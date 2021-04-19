package shim

import (
	"context"

	"github.com/containerd/containerd/runtime"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
)

type shimRecovery struct {
	id        string
	namespace string
}

// ID of the process
func (s *shimRecovery) ID() string {
	return s.id
}

// State returns the process state
func (s *shimRecovery) State(context.Context) (runtime.State, error) {
	return runtime.State{}, errors.New("")
}

// Kill signals a container
func (s *shimRecovery) Kill(context.Context, uint32, bool) error {
	return errors.New("")
}

// Pty resizes the processes pty/console
func (s *shimRecovery) ResizePty(context.Context, runtime.ConsoleSize) error {
	return errors.New("")
}

// CloseStdin closes the processes stdin
func (s *shimRecovery) CloseIO(context.Context) error {
	return errors.New("")
}

// Start the container's user defined process
func (s *shimRecovery) Start(context.Context) error {
	return errors.New("")
}

// Wait for the process to exit
func (s *shimRecovery) Wait(context.Context) (*runtime.Exit, error) {
	return nil, errors.New("")
}

// Delete deletes the process
func (s *shimRecovery) Delete(context.Context) (*runtime.Exit, error) {
	return nil, errors.New("")
}

// PID of the process
func (s *shimRecovery) PID() uint32 {
	return 0
}

// Namespace that the task exists in
func (s *shimRecovery) Namespace() string {
	return s.namespace
}

// Pause pauses the container process
func (s *shimRecovery) Pause(context.Context) error {
	return errors.New("")
}

// Resume unpauses the container process
func (s *shimRecovery) Resume(context.Context) error {
	return errors.New("")
}

// Exec adds a process into the container
func (s *shimRecovery) Exec(context.Context, string, runtime.ExecOpts) (runtime.Process, error) {
	return nil, errors.New("")
}

// Pids returns all pids
func (s *shimRecovery) Pids(context.Context) ([]runtime.ProcessInfo, error) {
	return nil, errors.New("")
}

// Checkpoint checkpoints a container to an image with live system data
func (s *shimRecovery) Checkpoint(context.Context, string, *types.Any) error {
	return errors.New("")
}

// Update sets the provided resources to a running task
func (s *shimRecovery) Update(context.Context, *types.Any) error {
	return errors.New("")
}

// Process returns a process within the task for the provided id
func (s *shimRecovery) Process(context.Context, string) (runtime.Process, error) {
	return nil, errors.New("")
}

// Stats returns runtime specific metrics for a task
func (s *shimRecovery) Stats(context.Context) (*types.Any, error) {
	return nil, errors.New("")
}
