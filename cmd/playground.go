package cmd

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"caminoclient/internal/config"
	"caminoclient/internal/logger"
	"caminoclient/internal/playground"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Execute() error {
	rootCmd := &cobra.Command{
		Use: "playground",
		Run: func(cmd *cobra.Command, args []string) {
			zapLogger, err := zap.NewDevelopment()
			if err != nil {
				log.Fatal(err)
			}
			zapSugaredLogger := zapLogger.Sugar()
			defer func() { _ = zapSugaredLogger.Sync() }()

			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			cfg, err := config.ReadConfig(ctx, zapSugaredLogger)
			if err != nil {
				return
			}

			if cfg.LogLevel == "info" {
				zapLogger, err = zap.NewProduction()
				if err != nil {
					log.Fatal(err)
				}
				_ = zapSugaredLogger.Sync()
				zapSugaredLogger = zapLogger.Sugar()
				defer func() { _ = zapSugaredLogger.Sync() }()
			}

			playground, err := playground.NewPlayground(
				ctx,
				logger.NewLoggerFromZap(zapSugaredLogger),
				cfg,
			)
			if err != nil {
				playground.Close(context.Background())
				return
			}

			playground.Run(ctx)
		},
	}
	if err := config.BindFlags(rootCmd); err != nil {
		return err
	}
	return rootCmd.Execute()
}
