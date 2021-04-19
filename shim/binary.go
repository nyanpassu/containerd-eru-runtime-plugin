/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package shim

import (
	"bytes"
	"context"
	"io"
	"os"
	gruntime "runtime"
	"strings"

	"github.com/containerd/containerd/events/exchange"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/pkg/timeout"
	"github.com/containerd/containerd/runtime"
	client "github.com/containerd/containerd/runtime/v2/shim"
	"github.com/containerd/containerd/runtime/v2/task"
	"github.com/containerd/ttrpc"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/nyanpassu/containerd-eru-runtime-plugin/bundle"
	"github.com/nyanpassu/containerd-eru-runtime-plugin/common"
)

func NewTask(
	ctx context.Context,
	bundle *bundle.Bundle,
	opts runtime.CreateOpts,
	containerdAddress, containerdTTRPCAddress string,
	events *exchange.Exchange,
	tasks common.Tasks,
) (runtime.Task, error) {
	topts := opts.TaskOptions
	if topts == nil {
		topts = opts.RuntimeOptions
	}
	b := shimBinary(ctx, bundle, opts.Runtime, containerdAddress, containerdTTRPCAddress, events, tasks)
	shim, err := b.start(ctx, topts, func() {
		log.G(ctx).WithField("id", bundle.ID).Infof("shim disconnected")
	})
	if err != nil {
		log.G(ctx).WithField("id", bundle.ID).Infof("shim start failed")
		return nil, err
	}
	t, err := shim.create(ctx, opts)
	if err != nil {
		log.G(ctx).WithField("id", bundle.ID).Infof("shim create failed")
		func() {
			dctx, cancel := timeout.WithContext(context.Background(), cleanupTimeout)
			defer cancel()
			_, errShim := shim.Delete(dctx)
			if errShim != nil {
				shim.Shutdown(dctx)
				shim.Close()
			}
		}()
		return nil, err
	}
	return t, nil
}

func shimBinary(
	ctx context.Context,
	bundle *bundle.Bundle,
	runtime, containerdAddress, containerdTTRPCAddress string,
	events *exchange.Exchange,
	tasks common.Tasks,
) *binary {
	return &binary{
		bundle:                 bundle,
		runtime:                runtime,
		containerdAddress:      containerdAddress,
		containerdTTRPCAddress: containerdTTRPCAddress,
		events:                 events,
		tasks:                  tasks,
	}
}

type binary struct {
	runtime                string
	containerdAddress      string
	containerdTTRPCAddress string
	bundle                 *bundle.Bundle
	events                 *exchange.Exchange
	tasks                  common.Tasks
}

func (b *binary) start(ctx context.Context, opts *types.Any, onClose func()) (_ *shim, err error) {
	log.G(ctx).Infof("binary start runtime = %s, id = %s", b.runtime, b.bundle.ID)

	args := []string{"-id", b.bundle.ID}
	if logrus.GetLevel() == logrus.DebugLevel {
		args = append(args, "-debug")
	}
	args = append(args, "start")

	cmd, err := client.Command(
		ctx,
		b.runtime,
		b.containerdAddress,
		b.containerdTTRPCAddress,
		b.bundle.Path,
		opts,
		args...,
	)
	if err != nil {
		return nil, err
	}
	// Windows needs a namespace when openShimLog
	ns, _ := namespaces.Namespace(ctx)
	shimCtx, cancelShimLog := context.WithCancel(namespaces.WithNamespace(context.Background(), ns))
	defer func() {
		if err != nil {
			cancelShimLog()
		}
	}()
	f, err := openShimLog(shimCtx, b.bundle, client.AnonDialer)
	if err != nil {
		return nil, errors.Wrap(err, "open shim log pipe")
	}
	defer func() {
		if err != nil {
			f.Close()
		}
	}()
	// open the log pipe and block until the writer is ready
	// this helps with synchronization of the shim
	// copy the shim's logs to containerd's output
	go func() {
		defer f.Close()
		_, err := io.Copy(os.Stderr, f)
		err = checkCopyShimLogError(ctx, err)
		if err != nil {
			log.G(ctx).WithError(err).Error("binary error copy shim log")
		}
	}()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "%s", out)
	}
	address := strings.TrimSpace(string(out))
	log.G(ctx).WithField("address", address).Info("connect shim")
	conn, err := client.Connect(address, client.AnonDialer)
	if err != nil {
		return nil, err
	}
	onCloseWithShimLog := func() {
		onClose()
		cancelShimLog()
		f.Close()
	}
	client := ttrpc.NewClient(conn, ttrpc.WithOnClose(onCloseWithShimLog))
	return &shim{
		bundle: b.bundle,
		client: client,
		task:   task.NewTaskClient(client),
		events: b.events,
		tasks:  b.tasks,
	}, nil
}

func (b *binary) delete(ctx context.Context) (*runtime.Exit, error) {
	log.G(ctx).Info("cleaning up dead shim")

	// Windows cannot delete the current working directory while an
	// executable is in use with it. For the cleanup case we invoke with the
	// default work dir and forward the bundle path on the cmdline.
	var bundlePath string
	if gruntime.GOOS != "windows" {
		bundlePath = b.bundle.Path
	}

	cmd, err := client.Command(ctx,
		b.runtime,
		b.containerdAddress,
		b.containerdTTRPCAddress,
		bundlePath,
		nil,
		"-id", b.bundle.ID,
		"-bundle", b.bundle.Path,
		"delete")
	if err != nil {
		return nil, err
	}
	var (
		out  = bytes.NewBuffer(nil)
		errb = bytes.NewBuffer(nil)
	)
	cmd.Stdout = out
	cmd.Stderr = errb
	if err := cmd.Run(); err != nil {
		return nil, errors.Wrapf(err, "%s", errb.String())
	}
	s := errb.String()
	if s != "" {
		log.G(ctx).Warnf("cleanup warnings %s", s)
	}
	var response task.DeleteResponse
	if err := response.Unmarshal(out.Bytes()); err != nil {
		return nil, err
	}
	if err := b.bundle.Delete(); err != nil {
		return nil, err
	}
	return &runtime.Exit{
		Status:    response.ExitStatus,
		Timestamp: response.ExitedAt,
		Pid:       response.Pid,
	}, nil
}
