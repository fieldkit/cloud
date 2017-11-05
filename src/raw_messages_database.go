package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

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
