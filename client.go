package main

import (
	"fmt"
	"net"

	"github.com/urfave/cli"

	"github.com/rancher/longhorn/rpc"
)

var (
	clientCmd = cli.Command{
		Name:  "client",
		Usage: "Start client",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "size, s",
				Usage: "Total data size to be send",
			},
			cli.IntFlag{
				Name:  "requestsize, r",
				Usage: "Size of each request",
			},
			cli.IntFlag{
				Name:  "queuedepth, q",
				Usage: "Number of request can be pending",
			},
		},
		Action: cmdStartClient,
	}
)

func cmdStartClient(c *cli.Context) {
	if err := doStartClient(c); err != nil {
		panic(err)
	}
}

func doStartClient(c *cli.Context) error {
	sockFile := c.GlobalString("socket")
	if sockFile == "" {
		return fmt.Errorf("Require unix domain socket location")
	}

	conn, err := net.Dial("unix", sockFile)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := rpc.NewClient(conn)

	buf := []byte("request")
	_, err = client.WriteAt(buf, 0)
	if err != nil {
		return err
	}
	fmt.Println("Send: ", string(buf))

	buf = make([]byte, 8)
	_, err = client.ReadAt(buf, 0)
	if err != nil {
		return err
	}
	fmt.Println("Receive: ", string(buf))

	return nil
}
