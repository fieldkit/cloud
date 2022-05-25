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
	memoryFiles := tests.NewInMemoryArchive(map[string][]byte{
		"/meta": files.Meta,
		"/data": files.Data,
	})
	handler := backend.NewIngestionReceivedHandler(e.DB, memoryFiles, logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}), publisher)

	queuedMeta, _, err := e.AddIngestion(user, "/meta", data.MetaTypeName, station.DeviceID, len(files.Meta))
	assert.NoError(err)

	queuedData, _, err := e.AddIngestion(user, "/data", data.DataTypeName, station.DeviceID, len(files.Data))
	assert.NoError(err)

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		QueuedID: queuedMeta.ID,
		UserID:   user.ID,
		Verbose:  true,
		Refresh:  true,
	}))

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		QueuedID: queuedData.ID,
		UserID:   user.ID,
		Verbose:  true,
		Refresh:  true,
	}))

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
									"internal": "<<PRESENCE>>",
									"name": "<<PRESENCE>>",
									"position": "<<PRESENCE>>",
									"sensors": [
										{ "key": "random0", "fullKey": "random-module-1.random0", "name": "random_0", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random1", "fullKey": "random-module-1.random1", "name": "random_1", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random2", "fullKey": "random-module-1.random2", "name": "random_2", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random3", "fullKey": "random-module-1.random3", "name": "random_3", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random4", "fullKey": "random-module-1.random4", "name": "random_4", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" }
									]
								},
								{
									"id": "<<PRESENCE>>",
									"flags": "<<PRESENCE>>",
									"fullKey": "<<PRESENCE>>",
									"hardwareId": "<<PRESENCE>>",
									"internal": "<<PRESENCE>>",
									"name": "<<PRESENCE>>",
									"position": "<<PRESENCE>>",
									"sensors": [
										{ "key": "random0", "fullKey": "random-module-2.random0", "name": "random_0", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random1", "fullKey": "random-module-2.random1", "name": "random_1", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random2", "fullKey": "random-module-2.random2", "name": "random_2", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random3", "fullKey": "random-module-2.random3", "name": "random_3", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random4", "fullKey": "random-module-2.random4", "name": "random_4", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random5", "fullKey": "random-module-2.random5", "name": "random_5", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random6", "fullKey": "random-module-2.random6", "name": "random_6", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random7", "fullKey": "random-module-2.random7", "name": "random_7", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random8", "fullKey": "random-module-2.random8", "name": "random_8", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" },
										{ "key": "random9", "fullKey": "random-module-2.random9", "name": "random_9", "ranges": null, "unitOfMeasure": "C", "reading": "<<PRESENCE>>" }
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
