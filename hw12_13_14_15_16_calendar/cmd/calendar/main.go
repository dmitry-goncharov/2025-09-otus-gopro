package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/config"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/server/http"
	storagefactory "github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/storage/factory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := config.NewConfig(configFile)
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

	calendar := app.NewApplication(log, storage)

	httpServer := internalhttp.NewServer(&config.HTTPServer, log, calendar)
	grpcServer := internalgrpc.NewServer(&config.GRPCServer, log, calendar)

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

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := httpServer.Stop(ctx); err != nil {
			log.Error("failed to stop http server: " + err.Error())
		}
	}()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		grpcServer.Stop(ctx)
	}()

	log.Info("calendar is running...")

	go func() {
		if err := httpServer.Start(ctx); err != nil {
			log.Error("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	go func() {
		if err := grpcServer.Start(ctx); err != nil {
			log.Error("failed to start grpc server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()
	<-ctx.Done()

	log.Info("calendar is stopped")
}
