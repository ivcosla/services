package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/skycoin/services/coin-api/cmd/cli/handler"
)

var (
	rpcaddr  = new(string)
	endpoint string
)

func main() {
	// App is a cli app
	hBTC := handler.NewBTC()
	hMulti := handler.NewMulti()
	httpServer := handler.NewServerHTTP()

	cliapp := cli.App{
		Commands: []cli.Command{
			{
				Name: "btc",
				Subcommands: cli.Commands{
					cli.Command{
						Name:   "generatekeys",
						Usage:  "Generate key pair",
						Action: hBTC.GenerateKeyPair,
					},
					cli.Command{
						Name:      "generateaddr",
						Usage:     "Generate BTC addr",
						ArgsUsage: "<publicKey>",
						Action:    hBTC.GenerateAddress,
					},
					cli.Command{
						Name:      "checkbalance",
						Usage:     "Check BTC balance",
						ArgsUsage: "<address>",
						Action:    hBTC.CheckBalance,
					},
					cli.Command{
						Name:      "checktxstatus",
						Usage:     "Check BTC transaction status",
						ArgsUsage: "<transaction id>",
						Action:    hBTC.CheckTransaction	,
					},
				},
				Before: func(c *cli.Context) error {
					endpoint = "btc"
					return nil
				},
			},
			{
				Name: "coin",
				Subcommands: cli.Commands{
					cli.Command{
						Name:   "generatekeys",
						Usage:  "Generate key pair",
						Action: hMulti.GenerateKeyPair,
					},
					cli.Command{
						Name:      "generateaddr",
						Usage:     "Generate BTC addr",
						ArgsUsage: "<publicKey>",
						Action:    hMulti.GenerateAddress,
					},
					cli.Command{
						Name:      "checkbalance",
						Usage:     "Check BTC balance",
						ArgsUsage: "<address>",
						Action:    hMulti.CheckBalance,
					},
				},
				Before: func(c *cli.Context) error {
					endpoint = "coin"
					return nil
				},
			},
			{
				Name: "server",
				Subcommands: cli.Commands{
					cli.Command{
						Name:      "start",
						Usage:     "Start HTTP Server",
						ArgsUsage: "<config_file>",
						Action:    httpServer.Start,
					},
					// cli.Command{
						// Name:   "stop",
						// Usage:  "Stop HTTP Server",
						// Action: httpServer.Stop,
					// },
				},
			},
		},
		Flags: []cli.Flag{
			cli.StringFlag{Name: "rpc", Destination: rpcaddr, Value: "localhost:12345"},
		},
	}
	err := cliapp.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
