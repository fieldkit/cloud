package backend

import (
	"context"

	"go.uber.org/zap"
)

type correlationIdType int

const (
	requestIdKey correlationIdType = iota
)

var rootLogger *zap.Logger

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	rootLogger = logger

	zap.RedirectStdLog(rootLogger)
}

func WithRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, requestIdKey, requestId)
}

func Logger(ctx context.Context) *zap.Logger {
	newLogger := rootLogger
	if ctx != nil {
		if ctxRqId, ok := ctx.Value(requestIdKey).(string); ok {
			newLogger = newLogger.With(zap.String("requestId", ctxRqId))
		}
	}
	return newLogger
}
