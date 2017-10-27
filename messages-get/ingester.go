package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func (i *MessageIngester) ApplySchema(nm *NormalizedMessage, ms *JsonMessageSchema) (err error) {
	if len(nm.ArrayValues) != len(ms.Fields) {
		return fmt.Errorf("%s: fields S=%v != M=%v", nm.SchemaId, len(ms.Fields), len(nm.ArrayValues))
	}

	return nil
}

func (i *MessageIngester) ApplySchemas(nm *NormalizedMessage, schemas []interface{}) (err error) {
	errors := make([]error, 0)

	if len(schemas) == 0 {
		errors = append(errors, fmt.Errorf("No matching schemas"))
	}

	for _, schema := range schemas {
		jsonSchema, ok := schema.(*JsonMessageSchema)
		if ok {
			applyErr := i.ApplySchema(nm, jsonSchema)
			if applyErr != nil {
				errors = append(errors, applyErr)
			} else {
				return
			}
		} else {
			log.Printf("Unexpected MessageSchema type: %v", schema)
		}
	}
	return fmt.Errorf("Schemas failed: %v", errors)
}

func (i *MessageIngester) HandleMessage(raw *RawMessage) error {
	rmd := RawMessageData{}
	err := json.Unmarshal([]byte(raw.Data), &rmd)
	if err != nil {
		return err
	}

	mp, err := IdentifyMessageProvider(&rmd)
	if err != nil {
		return err
	}
	if mp == nil {
		log.Printf("(%s)[Error]: No message provider: (UserAgent: %v) (ContentType: %s) Body: %s",
			rmd.Context.RequestId, rmd.Params.Headers.UserAgent, rmd.Params.Headers.ContentType, raw.Data)
		return nil
	}

	nm, err := mp.NormalizeMessage(&rmd)
	if err != nil {
		// log.Printf("%s(Error): %v", rmd.Context.RequestId, err)
	}
	if nm != nil {
		schemas, err := i.Schemas.LookupSchema(nm.SchemaId)
		if err != nil {
			return err
		}
		err = i.ApplySchemas(nm, schemas)
		if err != nil {
			if true {
				log.Printf("(%s)(%s)[Error]: %v %s", nm.MessageId, nm.SchemaId, err, nm.ArrayValues)
			} else {
				log.Printf("(%s)(%s)[Error]: %v", nm.MessageId, nm.SchemaId, err)
			}
		} else {
			if true {
				log.Printf("(%s)(%s)[Success]", nm.MessageId, nm.SchemaId)
			}
		}
	}

	return nil
}

type MessageIngester struct {
	Handler
	Schemas *SchemaRepository
}

func NewMessageIngester() *MessageIngester {
	return &MessageIngester{
		Schemas: &SchemaRepository{
			Map: make(map[SchemaId][]interface{}),
		},
	}
}
