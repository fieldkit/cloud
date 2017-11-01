package main

import (
	"log"
	"os"
)

func main() {
	ingester := NewMessageIngester()
	ingester.Schemas.AddLegacySchemas()

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
