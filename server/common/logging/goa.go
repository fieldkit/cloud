package logging

import (
	"github.com/goadesign/goa"
	"goa.design/goa/v3/middleware"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type adapter struct {
	logger *zap.Logger
}

func NewGoaAdapter(logger *zap.Logger) goa.LogAdapter {
	return &adapter{logger: logger.Named("goa")}
}

func (a *adapter) Info(msg string, data ...interface{}) {
	fields := toZapFields(data)
	a.getTaskedLogger(fields).Info(msg, *fields...)
}

func (a *adapter) Error(msg string, data ...interface{}) {
	fields := toZapFields(data)
	a.getTaskedLogger(fields).Error(msg, *fields...)
}

func (a *adapter) New(data ...interface{}) goa.LogAdapter {
	fields := toZapFields(data)
	return &adapter{logger: a.getTaskedLogger(fields).With(*fields...)}
}

func (a *adapter) getRequestId(fields *[]zapcore.Field) string {
	for _, f := range *fields {
		if f.Key == reqIDTagName {
			return f.String
		}
	}
	return ""
}

func (a *adapter) getTaskedLogger(fields *[]zapcore.Field) *zap.Logger {
	id := a.getRequestId(fields)
	if len(id) > 0 {
		return a.logger.With(zap.String(taskIDTagName, id))
	}
	return a.logger
}

type middlewareAdapter struct {
	logger *zap.Logger
}

func NewGoaMiddlewareAdapter(logger *zap.Logger) middleware.Logger {
	return &middlewareAdapter{logger: logger.Named("goa")}
}

func (a *middlewareAdapter) Log(keyvals ...interface{}) error {
	fields := toZapFields(keyvals)
	a.logger.Info("goa", *fields...)
	return nil
}
