package cli

import "github.com/urfave/cli/v2"

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "staking-cli"
	app.HelpName = AppName
	app.Commands = []*cli.Command{
		DepositCommand,
		BLSToExecutionCommand,
	}

	return app
}
