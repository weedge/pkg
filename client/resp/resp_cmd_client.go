package respclient

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"github.com/weedge/pkg/utils"
)

const (
	defaultBufSize = 4096
)

type SizeWriter int64

func (s *SizeWriter) Write(p []byte) (int, error) {
	*s += SizeWriter(len(p))
	return len(p), nil
}

type RespCmdClient struct {
	conn        net.Conn
	respReader  *RespReader
	respWriter  *RespWriter
	rBufferSize SizeWriter
	wBufferSize SizeWriter
	closed      atomic.Bool
}

func Connect(addr string) (*RespCmdClient, error) {
	return ConnectWithSize(addr, defaultBufSize, defaultBufSize)
}

func ConnectWithSize(addr string, readSize int, writeSize int) (*RespCmdClient, error) {
	conn, err := net.Dial(utils.GetProto(addr), addr)
	if err != nil {
		return nil, err
	}

	return NewRespCmdClientWithSize(conn, readSize, writeSize)
}

func NewRespCmdClient(conn net.Conn) (*RespCmdClient, error) {
	return NewRespCmdClientWithSize(conn, defaultBufSize, defaultBufSize)
}

func NewRespCmdClientWithSize(conn net.Conn, readSize int, writeSize int) (*RespCmdClient, error) {
	c := new(RespCmdClient)
	c.conn = conn
	br := bufio.NewReaderSize(io.TeeReader(c.conn, &c.rBufferSize), readSize)
	bw := bufio.NewWriterSize(io.MultiWriter(c.conn, &c.wBufferSize), writeSize)
	c.respReader = NewRespReader(br)
	c.respWriter = NewRespWriter(bw)
	c.closed.Store(false)

	return c, nil

}

// Close close net.Conn, tag closed
func (c *RespCmdClient) Close() {
	if c.closed.Load() {
		return
	}

	c.conn.Close()
	c.closed.Store(true)
}

// GetConn for set conn feature (eg: r/w deadline timeout)
func (c *RespCmdClient) GetConn() net.Conn {
	return c.conn
}

// DoWithStringArgs RESP command and receive the reply
func (c *RespCmdClient) DoWithStringArgs(args ...string) (interface{}, error) {
	if err := c.SendWithStringArgs(args...); err != nil {
		return nil, err
	}

	return c.Receive()
}

// SendWithStringArgs RESP command
func (c *RespCmdClient) SendWithStringArgs(args ...string) error {
	if err := c.WriteWithStringArgs(args...); err != nil {
		c.Close()
		return err
	}

	return nil
}

func (c *RespCmdClient) WriteWithStringArgs(args ...string) error {
	if args == nil {
		return errors.New("args cannot be nil")
	}

	if len(args) == 0 {
		return errors.New("args cannot be empty")
	}

	cmdResp := fmt.Sprintf("*%d\r\n", len(args))
	for _, arg := range args {
		cmdResp = cmdResp + fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg)
	}

	return c.WriteStringResp(cmdResp)
}

func (c *RespCmdClient) WriteStringResp(s string) error {
	_, err := c.respWriter.bw.WriteString(s)
	if err != nil {
		return err
	}
	return c.respWriter.Flush()
}

// Send RESP command and receive the reply
func (c *RespCmdClient) Do(cmd string, args ...interface{}) (interface{}, error) {
	if err := c.Send(cmd, args...); err != nil {
		return nil, err
	}

	return c.Receive()
}

// Send RESP command
func (c *RespCmdClient) Send(cmd string, args ...interface{}) error {
	if err := c.respWriter.WriteCommand(cmd, args...); err != nil {
		c.Close()
		return err
	}

	return nil
}

// Receive RESP reply
func (c *RespCmdClient) Receive() (interface{}, error) {
	if reply, err := c.respReader.Parse(); err != nil {
		c.Close()
		return nil, err
	} else {
		if e, ok := reply.(Error); ok {
			return reply, e
		} else {
			return reply, nil
		}
	}
}

// Receive RESP bulk string reply into writer w
func (c *RespCmdClient) ReceiveBulkTo(w io.Writer) error {
	err := c.respReader.ParseBulkTo(w)
	if err != nil {
		if _, ok := err.(Error); !ok {
			c.Close()
		}
	}
	return err
}
