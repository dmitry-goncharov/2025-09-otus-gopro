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
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout connection")
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		fmt.Println("Should be set host and port")
		os.Exit(1)
	}

	client := NewTelnetClient(net.JoinHostPort(args[0], args[1]), timeout, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		fmt.Println("Can't open connection:", err.Error())
		os.Exit(1)
	}

	defer func() {
		err := client.Close()
		if err != nil {
			fmt.Println("Can't close connection:", err.Error())
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	go func() {
		err := client.Send()
		if err != nil {
			fmt.Println("Can't send data:", err.Error())
		}
		cancel()
	}()

	go func() {
		err := client.Receive()
		if err != nil {
			fmt.Println("Can't receive data:", err.Error())
		}
		cancel()
	}()

	<-ctx.Done()
}
