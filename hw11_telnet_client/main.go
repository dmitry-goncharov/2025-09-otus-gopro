package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	conf, err := parseConf()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	client := NewTelnetClient(net.JoinHostPort(conf.host, conf.port), conf.timeout, os.Stdin, os.Stdout)

	err = client.Connect()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't open connection:", err.Error())
		os.Exit(1)
	}
	defer closeClient(client)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	go send(client, cancel)
	go receive(client)

	<-ctx.Done()
}

type conf struct {
	host    string
	port    string
	timeout time.Duration
}

func parseConf() (*conf, error) {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout connection")
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		return nil, fmt.Errorf("should be set host and port")
	}
	return &conf{
		host:    args[0],
		port:    args[1],
		timeout: timeout,
	}, nil
}

func send(client TelnetClient, cancel context.CancelFunc) {
	err := client.Send()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't send data:", err.Error())
		cancel()
	}
}

func receive(client TelnetClient) {
	err := client.Receive()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't receive data:", err.Error())
	}
}

func closeClient(client TelnetClient) {
	err := client.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't close connection:", err.Error())
	}
}
