package cli

import (
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/viwet/GoDepositCLI/version"
)

func NewApp() *cli.Command {
	app := new(cli.Command)

	app.EnableShellCompletion = true
	app.Usage = fmt.Sprintf("This CLI application generates deposit data for %s.", NetworkName)
	app.Flags = sharedFlags
	app.Name = "staking-cli"
	app.Version = version.Version()
	app.Commands = commands

	return app
}
