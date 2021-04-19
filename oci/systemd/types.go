package systemd

// UnitFile .
type UnitFile struct {
	ContainerID string
	Args        []string
}

// UnitStatus .
type UnitStatus int

const (
	// StatusCreated .
	StatusCreated UnitStatus = 1 << iota
	// StatusRunning .
	StatusRunning
	// StatusError .
	StatusError
)
