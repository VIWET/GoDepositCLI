package main

import (
	"log"
	"os"

	"github.com/viwet/GoDepositCLI/cli"
)

func main() {
	if err := cli.NewApp().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
