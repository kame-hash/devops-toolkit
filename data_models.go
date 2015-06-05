// data_models.go
package main

import (
	"errors"
	"time"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Logging  LoggingConfig  `yaml:"logging"`
	Metrics  MetricsConfig  `yaml:"metrics"`
	Health   HealthConfig   `yaml:"health"`
	Retry    RetryConfig    `yaml:"retry"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

type MetricsConfig struct {
	Port    int    `yaml:"port"`
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

type HealthConfig struct {
	Timeout time.Duration `yaml:"timeout"`
	Interval time.Duration `yaml:"interval"`
}

type RetryConfig struct {
	MaxAttempts int           `yaml:"maxAttempts"`
	Backoff     BackoffConfig `yaml:"backoff"`
}

type BackoffConfig struct {
	Initial    time.Duration `yaml:"initial"`
	Max         time.Duration `yaml:"max"`
	Multiplier  float64      `yaml:"multiplier"`
	Randomization float64      `yaml:"randomization"`
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

type HealthCheckResult struct {
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func (c *Config) Validate() error {
	if c.Server.Host == "" {
		return errors.New("server host is required")
	}
	if c.Logging.Level == "" {
		return errors.New("logging level is required")
	}
	return nil
}

func (c *Config) ToZapLogger() (*zap.Logger, error) {
	// create a new zap logger with the specified level and output
	return zap.NewProduction(zap.Level(c.Logging.Level), zap.Output(c.Logging.Output)), nil
}