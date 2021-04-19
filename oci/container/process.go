package container

// Process .
type Process interface {
	PID() int
	Close() error
}

type noProcess struct{}

func (p noProcess) PID() int     { return 0 }
func (p noProcess) Close() error { return nil }

func createInitProcess(id string) (Process, error) {
	return noProcess{}, nil
}
