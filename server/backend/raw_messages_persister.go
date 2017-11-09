package backend

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/O-C-R/sqlxcache"
	"io/ioutil"
	"log"
	"net/http"
)

type IncomingMessageContext struct {
	UserAgent string `json:"user-agent"`
	RequestId string `json:"request-id"`
}

type IncomingMessageParams struct {
	Headers     map[string][]string `json:"header"`
	QueryString map[string][]string `json:"querystring"`
}

type IncomingMessage struct {
	RawBody string                 `json:"body-raw"`
	Params  IncomingMessageParams  `json:"params"`
	Context IncomingMessageContext `json:"context"`
}

type RawMessageIngester struct {
	db *sqlxcache.DB
}

func NewRawMessageIngester(db *sqlxcache.DB) (*RawMessageIngester, error) {
	return &RawMessageIngester{
		db: db,
	}, nil
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

	m := IncomingMessage{
		RawBody: string(bodyBytes),
		Context: IncomingMessageContext{
			UserAgent: req.UserAgent(),
			RequestId: "",
		},
		Params: IncomingMessageParams{
			QueryString: queryString,
			Headers:     headers,
		},
	}

	data, err := json.Marshal(m)
	if err != nil {
		log.Printf("Error marshaling JSON: %s", err)
		return
	}

	rmi.db.ExecContext(context.TODO(), `INSERT INTO fieldkit.raw_message (time, origin_id, data) VALUES (now(), $1, $2)`, originId, string(data))
}
