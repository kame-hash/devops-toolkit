// service_test.go
package main

import (
	"context"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"devops-toolkit/config"
	"devops-toolkit/metrics"
	"devops-toolkit/service"
)

func TestService_LogParsing(t *testing.T) {
	// Create a test logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Create a test service
	srv := service.NewService(config.NewConfig(), metrics.NewPrometheusRegistry(), logger)

	// Test log parsing
	logLine := "2022-01-01 12:00:00 INFO this is a test log line"
	parsedLog, err := srv.ParseLog(context.Background(), logLine)
	if err != nil {
		t.Errorf("failed to parse log line: %v", err)
	}

	// Verify parsed log fields
	if parsedLog.Timestamp.IsZero() {
		t.Errorf("expected non-zero timestamp")
	}
	if parsedLog.Level != "INFO" {
		t.Errorf("expected log level to be INFO")
	}
}

func TestService_HealthCheck(t *testing.T) {
	// Create a test service
	srv := service.NewService(config.NewConfig(), metrics.NewPrometheusRegistry(), zap.NewNop())

	// Test health check
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	health, err := srv.HealthCheck(ctx)
	if err != nil {
		t.Errorf("health check failed: %v", err)
	}

	if health.Status != codes.OK {
		t.Errorf("expected health status to be OK")
	}
}