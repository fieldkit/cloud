package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/tests"
)

func TestGetMissingDataRecord(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/records/data/%d", 0), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForAdmin())
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusNotFound, rr.Code)
}

func TestGetMetaRecord(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	ar, err := e.AddMetaAndData(fd.Stations[0], fd.Owner, 5)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/records/meta/%d", ar.Meta.Records[0].ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForAdmin())
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"meta": {
			"id": "<<PRESENCE>>",
			"number": "<<PRESENCE>>",
			"provision_id": "<<PRESENCE>>",
			"raw": "<<PRESENCE>>",
			"pb": "<<PRESENCE>>",
			"time": "<<PRESENCE>>"
		}
	}`)
}

func TestGetDataRecord(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	ar, err := e.AddMetaAndData(fd.Stations[0], fd.Owner, 5)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/records/data/%d", ar.Data.Records[0].ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForAdmin())
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"meta": {
			"id": "<<PRESENCE>>",
			"number": "<<PRESENCE>>",
			"provision_id": "<<PRESENCE>>",
			"raw": "<<PRESENCE>>",
			"time": "<<PRESENCE>>",
			"pb": "<<PRESENCE>>"
		},
		"data": {
			"time": "<<PRESENCE>>",
			"id": "<<PRESENCE>>",
			"location": "<<PRESENCE>>",
			"meta_record_id": "<<PRESENCE>>",
			"number": "<<PRESENCE>>",
			"provision_id": "<<PRESENCE>>",
			"raw": "<<PRESENCE>>",
			"pb": "<<PRESENCE>>"
		}
	}`)
}

func TestGetDataRecordResolved(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	ar, err := e.AddMetaAndData(fd.Stations[0], fd.Owner, 5)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/records/data/%d/resolved", ar.Data.Records[0].ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForAdmin())
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"filtered": {
			"filters": "<<PRESENCE>>",
			"record": "<<PRESENCE>>"
		},
		"meta": {
			"id": "<<PRESENCE>>",
			"number": "<<PRESENCE>>",
			"time": "<<PRESENCE>>",
			"provision_id": "<<PRESENCE>>",
			"raw": "<<PRESENCE>>",
			"pb": "<<PRESENCE>>"
		},
		"data": {
			"time": "<<PRESENCE>>",
			"id": "<<PRESENCE>>",
			"location": "<<PRESENCE>>",
			"meta_record_id": "<<PRESENCE>>",
			"number": "<<PRESENCE>>",
			"provision_id": "<<PRESENCE>>",
			"raw": "<<PRESENCE>>",
			"pb": "<<PRESENCE>>"
		}
	}`)
}
