package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"

	"github.com/jacob-delgado/playground/pkg/config"
	"github.com/jacob-delgado/playground/pkg/grpc"
)

var rootCmd = &cobra.Command{
	Use:   "inventory",
	Short: "A sample application to serve files over grpc",
	Long: `A sample application that can be used
for experimentation with various tools.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Init()
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := cfg.Shutdown(context.Background()); err != nil {
				log.Panicf("Error shutting down tracer provider: %v", err)
			}
		}()

		zapLogger, err := zap.NewProduction()
		if err != nil {
			log.Panicf("could not initialize zap logger: %v", err)
		}
		logger := otelzap.New(zapLogger)
		defer func() {
			_ = logger.Sync()
		}()

		logger.Info("starting inventory server")

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)

		errCh := make(chan error)
		grpcServer := grpc.NewServer(logger)
		go grpcServer.Serve(errCh)

		select {
		case <-sigCh:
			logger.Info("exiting")
			return
		case err := <-errCh:
			if err != nil {
				logger.Panic("exiting due to error", zap.Error(err))
			}
		}
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
