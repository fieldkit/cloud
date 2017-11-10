package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/conservify/ingestion/ingestion"
	"github.com/nu7hatch/gouuid"
	"log"
	"os"
	"strings"
	"time"
)

const (
	UserAgent = "FieldKitBot/1.1"
)

type options struct {
	MessageId string
	Token     string
	Device    string
	Stream    string
}

func main() {
	o := options{}

	flag.StringVar(&o.MessageId, "id", "", "message id")
	flag.StringVar(&o.Token, "token", "", "token")
	flag.StringVar(&o.Device, "device", "", "device id")
	flag.StringVar(&o.Stream, "stream", "", "stream id")

	flag.Parse()

	if o.Device == "" {
		flag.Usage()
		os.Exit(1)
	}

	if o.MessageId == "" {
		id, err := uuid.NewV4()
		if err != nil {
			log.Fatalf("Error generating UUID: %s", err)
		}
		o.MessageId = id.String()
	}

	if o.Token == "" {
		id, err := uuid.NewV4()
		if err != nil {
			log.Fatalf("Error generating UUID: %s", err)
		}
		o.Token = id.String()
	}

	m := ingestion.HttpJsonMessage{
		Location: []float32{34.0647314, -118.2956819},
		Time:     time.Now().Unix(),
		Device:   o.Device,
		Stream:   o.Stream,
		Values: map[string]string{
			"cpu":  "89",
			"temp": "22",
		},
	}

	rawBody, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %s", err)
	}

	sqs := ingestion.SqsMessage{
		RawBody: string(rawBody),
		Params: ingestion.SqsMessageParams{
			QueryString: map[string]string{"token": o.Token},
			Headers: ingestion.SqsMessageHeaders{
				UserAgent:   UserAgent,
				ContentType: "application/vnd.fk.message+json",
			},
		},
		Context: ingestion.SqsMessageContext{
			RequestId: o.MessageId,
		},
	}

	WriteSql(&o, &sqs)
}

func WriteSql(o *options, sqs *ingestion.SqsMessage) {
	d, err := json.Marshal(sqs)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %s", err)
	}
	escaped := strings.Replace(string(d), "'", "\\'", -1)
	fmt.Printf("INSERT INTO messages_raw (sqs_id, hash, time, data) VALUES ('%s', '', now(), '%s') ON CONFLICT (sqs_id) DO UPDATE SET data = '%s', time = now();\n", o.MessageId, escaped, escaped)
}
