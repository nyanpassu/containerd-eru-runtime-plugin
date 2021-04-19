package command

import (
	"time"

	cli "github.com/urfave/cli"

	"github.com/nyanpassu/containerd-eru-runtime-plugin/oci/common"
)

// Events .
var Events = cli.Command{
	Name:  "events",
	Usage: "display container events such as OOM notifications, cpu, memory, and IO usage statistics",
	ArgsUsage: `<container-id>

Where "<container-id>" is the name for the instance of the container.`,
	Description: `The events command displays information about the container. By default the
information is displayed once every 5 seconds.`,
	Flags: []cli.Flag{
		&cli.DurationFlag{Name: "interval", Value: 5 * time.Second, Usage: "set the stats collection interval"},
		&cli.BoolFlag{Name: "stats", Usage: "display the container's stats then exit"},
	},
	Action: func(context *cli.Context) error {
		return common.ErrNotImplemented
	},
}
