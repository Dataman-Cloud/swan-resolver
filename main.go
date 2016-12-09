package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"golang.org/x/net/context"
)

func ServerCommand() cli.Command {
	return cli.Command{
		Name:      "server",
		Usage:     "start a dns proxy server",
		ArgsUsage: "[name]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "domain",
				Value: "swan",
				Usage: "default doamin prefix",
			},

			cli.StringFlag{
				Name:  "listener",
				Value: "0.0.0.0",
				Usage: "default ip addr",
			},

			cli.IntFlag{
				Name:  "port",
				Value: 53,
				Usage: "default port",
			},
		},
		Action: func(c *cli.Context) error {
			resolver := NewResolver(NewConfig(c))
			resolver.Start(context.Background())

			return nil
		},
	}
}

func main() {
	resolver := cli.NewApp()
	resolver.Name = "swan-resolver"
	resolver.Usage = "command-line client for resolver"
	resolver.Version = "0.1"
	resolver.Copyright = "(c) 2016 Dataman Cloud"

	resolver.Commands = []cli.Command{
		ServerCommand(),
	}

	if err := resolver.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
