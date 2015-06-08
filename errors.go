// pkg/errors/errors.go
package errors

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ErrConfigInvalid is returned when the configuration is invalid
var ErrConfigInvalid = errors.New("configuration is invalid")

// ErrHealthCheckFailed is returned when a health check fails
var ErrHealthCheckFailed = errors.New("health check failed")

// ErrLogParseFailed is returned when log parsing fails
var ErrLogParseFailed = errors.New("log parse failed")

// ErrMetricPushFailed is returned when pushing metrics to Prometheus fails
var ErrMetricPushFailed = errors.New("metric push failed")

// ErrRetryExceeded is returned when the retry count is exceeded
var ErrRetryExceeded = errors.New("retry count exceeded")

// WrapError wraps an error with additional context
func WrapError(err error, msg string) error {
	return errors.Wrap(err, msg)
}

// WrapErrorWithFields wraps an error with additional context and fields
func WrapErrorWithFields(err error, msg string, fields []zapcore.Field) error {
	return errors.Wrap(err, msg)
}

// NewHTTPError returns an error for an HTTP request
func NewHTTPError(code int, msg string) error {
	return fmt.Errorf("HTTP %d: %s", code, msg)
}

// IsHTTPError checks if an error is an HTTP error
func IsHTTPError(err error) bool {
	var httpErr *errors.Error
	if errors.As(err, &httpErr) {
		_, ok := httpErr.(*http.HTTPError)
		return ok
	}
	return false
}
func getLogLevel(err error) zapcore.Level {
	if err == nil {
		return zapcore.DebugLevel
	}
	if IsHTTPError(err) {
		httpErr, ok := err.(*errors.Error)
		if ok {
			httpCode := httpErr.Code
			if httpCode >= http.StatusInternalServerError {
				return zapcore.ErrorLevel
			}
		}
	}
	return zapcore.WarnLevel
}
func logError(logger *zap.Logger, err error) {
	logLevel := getLogLevel(err)
	logger.Check(logLevel, err.Error()).Write()
}