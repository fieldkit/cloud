package main

import (
	"log"
)

func main() {
	ingester := NewMessageIngester()
	ingester.Schemas.AddLegacySchemas()

	err := ProcessRawMessages(ingester)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
