package errors

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type StructuredError struct {
	source error
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
		source: actualCause,
		fields: fields,
	}
}

func (e StructuredError) Error() string {
	return e.source.Error()
}

func (e StructuredError) Cause() error {
	return e.source
}

func (e StructuredError) Unwrap() error {
	return e.source
}

func (e StructuredError) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, f := range e.fields {
		f.AddTo(enc)
	}
	ToField(e.source).AddTo(enc)
	return nil
}

func ToField(err error) zap.Field {
	return ToNamedField("cause", err)
}

func ToNamedField(key string, err error) zap.Field {
	if err == nil {
		return zap.Skip()
	}
	if e, ok := err.(zapcore.ObjectMarshaler); ok {
		return zap.Object(key, e)
	}
	return zap.NamedError(key, err)
}
