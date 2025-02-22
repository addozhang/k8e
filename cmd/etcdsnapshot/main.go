package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/xiaods/k8e/pkg/cli/cmds"
	"github.com/xiaods/k8e/pkg/cli/etcdsnapshot"
	"github.com/xiaods/k8e/pkg/configfilearg"
)

func main() {
	app := cmds.NewApp()
	app.Commands = []cli.Command{
		cmds.NewEtcdSnapshotCommand(etcdsnapshot.Run,
			cmds.NewEtcdSnapshotSubcommands(
				etcdsnapshot.Delete,
				etcdsnapshot.List,
				etcdsnapshot.Prune,
				etcdsnapshot.Run),
		),
	}

	if err := app.Run(configfilearg.MustParse(os.Args)); err != nil {
		logrus.Fatal(err)
	}
}
