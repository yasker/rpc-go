package main

import (
	"fmt"
	"net"
	"reflect"
	"time"

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
	requestSize := c.Int("requestsize")
	if requestSize == 0 {
		requestSize = 4096
	}

	conn, err := net.Dial("unix", sockFile)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := rpc.NewClient(conn)

	return startTest(client, requestSize)
}

func startTest(client *rpc.Client, requestSize int) error {
	buf := GetRandomStringBytes(SampleSize)
	fmt.Println("Sample ready")

	writeStart := time.Now()
	for offset := 0; offset < SampleSize; offset += requestSize {
		_, err := client.WriteAt(buf[offset:offset+requestSize], int64(offset))
		if err != nil {
			return err
		}
	}
	writeStop := time.Now()
	writeDuration := writeStop.Sub(writeStart)
	fmt.Println("Write done in ", writeDuration)
	fmt.Printf("Write bandwidth %.2f MB/s\n", (float64(SampleSize/(1024*1024)) / writeDuration.Seconds()))

	tmpBuf := make([]byte, requestSize)
	readStart := time.Now()
	for offset := 0; offset < SampleSize; offset += requestSize {
		_, err := client.ReadAt(tmpBuf, int64(offset))
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(tmpBuf, buf[offset:offset+requestSize]) {
			return fmt.Errorf("Inconsistent at offset %v", offset)
		}
	}
	readStop := time.Now()
	readDuration := readStop.Sub(readStart)
	fmt.Println("Read done in ", readDuration)
	fmt.Printf("Read bandwidth %.2f MB/s\n", (float64(SampleSize/(1024*1024)) / readDuration.Seconds()))

	return nil
}
