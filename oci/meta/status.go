package meta

import (
	"github.com/nyanpassu/containerd-eru-runtime-plugin/oci/common"
	"github.com/nyanpassu/containerd-eru-runtime-plugin/oci/systemd"
	"github.com/pkg/errors"
)

// Status .
type Status interface {
	GetStatus(id string) (string, error)
}

// NewStatus .
func NewStatus(unitManager systemd.UnitManager) Status {
	return &status{unitManager: unitManager}
}

type status struct {
	unitManager systemd.UnitManager
}

func (s *status) GetStatus(id string) (string, error) {
	unit, exists, err := s.unitManager.GetUnit(id)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", errors.New("unit file not exists")
	}
	status, err := unit.Status()
	if err != nil {
		return "", err
	}
	switch status {
	case systemd.StatusCreated:
		return common.StatusContainerCreated, nil
	case systemd.StatusError:
		return common.StatusContainerPaused, nil
	case systemd.StatusRunning:
		return common.StatusContainerRunning, nil
	}
	return common.StatusContainerStopped, nil
}
