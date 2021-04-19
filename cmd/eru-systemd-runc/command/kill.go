package command

import (
	cli "github.com/urfave/cli"

	"github.com/nyanpassu/containerd-eru-runtime-plugin/oci/common"
)

// Kill .
var Kill = cli.Command{
	Name:  "kill",
	Usage: "kill sends the specified signal (default: SIGTERM) to the container's init process",
	ArgsUsage: `<container-id> [signal]

Where "<container-id>" is the name for the instance of the container and
"[signal]" is the signal to be sent to the init process.

EXAMPLE:
For example, if the container id is "ubuntu01" the following will send a "KILL"
signal to the init process of the "ubuntu01" container:

       # runc kill ubuntu01 KILL`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "all, a",
			Usage: "send the specified signal to all processes inside the container",
		},
	},
	Action: func(context *cli.Context) error {
		return common.ErrNotImplemented
	},
}
