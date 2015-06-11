// middleware/http_interceptor.go
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/spf13/cobra"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
}

func httpInterceptor(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// start timer for request latency metric
		startTime := time.Now()
		// set rate limiter for incoming requests
		limiter := rate.NewLimiter(rate.Every(10*time.Millisecond), 100)
		if !limiter.Allow() {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// call original handler
		handler(w, r)

		// calculate request latency
		latency := time.Since(startTime)
		// increment prometheus metrics
		opsProcessed.Inc()
		opsLatency.Observe(float64(latency.Seconds()))
	}
}

var opsProcessed = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "ops_processed_total",
		Help: "The total number of operations processed",
	},
)

var opsLatency = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name:    "ops_latency_seconds",
		Help:    "The latency of operations in seconds",
		Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 5},
	},
)

func main() {
	rootCmd := &cobra.Command{
		Use: "devops-toolkit",
		Run: func(cmd *cobra.Command, args []string) {
			http.Handle("/metrics", promhttp.Handler())
			http.HandleFunc("/", httpInterceptor(func(w http.ResponseWriter, r *http.Request) {
				logger.Info("incoming request", zap.String("url", r.URL.String()))
			}))
			log.Fatal(http.ListenAndServe(":8080", nil))
		},
	}
	rootCmd.Execute()
}