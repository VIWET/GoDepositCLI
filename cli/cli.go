package cli

import (
	"github.com/urfave/cli/v3"
	"github.com/viwet/GoDepositCLI/version"
)

func NewApp() *cli.Command {
	app := new(cli.Command)

	app.Flags = []cli.Flag{NonInteractiveFlag, EngineWorkersFlag}
	app.Name = "staking-cli"
	app.Version = version.Version()
	app.Commands = commands
	return app
}
