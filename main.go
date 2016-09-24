package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "RPC Go benchmark"
	app.Usage = "Benchmark tool for a certain RPC implemenation in Go"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "socket",
			Value: "/tmp/rpc.sock",
			Usage: "Specify unix domain socket for communication between server and client",
		},
	}
	app.Commands = []cli.Command{
		serverCmd,
		clientCmd,
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(fmt.Errorf("Error when executing command: %v", err))
	}
}
