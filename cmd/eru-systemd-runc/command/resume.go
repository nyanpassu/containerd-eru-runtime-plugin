package command

import (
	cli "github.com/urfave/cli"

	"github.com/nyanpassu/containerd-eru-runtime-plugin/oci/common"
)

// Resume .
var Resume = cli.Command{
	Name:  "resume",
	Usage: "resumes all processes that have been previously paused",
	ArgsUsage: `<container-id>

Where "<container-id>" is the name for the instance of the container to be
resumed.`,
	Description: `The resume command resumes all processes in the instance of the container.

Use runc list to identify instances of containers and their current status.`,
	Action: func(context *cli.Context) error {
		return common.ErrNotImplemented
	},
}
