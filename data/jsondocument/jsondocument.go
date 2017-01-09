package jsondocument

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	arrayIndexRegexp          *regexp.Regexp
	jsonPointerSyntaxError    = errors.New("JSON pointer syntax error")
	jsonPointerReferenceError = errors.New("JSON pointer reference error")
	jsonDelimiterError        = errors.New("unknown JSON delimiter")
	jsonTypeError             = errors.New("unexpected JSON type")
)

func init() {
	arrayIndexRegexp = regexp.MustCompile(`^(?:0|(?:[1-9]\d*))$`)
}

func arrayIndex(index string) int {
	uint64Index, err := strconv.ParseUint(index, 10, 64)
	if err != nil {
		panic(err)
	}

	return int(uint64Index)
}

type Document struct {
	data interface{}
}

func (d *Document) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.data)
}

func (d *Document) UnmarshalJSON(data []byte) error {
	token, err := json.NewDecoder(bytes.NewReader(data)).Token()
	if err != nil {
		return err
	}

	if delim, ok := token.(json.Delim); ok {
		switch delim {
		case '[':
			array := []*Document{}
			if err := json.Unmarshal(data, &array); err != nil {
				return err
			}

			d.data = array
			return nil
		case '{':
			object := map[string]*Document{}
			if err := json.Unmarshal(data, &object); err != nil {
				return err
			}

			d.data = object
			return nil
		}

		return jsonDelimiterError
	}

	return json.Unmarshal(data, &d.data)
}

func Array() *Document {
	return &Document{[]*Document{}}
}

func Object() *Document {
	return &Document{map[string]*Document{}}
}

func String(value string) *Document {
	return &Document{value}
}

func Number(value float64) *Document {
	return &Document{value}
}

func Boolean(value bool) *Document {
	return &Document{value}
}

func Null() *Document {
	return &Document{}
}

func (d *Document) String() (string, error) {
	value, ok := d.data.(string)
	if !ok {
		return "", jsonTypeError
	}

	return value, nil
}

func (d *Document) Number() (float64, error) {
	value, ok := d.data.(float64)
	if !ok {
		return .0, jsonTypeError
	}

	return value, nil
}

func (d *Document) Boolean() (bool, error) {
	value, ok := d.data.(bool)
	if !ok {
		return false, jsonTypeError
	}

	return value, nil
}

func (d *Document) ArrayIndex(index int) (*Document, error) {
	array, ok := d.data.([]*Document)
	if !ok {
		return nil, jsonTypeError
	}

	if index >= len(array) {
		return nil, jsonTypeError
	}

	return array[index], nil
}

func (d *Document) Copy() *Document {
	var c *Document
	switch data := d.data.(type) {
	case string, float64, bool, nil:
		c = &Document{data}
	case []*Document:
		newArray := make([]*Document, len(data))
		for i, value := range data {
			newArray[i] = value.Copy()
		}

		c = &Document{newArray}
	case map[string]*Document:
		newObject := make(map[string]*Document)
		for key, value := range data {
			newObject[key] = value.Copy()
		}

		c = &Document{newObject}
	}

	return c
}

func parsePointer(pointer string) ([]string, error) {
	if pointer == "" {
		return []string{}, nil
	}

	pointerComponents := strings.Split(pointer, "/")
	if pointerComponents[0] != "" {
		return nil, jsonPointerSyntaxError
	}

	return pointerComponents[1:], nil
}

func (d *Document) SetDocument(pointer string, value *Document) error {
	references, err := parsePointer(pointer)
	if err != nil {
		return err
	}

	document := d
	for _, referenceComponent := range references {

		// A dash indicates the reference is to the position after the last
		// element of the array.
		if referenceComponent == "-" {
			if document.data == nil {
				document.data = make([]*Document, 0, 1)
			}

			if array, ok := document.data.([]*Document); ok {
				newDocument := Null()
				document.data = append(array, newDocument)
				document = newDocument
				continue
			}

			return jsonPointerReferenceError
		}

		if arrayIndexRegexp.MatchString(referenceComponent) {
			index := arrayIndex(referenceComponent)
			if document.data == nil {
				array := make([]*Document, index+1)
				document.data = array
				array[index] = Null()
				document = array[index]
				continue
			}

			if array, ok := document.data.([]*Document); ok {
				if len(array) <= index {
					array = append(array, make([]*Document, index-len(array)+1)...)
					document.data = array
				}

				if array[index] == nil {
					array[index] = Null()
				}

				document = array[index]
				continue
			}

			return jsonPointerReferenceError
		}

		if document.data == nil {
			document.data = map[string]*Document{}
		}

		if object, ok := document.data.(map[string]*Document); ok {
			if keyValue, ok := object[referenceComponent]; ok && keyValue != nil {
				document = keyValue
				continue
			}

			object[referenceComponent] = Null()
			document = object[referenceComponent]
			continue
		}

		return jsonPointerReferenceError
	}

	document.data = value.Copy().data
	return nil
}

func (d *Document) Document(pointer string) (*Document, error) {
	references, err := parsePointer(pointer)
	if err != nil {
		return nil, err
	}

	document := d
	for _, referenceComponent := range references {
		if arrayIndexRegexp.MatchString(referenceComponent) {
			array, ok := document.data.([]*Document)
			if !ok {
				return nil, jsonPointerReferenceError
			}

			index := arrayIndex(referenceComponent)
			if index >= len(array) {
				return nil, jsonPointerReferenceError
			}

			document = array[index]
		} else {
			object, ok := document.data.(map[string]*Document)
			if !ok {
				return nil, jsonPointerReferenceError
			}

			value, ok := object[referenceComponent]
			if !ok {
				return nil, jsonPointerReferenceError
			}

			document = value
		}
	}

	return document, nil
}
