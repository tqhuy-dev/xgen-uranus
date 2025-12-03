package interceptors

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tqhuy-dev/xgen-uranus/common"
	"go.uber.org/zap"
)

type httpLoggingOptions struct {
	appName string
}

type HttpLoggingOption func(*httpLoggingOptions)

func WithHttpAppName(appName string) HttpLoggingOption {
	return func(o *httpLoggingOptions) {
		o.appName = appName
	}
}

// HttpZapLogMiddleware is a Gin middleware that logs HTTP requests using zap
func HttpZapLogMiddleware(logger *zap.Logger, opts ...HttpLoggingOption) gin.HandlerFunc {
	o := &httpLoggingOptions{}
	for _, opt := range opts {
		opt(o)
	}

	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method

		// Process request
		c.Next()

		// After request - log the response
		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		correlationId, _ := c.Get(common.CorrelationIdKey)
		correlationIdStr, _ := correlationId.(string)

		fields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Float32("duration_ms", float32(duration.Nanoseconds()/1000)/1000),
			zap.String("app_name", o.appName),
			zap.String("correlation_id", correlationIdStr),
		}

		if query != "" {
			fields = append(fields, zap.String("query", query))
		}

		// Log errors if any
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("error", c.Errors.String()))
		}

		// Determine log level based on status code
		msg := "HTTP request completed"
		switch {
		case statusCode >= 500:
			logger.Error(msg, fields...)
		case statusCode >= 400:
			logger.Warn(msg, fields...)
		default:
			logger.Info(msg, fields...)
		}
	}
}
