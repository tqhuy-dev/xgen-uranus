package interceptors

import (
	"context"
	"path"
	"time"

	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/tqhuy-dev/xgen-uranus/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// DefaultCodeToLevel is the default implementation of gRPC return codes and interceptor log level for server side.
func DefaultCodeToLevel(code codes.Code) zapcore.Level {
	switch code {
	case codes.OK:
		return zap.InfoLevel
	case codes.Canceled:
		return zap.InfoLevel
	case codes.Unknown:
		return zap.ErrorLevel
	case codes.InvalidArgument:
		return zap.InfoLevel
	case codes.DeadlineExceeded:
		return zap.WarnLevel
	case codes.NotFound:
		return zap.InfoLevel
	case codes.AlreadyExists:
		return zap.InfoLevel
	case codes.PermissionDenied:
		return zap.WarnLevel
	case codes.Unauthenticated:
		return zap.InfoLevel // unauthenticated requests can happen
	case codes.ResourceExhausted:
		return zap.WarnLevel
	case codes.FailedPrecondition:
		return zap.WarnLevel
	case codes.Aborted:
		return zap.WarnLevel
	case codes.OutOfRange:
		return zap.WarnLevel
	case codes.Unimplemented:
		return zap.ErrorLevel
	case codes.Internal:
		return zap.ErrorLevel
	case codes.Unavailable:
		return zap.WarnLevel
	case codes.DataLoss:
		return zap.ErrorLevel
	default:
		return zap.ErrorLevel
	}
}

type options struct {
	levelFunc       CodeToLevel
	shouldLog       grpc_logging.Decider
	codeFunc        grpc_logging.ErrorToCode
	durationFunc    DurationToField
	messageFunc     MessageProducer
	timestampFormat string
	appName         string
}

type Option func(*options)
type MessageProducer func(ctx context.Context, msg string, level zapcore.Level, code codes.Code, err error, duration zapcore.Field, appName string)
type DurationToField func(duration time.Duration) zapcore.Field

func WithAppName(appName string) Option {
	return func(o *options) {
		o.appName = appName
	}
}

type CodeToLevel func(code codes.Code) zapcore.Level

// DefaultDurationToField is the default implementation of converting request duration to a Zap field.
var DefaultDurationToField = DurationToTimeMillisField

// DurationToTimeMillisField converts the duration to milliseconds and uses the key `grpc.time_ms`.
func DurationToTimeMillisField(duration time.Duration) zapcore.Field {
	return zap.Float32("duration_ms", durationToMilliseconds(duration))
}

func durationToMilliseconds(duration time.Duration) float32 {
	return float32(duration.Nanoseconds()/1000) / 1000
}

var (
	defaultOptions = &options{
		levelFunc:       DefaultCodeToLevel,
		shouldLog:       grpc_logging.DefaultDeciderMethod,
		codeFunc:        grpc_logging.DefaultErrorToCode,
		durationFunc:    DefaultDurationToField,
		messageFunc:     DefaultMessageProducer,
		timestampFormat: time.RFC3339,
	}
)

// DefaultMessageProducer writes the default message
func DefaultMessageProducer(ctx context.Context, msg string, level zapcore.Level, code codes.Code, err error, duration zapcore.Field, appName string) {
	// re-extract logger from newCtx, as it may have extra fields that changed in the holder.
	correlationId, _ := ctx.Value(common.CorrelationIdKey).(string)
	fields := []zapcore.Field{
		zap.String("code", code.String()),
		zap.String("app_name", appName),
		zap.String("correlation_id", correlationId),
		duration,
	}
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	ctxzap.Extract(ctx).Check(level, msg).Write(fields...)
}

func evaluateServerOpt(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	optCopy.levelFunc = grpc_zap.DefaultCodeToLevel
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

func ZapLogInterceptor(logger *zap.Logger, opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateServerOpt(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		newCtx := newLoggerForCall(ctx, logger, info.FullMethod, startTime, o.timestampFormat)

		resp, err := handler(newCtx, req)
		if !o.shouldLog(info.FullMethod, err) {
			return resp, err
		}
		code := o.codeFunc(err)
		level := o.levelFunc(code)
		duration := o.durationFunc(time.Since(startTime))

		o.messageFunc(newCtx, "finished unary call with code "+code.String(), level, code, err, duration, o.appName)
		return resp, err
	}
}

func newLoggerForCall(ctx context.Context, logger *zap.Logger, fullMethodString string, start time.Time, timestampFormat string) context.Context {
	var f []zapcore.Field
	f = append(f, zap.String("start_time", start.Format(timestampFormat)))
	//if d, ok := ctx.Deadline(); ok {
	//	f = append(f, zap.String("grpc.request.deadline", d.Format(timestampFormat)))
	//}
	callLog := logger.With(append(f, serverCallFields(fullMethodString)...)...)
	return ctxzap.ToContext(ctx, callLog)
}

func serverCallFields(fullMethodString string) []zapcore.Field {
	method := path.Base(fullMethodString)
	return []zapcore.Field{
		zap.String("method", method),
	}
}
