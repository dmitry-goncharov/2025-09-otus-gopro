package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetClientSimple{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
		conn:    nil,
	}
}

type TelnetClientSimple struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

var errNoConnection = fmt.Errorf("no connection")

func (c *TelnetClientSimple) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", c.address)
	if err != nil {
		return fmt.Errorf("error connection: %w", err)
	}
	c.conn = conn
	return nil
}

func (c *TelnetClientSimple) Close() error {
	if c.conn == nil {
		return errNoConnection
	}
	return c.conn.Close()
}

func (c *TelnetClientSimple) Send() error {
	if c.conn == nil {
		return errNoConnection
	}
	_, err := io.Copy(c.conn, c.in)
	if err != nil {
		return fmt.Errorf("error sending data: %w", err)
	}
	return nil
}

func (c *TelnetClientSimple) Receive() error {
	if c.conn == nil {
		return errNoConnection
	}
	_, err := io.Copy(c.out, c.conn)
	if err != nil {
		return fmt.Errorf("error receiving data: %w", err)
	}
	return nil
}
