package command

import (
	"encoding/json"
	"os"

	cli "github.com/urfave/cli"

	"github.com/nyanpassu/containerd-eru-runtime-plugin/oci/common"
	"github.com/nyanpassu/containerd-eru-runtime-plugin/oci/container"
)

// State .
var State = cli.Command{
	Name:  "state",
	Usage: "output the state of a container",
	ArgsUsage: `<container-id>

Where "<container-id>" is your name for the instance of the container.`,
	Description: `The state command outputs current state information for the
instance of a container.`,
	Action: func(context *cli.Context) error {
		var (
			err error
			c   container.Container
			s   common.ContainerState
		)
		if c, err = getContainer(context); err != nil {
			return err
		}
		if s, err = c.State(); err != nil {
			return err
		}
		data, err := json.MarshalIndent(s, "", "  ")
		if err != nil {
			return err
		}
		os.Stdout.Write(data)
		return nil
	},
}
