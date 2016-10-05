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

	SampleSize = 100 * 1024 * 1024
	//SampleSize = 64 * 1024
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

type Processor struct {
	Sample []byte
}

func handleConnection(c net.Conn) {
	processor := &Processor{}
	processor.Sample = make([]byte, SampleSize, SampleSize)

	server := rpc.NewServer(c, processor)
	if err := server.Handle(); err != nil && err != io.EOF {
		panic(fmt.Sprintf("Fail to start server due to %v", err))
	} else if err == io.EOF {
		fmt.Println("Connection closed")
	}
}

func (p *Processor) ReadAt(buf []byte, off int64) (n int, err error) {
	n = copy(buf, p.Sample[off:off+int64(len(buf))])
	if n != len(buf) {
		return 0, fmt.Errorf("Fail to copy completely")
	}
	//fmt.Println("Received read at: len, buf", off, len(buf), string(buf[:16]))
	return n, nil
}

func (p *Processor) WriteAt(buf []byte, off int64) (n int, err error) {
	//fmt.Println("Received write at: len, buf", off, len(buf), string(buf[:16]))
	n = copy(p.Sample[off:], buf)
	if n != len(buf) {
		return 0, fmt.Errorf("Fail to copy completely")
	}
	return n, nil
}
