package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/urfave/cli"

	"github.com/rancher/longhorn/rpc"
)

var (
	serverCmd = cli.Command{
		Name:   "server",
		Usage:  "Start server",
		Action: cmdStartServer,
	}

	context []byte
)

func cmdStartServer(c *cli.Context) {
	if err := doStartServer(c); err != nil {
		panic(err)
	}
}

func doStartServer(c *cli.Context) error {
	sockFile := c.GlobalString("socket")
	if sockFile == "" {
		return fmt.Errorf("Empty socket address")
	}

	os.Remove(sockFile)

	context = []byte(GetRandomString(128 * 1024 * 1024))

	ln, err := net.Listen("unix", sockFile)
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go handleConnection(conn)
	}
}

type Processor struct{}

func handleConnection(c net.Conn) {
	processor := &Processor{}
	server := rpc.NewServer(c, processor)
	if err := server.Handle(); err != nil && err != io.EOF {
		panic(fmt.Sprintf("Fail to start server due to %v", err))
	} else if err == io.EOF {
		fmt.Println("Connection closed")
	}
}

func (p *Processor) ReadAt(buf []byte, off int64) (n int, err error) {
	if off == 0 {
		copy(buf, context[:len(buf)])
		return len(buf), nil
	}
	return 0, fmt.Errorf("No data for read")
}

func (p *Processor) WriteAt(buf []byte, off int64) (n int, err error) {
	if off == 0 {
		return len(buf), nil
	}
	return 0, fmt.Errorf("No data for write")
}
