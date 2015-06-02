// core/service.go
package core

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"devops-toolkit/config"
	"devops-toolkit/logparser"
	"devops-toolkit/metrics"
	"devops-toolkit/retry"
)

type Service struct {
	logger    *zap.Logger
	config    *config.Config
	parser    *logparser.LogParser
	retry     *retry.Retry
	metrics   *metrics.Metrics
	health    *healthChecker
}

func NewService(l *zap.Logger, c *config.Config, p *logparser.LogParser, r *retry.Retry, m *metrics.Metrics, h *healthChecker) *Service {
	return &Service{
		logger:    l,
		config:    c,
		parser:    p,
		retry:     r,
		metrics:   m,
		health:    h,
	}
}

func (s *Service) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return s.startLogParsing(ctx)
	})
	g.Go(func() error {
		return s.startHealthChecks(ctx)
	})
	g.Go(func() error {
		return s.startMetricsServer(ctx)
	})
	return g.Wait()
}

func (s *Service) startLogParsing(ctx context.Context) error {
	s.logger.Info("starting log parsing")
	// parse logs every 10 seconds
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			// parse logs and update metrics
			logs, err := s.parser.ParseLogs()
			if err != nil {
				s.logger.Error("failed to parse logs", zap.Error(err))
				continue
			}
			s.metrics.UpdateLogCount(logs)
		}
	}
}

func (s *Service) startHealthChecks(ctx context.Context) error {
	s.logger.Info("starting health checks")
	// perform health checks every minute
	t := time.NewTicker(1 * time.Minute)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			// perform health checks and update metrics
			err := s.health.Check()
			if err != nil {
				s.logger.Error("health check failed", zap.Error(err))
				s.metrics.UpdateHealthStatus(err)
			} else {
				s.metrics.UpdateHealthStatus(nil)
			}
		}
	}
}

func (s *Service) startMetricsServer(ctx context.Context) error {
	s.logger.Info("starting metrics server")
	http.Handle("/metrics", s.metrics.Handler())
	return http.ListenAndServe(":8080", nil)
}