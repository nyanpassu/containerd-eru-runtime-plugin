package command

import (
	cli "github.com/urfave/cli"

	"github.com/nyanpassu/containerd-eru-runtime-plugin/oci/common"
)

// Update .
var Update = cli.Command{
	Name:      "update",
	Usage:     "update container resource constraints",
	ArgsUsage: `<container-id>`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "resources, r",
			Value: "",
			Usage: `path to the file containing the resources to update or '-' to read from the standard input

The accepted format is as follow (unchanged values can be omitted):

{
  "memory": {
    "limit": 0,
    "reservation": 0,
    "swap": 0,
    "kernel": 0,
    "kernelTCP": 0
  },
  "cpu": {
    "shares": 0,
    "quota": 0,
    "period": 0,
    "realtimeRuntime": 0,
    "realtimePeriod": 0,
    "cpus": "",
    "mems": ""
  },
  "blockIO": {
    "weight": 0
  }
}

Note: if data is to be read from a file or the standard input, all
other options are ignored.
`,
		},
		&cli.IntFlag{
			Name:  "blkio-weight",
			Usage: "Specifies per cgroup weight, range is from 10 to 1000",
		},
		&cli.StringFlag{
			Name:  "cpu-period",
			Usage: "CPU CFS period to be used for hardcapping (in usecs). 0 to use system default",
		},
		&cli.StringFlag{
			Name:  "cpu-quota",
			Usage: "CPU CFS hardcap limit (in usecs). Allowed cpu time in a given period",
		},
		&cli.StringFlag{
			Name:  "cpu-share",
			Usage: "CPU shares (relative weight vs. other containers)",
		},
		&cli.StringFlag{
			Name:  "cpu-rt-period",
			Usage: "CPU realtime period to be used for hardcapping (in usecs). 0 to use system default",
		},
		&cli.StringFlag{
			Name:  "cpu-rt-runtime",
			Usage: "CPU realtime hardcap limit (in usecs). Allowed cpu time in a given period",
		},
		&cli.StringFlag{
			Name:  "cpuset-cpus",
			Usage: "CPU(s) to use",
		},
		&cli.StringFlag{
			Name:  "cpuset-mems",
			Usage: "Memory node(s) to use",
		},
		&cli.StringFlag{
			Name:  "kernel-memory",
			Usage: "Kernel memory limit (in bytes)",
		},
		&cli.StringFlag{
			Name:  "kernel-memory-tcp",
			Usage: "Kernel memory limit (in bytes) for tcp buffer",
		},
		&cli.StringFlag{
			Name:  "memory",
			Usage: "Memory limit (in bytes)",
		},
		&cli.StringFlag{
			Name:  "memory-reservation",
			Usage: "Memory reservation or soft_limit (in bytes)",
		},
		&cli.StringFlag{
			Name:  "memory-swap",
			Usage: "Total memory usage (memory + swap); set '-1' to enable unlimited swap",
		},
		&cli.IntFlag{
			Name:  "pids-limit",
			Usage: "Maximum number of pids allowed in the container",
		},
		&cli.StringFlag{
			Name:  "l3-cache-schema",
			Usage: "The string of Intel RDT/CAT L3 cache schema",
		},
		&cli.StringFlag{
			Name:  "mem-bw-schema",
			Usage: "The string of Intel RDT/MBA memory bandwidth schema",
		},
	},
	Action: func(context *cli.Context) error {
		return common.ErrNotImplemented
	},
}
