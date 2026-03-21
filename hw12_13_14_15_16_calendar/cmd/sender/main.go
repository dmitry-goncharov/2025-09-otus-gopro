package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/cmd"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/config"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/logger"
	queuefactory "github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/queue/factory"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/service/sender"
	storagefactory "github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/storage/factory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/sender/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == cmd.VERSION {
		cmd.PrintVersion()
		return
	}

	config, err := config.NewSenderConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting config: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.New(config.Logger.Level, config.Logger.Source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating logger: %v\n", err)
		os.Exit(1)
	}

	storage, err := storagefactory.NewStorage(&config.Storage, log)
	if err != nil {
		log.Error("error creating storage: " + err.Error())
		os.Exit(1)
	}

	messageQueue, err := queuefactory.NewMessageQueue(log, &config.Queue)
	if err != nil {
		log.Error("error creating message queue: " + err.Error())
		os.Exit(1)
	}
	defer func() {
		err := messageQueue.Close()
		if err != nil {
			log.Error("error closing message queue: " + err.Error())
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	err = storage.Connect(ctx)
	if err != nil {
		log.Error("error connecting to storage: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
	defer func() {
		err := storage.Close(ctx)
		if err != nil {
			log.Error("failed to close storage: " + err.Error())
		}
	}()

	messages, err := messageQueue.Consume(ctx)
	if err != nil {
		log.Error("error consuming messages: " + err.Error())
		os.Exit(1)
	}

	sender := sender.NewLogSender(log, storage)
	sender.Send(ctx, messages)

	log.Info("sender is running...")

	<-ctx.Done()

	log.Info("sender is stopped")
}
