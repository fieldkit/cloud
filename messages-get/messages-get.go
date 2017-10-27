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
	RequestId string `json:"request-id"`
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

	rows, err := db.Query("SELECT sqs_id, data FROM messages_raw ORDER BY time")
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		raw := &RawMessage{}
		rows.Scan(&raw.SqsId, &raw.Data)
		if false {
			fmt.Printf("DATA: %s\n", raw.Data)
		}
		err = h.HandleMessage(raw)
		if err != nil {
			return err
		}
	}

	return nil
}

type MessageId string

type SchemaId string

type NormalizedMessage struct {
	MessageId   MessageId
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
		MessageId:   MessageId(rmd.Context.RequestId),
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
		stationNameRe := regexp.MustCompile("[A-Z][A-Z]")
		trimmed := strings.TrimSpace(string(bytes))
		fields := strings.Split(trimmed, ",")
		maybeStationName := fields[2]
		if stationNameRe.MatchString(maybeStationName) {
			nm = &NormalizedMessage{
				MessageId:   MessageId(rmd.Context.RequestId),
				SchemaId:    SchemaId(fmt.Sprintf("%s-%s", serial, maybeStationName)),
				Time:        0,
				ArrayValues: fields,
			}
		} else {
			return nil, fmt.Errorf("Invalid station name: %s", maybeStationName)
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
			MessageId:   MessageId(rmd.Context.RequestId),
			SchemaId:    SchemaId(fmt.Sprintf("%s-%d", serial, id)),
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

func (nm *NormalizedMessage) ApplySchema(ms *JsonMessageSchema) (err error) {
	if len(nm.ArrayValues) != len(ms.Fields) {
		return fmt.Errorf("%s: fields S=%v != M=%v", nm.SchemaId, len(ms.Fields), len(nm.ArrayValues))
	}

	return nil
}

func (nm *NormalizedMessage) ApplySchemas(schemas []interface{}) (err error) {
	errors := make([]error, 0)

	if len(schemas) == 0 {
		errors = append(errors, fmt.Errorf("No matching schemas"))
	}

	for _, schema := range schemas {
		jsonSchema, ok := schema.(*JsonMessageSchema)
		if ok {
			applyErr := nm.ApplySchema(jsonSchema)
			if applyErr != nil {
				errors = append(errors, applyErr)
			} else {
				return
			}
		} else {
			log.Printf("Unexpected MessageSchema type: %v", schema)
		}
	}
	return fmt.Errorf("Schemas failed: %v", errors)
}

func ShouldSkip(nm *NormalizedMessage) bool {
	return false
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
		log.Printf("(%s)[Error]: No message provider: (UserAgent: %v) (ContentType: %s) Body: %s",
			rmd.Context.RequestId, rmd.Params.Headers.UserAgent, rmd.Params.Headers.ContentType, raw.Data)
		return nil
	}

	nm, err := mp.NormalizeMessage(&rmd)
	if err != nil {
		// log.Printf("%s(Error): %v", rmd.Context.RequestId, err)
	}
	if nm != nil {
		if ShouldSkip(nm) {

		} else {
			schemas, err := i.Schemas.LookupSchema(nm.SchemaId)
			if err != nil {
				return err
			}
			err = nm.ApplySchemas(schemas)
			if err != nil {
				if true {
					log.Printf("(%s)(%s)[Error]: %v %s", nm.MessageId, nm.SchemaId, err, nm.ArrayValues)
				} else {
					log.Printf("(%s)(%s)[Error]: %v", nm.MessageId, nm.SchemaId, err)
				}
			} else {
				if true {
					log.Printf("(%s)(%s)[Success]", nm.MessageId, nm.SchemaId)
				}
			}
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
			Map: make(map[SchemaId][]interface{}),
		},
	}
}

func main() {
	ingester := NewMessageIngester()
	ingester.Schemas.DefineSchema("(\\d+)-AT", &JsonMessageSchema{
		Fields: []JsonSchemaField{
			JsonSchemaField{Name: "Time", Type: "uint32"},
			JsonSchemaField{Name: "Station", Type: "string"},
			JsonSchemaField{Name: "Type", Type: "string"},
			JsonSchemaField{Name: "Battery Voltage", Type: "float32"},
			JsonSchemaField{Name: "Battery Percentage", Type: "float32"},

			JsonSchemaField{Name: "Atlas Battery", Type: "float32"},
			JsonSchemaField{Name: "Orp", Type: "float32"},
			JsonSchemaField{Name: "DO", Type: "float32"},
			JsonSchemaField{Name: "pH", Type: "float32"},
			JsonSchemaField{Name: "EC 1", Type: "float32"},
			JsonSchemaField{Name: "EC 2", Type: "float32"},
			JsonSchemaField{Name: "EC 3", Type: "float32"},
			JsonSchemaField{Name: "EC 4", Type: "float32"},
			JsonSchemaField{Name: "Water Temp", Type: "float32"},
			JsonSchemaField{Name: "Ignore", Type: "float32"},
			JsonSchemaField{Name: "Ignore", Type: "float32"},
			JsonSchemaField{Name: "Ignore", Type: "float32"},
		},
	})

	ingester.Schemas.DefineSchema("(\\d+)-WE", &JsonMessageSchema{
		Fields: []JsonSchemaField{
			JsonSchemaField{Name: "Time", Type: "uint32"},
			JsonSchemaField{Name: "Station", Type: "string"},
			JsonSchemaField{Name: "Type", Type: "string"},
			JsonSchemaField{Name: "Battery Voltage", Type: "float32"},
			JsonSchemaField{Name: "Battery Percentage", Type: "float32"},
			JsonSchemaField{Name: "Old Voltage", Type: "float32"},

			JsonSchemaField{Name: "Temp", Type: "float32"},
			JsonSchemaField{Name: "Humidity", Type: "float32"},
			JsonSchemaField{Name: "Pressure", Type: "float32"},
			JsonSchemaField{Name: "WindSpeed2", Type: "float32"},
			JsonSchemaField{Name: "WindDir2", Type: "float32"},
			JsonSchemaField{Name: "WindGust10", Type: "float32"},
			JsonSchemaField{Name: "WindGustDir10", Type: "float32"},
			JsonSchemaField{Name: "DailyRain", Type: "float32"},
		},
	})

	ingester.Schemas.DefineSchema("(\\d+)-SO", &JsonMessageSchema{
		Fields: []JsonSchemaField{
			JsonSchemaField{Name: "Time", Type: "uint32"},
			JsonSchemaField{Name: "Station", Type: "string"},
			JsonSchemaField{Name: "Type", Type: "string"},
			JsonSchemaField{Name: "Battery Voltage", Type: "float32"},
			JsonSchemaField{Name: "Battery Percentage", Type: "float32"},

			JsonSchemaField{Name: "Sonar Battery", Type: "float32"},
			JsonSchemaField{Name: "Depth 1", Type: "float32"},
			JsonSchemaField{Name: "Depth 2", Type: "float32"},
			JsonSchemaField{Name: "Depth 3", Type: "float32"},
			JsonSchemaField{Name: "Depth 4", Type: "float32"},
			JsonSchemaField{Name: "Depth 5", Type: "float32"},
		},
	})

	ingester.Schemas.DefineSchema("(\\d+)-LO", &JsonMessageSchema{
		Fields: []JsonSchemaField{
			JsonSchemaField{Name: "Time", Type: "uint32"},
			JsonSchemaField{Name: "Station", Type: "string"},
			JsonSchemaField{Name: "Type", Type: "string"},
			JsonSchemaField{Name: "Battery Voltage", Type: "float32"},
			JsonSchemaField{Name: "Battery Percentage", Type: "float32"},

			JsonSchemaField{Name: "Latitude", Type: "float32"},
			JsonSchemaField{Name: "Longitude", Type: "float32"},
			JsonSchemaField{Name: "Altitude", Type: "float32"},
			JsonSchemaField{Name: "Uptime", Type: "float32"},
		},
	})

	ingester.Schemas.DefineSchema("(\\d+)-ST", &JsonMessageSchema{
		Fields: []JsonSchemaField{
			JsonSchemaField{Name: "Time", Type: "uint32"},
			JsonSchemaField{Name: "Station", Type: "string"},
			JsonSchemaField{Name: "Type", Type: "string"},
			JsonSchemaField{Name: "Battery Voltage", Type: "float32"},
			JsonSchemaField{Name: "Battery Percentage", Type: "float32"},

			JsonSchemaField{Name: "SD Available", Type: "bool"},
			JsonSchemaField{Name: "GPS Fix", Type: "bool"},
			JsonSchemaField{Name: "Battery Sleep Time", Type: "uint32"},
			JsonSchemaField{Name: "Tx Failures", Type: "uint32"},
			JsonSchemaField{Name: "Tx Skips", Type: "uint32"},
			JsonSchemaField{Name: "Weather Readings Received", Type: "uint32"},
			JsonSchemaField{Name: "Atlas Packets Rx", Type: "uint32"},
			JsonSchemaField{Name: "Sonar Packets Rx", Type: "uint32"},
			JsonSchemaField{Name: "Dead For", Type: "uint32"},

			JsonSchemaField{Name: "Avg Tx Time", Type: "float32"},

			JsonSchemaField{Name: "Idle Interval", Type: "uint32"},
			JsonSchemaField{Name: "Airwaves Interval", Type: "uint32"},
			JsonSchemaField{Name: "Weather Interval", Type: "uint32"},
			JsonSchemaField{Name: "Weather Start", Type: "uint32"},
			JsonSchemaField{Name: "Weather Ignore", Type: "uint32"},
			JsonSchemaField{Name: "Weather Off", Type: "uint32"},
			JsonSchemaField{Name: "Weather Reading", Type: "uint32"},

			JsonSchemaField{Name: "Tx A Offset", Type: "uint32"},
			JsonSchemaField{Name: "Tx A Interval", Type: "uint32"},
			JsonSchemaField{Name: "Tx B Offset", Type: "uint32"},
			JsonSchemaField{Name: "Tx B Interval", Type: "uint32"},
			JsonSchemaField{Name: "Tx C Offset", Type: "uint32"},
			JsonSchemaField{Name: "Tx C Interval", Type: "uint32"},

			JsonSchemaField{Name: "Uptime", Type: "uint32"},
		},
	})

	ingester.Schemas.DefineSchema("(\\d+)-ST", &JsonMessageSchema{
		Fields: []JsonSchemaField{
			JsonSchemaField{Name: "Time", Type: "uint32"},
			JsonSchemaField{Name: "Station", Type: "string"},
			JsonSchemaField{Name: "Type", Type: "string"},
			JsonSchemaField{Name: "Battery Voltage", Type: "float32"},
			JsonSchemaField{Name: "Battery Percentage", Type: "float32"},

			JsonSchemaField{Name: "SD Available", Type: "bool"},
			JsonSchemaField{Name: "GPS Fix", Type: "bool"},
			JsonSchemaField{Name: "Battery Sleep Time", Type: "uint32"},
			JsonSchemaField{Name: "Deep Sleep Time", Type: "uint32", Optional: true},
			JsonSchemaField{Name: "Tx Failures", Type: "uint32"},
			JsonSchemaField{Name: "Tx Skips", Type: "uint32"},
			JsonSchemaField{Name: "Weather Readings Received", Type: "uint32"},
			JsonSchemaField{Name: "Atlas Packets Rx", Type: "uint32"},
			JsonSchemaField{Name: "Sonar Packets Rx", Type: "uint32"},
			JsonSchemaField{Name: "Dead For", Type: "uint32"},

			JsonSchemaField{Name: "Avg Tx Time", Type: "float32"},

			JsonSchemaField{Name: "Idle Interval", Type: "uint32"},
			JsonSchemaField{Name: "Airwaves Interval", Type: "uint32"},
			JsonSchemaField{Name: "Weather Interval", Type: "uint32"},
			JsonSchemaField{Name: "Weather Start", Type: "uint32"},
			JsonSchemaField{Name: "Weather Ignore", Type: "uint32"},
			JsonSchemaField{Name: "Weather Off", Type: "uint32"},
			JsonSchemaField{Name: "Weather Reading", Type: "uint32"},

			JsonSchemaField{Name: "Tx A Offset", Type: "uint32"},
			JsonSchemaField{Name: "Tx A Interval", Type: "uint32"},
			JsonSchemaField{Name: "Tx B Offset", Type: "uint32"},
			JsonSchemaField{Name: "Tx B Interval", Type: "uint32"},
			JsonSchemaField{Name: "Tx C Offset", Type: "uint32"},
			JsonSchemaField{Name: "Tx C Interval", Type: "uint32"},

			JsonSchemaField{Name: "Uptime", Type: "uint32"},
		},
	})

	ingester.Schemas.DefineSchema("(\\d+)-1", &JsonMessageSchema{
		Fields: []JsonSchemaField{
			JsonSchemaField{Name: "Latitude", Type: "float32"},
			JsonSchemaField{Name: "Longitude", Type: "float32"},
			JsonSchemaField{Name: "Altitude", Type: "float32"},
			JsonSchemaField{Name: "Temperature", Type: "float32"},
			JsonSchemaField{Name: "Humidity", Type: "float32"},
			JsonSchemaField{Name: "Battery Percentage", Type: "float32"},
			JsonSchemaField{Name: "Battery Voltage", Type: "float32"},
			JsonSchemaField{Name: "Uptime", Type: "uint32"},
		},
	})

	err := ProcessRawMessages(ingester)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

type SchemaRepository struct {
	Map map[SchemaId][]interface{}
}

func (sr *SchemaRepository) DefineSchema(id SchemaId, ms interface{}) (err error) {
	if sr.Map[id] == nil {
		sr.Map[id] = make([]interface{}, 0)
	}
	sr.Map[id] = append(sr.Map[id], ms)
	return
}

func (sr *SchemaRepository) LookupSchema(id SchemaId) (ms []interface{}, err error) {
	ms = make([]interface{}, 0)

	// TODO: This will become very slow, just to prove the concept.
	for key, value := range sr.Map {
		re := regexp.MustCompile(string(key))
		if re.MatchString(string(id)) {
			ms = append(ms, value...)
		}
	}
	return
}

type MessageSchema struct {
}

type JsonSchemaField struct {
	Name     string
	Type     string
	Optional bool
}

type JsonMessageSchema struct {
	MessageSchema
	Fields []JsonSchemaField
}

type TwitterMessageSchema struct {
	MessageSchema
}
