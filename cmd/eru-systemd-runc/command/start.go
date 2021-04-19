package command

import (
	cli "github.com/urfave/cli"

	"github.com/nyanpassu/containerd-eru-runtime-plugin/oci/container"
)

// Start .
var Start = cli.Command{
	Name:  "start",
	Usage: "executes the user defined process in a created container",
	ArgsUsage: `<container-id>

Where "<container-id>" is your name for the instance of the container that you
are starting. The name you provide for the container instance must be unique on
your host.`,
	Description: `The start command executes the user defined process in a created container.`,
	Action: func(context *cli.Context) error {
		var (
			err error
			c   container.Container
		)
		if c, err = getContainer(context); err != nil {
			return err
		}
		if err = c.Start(); err != nil {
			return err
		}
		return nil
	},
}
