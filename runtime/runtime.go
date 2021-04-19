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

package runtime

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/events/exchange"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/runtime"

	"github.com/nyanpassu/containerd-eru-runtime-plugin/bundle"
	"github.com/nyanpassu/containerd-eru-runtime-plugin/common"
	"github.com/nyanpassu/containerd-eru-runtime-plugin/shim"
)

// Config for the v2 runtime
type Config struct {
	// Supported platforms
	Platforms []string `toml:"platforms"`
}

// New task manager for v2 shims
func New(
	ctx context.Context,
	root, state, containerdAddress, containerdTTRPCAddress string,
	events *exchange.Exchange,
	cs containers.Store,
) (runtime.PlatformRuntime, error) {
	for _, d := range []string{root, state} {
		if err := os.MkdirAll(d, 0711); err != nil {
			return nil, err
		}
	}
	m := &TaskManager{
		root:                   root,
		state:                  state,
		containerdAddress:      containerdAddress,
		containerdTTRPCAddress: containerdTTRPCAddress,
		tasks:                  newTaskList(),
		events:                 events,
		containers:             cs,
	}
	if err := m.loadExistingTasks(ctx); err != nil {
		return nil, err
	}
	return m, nil
}

// TaskManager manages v2 shim's and their tasks
type TaskManager struct {
	root                   string
	state                  string
	containerdAddress      string
	containerdTTRPCAddress string

	tasks      *taskList
	events     *exchange.Exchange
	containers containers.Store
}

// ID of the task manager
func (m *TaskManager) ID() string {
	return common.RuntimeName
}

// Create a new task
func (m *TaskManager) Create(ctx context.Context, id string, opts runtime.CreateOpts) (runtime.Task, error) {
	bundle, err := bundle.NewBundle(ctx, m.root, m.state, id, opts.Spec.Value)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			bundle.Delete()
		}
	}()
	err = bundle.SaveOpts(ctx, opts)
	if err != nil {
		return nil, err
	}
	topts := opts.TaskOptions
	if topts == nil {
		topts = opts.RuntimeOptions
	}

	t, err := shim.NewTask(ctx, bundle, opts, m.containerdAddress, m.containerdTTRPCAddress, m.events, m.tasks)
	if err != nil {
		return nil, err
	}
	m.tasks.Add(ctx, t)
	return t, nil
}

// Get a specific task
func (m *TaskManager) Get(ctx context.Context, id string) (runtime.Task, error) {
	return m.tasks.Get(ctx, id)
}

// Tasks lists all tasks
func (m *TaskManager) Tasks(ctx context.Context, all bool) ([]runtime.Task, error) {
	return m.tasks.GetAll(ctx, all)
}

func (m *TaskManager) loadExistingTasks(ctx context.Context) error {
	nsDirs, err := ioutil.ReadDir(m.state)
	if err != nil {
		return err
	}
	for _, nsd := range nsDirs {
		if !nsd.IsDir() {
			continue
		}
		ns := nsd.Name()
		// skip hidden directories
		if len(ns) > 0 && ns[0] == '.' {
			continue
		}
		log.G(ctx).WithField("namespace", ns).Debug("loading tasks in namespace")
		if err := m.loadTasks(namespaces.WithNamespace(ctx, ns)); err != nil {
			log.G(ctx).WithField("namespace", ns).WithError(err).Error("loading tasks in namespace")
			continue
		}
		if err := m.cleanupWorkDirs(namespaces.WithNamespace(ctx, ns)); err != nil {
			log.G(ctx).WithField("namespace", ns).WithError(err).Error("cleanup working directory in namespace")
			continue
		}
	}
	return nil
}

func (m *TaskManager) loadTasks(ctx context.Context) error {
	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil {
		return err
	}
	shimDirs, err := ioutil.ReadDir(filepath.Join(m.state, ns))
	if err != nil {
		return err
	}
	for _, sd := range shimDirs {
		if !sd.IsDir() {
			continue
		}
		id := sd.Name()
		// skip hidden directories
		if len(id) > 0 && id[0] == '.' {
			continue
		}
		bundle, err := bundle.LoadBundle(ctx, m.state, id)
		if err != nil {
			// fine to return error here, it is a programmer error if the context
			// does not have a namespace
			return err
		}
		// fast path
		bf, err := ioutil.ReadDir(bundle.Path)
		if err != nil {
			bundle.Delete()
			log.G(ctx).WithError(err).Errorf("fast path read bundle path for %s", bundle.Path)
			continue
		}
		if len(bf) == 0 {
			bundle.Delete()
			continue
		}
		err = m.checkContainerExist(ctx, id)
		if err != nil {
			log.G(ctx).WithError(err).Errorf("loading container %s", id)
			if err := mount.UnmountAll(filepath.Join(bundle.Path, "rootfs"), 0); err != nil {
				log.G(ctx).WithError(err).Errorf("forceful unmount of rootfs %s", id)
			}
			bundle.Delete()
			continue
		}
		m.tasks.Add(ctx, shim.LoadTaskFromBundle(ctx, bundle, m.containerdAddress, m.containerdTTRPCAddress, m.events, m.tasks, func() {}))
	}
	return nil
}

func (m *TaskManager) checkContainerExist(ctx context.Context, id string) error {
	_, err := m.containers.Get(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *TaskManager) addNewTaskFromSupplier(ctx context.Context, supplier func(onClose func(runtime.Task)) (runtime.Task, error)) error {
	task, err := supplier(func(task runtime.Task) {
		log.G(ctx).WithField("id", task.ID()).Info("shim disconnected")
		m.tasks.Delete(ctx, task)
	})
	if err != nil {
		return err
	}
	m.tasks.Add(ctx, task)
	return nil
}

func (m *TaskManager) cleanupWorkDirs(ctx context.Context) error {
	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil {
		return err
	}
	dirs, err := ioutil.ReadDir(filepath.Join(m.root, ns))
	if err != nil {
		return err
	}
	for _, d := range dirs {
		// if the task was not loaded, cleanup and empty working directory
		// this can happen on a reboot where /run for the bundle state is cleaned up
		// but that persistent working dir is left
		if _, err := m.tasks.Get(ctx, d.Name()); err != nil {
			path := filepath.Join(m.root, ns, d.Name())
			if err := os.RemoveAll(path); err != nil {
				log.G(ctx).WithError(err).Errorf("cleanup working dir %s", path)
			}
		}
	}
	return nil
}
