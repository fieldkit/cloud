package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/common/jobs"
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

	assert.NoError(e.AddRandomSensors())

	fd, err := e.AddStations(1)
	assert.NoError(err)

	user := fd.Owner
	station := fd.Stations[0]

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload, err := json.Marshal(
		struct {
			DeviceID string `json:"deviceId"`
			Name     string `json:"name"`
		}{
			DeviceID: station.DeviceIDHex(),
			Name:     station.Name,
		},
	)
	assert.NoError(err)

	req, _ := http.NewRequest("POST", "/user/stations", bytes.NewReader(payload))
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(user))
	rr := tests.ExecuteRequest(req, api)

	files, err := e.NewFilePair(4, 16)
	assert.NoError(err)

	publisher := jobs.NewDevNullMessagePublisher()
	mc := jobs.NewMessageContext(publisher, nil)
	memoryFiles := tests.NewInMemoryArchive(map[string][]byte{
		"/meta": files.Meta,
		"/data": files.Data,
	})
	handler := backend.NewIngestionReceivedHandler(e.DB, e.DbPool, memoryFiles, logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}), publisher, nil)

	queuedMeta, _, err := e.AddIngestion(user, "/meta", data.MetaTypeName, station.DeviceID, len(files.Meta))
	assert.NoError(err)

	queuedData, _, err := e.AddIngestion(user, "/data", data.DataTypeName, station.DeviceID, len(files.Data))
	assert.NoError(err)

	log := Logger(e.Ctx).Sugar()

	log.Infow("first test")

	assert.NoError(handler.Start(e.Ctx, &messages.IngestionReceived{
		QueuedID: queuedMeta.ID,
		UserID:   user.ID,
		Verbose:  true,
	}, mc))

	log.Infow("second test")

	assert.NoError(handler.Start(e.Ctx, &messages.IngestionReceived{
		QueuedID: queuedData.ID,
		UserID:   user.ID,
		Verbose:  true,
	}, mc))

	log.Infow("end ingestion")

	req, _ = http.NewRequest("GET", "/user/stations", nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(user))
	rr = tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"stations": [
			{
				"id": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"lastReadingAt": "<<PRESENCE>>",
				"model": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"interestingness": "<<PRESENCE>>",
				"attributes": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": null,
				"readOnly": "<<PRESENCE>>",
				"location": "<<PRESENCE>>",
				"configurations": {
					"all": [
						{
							"id": "<<PRESENCE>>",
							"provisionId": "<<PRESENCE>>",
							"time": "<<PRESENCE>>",
							"modules": [
								{
									"id": "<<PRESENCE>>",
									"flags": "<<PRESENCE>>",
									"fullKey": "<<PRESENCE>>",
									"hardwareId": "<<PRESENCE>>",
									"hardwareIdBase64": "<<PRESENCE>>",
									"internal": "<<PRESENCE>>",
									"name": "<<PRESENCE>>",
									"position": "<<PRESENCE>>",
									"sensors": [
										{ "key": "random0", "fullKey": "fk.random-module-1.random0", "name": "random0", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random1", "fullKey": "fk.random-module-1.random1", "name": "random1", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random2", "fullKey": "fk.random-module-1.random2", "name": "random2", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random3", "fullKey": "fk.random-module-1.random3", "name": "random3", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random4", "fullKey": "fk.random-module-1.random4", "name": "random4", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" }
									]
								},
								{
									"id": "<<PRESENCE>>",
									"flags": "<<PRESENCE>>",
									"fullKey": "<<PRESENCE>>",
									"hardwareId": "<<PRESENCE>>",
									"hardwareIdBase64": "<<PRESENCE>>",
									"internal": "<<PRESENCE>>",
									"name": "<<PRESENCE>>",
									"position": "<<PRESENCE>>",
									"sensors": [
										{ "key": "random0", "fullKey": "fk.random-module-2.random0", "name": "random0", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random1", "fullKey": "fk.random-module-2.random1", "name": "random1", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random2", "fullKey": "fk.random-module-2.random2", "name": "random2", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random3", "fullKey": "fk.random-module-2.random3", "name": "random3", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random4", "fullKey": "fk.random-module-2.random4", "name": "random4", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random5", "fullKey": "fk.random-module-2.random5", "name": "random5", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random6", "fullKey": "fk.random-module-2.random6", "name": "random6", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random7", "fullKey": "fk.random-module-2.random7", "name": "random7", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random8", "fullKey": "fk.random-module-2.random8", "name": "random8", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random9", "fullKey": "fk.random-module-2.random9", "name": "random9", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" }
									]
								}
							]
						}
					]
				}
			}
		]
	}`)
}
