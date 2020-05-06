package ingester

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"

	"github.com/bxcodec/faker/v3"

	"github.com/fieldkit/cloud/server/common"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/tests"
)

func TestAnyUnauthenticatedGetReturnsUnauthorized(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	handler, _, err := NewTestableIngester(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", "/why-would-this-path-be-handled", nil)
	rr := tests.ExecuteRequest(req, handler)

	assert.Equal(http.StatusUnauthorized, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `{
		"id": "<<PRESENCE>>",
		"code": "jwt_security_error",
		"detail": "<<PRESENCE>>",
		"status": "<<PRESENCE>>"
	}`)
}

func TestAnyInvalidAuthorization(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	handler, _, err := NewTestableIngester(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", "/why-would-this-path-be-handled", nil)
	req.Header.Add("Authorization", "Bearer TOTALLY-INVALID")
	rr := tests.ExecuteRequest(req, handler)

	assert.Equal(http.StatusUnauthorized, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `{
		"id": "<<PRESENCE>>",
		"code": "jwt_security_error",
		"detail": "<<PRESENCE>>",
		"status": "<<PRESENCE>>"
	}`)
}

func TestAnyGetReturnsStatusOK(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	handler, user, err := NewTestableIngester(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", "/why-would-this-path-be-handled", nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(user))
	rr := tests.ExecuteRequest(req, handler)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `{}`)
}

func TestIngestWithoutContentType(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	handler, user, err := NewTestableIngester(e)
	assert.NoError(err)

	data, err := NewRandomData(1024)
	assert.NoError(err)

	req, _ := http.NewRequest("POST", "/", nil)
	assert.NoError(AddIngestionHeaders(e, req, user, data))
	req.Header.Del(common.ContentTypeHeaderName)
	rr := tests.ExecuteRequest(req, handler)

	assert.Equal(http.StatusBadRequest, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"message": "invalid headers",
		"headers": "<<PRESENCE>>"
	}`)
}

func TestIngestWithoutContentLength(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	handler, user, err := NewTestableIngester(e)
	assert.NoError(err)

	data, err := NewRandomData(1024)
	assert.NoError(err)

	req, _ := http.NewRequest("POST", "/", nil)
	assert.NoError(AddIngestionHeaders(e, req, user, data))
	req.Header.Del(common.ContentLengthHeaderName)
	rr := tests.ExecuteRequest(req, handler)

	assert.Equal(http.StatusBadRequest, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"message": "invalid headers",
		"headers": "<<PRESENCE>>"
	}`)
}

func TestIngestWithZeroContentLength(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	handler, user, err := NewTestableIngester(e)
	assert.NoError(err)

	req, _ := http.NewRequest("POST", "/", nil)
	assert.NoError(AddIngestionHeaders(e, req, user, []byte{}))
	rr := tests.ExecuteRequest(req, handler)

	assert.Equal(http.StatusBadRequest, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"message": "invalid headers",
		"headers": "<<PRESENCE>>"
	}`)
}

func TestIngestWithoutDeviceID(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	handler, user, err := NewTestableIngester(e)
	assert.NoError(err)

	data, err := NewRandomData(1024)
	assert.NoError(err)

	req, _ := http.NewRequest("POST", "/", nil)
	assert.NoError(AddIngestionHeaders(e, req, user, data))
	req.Header.Del(common.FkDeviceIdHeaderName)
	rr := tests.ExecuteRequest(req, handler)

	assert.Equal(http.StatusBadRequest, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"message": "invalid headers",
		"headers": "<<PRESENCE>>"
	}`)
}

func TestIngestWithoutGeneration(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	handler, user, err := NewTestableIngester(e)
	assert.NoError(err)

	data, err := NewRandomData(1024)
	assert.NoError(err)

	req, _ := http.NewRequest("POST", "/", nil)
	assert.NoError(AddIngestionHeaders(e, req, user, data))
	req.Header.Del(common.FkGenerationHeaderName)
	rr := tests.ExecuteRequest(req, handler)

	assert.Equal(http.StatusBadRequest, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"message": "invalid headers",
		"headers": "<<PRESENCE>>"
	}`)
}

func TestIngestWithoutBlocks(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	handler, user, err := NewTestableIngester(e)
	assert.NoError(err)

	data, err := NewRandomData(1024)
	assert.NoError(err)

	req, _ := http.NewRequest("POST", "/", nil)
	assert.NoError(AddIngestionHeaders(e, req, user, data))
	req.Header.Del(common.FkBlocksHeaderName)
	rr := tests.ExecuteRequest(req, handler)

	assert.Equal(http.StatusBadRequest, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"message": "invalid headers",
		"headers": "<<PRESENCE>>"
	}`)
}

func TestIngestWithoutType(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	handler, user, err := NewTestableIngester(e)
	assert.NoError(err)

	data, err := NewRandomData(1024)
	assert.NoError(err)

	req, _ := http.NewRequest("POST", "/", nil)
	assert.NoError(AddIngestionHeaders(e, req, user, data))
	req.Header.Del(common.FkTypeHeaderName)
	rr := tests.ExecuteRequest(req, handler)

	assert.Equal(http.StatusBadRequest, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"message": "invalid headers",
		"headers": "<<PRESENCE>>"
	}`)
}

func TestIngestSuccessfully(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	handler, user, err := NewTestableIngester(e)
	assert.NoError(err)

	data, err := NewRandomData(1024)
	assert.NoError(err)

	req, _ := http.NewRequest("POST", "/", nil)
	assert.NoError(AddIngestionHeaders(e, req, user, data))
	rr := tests.ExecuteRequest(req, handler)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"id": "<<PRESENCE>>",
		"upload_id": "<<PRESENCE>>"
	}`)

	body := IngestionSuccessful{}
	assert.NoError(json.Unmarshal(rr.Body.Bytes(), &body))

	found := 0
	assert.NoError(e.DB.GetContext(e.Ctx, &found, "SELECT COUNT(*) FROM fieldkit.ingestion WHERE id = $1", body.ID))
	assert.Equal(found, 1)
}

func AddIngestionHeaders(e *tests.TestEnv, r *http.Request, user *data.User, data []byte) error {
	r.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(user))
	r.Header.Add("Content-Type", common.FkDataBinaryContentType)
	r.Header.Add("Content-Length", fmt.Sprintf("%d", len(data)))
	r.Header.Add("Fk-DeviceId", faker.UUIDDigit())
	r.Header.Add("Fk-Generation", faker.UUIDDigit())
	r.Header.Add("Fk-Type", "data")
	r.Header.Add("Fk-Blocks", "1,2")

	if data != nil {
		r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	}

	return nil
}

func NewRandomData(n int) ([]byte, error) {
	data := make([]byte, n)
	actual, err := rand.Read(data)
	if actual != n {
		return nil, fmt.Errorf("unexpected random byte read")
	}
	return data, err
}
