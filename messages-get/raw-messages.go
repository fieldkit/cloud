package main

import (
	"database/sql"
	_ "fmt"
	_ "github.com/lib/pq"
	"os"
)

type RawMessage struct {
	SqsId string
	Data  string
}

type HandlerFunc func(raw *RawMessage) error

func (f HandlerFunc) HandleMessage(raw *RawMessage) error {
	return f(raw)
}

type Handler interface {
	HandleMessage(raw *RawMessage) error
}

type RawMessageHeaders struct {
	UserAgent   string `json:"User-Agent"`
	ContentType string `json:"Content-Type"`
}

type RawMessageParams struct {
	Headers RawMessageHeaders `json:"header"`
}

type RawMessageContext struct {
	UserAgent string `json:"user-agent"`
	RequestId string `json:"request-id"`
}

type RawMessageData struct {
	RawBody string            `json:"body-raw"`
	Params  RawMessageParams  `json:"params"`
	Context RawMessageContext `json:"context"`
}

func ProcessRawMessages(h Handler) error {
	databaseHostname := os.Getenv("DATABASE_HOSTNAME")
	databaseUser := os.Getenv("DATABASE_USER")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databaseName := os.Getenv("DATABASE_NAME")

	connectionString := "postgres://" + databaseUser + ":" + databasePassword + "@" + databaseHostname + "/" + databaseName + "?sslmode=disable"

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	defer db.Close()

	rows, err := db.Query("SELECT sqs_id, data FROM messages_raw ORDER BY time")
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		raw := &RawMessage{}
		rows.Scan(&raw.SqsId, &raw.Data)
		err = h.HandleMessage(raw)
		if err != nil {
			return err
		}
	}

	return nil
}
