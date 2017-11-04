package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"log"
	"os"
	"strings"
	"time"
)

const (
	UserAgent = "FieldKitBot/1.1"
)

// From http_provider.go
type HttpJsonMessage struct {
	Location []float32         `json:"location"`
	Time     int64             `json:"time"`
	Device   string            `json:"device"`
	Stream   string            `json:"stream"`
	Values   map[string]string `json:"values"`
}

type IncomingSqsMessageParams struct {
	Path        map[string]string `json:"path"`
	QueryString map[string]string `json:"querystring"`
	Header      map[string]string `json:"header"`
}

type IncomingSqsMessage struct {
	BodyRaw string                   `json:"body-raw"`
	Params  IncomingSqsMessageParams `json:"params"`
	Context map[string]string        `json:"context"`
}

type options struct {
	MessageId string
	Device    string
	Stream    string
}

func main() {
	o := options{}

	flag.StringVar(&o.MessageId, "id", "", "message id")
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

	m := HttpJsonMessage{
		Location: []float32{34.0647314, -118.2956819},
		Time:     time.Now().Unix(),
		Device:   o.Device,
		Stream:   o.Stream,
		Values: map[string]string{
			"cpu":  "89",
			"temp": "22",
		},
	}

	bodyRaw, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %s", err)
	}

	sqs := IncomingSqsMessage{
		BodyRaw: string(bodyRaw),
		Params: IncomingSqsMessageParams{
			Path:        map[string]string{},
			QueryString: map[string]string{},
			Header: map[string]string{
				"User-Agent":          UserAgent,
				"FieldKit-Message-Id": o.MessageId,
				"Content-Type":        "application/vnd.fk.message+json",
			},
		},
		Context: map[string]string{
			"Request-Id": o.MessageId,
			"User-Agent": UserAgent,
		},
	}

	WriteSql(&o, &sqs)
}

func WriteSql(o *options, sqs *IncomingSqsMessage) {
	d, err := json.Marshal(sqs)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %s", err)
	}
	escaped := strings.Replace(string(d), "'", "\\'", -1)
	fmt.Printf("INSERT INTO messages_raw (sqs_id, hash, time, data) VALUES ('%s', '', now(), '%s') ON CONFLICT (sqs_id) DO UPDATE SET data = '%s', time = now();\n", o.MessageId, escaped, escaped)
}
