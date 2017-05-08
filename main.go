package main

import (
	"github.com/Dataman-Cloud/swan-resolver/nameserver"

	"github.com/Sirupsen/logrus"
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
				Value: "swan.com",
				Usage: "default doamin prefix",
			},

			cli.StringFlag{
				Name:  "listen-addr",
				Value: "0.0.0.0:53",
				Usage: "default listen addr",
			},

			cli.StringFlag{
				Name:  "log-level",
				Value: "debug",
			},
		},
		Action: func(c *cli.Context) error {
			level, _ := logrus.ParseLevel(c.String("log-level"))
			logrus.SetLevel(level)

			resolver := nameserver.NewResolver(nameserver.NewConfig(c))
			go func() {
				a := nameserver.RecordChangeEvent{
					Change:  "add",
					Type:    nameserver.A,
					Ip:      "192.168.1.1",
					Cluster: "cluster",
					RunAs:   "xcm",
					AppName: "nginx",
					SlotID:  "1",
				}
				resolver.RecordChangeChan <- &a

				a1 := nameserver.RecordChangeEvent{
					Change:  "add",
					Type:    nameserver.A,
					Ip:      "192.168.1.2",
					Cluster: "cluster",
					RunAs:   "xcm",
					AppName: "nginx",
					SlotID:  "0",
				}
				resolver.RecordChangeChan <- &a1

				srv := nameserver.RecordChangeEvent{
					Change:  "add",
					Type:    nameserver.SRV ^ nameserver.A,
					Ip:      "192.168.1.3",
					Port:    "1234",
					Cluster: "cluster",
					RunAs:   "xcm",
					AppName: "nginx",
					SlotID:  "2",
				}
				resolver.RecordChangeChan <- &srv

				srv1 := nameserver.RecordChangeEvent{
					Change:  "add",
					Type:    nameserver.SRV ^ nameserver.A,
					Ip:      "192.168.1.4",
					Port:    "1235",
					Cluster: "cluster",
					RunAs:   "xcm",
					AppName: "nginx",
					SlotID:  "3",
				}
				resolver.RecordChangeChan <- &srv1

				proxy1 := nameserver.RecordChangeEvent{
					Change:  "add",
					Type:    nameserver.A,
					Ip:      "192.168.1.5",
					IsProxy: true,
				}
				resolver.RecordChangeChan <- &proxy1

				proxy2 := nameserver.RecordChangeEvent{
					Change:  "add",
					Type:    nameserver.A,
					Ip:      "192.168.1.6",
					IsProxy: true,
				}
				resolver.RecordChangeChan <- &proxy2

				da1 := nameserver.RecordChangeEvent{
					Change: "del",
					Type:   nameserver.A,
					Ip:     "192.168.1.2",
				}
				resolver.RecordChangeChan <- &da1
			}()

			started := make(chan bool, 1)
			return resolver.Start(context.Background(), started)
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

	resolver.RunAndExitOnError()
}
