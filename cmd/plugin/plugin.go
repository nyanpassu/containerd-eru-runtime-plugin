package main

import (
	"context"
	"log"

	containerdRuntime "github.com/containerd/containerd/runtime"
	"github.com/containerd/containerd/runtime/proxy"

	"github.com/nyanpassu/containerd-eru-runtime-plugin/runtime"
)

func main() {
	rt := newRuntime()
	proxy.RunPlatformRuntimeService(rt)
}

func newRuntime() containerdRuntime.PlatformRuntime {
	root := ""
	state := ""
	containerdAddress := ""
	containerdTTRPCAddress := ""
	rt, err := runtime.New(context.Background(), root, state, containerdAddress, containerdTTRPCAddress, nil, nil)
	if err != nil {
		log.Fatalln(err)
	}
	return rt
}
