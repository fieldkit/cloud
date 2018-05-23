package errors

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type StructuredError struct {
	error
	fields []zapcore.Field
}

func Structured(causeOrMessage interface{}, raw ...interface{}) *StructuredError {
	var actualCause error

	if c, ok := causeOrMessage.(error); ok {
		actualCause = c
	} else {
		actualCause = fmt.Errorf(causeOrMessage.(string))
	}

	fields := make([]zapcore.Field, len(raw)/2)
	for i := 0; i < len(raw)/2; i += 1 {
		key := raw[i*2].(string)
		fields[i] = zap.Any(key, raw[i*2+1])
	}
	return &StructuredError{
		error:  actualCause,
		fields: fields,
	}
}

func (se *StructuredError) Cause() error {
	return se.error
}

func (se *StructuredError) Fields() []zapcore.Field {
	return se.fields
}
