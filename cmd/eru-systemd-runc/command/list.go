package command

import (
	cli "github.com/urfave/cli"

	"github.com/nyanpassu/containerd-eru-runtime-plugin/oci/common"
)

const formatOptions = `table or json`

// List .
var List = cli.Command{
	Name:  "list",
	Usage: "lists containers started by runc with the given root",
	ArgsUsage: `

Where the given root is specified via the global option "--root"
(default: "/run/runc").

EXAMPLE 1:
To list containers created via the default "--root":
       # runc list

EXAMPLE 2:
To list containers created using a non-default value for "--root":
       # runc --root value list`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "format, f",
			Value: "table",
			Usage: `select one of: ` + formatOptions,
		},
		&cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "display only container IDs",
		},
	},
	Action: func(context *cli.Context) error {
		return common.ErrNotImplemented
	},
}
