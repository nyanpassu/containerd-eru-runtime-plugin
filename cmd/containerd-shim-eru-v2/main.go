package main

import (
	"github.com/containerd/containerd/runtime/v2/shim"

	eruShim "github.com/nyanpassu/containerd-eru-runtime-plugin/v2/shim"
)

func main() {
	// init and execute the shim
	shim.Run("io.containerd.systemd.v1", eruShim.New)
}
