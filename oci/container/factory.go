package container

import (
	specsGo "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"

	"github.com/nyanpassu/containerd-eru-runtime-plugin/oci/meta"
	"github.com/nyanpassu/containerd-eru-runtime-plugin/oci/systemd"
)

// Factory .
type Factory interface {
	CreateContainer(id string, spec *specsGo.Spec) (Container, error)
	GetContainer(id string) (Container, error)
}

// NewFactory .
func NewFactory(m meta.Meta, status meta.Status, systemd systemd.UnitManager) (Factory, error) {
	return &factory{
		meta:    m,
		status:  status,
		systemd: systemd,
	}, nil
}

type factory struct {
	meta    meta.Meta
	status  meta.Status
	systemd systemd.UnitManager
}

func (f *factory) CreateContainer(id string, spec *specsGo.Spec) (Container, error) {
	if err := f.meta.CreateContainer(meta.Container{
		ID: id,
	}); err != nil {
		return nil, err
	}
	if err := f.systemd.Create(systemd.UnitFile{
		ContainerID: id,
		Args:        spec.Process.Args,
	}); err != nil {
		if errDel := f.meta.DeleteContainer(id); errDel != nil {
			return nil, errors.Wrap(err, errDel.Error())
		}
		return nil, err
	}
	return newContainer(id, f.meta, f.status, f.systemd), nil
}

func (f *factory) GetContainer(id string) (Container, error) {
	_, err := f.meta.GetContainer(id)
	if err != nil {
		return nil, err
	}
	return newContainer(id, f.meta, f.status, f.systemd), nil
}
