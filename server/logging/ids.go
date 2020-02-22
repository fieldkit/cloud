package logging

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"sync/atomic"
)

func MakeShortID() string {
	b := make([]byte, 6)
	io.ReadFull(rand.Reader, b)
	return base64.StdEncoding.EncodeToString(b)
}

type IdGenerator struct {
	id     int64
	prefix string
}

// algorithm taken from goa RequestId middleware
// algorithm taken from https://github.com/zenazn/goji/blob/master/web/middleware/request_id.go#L44-L50
func MakeCommonPrefix() string {
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}
	return string(b64[0:10])
}

func NewIdGenerator() *IdGenerator {
	return &IdGenerator{
		id:     0,
		prefix: MakeCommonPrefix(),
	}
}

func (g *IdGenerator) Generate() string {
	return fmt.Sprintf("%s-%d", g.prefix, atomic.AddInt64(&g.id, 1))
}
