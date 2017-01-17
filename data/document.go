package data

import (
	"github.com/O-C-R/auth/id"

	"github.com/O-C-R/fieldkit/data/jsondocument"
)

type Document struct {
	ID        id.ID                  `json:"id" db:"id"`
	MessageID id.ID                  `json:"message_id" db:"message_id"`
	RequestID id.ID                  `json:"request_id" db:"request_id"`
	InputID   id.ID                  `json:"input_id" db:"input_id"`
	Data      *jsondocument.Document `json:"data" db:"data"`
}

func NewDocument(messageID, requestID, inputID id.ID) (*Document, error) {
	documentID, err := id.New()
	if err != nil {
		return nil, err
	}

	return &Document{
		ID:        documentID,
		MessageID: messageID,
		RequestID: requestID,
		InputID:   inputID,
	}, nil
}

type Message struct {
	ID        id.ID                  `json:"id" db:"id"`
	RequestID id.ID                  `json:"request_id" db:"request_id"`
	InputID   id.ID                  `json:"input_id" db:"input_id"`
	Data      *jsondocument.Document `json:"data" db:"data"`
}

func NewMessage(requestID, inputID id.ID) (*Message, error) {
	messageID, err := id.New()
	if err != nil {
		return nil, err
	}

	return &Message{
		ID:        messageID,
		RequestID: requestID,
		InputID:   inputID,
	}, nil
}
