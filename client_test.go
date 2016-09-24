package main

import (
	"net"
	"testing"

	. "gopkg.in/check.v1"

	"github.com/rancher/longhorn/rpc"
)

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct {
}

var _ = Suite(&TestSuite{})

var (
	client *rpc.Client
)

func (s *TestSuite) SetUpTest(c *C) {
	socket := "/tmp/rpc.sock"

	conn, err := net.Dial("unix", socket)
	if err != nil {
		panic(err)
	}
	client = rpc.NewClient(conn)
}

func (s *TestSuite) TearDownTest(c *C) {
	client.Close()
}

func (s *TestSuite) doClientRead(size int, c *C) {
	buf := []byte(GetRandomString(size))
	for i := 0; i < c.N; i++ {
		_, err := client.ReadAt(buf, 0)
		if err != nil {
			panic("Fail reading")
		}
	}

}

func (s *TestSuite) doClientWrite(size int, c *C) {
	buf := []byte(GetRandomString(size))
	for i := 0; i < c.N; i++ {
		_, err := client.WriteAt(buf, 0)
		if err != nil {
			panic("Fail writing")
		}
	}

}

func (s *TestSuite) BenchmarkClientRead1byte(c *C) {
	size := 1

	s.doClientRead(size, c)
}

func (s *TestSuite) BenchmarkClientRead4096byte(c *C) {
	size := 4 * 4096

	s.doClientRead(size, c)
}

func (s *TestSuite) BenchmarkClientWrite1byte(c *C) {
	size := 1

	s.doClientWrite(size, c)
}

func (s *TestSuite) BenchmarkClientWrite4096byte(c *C) {
	size := 4 * 4096

	s.doClientWrite(size, c)
}
