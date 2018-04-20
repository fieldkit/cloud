/*
Package goazap contains an adapter that makes it possible to configure goa so it uses zap
as logger backend.
Usage:
    logger, err := zap.NewProduction()
	...
    // Initialize logger handler using zap package
    service.WithLogger(goazap.New(logger))
    // ... Proceed with configuring and starting the goa service

    // In handlers:
    goazap.Logger(ctx).Info("foo")
*/
package goazap

import (
	"github.com/goadesign/goa"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/context"
)

type adapter struct {
	Logger *zap.Logger
}

// New wraps a zap logger into a goa logger adapter.
func New(logger *zap.Logger) goa.LogAdapter {
	return &adapter{Logger: logger}
}

// Logger returns the zap logger stored in the given context if any, nil otherwise.
func Logger(ctx context.Context) *zap.Logger {
	logger := goa.ContextLogger(ctx)
	if a, ok := logger.(*adapter); ok {
		return a.Logger
	}
	return nil
}

// Info logs informational messages using zap.
func (a *adapter) Info(msg string, data ...interface{}) {
	fields := data2fields(data)
	a.Logger.Info(msg, *fields...)
}

// Error logs error messages using zap.
func (a *adapter) Error(msg string, data ...interface{}) {
	fields := data2fields(data)
	a.Logger.Error(msg, *fields...)
}

// New creates a new logger given a context.
func (a *adapter) New(data ...interface{}) goa.LogAdapter {
	fields := data2fields(data)
	return &adapter{Logger: a.Logger.With(*fields...)}
}

func data2fields(keyvals []interface{}) *[]zapcore.Field {
	n := (len(keyvals) + 1) / 2
	fields := make([]zapcore.Field, n)

	fi := 0
	for i := 0; i < len(keyvals); i += 2 {
		if key, ok := keyvals[i].(string); ok {
			if i+1 < len(keyvals) {
				v := keyvals[i+1]
				fields[fi] = zap.Any(key, v)
			}
		} else {
			fields[fi] = zap.Skip()
		}
		fi = fi + 1
	}
	return &fields
}
