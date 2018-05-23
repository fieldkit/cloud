package logging

import (
	"go.uber.org/zap/zapcore"
)

func NewStructuredErrorsCore(c zapcore.Core) zapcore.Core {
	return &structuredErrorsCore{c}
}

type causedErr interface {
	error
	Cause() error
}

type structuredErr interface {
	error
	Cause() error
	Fields() []zapcore.Field
}

type structuredErrorsCore struct {
	zapcore.Core
}

func (c *structuredErrorsCore) With(fields []zapcore.Field) zapcore.Core {
	fields = extractErrorFields(fields)
	return &structuredErrorsCore{
		c.Core.With(fields),
	}
}

func (c *structuredErrorsCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	if !hasFieldsToExtract(fields) {
		return c.Core.Write(ent, fields)
	}
	fields = extractErrorFields(fields)
	return c.Core.Write(ent, fields)
}

func (c *structuredErrorsCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	// HACK: This breaks sampling cores. See https://github.com/uber-go/zap/issues/529
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func hasFieldsToExtract(fields []zapcore.Field) bool {
	for _, field := range fields {
		if field.Type != zapcore.ErrorType {
			continue
		}
		if _, ok := field.Interface.(structuredErr); ok {
			return true
		}
	}
	return false
}

func extractErrorFields(fields []zapcore.Field) []zapcore.Field {
	copied := false

	for i, field := range fields {
		if field.Type != zapcore.ErrorType {
			continue
		}

		plainError := field.Interface.(error)

		for {
			if structured, ok := plainError.(structuredErr); ok {
				if !copied {
					copied = true
					oldFields := fields
					fields = make([]zapcore.Field, len(fields))
					copy(fields, oldFields)
				}

				fields = append(fields, structured.Fields()...)

				plainError = structured.Cause()
			} else if causedErr, ok := plainError.(causedErr); ok {
				plainError = causedErr.Cause()
			} else {
				break
			}

			if plainError == nil {
				break
			}
		}

		fields[i] = field
	}
	return fields
}
