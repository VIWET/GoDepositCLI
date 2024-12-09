package cli

import "github.com/urfave/cli/v2"

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "staking-cli"
	app.Commands = []*cli.Command{DepositCommand}
	return app
}
