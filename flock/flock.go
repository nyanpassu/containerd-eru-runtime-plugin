package flock

import (
	"os"
	"syscall"
)

func AquireLock(bundlePath string) (func() error, error) {
	var (
		f   *os.File
		err error
	)
	if f, err = os.Create(bundlePath + "/shimlock"); err != nil {
		return nil, err
	}
	fd := int(f.Fd())
	if err = syscall.Flock(fd, syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		return nil, err
	}
	return func() error {
		return syscall.Flock(fd, syscall.LOCK_UN|syscall.LOCK_NB)
	}, nil
}
