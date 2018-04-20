package backend

import (
	"context"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/google/uuid"
)

type correlationIdType int

const (
	facilityKey correlationIdType = iota
	taskIdKey   correlationIdType = iota
)

var rootLogger *zap.Logger

func init() {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	rootLogger = logger

	zap.RedirectStdLog(rootLogger)
}

func WithFacility(ctx context.Context, facility string) context.Context {
	return context.WithValue(ctx, facilityKey, facility)
}

func WithTaskId(ctx context.Context, taskId string) context.Context {
	return context.WithValue(ctx, taskIdKey, taskId)
}

func WithNewTaskId(ctx context.Context) context.Context {
	id, err := uuid.NewRandom()
	if err != nil {
		return ctx
	}
	return WithTaskId(ctx, id.String())
}

func Logger(ctx context.Context) *zap.Logger {
	newLogger := rootLogger
	if ctx != nil {
		name := make([]string, 0)
		if ctxFacility, ok := ctx.Value(facilityKey).(string); ok {
			newLogger = newLogger.With(zap.String("facility", ctxFacility))
			name = append(name, ctxFacility)
		}
		if ctxTaskId, ok := ctx.Value(taskIdKey).(string); ok {
			newLogger = newLogger.With(zap.String("taskId", ctxTaskId))
			name = append(name, ctxTaskId)
		}
		if len(name) > 0 {
			newLogger = newLogger.Named(strings.Join(name, "."))
		}
	}
	return newLogger
}
