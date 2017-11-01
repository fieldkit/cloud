package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"net/url"
	"os"
	"strings"
)

type RawMessageRow struct {
	SqsId string
	Data  string
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

type RawMessage struct {
	Row  *RawMessageRow
	Data *RawMessageData
	Form *url.Values
}

type HandlerFunc func(raw *RawMessage) error

func (f HandlerFunc) HandleMessage(raw *RawMessage) error {
	return f(raw)
}

type Handler interface {
	HandleMessage(raw *RawMessage) error
}

func CreateRawMessageFromRow(row *RawMessageRow) (raw *RawMessage, err error) {
	rmd := RawMessageData{}
	err = json.Unmarshal([]byte(row.Data), &rmd)
	if err != nil {
		return nil, err
	}

	if strings.Contains(rmd.Params.Headers.ContentType, FormUrlEncodedMimeType) {
		form, err := url.ParseQuery(rmd.RawBody)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse message content.")
		}

		raw = &RawMessage{
			Row:  row,
			Data: &rmd,
			Form: &form,
		}

		return raw, nil
	}

	return nil, fmt.Errorf("Unexpected ContentType: %s", rmd.Params.Headers.ContentType)
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
		row := &RawMessageRow{}
		rows.Scan(&row.SqsId, &row.Data)

		raw, err := CreateRawMessageFromRow(row)
		if err != nil {
			return err
		}

		err = h.HandleMessage(raw)
		if err != nil {
			return err
		}
	}

	return nil
}
