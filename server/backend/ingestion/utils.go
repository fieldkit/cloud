package ingestion

import (
	"fmt"
	"unicode"
)

// ToSnake convert the given string to snake case following the Golang format:
// acronyms are converted to lower-case and preceded by an underscore.
func ToSnake(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		if !unicode.IsSpace(runes[i]) {
			out = append(out, unicode.ToLower(runes[i]))
		}
	}

	return string(out)
}

type IngestError struct {
	Cause    error
	Critical bool
}

func (e *IngestError) Error() string {
	return e.Cause.Error()
}

func NewErrorf(critical bool, f string, a ...interface{}) *IngestError {
	return &IngestError{
		Cause:    fmt.Errorf(f, a...),
		Critical: critical,
	}
}

func NewError(cause error) *IngestError {
	if err, ok := cause.(*IngestError); ok {
		return &IngestError{
			Cause:    err,
			Critical: err.Critical,
		}
	}
	return &IngestError{
		Cause:    cause,
		Critical: true,
	}

}
