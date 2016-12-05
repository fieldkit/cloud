package data

import (
	"encoding/json"
	"testing"
)

const testErrorsJSONString = `{"key":["message","message"]}`

func TestErrrors(t *testing.T) {
	e := NewErrors()
	e.Error("key", "message")
	e.Error("key", "message")

	errorsJson, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	if errorsJSONString := string(errorsJson); errorsJSONString != testErrorsJSONString {
		t.Errorf("incorrect JSON string\n%s", errorsJSONString)
	}
}
