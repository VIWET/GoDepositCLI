package cli

import (
	"github.com/urfave/cli/v2"
	"github.com/viwet/GoDepositCLI/version"
)

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "staking-cli"
	app.Version = version.Version()
	app.Commands = commands
	return app
}
