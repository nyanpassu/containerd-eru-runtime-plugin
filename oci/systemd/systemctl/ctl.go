package systemctl

import (
	"os/exec"
)

// Start .
func Start(name string) error {
	return exec.Command("systemctl", "start", name).Run()
}

// Stop .
func Stop(name string) error {
	return exec.Command("systemctl", "stop", name).Run()
}

// Enable .
func Enable(name string) error {
	return exec.Command("systemctl", "enable", name).Run()
}

// Disable .
func Disable(name string) error {
	return exec.Command("systemctl", "disable", name).Run()
}
