package jsondocument

import (
	"math"
	"testing"
)

const (
	testString  = "string value"
	testNumber  = math.Pi
	testBoolean = true
)

var (
	testPointers = []string{
		"",
		"/",
		"/0",
		"/key",
		"/0/0",
		"/0/key",
		"/key/0",
		"/key/key",
	}
)

func TestStringDocument(t *testing.T) {
	document := String(testString)
	data, err := document.String()
	if err != nil {
		t.Fatal(err)
	}

	if data != testString {
		t.Error("unexpected string value")
	}
}

func TestNumberDocument(t *testing.T) {
	document := Number(testNumber)
	data, err := document.Number()
	if err != nil {
		t.Fatal(err)
	}

	if data != testNumber {
		t.Error("unexpected number value")
	}
}

func TestBooleanDocument(t *testing.T) {
	document := Boolean(testBoolean)
	data, err := document.Boolean()
	if err != nil {
		t.Fatal(err)
	}

	if data != testBoolean {
		t.Error("unexpected boolean value")
	}
}

func TestDocument(t *testing.T) {
	for _, pointer := range testPointers {
		document := &Document{}
		if err := document.SetDocument(pointer, String(testString)); err != nil {
			t.Error(err)
			continue
		}

		returnedDocument, err := document.Document(pointer)
		if err != nil {
			t.Error(err)
			continue
		}

		data, err := returnedDocument.String()
		if err != nil {
			t.Error(err)
			continue
		}

		if data != testString {
			t.Errorf("unexpected string, %s", data)
		}
	}
}
