package main

import (
	"database/sql"
	"fmt"
	"github.com/fieldkit/cloud/server/ingestion"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	sr := ingestion.NewInMemorySchemas()
	ingestion.AddLegacySchemas(sr)
	streams := ingestion.NewInMemoryStreams()
	ingester := ingestion.NewMessageIngester(sr, streams)

	dbOptions := MessageDatabaseOptions{
		Hostname: os.Getenv("DATABASE_HOSTNAME"),
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Database: os.Getenv("DATABASE_NAME"),
	}

	err := ProcessRawMessages(&dbOptions, ingester)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Done: %+v {Failures:%d}", ingester.Statistics, ingester.Statistics.Processed-ingester.Statistics.Successes)
}

type MessageDatabaseOptions struct {
	Hostname string
	User     string
	Password string
	Database string
}

func ProcessRawMessages(o *MessageDatabaseOptions, h ingestion.RawMessageHandler) error {
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
		row := &ingestion.RawMessageRow{}
		rows.Scan(&row.Id, &row.Data, &row.Time)

		raw, err := ingestion.CreateRawMessageFromRow(row)
		if err != nil {
			return fmt.Errorf("(%s)[Error] %v", row.Id, err)
		}

		err = h.HandleMessage(raw)
		if err != nil {
			return fmt.Errorf("(%s)[Error] %v", row.Id, err)
		}
	}

	return nil
}
