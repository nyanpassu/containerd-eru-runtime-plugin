package systemd

import (
	"fmt"
	"os"
	"strings"

	"github.com/nyanpassu/systemd-oci/utils"
)

// FileManager .
type FileManager interface {
	GenerateSystemdUnitFile(id string, unitName string, args []string) error
	UnitFileExists(unitName string) (bool, error)
	RemoveSystemdUnitFile(unitName string) error
}

// NewFileManager .
func NewFileManager(systemdFileBasePath string) FileManager {
	return &systemdFileManager{basePath: systemdFileBasePath}
}

type systemdFileManager struct {
	basePath string
}

func (m *systemdFileManager) GenerateSystemdUnitFile(id string, unitName string, args []string) error {
	f, err := os.OpenFile(m.getAbsPath(unitName), os.O_RDWR|os.O_CREATE|os.O_SYNC, 0644)
	if err != nil {
		return err
	}

	if _, err := f.WriteString(
		fmt.Sprintf(`
[Unit]
Description=%s

[Service]
Type=forking
ExecStart=eru-systemd-rununit %s %s

[Install]
WantedBy=multi-user.target
		`,
			unitName,
			id,
			strings.Join(args, " ")),
	); err != nil {
		return err
	}
	return f.Close()
}

func (m *systemdFileManager) RemoveSystemdUnitFile(unitName string) error {
	exists, err := m.UnitFileExists(unitName)
	if err != nil {
		return err
	}
	if exists {
		return os.Remove(m.getAbsPath(unitName))
	}
	return nil
}

func (m *systemdFileManager) UnitFileExists(unitName string) (bool, error) {
	return utils.FileExists(m.getAbsPath(unitName))
}

func (m *systemdFileManager) getAbsPath(unitName string) string {
	return m.basePath + "/" + unitName + ".service"
}
