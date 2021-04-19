package command

import (
	cli "github.com/urfave/cli"

	"github.com/nyanpassu/containerd-eru-runtime-plugin/oci/common"
)

// Checkpoint .
var Checkpoint = cli.Command{
	Name:  "checkpoint",
	Usage: "checkpoint a running container",
	ArgsUsage: `<container-id>

Where "<container-id>" is the name for the instance of the container to be
checkpointed.`,
	Description: `The checkpoint command saves the state of the container instance.`,
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "image-path", Value: "", Usage: "path for saving criu image files"},
		&cli.StringFlag{Name: "work-path", Value: "", Usage: "path for saving work files and logs"},
		&cli.StringFlag{Name: "parent-path", Value: "", Usage: "path for previous criu image files in pre-dump"},
		&cli.BoolFlag{Name: "leave-running", Usage: "leave the process running after checkpointing"},
		&cli.BoolFlag{Name: "tcp-established", Usage: "allow open tcp connections"},
		&cli.BoolFlag{Name: "ext-unix-sk", Usage: "allow external unix sockets"},
		&cli.BoolFlag{Name: "shell-job", Usage: "allow shell jobs"},
		&cli.BoolFlag{Name: "lazy-pages", Usage: "use userfaultfd to lazily restore memory pages"},
		&cli.IntFlag{Name: "status-fd", Value: -1, Usage: "criu writes \\0 to this FD once lazy-pages is ready"},
		&cli.StringFlag{Name: "page-server", Value: "", Usage: "ADDRESS:PORT of the page server"},
		&cli.BoolFlag{Name: "file-locks", Usage: "handle file locks, for safety"},
		&cli.BoolFlag{Name: "pre-dump", Usage: "dump container's memory information only, leave the container running after this"},
		&cli.StringFlag{Name: "manage-cgroups-mode", Value: "", Usage: "cgroups mode: 'soft' (default), 'full' and 'strict'"},
		&cli.StringSliceFlag{Name: "empty-ns", Usage: "create a namespace, but don't restore its properties"},
		&cli.BoolFlag{Name: "auto-dedup", Usage: "enable auto deduplication of memory images"},
	},
	Action: func(context *cli.Context) error {
		return common.ErrNotImplemented
	},
}
