package logging

import (
	"context"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ServiceTrace(ctx context.Context) []string {
	if ctxServiceTrace, ok := ctx.Value(serviceTraceKey).([]string); ok {
		return ctxServiceTrace
	}
	return []string{}
}

func WithDeviceId(ctx context.Context, deviceId string) context.Context {
	if deviceId == "" {
		return ctx
	}
	return context.WithValue(ctx, deviceIdKey, deviceId)
}

func PushServiceTrace(ctx context.Context, value ...string) context.Context {
	existing := ServiceTrace(ctx)
	return context.WithValue(ctx, serviceTraceKey, append(existing, value...))
}

func WithFacility(ctx context.Context, facility string) context.Context {
	return PushServiceTrace(context.WithValue(ctx, facilityKey, facility), facility)
}

func WithTaskId(ctx context.Context, taskId string) context.Context {
	return PushServiceTrace(context.WithValue(ctx, taskIdKey, taskId), taskId)
}

func WithNewTaskId(ctx context.Context, g *IdGenerator) context.Context {
	taskId := g.Generate()
	return PushServiceTrace(context.WithValue(ctx, taskIdKey, taskId), taskId)
}

func Logger(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return rootLogger
	}

	newLogger := rootLogger
	if ctxDeviceId, ok := ctx.Value(deviceIdKey).(string); ok {
		newLogger = newLogger.With(zap.String(deviceIdTagName, ctxDeviceId))
	}
	if ctxFacility, ok := ctx.Value(facilityKey).(string); ok {
		newLogger = newLogger.With(zap.String(facilityTagName, ctxFacility))
	}
	if ctxTaskId, ok := ctx.Value(taskIdKey).(string); ok {
		newLogger = newLogger.With(zap.String(taskIdTagName, ctxTaskId))
	}
	if ctxHandler, ok := ctx.Value(handlerKey).(string); ok {
		newLogger = newLogger.With(zap.String(handlerTagName, ctxHandler))
	}
	if ctxQueue, ok := ctx.Value(queueKey).(string); ok {
		newLogger = newLogger.With(zap.String(queueTagName, ctxQueue))
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
