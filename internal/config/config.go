package config

import (
	"context"

	"github.com/ava-labs/avalanchego/utils/wrappers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	configFileName = "camino-client"

	configFlagKey = "config"

	logLevelKey = "log_level"
)

func BindFlags(cmd *cobra.Command) error {
	cmd.PersistentFlags().String(configFlagKey, ".", "path to config file dir")

	cmd.PersistentFlags().String(logLevelKey, ".", "log_level")

	errs := wrappers.Errs{}
	errs.Add(
		viper.BindPFlag(configFlagKey, cmd.PersistentFlags().Lookup(configFlagKey)),

		viper.BindPFlag(logLevelKey, cmd.PersistentFlags().Lookup(logLevelKey)),
	)
	return errs.Err
}

type Config struct {
	LogLevel string `mapstructure:"log_level"`
}

func ReadConfig(ctx context.Context, logger *zap.SugaredLogger) (*Config, error) {
	logger.Debug("Reading config...")
	viper.SetConfigName(configFileName)
	viper.AddConfigPath(".")
	viper.AddConfigPath(viper.GetString(configFlagKey)) // must already be bound from flag

	if err := viper.ReadInConfig(); err != nil {
		logger.Info(err)
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		logger.Error(err)
		return nil, err
	}

	return cfg, nil
}
