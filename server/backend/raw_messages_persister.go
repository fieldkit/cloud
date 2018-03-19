package backend

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/ingestion"
	"github.com/fieldkit/cloud/server/data"
)

type IncomingMessageContext struct {
	UserAgent string `json:"user-agent"`
	RequestId string `json:"request-id"`
}

type IncomingMessageParams struct {
	Headers     map[string]string `json:"header"`
	QueryString map[string]string `json:"querystring"`
}

type IncomingMessage struct {
	RawBody []byte                 `json:"body-raw"`
	Params  IncomingMessageParams  `json:"params"`
	Context IncomingMessageContext `json:"context"`
}

type RawMessageIngester struct {
	incoming    chan *ingestion.RawMessageRow
	ingester    *ingestion.MessageIngester
	backend     *Backend
	db          *sqlxcache.DB
	recordAdder *RecordAdder
}

func NewRawMessageIngester(b *Backend) (rmi *RawMessageIngester, err error) {
	incoming := make(chan *ingestion.RawMessageRow, 100)

	sr := NewDatabaseSchemas(b.db)
	streams := NewDatabaseStreams(b.db)
	ingester := ingestion.NewMessageIngester(sr, streams)

	rmi = &RawMessageIngester{
		incoming:    incoming,
		ingester:    ingester,
		backend:     b,
		db:          b.db,
		recordAdder: NewRecordAdder(b),
	}

	go backgroundIngestion(rmi)

	return
}

type RecordAdder struct {
	backend *Backend
}

func NewRecordAdder(backend *Backend) *RecordAdder {
	return &RecordAdder{
		backend: backend,
	}
}

func (da *RecordAdder) EmitSourceChanged(id int64) {
	da.backend.SourceChanges <- SourceChange{
		SourceID: id,
	}
}

func (da *RecordAdder) AddRecord(ctx context.Context, im *ingestion.IngestedMessage) error {
	// TODO: Not terribly happy with this hack.
	ids := im.Schema.Ids.(DatabaseIds)
	d := data.Record{
		SchemaID:  int32(ids.SchemaID),
		SourceID:  int32(ids.DeviceID),
		TeamID:    nil,
		UserID:    nil,
		Timestamp: *im.Time,
		Location:  data.NewLocation(im.Location.Coordinates),
		Fixed:     im.Fixed,
		Visible:   true,
	}
	d.SetData(im.Fields)
	return da.backend.AddRecord(ctx, &d)
}

func backgroundIngestion(rmi *RawMessageIngester) {
	for row := range rmi.incoming {
		raw, err := ingestion.CreateRawMessageFromRow(row)
		if err != nil {
			log.Printf("(%s)[Error] %v", row.Id, err)
			log.Printf("%s", row.Data)
		} else {
			im, pm, err := rmi.ingester.Ingest(context.TODO(), raw)
			if err != nil {
				if pm != nil {
					log.Printf("(%s)(%s)[Error]: %v %s", pm.MessageId, pm.SchemaId, err, pm.ArrayValues)
				} else {
					log.Printf("(%s)[Error] %v", row.Id, err)
				}
				if true {
					log.Printf("RawMessage: contentType=%s queryString=%v", raw.ContentType, raw.QueryString)
				}
			} else {
				rmi.recordAdder.AddRecord(context.TODO(), im)
				log.Printf("(%s)(%s)[Success]", pm.MessageId, pm.SchemaId)
			}
		}
	}
}

// TODO: Eventually I'd like to see this go away. It's a relic from some Request Template stuff we had in Lambda.
func flatten(m map[string][]string) map[string]string {
	f := make(map[string]string)
	for k, v := range m {
		f[k] = v[0]
	}
	return f
}

func generateOriginId() (string, error) {
	data := make([]byte, 10)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha256.Sum256(data)), nil
}

func (rmi *RawMessageIngester) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading HTTP body: %v", err)
		return
	}

	headers := req.Header
	queryString := req.URL.Query()
	originId, err := generateOriginId()
	if err != nil {
		log.Printf("Error generating OriginId: %v", err)
		return
	}

	m := &IncomingMessage{
		RawBody: bodyBytes,
		Context: IncomingMessageContext{
			UserAgent: req.UserAgent(),
			RequestId: originId,
		},
		Params: IncomingMessageParams{
			QueryString: flatten(queryString),
			Headers:     flatten(headers),
		},
	}

	data, err := json.Marshal(m)
	if err != nil {
		log.Printf("Error marshaling JSON: %s", err)
		return
	}

	time := time.Now()

	rmi.db.ExecContext(context.TODO(), `INSERT INTO fieldkit.raw_message (time, origin_id, data) VALUES ($1, $2, $3)`, time, originId, string(data))

	rmi.incoming <- &ingestion.RawMessageRow{
		Time: uint64(time.Unix()),
		Id:   originId,
		Data: string(data),
	}
}
