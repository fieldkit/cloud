package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/common/logging"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/messages"
	"github.com/fieldkit/cloud/server/tests"
)

func TestQueryStationWithConfigurations(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	user := fd.Owner
	station := fd.Stations[0]

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload, err := json.Marshal(
		struct {
			DeviceID   string                 `json:"device_id"`
			Name       string                 `json:"name"`
			StatusJSON map[string]interface{} `json:"status_json"`
		}{
			DeviceID:   station.DeviceIDHex(),
			Name:       station.Name,
			StatusJSON: make(map[string]interface{}),
		},
	)
	assert.NoError(err)

	req, _ := http.NewRequest("POST", "/stations", bytes.NewReader(payload))
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(user))
	rr := tests.ExecuteRequest(req, api)

	files, err := e.NewFilePair(4, 16)
	assert.NoError(err)

	memoryFiles := tests.NewInMemoryArchive(map[string][]byte{
		"/meta": files.Meta,
		"/data": files.Data,
	})
	handler := backend.NewIngestionReceivedHandler(e.DB, memoryFiles, logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}))

	queuedMeta, _, err := e.AddIngestion(user, "/meta", data.MetaTypeName, station.DeviceID, len(files.Meta))
	assert.NoError(err)

	queuedData, _, err := e.AddIngestion(user, "/data", data.DataTypeName, station.DeviceID, len(files.Data))
	assert.NoError(err)

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		QueuedID: queuedMeta.ID,
		UserID:   user.ID,
		Verbose:  true,
	}))

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		QueuedID: queuedData.ID,
		UserID:   user.ID,
		Verbose:  true,
	}))

	req, _ = http.NewRequest("GET", "/stations", nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(user))
	rr = tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"stations": [
			{
				"id": "<<PRESENCE>>",
				"images": "<<PRESENCE>>",
				"updated": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"device_id": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": "<<PRESENCE>>",
				"read_only": "<<PRESENCE>>",
				"configurations": {
					"all": [
						{
							"id": "<<PRESENCE>>",
							"provision_id": "<<PRESENCE>>",
							"time": "<<PRESENCE>>",
							"modules": [
								{
									"id": "<<PRESENCE>>",
									"flags": "<<PRESENCE>>",
									"hardware_id": "<<PRESENCE>>",
									"internal": "<<PRESENCE>>",
									"name": "<<PRESENCE>>",
									"position": "<<PRESENCE>>",
									"sensors": [
										{ "key": "", "name": "random_0", "ranges": null, "unit_of_measure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "", "name": "random_1", "ranges": null, "unit_of_measure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "", "name": "random_2", "ranges": null, "unit_of_measure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "", "name": "random_3", "ranges": null, "unit_of_measure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "", "name": "random_4", "ranges": null, "unit_of_measure": "C", "reading": "<<PRESENCE>>" }
									]
								},
								{
									"id": "<<PRESENCE>>",
									"flags": "<<PRESENCE>>",
									"hardware_id": "<<PRESENCE>>",
									"internal": "<<PRESENCE>>",
									"name": "<<PRESENCE>>",
									"position": "<<PRESENCE>>",
									"sensors": [
										{ "key": "", "name": "random_0", "ranges": null, "unit_of_measure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "", "name": "random_1", "ranges": null, "unit_of_measure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "", "name": "random_2", "ranges": null, "unit_of_measure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "", "name": "random_3", "ranges": null, "unit_of_measure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "", "name": "random_4", "ranges": null, "unit_of_measure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "", "name": "random_5", "ranges": null, "unit_of_measure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "", "name": "random_6", "ranges": null, "unit_of_measure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "", "name": "random_7", "ranges": null, "unit_of_measure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "", "name": "random_8", "ranges": null, "unit_of_measure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "", "name": "random_9", "ranges": null, "unit_of_measure": "C", "reading": "<<PRESENCE>>" }
									]
								}
							]
						},
						{
							"id": "<<PRESENCE>>",
							"provision_id": "<<PRESENCE>>",
							"time": "<<PRESENCE>>",
							"modules": [ "<<PRESENCE>>", "<<PRESENCE>>" ]
						},
						{
							"id": "<<PRESENCE>>",
							"provision_id": "<<PRESENCE>>",
							"time": "<<PRESENCE>>",
							"modules": [ "<<PRESENCE>>", "<<PRESENCE>>" ]
						},
						{
							"id": "<<PRESENCE>>",
							"provision_id": "<<PRESENCE>>",
							"time": "<<PRESENCE>>",
							"modules": [ "<<PRESENCE>>", "<<PRESENCE>>" ]
						}
					]
				}
			}
		]
	}`)
}
