// config.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var cfg Config

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Logging LoggingConfig `yaml:"logging"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}

func initConfig(cmd *cobra.Command) error {
	cfgFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	if cfgFile == "" {
		return errors.New("config file not provided")
	}

	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

func setupLogger() (*zap.Logger, error) {
	level, err := zap.ParseAtomicLevel(cfg.Logging.Level)
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
	}

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level = level

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return logger, nil
}

func getConfig() Config {
	return cfg
}

func getServerPort() int {
	return cfg.Server.Port
}

func getLogLevel() string {
	return cfg.Logging.Level
}