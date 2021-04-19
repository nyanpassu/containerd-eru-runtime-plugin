package systemd

import "github.com/nyanpassu/systemd-oci/systemd/systemctl"

// UnitManager .
type UnitManager interface {
	GetUnit(string) (Unit, bool, error)
	Create(UnitFile) error
}

// NewUnitManager .
func NewUnitManager(fileManager FileManager) (UnitManager, error) {
	return &unitManager{fileManager: fileManager}, nil
}

type unitManager struct {
	fileManager FileManager
}

func (m *unitManager) GetUnit(id string) (Unit, bool, error) {
	unitName := generateSystemdUnitName(id)
	exists, err := m.fileManager.UnitFileExists(unitName)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	return NewUnit(id, unitName, m.fileManager), true, nil
}

func (m *unitManager) Create(unit UnitFile) error {
	unitName := generateSystemdUnitName(unit.ContainerID)
	if err := m.fileManager.GenerateSystemdUnitFile(unit.ContainerID, unitName, unit.Args); err != nil {
		return err
	}
	return systemctl.Enable(unitName)
}
