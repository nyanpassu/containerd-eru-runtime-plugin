package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli"

	"github.com/nyanpassu/containerd-eru-runtime-plugin/cmd/eru-systemd-runc/command"
)

func run(cli *cli.Context) error {
	return nil
}

func main() {
	app := &cli.App{
		Name:    "Barrel",
		Usage:   "Dockerd with calico fixed IP feature",
		Action:  run,
		Version: "demo",
		Commands: []cli.Command{
			command.Checkpoint,
			command.Create,
			command.Delete,
			command.Events,
			command.Exec,
			command.Init,
			command.Kill,
			command.List,
			command.Pause,
			command.Ps,
			command.Restore,
			command.Resume,
			command.Run,
			command.Spec,
			command.Start,
			command.State,
			command.Update,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
