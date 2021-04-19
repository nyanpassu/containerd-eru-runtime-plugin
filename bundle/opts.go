package bundle

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/containerd/containerd/runtime"
)

const optsFileName = "opts"

func (b *Bundle) SaveOpts(ctx context.Context, opts runtime.CreateOpts) error {
	content, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(filepath.Join(b.Path, optsFileName), content, 0666); err != nil {
		return err
	}
	return nil
}

func (b *Bundle) LoadOpts(ctx context.Context) (runtime.CreateOpts, error) {
	var (
		opts runtime.CreateOpts
	)
	f, err := os.Open(filepath.Join(b.Path, optsFileName))
	if err != nil {
		return runtime.CreateOpts{}, err
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return runtime.CreateOpts{}, err
	}
	if err = json.Unmarshal(content, &opts); err != nil {
		return runtime.CreateOpts{}, err
	}
	return opts, nil
}
