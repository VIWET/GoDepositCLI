package main

import (
	"context"
	"log"
	"os"

	"github.com/viwet/GoDepositCLI/cli"
)

func main() {
	ctx := context.Background()

	if err := cli.NewApp().Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
