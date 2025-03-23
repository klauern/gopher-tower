package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/klauern/gopher-tower/cmd/cli/commands"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "gopher-cli",
		Usage: "A CLI tool for interacting with the Gopher Tower API",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "server",
				Aliases: []string{"s"},
				Value:   "http://localhost:8080",
				Usage:   "API server URL",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "status",
				Usage: "Check the status of the Gopher Tower API server",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					serverURL := cmd.Root().String("server")
					fmt.Printf("Checking status of server at %s...\n", serverURL)
					// TODO: Implement actual status check
					return nil
				},
			},
			commands.JobsCommand(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
