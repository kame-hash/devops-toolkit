// utils/utils.go
package utils

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction(zapcore.InfoLevel)
	if err != nil {
		log.Fatal(err)
	}
}

func getFilePath(filename string) (string, error) {
	// if the file is absolute, return it
	if filepath.IsAbs(filename) {
		return filename, nil
	}
	// try to find the file in the current working directory
	if _, err := os.Stat(filename); err == nil {
		return filepath.Abs(filename)
	}
	// try to find the file in the config directory
	configDir := os.Getenv("CONFIG_DIR")
	if configDir == "" {
		return "", errors.New("CONFIG_DIR environment variable is not set")
	}
	return filepath.Join(configDir, filename), nil
}

func readFile(filename string) ([]byte, error) {
	filePath, err := getFilePath(filename)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(filePath)
}

func prometheusRegisterMetric(metric prometheus.Collector) {
	if err := prometheus.Register(metric); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			logger.Warn("metric already registered", zap.Any("metric", are.Metric))
		} else {
			logger.Error("failed to register metric", zap.Error(err))
		}
	}
}