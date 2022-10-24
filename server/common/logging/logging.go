package logging

import (
	"context"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func OnlyLogIf(logger *zap.SugaredLogger, verbose bool) *zap.SugaredLogger {
	if !verbose {
		return zap.NewNop().Sugar()
	}
	return logger.With("verbose", true)
}

func ServiceTrace(ctx context.Context) []string {
	if ctxServiceTrace, ok := ctx.Value(serviceTraceKey).([]string); ok {
		return ctxServiceTrace
	}
	return []string{}
}

func WithDeviceID(ctx context.Context, deviceID string) context.Context {
	if deviceID == "" {
		return ctx
	}
	return context.WithValue(ctx, deviceIDKey, deviceID)
}

const (
	MaximumServiceTrace = 10
)

func PushServiceTrace(ctx context.Context, value ...string) context.Context {
	if !IsServiceTraceEnabled() {
		return ctx
	}

	existing := ServiceTrace(ctx)
	narrowed := existing
	for _, part := range value {
		if len(narrowed) >= MaximumServiceTrace {
			narrowed := make([]string, MaximumServiceTrace)
			narrowed[0] = existing[0]
			narrowed[1] = existing[1]
			for i := 2; i <= MaximumServiceTrace-3; i += 1 {
				narrowed[i] = existing[i+1]
			}
			narrowed[MaximumServiceTrace-1] = part
		} else {
			narrowed = append(narrowed, part)
		}
	}
	return context.WithValue(ctx, serviceTraceKey, narrowed)
}

func WithFacility(ctx context.Context, facility string) context.Context {
	return PushServiceTrace(context.WithValue(ctx, facilityKey, facility), facility)
}

func WithTaskID(ctx context.Context, taskID string) context.Context {
	return PushServiceTrace(context.WithValue(ctx, taskIDKey, taskID), taskID)
}

func WithNewTaskID(ctx context.Context, g *IdGenerator) context.Context {
	taskID := g.Generate()
	return PushServiceTrace(context.WithValue(ctx, taskIDKey, taskID), taskID)
}

func WithUserID(ctx context.Context, userID int32) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func FindTaskID(ctx context.Context) string {
	return ctx.Value(taskIDKey).(string)
}

func Logger(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return rootLogger
	}

	newLogger := rootLogger
	if ctxDeviceID, ok := ctx.Value(deviceIDKey).(string); ok {
		newLogger = newLogger.With(zap.String(deviceIDTagName, ctxDeviceID))
	}
	if ctxFacility, ok := ctx.Value(facilityKey).(string); ok {
		newLogger = newLogger.With(zap.String(facilityTagName, ctxFacility))
	}
	if ctxTaskID, ok := ctx.Value(taskIDKey).(string); ok {
		newLogger = newLogger.With(zap.String(taskIDTagName, ctxTaskID))
	}
	if ctxHandler, ok := ctx.Value(handlerKey).(string); ok {
		newLogger = newLogger.With(zap.String(handlerTagName, ctxHandler))
	}
	if ctxQueue, ok := ctx.Value(queueKey).(string); ok {
		newLogger = newLogger.With(zap.String(queueTagName, ctxQueue))
	}
	if ctxUserID, ok := ctx.Value(userIDKey).(string); ok {
		newLogger = newLogger.With(zap.String(userIDTagName, ctxUserID))
	}
	if ctxServiceTrace, ok := ctx.Value(serviceTraceKey).([]string); ok {
		newLogger = newLogger.With(zap.String(serviceTraceTagName, strings.Join(ctxServiceTrace, " ")))
	}
	return newLogger
}

func toZapFields(data []interface{}) *[]zapcore.Field {
	n := (len(data) + 1) / 2
	fields := make([]zapcore.Field, n)

	fi := 0
	for i := 0; i < len(data); i += 2 {
		if key, ok := data[i].(string); ok {
			if i+1 < len(data) {
				v := data[i+1]
				fields[fi] = zap.Any(key, v)
			}
		} else {
			fields[fi] = zap.Skip()
		}
		fi = fi + 1
	}
	return &fields
}
