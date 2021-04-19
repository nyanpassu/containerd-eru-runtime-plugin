package systemd

import (
	"os"
	"syscall"

	"github.com/nyanpassu/systemd-oci/systemd/systemctl"
)

// Unit .
type Unit interface {
	Status() (UnitStatus, error)
	Start() error
	Stop() error
	Delete() error
}

// NewUnit .
func NewUnit(id string, unitName string, fileManager FileManager) Unit {
	return &unit{
		id:          id,
		unitName:    unitName,
		fileManager: fileManager,
	}
}

type unit struct {
	id          string
	unitName    string
	fileManager FileManager
}

func (u *unit) Status() (UnitStatus, error) {
	pid, exists, err := u.getPid()
	if err != nil {
		return StatusError, err
	}
	if !exists {
		return StatusCreated, nil
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return StatusError, err
	}
	if err := process.Signal(syscall.Signal(0)); err != nil {
		return StatusError, nil
	}
	return StatusRunning, nil
}

func (u *unit) Start() error {
	return systemctl.Start(u.unitName)
}

func (u *unit) Stop() error {
	return systemctl.Stop(u.unitName)
}

func (u *unit) Delete() error {
	if err := systemctl.Stop(u.unitName); err != nil {
		return err
	}
	if err := systemctl.Disable(u.unitName); err != nil {
		return err
	}
	return u.fileManager.RemoveSystemdUnitFile(u.unitName)
}

func (u *unit) getPid() (int, bool, error) {
	return ReadPid(u.id)
}

func generateSystemdUnitName(id string) string {
	return "systemdoci-" + id
}
