package log

import (
	"context"

	"github.com/jackc/pgx/v5/tracelog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type PgxLogAdapter struct{}

func (PgxLogAdapter) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	logger := FromContext(ctx).Desugar().
		WithOptions(zap.AddCallerSkip(5)).
		WithOptions(zap.AddStacktrace(zapcore.FatalLevel))

	fields := make([]zapcore.Field, 0, len(data))
	for k, v := range data {
		if k == "time" {
			// below renaming is required by GKE to properly display SQL queries
			k = "queryTime"
		}
		fields = append(fields, zap.Reflect(k, v))
	}

	switch level {
	case tracelog.LogLevelTrace:
		logger.Debug(msg, append(fields, zap.Stringer("PGX_LOG_LEVEL", level))...)
	case tracelog.LogLevelDebug:
		logger.Debug(msg, fields...)
	case tracelog.LogLevelInfo:
		logger.Info(msg, fields...)
	case tracelog.LogLevelWarn:
		logger.Warn(msg, fields...)
	case tracelog.LogLevelError:
		logger.Error(msg, fields...)
	default:
		logger.Error(msg, append(fields, zap.Stringer("PGX_LOG_LEVEL", level))...)
	}
}
