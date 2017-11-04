package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"net/url"
	"strings"
)

type RawMessageRow struct {
	SqsId string
	Data  string
	Time  uint64
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
		return nil, fmt.Errorf("Malformed RawMessage: %s", row.Data)
	}

	if rmd.Context.RequestId == "" {
		return nil, fmt.Errorf("Malformed RawMessage: Context missing RequestId.")
	}

	if strings.Contains(rmd.Params.Headers.ContentType, FormUrlEncodedMimeType) {
		form, err := url.ParseQuery(rmd.RawBody)
		if err != nil {
			return nil, fmt.Errorf("Malformed RawMessage: RawBody invalid for %s.", rmd.Params.Headers.ContentType)
		}

		raw = &RawMessage{
			Row:  row,
			Data: &rmd,
			Form: &form,
		}

		return raw, nil
	}

	if strings.Contains(rmd.Params.Headers.ContentType, JsonMimeType) {
		raw = &RawMessage{
			Row:  row,
			Data: &rmd,
			Form: &url.Values{},
		}

		return raw, nil
	}

	return nil, fmt.Errorf("Unexpected ContentType: %s", rmd.Params.Headers.ContentType)
}

type MessageDatabaseOptions struct {
	Hostname string
	User     string
	Password string
	Database string
}

func ProcessRawMessages(o *MessageDatabaseOptions, h Handler) error {
	connectionString := "postgres://" + o.User + ":" + o.Password + "@" + o.Hostname + "/" + o.Database + "?sslmode=disable"

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	defer db.Close()

	rows, err := db.Query("SELECT sqs_id, data, time FROM messages_raw ORDER BY time")
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		row := &RawMessageRow{}
		rows.Scan(&row.SqsId, &row.Data, &row.Time)

		// fmt.Printf("%s,%s\n", row.SqsId, row.Data)

		raw, err := CreateRawMessageFromRow(row)
		if err != nil {
			return fmt.Errorf("(%s)[Error] %v", row.SqsId, err)
		}

		err = h.HandleMessage(raw)
		if err != nil {
			return fmt.Errorf("(%s)[Error] %v", row.SqsId, err)
		}
	}

	return nil
}
