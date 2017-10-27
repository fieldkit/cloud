package main

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"log"
	"math"
	"net/url"
	"os"
	"regexp"
	_ "strconv"
	"strings"
	"unicode"

	"github.com/golang/protobuf/proto"
)

type RawMessageHeaders struct {
	UserAgent   string `json:"User-Agent"`
	ContentType string `json:"Content-Type"`
}

type RawMessageParams struct {
	Headers RawMessageHeaders `json:"header"`
}

type RawMessageContext struct {
	UserAgent string `json:"user-agent"`
}

type RawMessageData struct {
	RawBody string            `json:"body-raw"`
	Params  RawMessageParams  `json:"params"`
	Context RawMessageContext `json:"context"`
}

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

	rows, err := db.Query("SELECT sqs_id, data FROM messages_raw")
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

type SchemaId string

type NormalizedMessage struct {
	SchemaId    SchemaId
	Time        uint64
	ArrayValues []string
}

type MessageProvider interface {
	NormalizeMessage(rmd *RawMessageData) (nm *NormalizedMessage, err error)
}

type MessageProviderBase struct {
	Form url.Values
}

type ParticleMessageProvider struct {
	MessageProviderBase
}

func (i *ParticleMessageProvider) NormalizeMessage(rmd *RawMessageData) (nm *NormalizedMessage, err error) {
	return nil, nil
}

type TwilioMessageProvider struct {
	MessageProviderBase
}

func (i *TwilioMessageProvider) NormalizeMessage(rmd *RawMessageData) (nm *NormalizedMessage, err error) {
	trimmed := strings.TrimSpace(i.Form.Get("Body"))
	fields := strings.Split(trimmed, ",")
	if len(fields) < 2 {
		return nil, nil
	}
	nm = &NormalizedMessage{
		SchemaId:    SchemaId(i.Form.Get("From") + "-" + fields[1]),
		ArrayValues: fields,
	}
	return nil, nil
}

type RockBlockMessageProvider struct {
	MessageProviderBase
}

func (i *RockBlockMessageProvider) NormalizeMessage(rmd *RawMessageData) (nm *NormalizedMessage, err error) {
	// Data from RockBlock is always hex encoded.
	data := i.Form.Get("data")
	if len(data) == 0 {
		// TODO: We should annotate incoming messages with information about their failure for logging/debugging.
		// debug.Log(MessageId, Bla bla)
		return
	}

	bytes, err := hex.DecodeString(data)
	if err != nil {
		return
	}

	serial := i.Form.Get("serial")

	if unicode.IsPrint(rune(bytes[0])) {
		trimmed := strings.TrimSpace(string(bytes))
		fields := strings.Split(trimmed, ",")
		nm = &NormalizedMessage{
			SchemaId:    SchemaId(serial + "-" + fields[2]),
			Time:        0,
			ArrayValues: fields,
		}
	} else {
		// This is a protobuf message or some other kind of similar low level binary.
		buffer := proto.NewBuffer(bytes)

		// NOTE: Right now we're only dealing with the binary format we
		// came up with during the 'NatGeo demo phase' Eventually this
		// will be a property message we can just slurp up. Though, maybe this
		// will be a great RB message going forward?
		id, err := buffer.DecodeVarint()
		if err != nil {
			return nil, err
		}
		time, err := buffer.DecodeVarint()
		if err != nil {
			return nil, err
		}

		// HACK
		values := make([]string, 0)

		for {
			f64, err := buffer.DecodeFixed32()
			if err == io.ErrUnexpectedEOF {
				break
			} else if err != nil {
				return nil, err
			}

			value := math.Float32frombits(uint32(f64))
			values = append(values, fmt.Sprintf("%f", value))
		}

		nm = &NormalizedMessage{
			SchemaId:    SchemaId(serial + "-" + string(id)),
			Time:        time,
			ArrayValues: values,
		}
	}
	return
}

func IdentifyMessageProvider(rmd *RawMessageData) (t MessageProvider, err error) {
	if strings.Contains(rmd.Params.Headers.ContentType, "x-www-form-urlencoded") {
		form, err := url.ParseQuery(rmd.RawBody)
		if err != nil {
			return nil, nil
		}

		if form.Get("device_type") == "ROCKBLOCK" {
			t = &RockBlockMessageProvider{
				MessageProviderBase: MessageProviderBase{Form: form},
			}
		} else if form.Get("coreid") != "" {
			t = &ParticleMessageProvider{
				MessageProviderBase: MessageProviderBase{Form: form},
			}
		} else if form.Get("SmsSid") != "" {
			t = &TwilioMessageProvider{
				MessageProviderBase: MessageProviderBase{Form: form},
			}
		}
	}

	return t, nil

}

func (i *MessageIngester) HandleMessage(raw *RawMessage) error {
	rmd := RawMessageData{}
	err := json.Unmarshal([]byte(raw.Data), &rmd)
	if err != nil {
		return err
	}

	mp, err := IdentifyMessageProvider(&rmd)
	if err != nil {
		return err
	}
	if mp == nil {
		fmt.Printf("Unknown Message: (UserAgent: %v) (ContentType: %s) Body: %s\n",
			rmd.Params.Headers.UserAgent, rmd.Params.Headers.ContentType, raw.Data)
		return nil
	}

	nm, err := mp.NormalizeMessage(&rmd)
	if err != nil {
		return err
	}
	if nm != nil {
		schema, err := i.Schemas.LookupSchema(nm.SchemaId)
		if err != nil {
			return err
		}
		if schema != nil {
			fmt.Printf("%v %v %v\n", nm.SchemaId, *nm, schema)
		}
	}

	return nil
}

type MessageIngester struct {
	Handler
	Schemas *SchemaRepository
}

func NewMessageIngester() *MessageIngester {
	return &MessageIngester{
		Schemas: &SchemaRepository{
			Map: make(map[SchemaId]interface{}),
		},
	}
}

func main() {
	ingester := NewMessageIngester()
	ingester.Schemas.DefineSchema("(\\d+)-ST", &JsonMessageSchema{
		Fields: []JsonSchemaField{
			JsonSchemaField{Name: "Battery Voltage", Type: "float32"},
			JsonSchemaField{Name: "Battery Percentage", Type: "float32"},

			JsonSchemaField{Name: "SD Available", Type: "float32"},
			JsonSchemaField{Name: "GPS Fix", Type: "float32"},
			JsonSchemaField{Name: "Battery Sleep Time", Type: "float32"},
			JsonSchemaField{Name: "Deep Sleep Time", Type: "float32"},
			JsonSchemaField{Name: "Tx Failures", Type: "float32"},
			JsonSchemaField{Name: "Tx Skips", Type: "float32"},
			JsonSchemaField{Name: "Weather Readings Received", Type: "float32"},
			JsonSchemaField{Name: "Atlas Packets Rx", Type: "float32"},
			JsonSchemaField{Name: "Sonar Packets Rx", Type: "float32"},
			JsonSchemaField{Name: "Dead For", Type: "float32"},

			JsonSchemaField{Name: "Avg Tx Time", Type: "float32"},

			JsonSchemaField{Name: "Idle Interval", Type: "float32"},
			JsonSchemaField{Name: "Airwaves Interval", Type: "float32"},
			JsonSchemaField{Name: "Weather Interval", Type: "float32"},
			JsonSchemaField{Name: "Weather Start", Type: "float32"},
			JsonSchemaField{Name: "Weather Ignore", Type: "float32"},
			JsonSchemaField{Name: "Weather Off", Type: "float32"},
			JsonSchemaField{Name: "Weather Reading", Type: "float32"},

			JsonSchemaField{Name: "Tx A Offset", Type: "float32"},
			JsonSchemaField{Name: "Tx A Interval", Type: "float32"},
			JsonSchemaField{Name: "Tx B Offset", Type: "float32"},
			JsonSchemaField{Name: "Tx B Interval", Type: "float32"},
			JsonSchemaField{Name: "Tx C Offset", Type: "float32"},
			JsonSchemaField{Name: "Tx C Interval", Type: "float32"},
		},
	})

	err := ProcessRawMessages(ingester)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

type SchemaRepository struct {
	Map map[SchemaId]interface{}
}

func (sr *SchemaRepository) DefineSchema(id SchemaId, ms interface{}) (err error) {
	sr.Map[id] = ms
	return
}

func (sr *SchemaRepository) LookupSchema(id SchemaId) (ms interface{}, err error) {
	// TODO: This will become very slow, just to prove the concept.
	for key, value := range sr.Map {
		re := regexp.MustCompile(string(key))
		if re.MatchString(string(id)) {
			return value, nil
		}
	}
	return
}

type MessageSchema struct {
}

type JsonSchemaField struct {
	Name string
	Type string
}

type JsonMessageSchema struct {
	MessageSchema
	Fields []JsonSchemaField
}

type TwitterMessageSchema struct {
	MessageSchema
}
