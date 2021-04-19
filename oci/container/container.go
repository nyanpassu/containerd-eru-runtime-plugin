package container

import (
	"errors"

	"github.com/nyanpassu/containerd-eru-runtime-plugin/oci/common"
	"github.com/nyanpassu/containerd-eru-runtime-plugin/oci/meta"
	"github.com/nyanpassu/containerd-eru-runtime-plugin/oci/systemd"
)

// Container .
type Container interface {
	Init() (int, error)
	Start() error
	State() (common.ContainerState, error)
	Kill() error
	Delete() error
	Pause() error
	Resume() error
}

func newContainer(
	id string,
	meta meta.Meta,
	statusMeta meta.Status,
	systemd systemd.UnitManager,
) *container {
	return &container{
		id:         id,
		meta:       meta,
		statusMeta: statusMeta,
		systemd:    systemd,
	}
}

type container struct {
	id         string
	process    Process
	meta       meta.Meta
	statusMeta meta.Status
	systemd    systemd.UnitManager
}

func (c *container) Init() (int, error) {
	process, err := createInitProcess(c.id)
	if err != nil {
		return 0, err
	}
	c.process = process
	return process.PID(), nil
}

func (c *container) Start() error {
	var (
		err  error
		unit systemd.Unit
	)
	if unit, err = c.getSystemdUnit(); err != nil {
		return err
	}
	return unit.Start()
}

func (c *container) State() (common.ContainerState, error) {
	status, err := c.statusMeta.GetStatus(c.id)
	if err != nil {
		return common.ContainerState{}, err
	}
	return common.ContainerState{
		Version: common.OCIVersion,
		ID:      c.id,
		Status:  status,
	}, nil
}

func (c *container) Kill() error {
	var (
		err  error
		unit systemd.Unit
	)
	if unit, err = c.getSystemdUnit(); err != nil {
		return err
	}
	if err = unit.Stop(); err != nil {
		return err
	}

	if err := c.closeInitProcess(); err != nil {
		return err
	}
	return nil
}

func (c *container) Delete() error {
	var (
		err    error
		unit   systemd.Unit
		exists bool
	)
	if err = c.closeInitProcess(); err != nil {
		return err
	}
	if unit, exists, err = c.systemd.GetUnit(c.id); err != nil {
		return err
	}
	if exists {
		if err = unit.Delete(); err != nil {
			return err
		}
	}
	if err := c.meta.DeleteContainer(c.id); err != nil {
		return err
	}
	return nil
}

func (c *container) Pause() error {
	var (
		err  error
		unit systemd.Unit
	)
	if unit, err = c.getSystemdUnit(); err != nil {
		return err
	}
	return unit.Stop()
}

func (c *container) Resume() error {
	var (
		err  error
		unit systemd.Unit
	)
	if unit, err = c.getSystemdUnit(); err != nil {
		return err
	}
	return unit.Start()
}

func (c *container) getSystemdUnit() (systemd.Unit, error) {
	var (
		err    error
		unit   systemd.Unit
		exists bool
	)
	if unit, exists, err = c.systemd.GetUnit(c.id); err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("systemd unit file not exists")
	}
	return unit, nil
}

func (c *container) closeInitProcess() error {
	if c.process != nil {
		if err := c.process.Close(); err != nil {
			return err
		}
		c.process = nil
	}
	return nil
}
