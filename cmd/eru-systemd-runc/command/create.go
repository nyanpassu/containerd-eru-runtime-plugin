package command

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	cli "github.com/urfave/cli"

	"github.com/nyanpassu/containerd-eru-runtime-plugin/oci/common"
	"github.com/nyanpassu/containerd-eru-runtime-plugin/oci/container"
	"github.com/nyanpassu/containerd-eru-runtime-plugin/oci/utils"
)

// Create .
var Create = cli.Command{
	Name:  "create",
	Usage: "create a container",
	ArgsUsage: `<container-id>

Where "<container-id>" is your name for the instance of the container that you
are starting. The name you provide for the container instance must be unique on
your host.`,
	Description: `The create command creates an instance of a container for a bundle. The bundle
is a directory with a specification file named "` + common.SpecConfig + `" and a root
filesystem.

The specification file includes an args parameter. The args parameter is used
to specify command(s) that get run when the container is started. To change the
command(s) that get executed on start, edit the args parameter of the spec. See
"runc spec --help" for more explanation.`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "bundle, b",
			Value: "",
			Usage: `path to the root of the bundle directory, defaults to the current directory`,
		},
		&cli.StringFlag{
			Name:  "console-socket",
			Value: "",
			Usage: "path to an AF_UNIX socket which will receive a file descriptor referencing the master end of the console's pseudoterminal",
		},
		&cli.StringFlag{
			Name:  "pid-file",
			Value: "",
			Usage: "specify the file to write the process id to",
		},
		&cli.BoolFlag{
			Name:  "no-pivot",
			Usage: "do not use pivot root to jail process inside rootfs.  This should be used whenever the rootfs is on top of a ramdisk",
		},
		&cli.BoolFlag{
			Name:  "no-new-keyring",
			Usage: "do not create a new session keyring for the container.  This will cause the container to inherit the calling processes session key",
		},
		&cli.IntFlag{
			Name:  "preserve-fds",
			Usage: "Pass N additional file descriptors to the container (stdio + $LISTEN_FDS + N in total)",
		},
	},
	Action: func(context *cli.Context) error {
		var (
			err         error
			c           container.Container
			initPID     int
			pidFilePath string
		)
		if pidFilePath, err = getPidFilePath(context); err != nil {
			return err
		}
		if c, err = createContainer(context); err != nil {
			return err
		}
		if initPID, err = c.Init(); err != nil {
			return err
		}
		if pidFilePath != "" {
			if err := utils.WritePidFile(pidFilePath, initPID); err != nil {
				if errDel := c.Delete(); errDel != nil {
					return errors.Wrap(err, errDel.Error())
				}
				return err
			}
		}
		os.Exit(0)
		return nil
	},
}

func getPidFilePath(context *cli.Context) (string, error) {
	pidFile := context.String("pid-file")
	if pidFile == "" {
		return "", nil
	}

	// convert pid-file to an absolute path so we can write to the right
	// file after chdir to bundle
	pidFile, err := filepath.Abs(pidFile)
	if err != nil {
		return "", err
	}
	return pidFile, nil
}
