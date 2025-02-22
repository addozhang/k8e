package cmds

import (
	"github.com/urfave/cli"
)

const SecretsEncryptCommand = "secrets-encrypt"

var EncryptFlags = []cli.Flag{
	DataDirFlag,
	ServerToken,
}

func NewSecretsEncryptCommand(action func(*cli.Context) error, subcommands []cli.Command) cli.Command {
	return cli.Command{
		Name:            SecretsEncryptCommand,
		Usage:           "Control secrets encryption and keys rotation",
		SkipFlagParsing: false,
		SkipArgReorder:  true,
		Action:          action,
		Subcommands:     subcommands,
	}
}

func NewSecretsEncryptSubcommands(status, enable, disable, prepare, rotate, reencrypt func(ctx *cli.Context) error) []cli.Command {
	return []cli.Command{
		{
			Name:            "status",
			Usage:           "Print current status of secrets encryption",
			SkipFlagParsing: false,
			SkipArgReorder:  true,
			Action:          status,
			Flags:           EncryptFlags,
		},
		{
			Name:            "enable",
			Usage:           "Enable secrets encryption",
			SkipFlagParsing: false,
			SkipArgReorder:  true,
			Action:          enable,
			Flags:           EncryptFlags,
		},
		{
			Name:            "disable",
			Usage:           "Disable secrets encryption",
			SkipFlagParsing: false,
			SkipArgReorder:  true,
			Action:          disable,
			Flags:           EncryptFlags,
		},
		{
			Name:            "prepare",
			Usage:           "Prepare for encryption keys rotation",
			SkipFlagParsing: false,
			SkipArgReorder:  true,
			Action:          prepare,
			Flags: append(EncryptFlags, &cli.BoolFlag{
				Name:        "f,force",
				Usage:       "Force preparation.",
				Destination: &ServerConfig.EncryptForce,
			}),
		},
		{
			Name:            "rotate",
			Usage:           "Rotate secrets encryption keys",
			SkipFlagParsing: false,
			SkipArgReorder:  true,
			Action:          rotate,
			Flags: append(EncryptFlags, &cli.BoolFlag{
				Name:        "f,force",
				Usage:       "Force key rotation.",
				Destination: &ServerConfig.EncryptForce,
			}),
		},
		{
			Name:            "reencrypt",
			Usage:           "Reencrypt all data with new encryption key",
			SkipFlagParsing: false,
			SkipArgReorder:  true,
			Action:          reencrypt,
			Flags: append(EncryptFlags,
				&cli.BoolFlag{
					Name:        "f,force",
					Usage:       "Force secrets reencryption.",
					Destination: &ServerConfig.EncryptForce,
				},
				&cli.BoolFlag{
					Name:        "skip",
					Usage:       "Skip removing old key",
					Destination: &ServerConfig.EncryptSkip,
				}),
		},
	}
}
