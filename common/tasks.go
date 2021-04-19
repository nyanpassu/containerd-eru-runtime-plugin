package common

import (
	"context"

	"github.com/containerd/containerd/runtime"
)

type Tasks interface {
	Add(context.Context, runtime.Task) error
	Delete(context.Context, runtime.Task)
}
