package cmds

import (
	"github.com/urfave/cli"
	"github.com/xiaods/k8e/pkg/version"
)

const CertCommand = "certificate"

var (
	ServicesList     cli.StringSlice
	CertCommandFlags = []cli.Flag{
		DebugFlag,
		ConfigFlag,
		LogFile,
		AlsoLogToStderr,
		cli.StringFlag{
			Name:        "data-dir,d",
			Usage:       "(data) Folder to hold state default /var/lib/" + version.Program + " or ${HOME}/." + version.Program + "/ if not root",
			Destination: &ServerConfig.DataDir,
		},
		cli.StringSliceFlag{
			Name:  "service,s",
			Usage: "List of services to rotate certificates for. Options include (admin, api-server, controller-manager, scheduler, " + version.Program + "-controller, " + version.Program + "-server, cloud-controller, etcd, auth-proxy, kubelet, kube-proxy)",
			Value: &ServicesList,
		},
	}
)

func NewCertCommand(subcommands []cli.Command) cli.Command {
	return cli.Command{
		Name:            CertCommand,
		Usage:           "Certificates management",
		SkipFlagParsing: false,
		SkipArgReorder:  true,
		Subcommands:     subcommands,
		Flags:           CertCommandFlags,
	}
}

func NewCertSubcommands(rotate func(ctx *cli.Context) error) []cli.Command {
	return []cli.Command{
		{
			Name:            "rotate",
			Usage:           "Certificate rotation",
			SkipFlagParsing: false,
			SkipArgReorder:  true,
			Action:          rotate,
			Flags:           CertCommandFlags,
		},
	}
}
