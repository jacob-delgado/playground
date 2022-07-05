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
	pgrpc "github.com/jacob-delgado/playground/pkg/grpc"
	"github.com/jacob-delgado/playground/pkg/http"
)

var rootCmd = &cobra.Command{
	Use:   "playground",
	Short: "A sample application to serve files over http/grpc/xds",
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

		logger.Info("starting playground server")

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)

		errCh := make(chan error)
		server := http.NewServer(logger)
		go server.Serve(errCh)

		pgServer := pgrpc.NewServer(logger)
		go pgServer.Serve(errCh)

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
