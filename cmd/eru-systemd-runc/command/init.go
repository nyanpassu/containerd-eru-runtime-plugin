package command

import (
	cli "github.com/urfave/cli"

	"github.com/nyanpassu/containerd-eru-runtime-plugin/oci/common"
)

// Init .
var Init = cli.Command{
	Name:  "init",
	Usage: `initialize the namespaces and launch the process (do not call it outside of runc)`,
	Action: func(context *cli.Context) error {
		return common.ErrNotImplemented
	},
}
