package backend

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/common/logging"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/messages"
	"github.com/fieldkit/cloud/server/tests"
)

func TestIngestionReceivedNoSuchIngestion(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	publisher := jobs.NewDevNullMessagePublisher()
	mc := jobs.NewMessageContext(publisher, nil)
	handler := NewIngestionReceivedHandler(e.DB, e.DbPool, tests.NewInMemoryArchive(map[string][]byte{}), logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}), publisher, nil)

	err = handler.Start(e.Ctx, &messages.IngestionReceived{
		UserID:   user.ID,
		QueuedID: int64(30342),
		Verbose:  true,
	}, mc)

	assert.Errorf(err, "queued ingestion missing: %d", 30342)
}

func TestIngestionReceivedCorruptedFile(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	deviceID := fd.Stations[0].DeviceID

	randomData, err := e.NewRandomData(1024)
	assert.NoError(err)

	publisher := jobs.NewDevNullMessagePublisher()
	mc := jobs.NewMessageContext(publisher, nil)
	files := tests.NewInMemoryArchive(map[string][]byte{
		"/file": []byte{},
	})
	handler := NewIngestionReceivedHandler(e.DB, e.DbPool, files, logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}), publisher, nil)

	queued, _, err := e.AddIngestion(user, "/file", data.MetaTypeName, deviceID, len(randomData))
	assert.NoError(err)

	assert.NoError(handler.Start(e.Ctx, &messages.IngestionReceived{
		UserID:   user.ID,
		QueuedID: queued.ID,
		Verbose:  true,
	}, mc))
}
func TestIngestionReceivedMetaOnly(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	deviceID := fd.Stations[0].DeviceID

	files, err := e.NewFilePair(1, 16)
	assert.NoError(err)

	publisher := jobs.NewDevNullMessagePublisher()
	mc := jobs.NewMessageContext(publisher, nil)
	memoryFiles := tests.NewInMemoryArchive(map[string][]byte{
		"/meta": files.Meta,
		"/data": files.Data,
	})
	handler := NewIngestionReceivedHandler(e.DB, e.DbPool, memoryFiles, logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}), publisher, nil)

	queued, _, err := e.AddIngestion(user, "/meta", data.MetaTypeName, deviceID, len(files.Meta))
	assert.NoError(err)

	assert.NoError(handler.Start(e.Ctx, &messages.IngestionReceived{
		UserID:   user.ID,
		QueuedID: queued.ID,
		Verbose:  true,
	}, mc))
}

func TestIngestionReceivedMetaAndData(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	assert.NoError(e.AddRandomSensors())

	fd, err := e.AddStations(1)
	assert.NoError(err)

	deviceID := fd.Stations[0].DeviceID

	files, err := e.NewFilePair(1, 16)
	assert.NoError(err)

	publisher := jobs.NewDevNullMessagePublisher()
	mc := jobs.NewMessageContext(publisher, nil)
	memoryFiles := tests.NewInMemoryArchive(map[string][]byte{
		"/meta": files.Meta,
		"/data": files.Data,
	})
	handler := NewIngestionReceivedHandler(e.DB, e.DbPool, memoryFiles, logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}), publisher, nil)

	queuedMeta, _, err := e.AddIngestion(user, "/meta", data.MetaTypeName, deviceID, len(files.Meta))
	assert.NotNil(queuedMeta)
	assert.NoError(err)

	queuedData, _, err := e.AddIngestion(user, "/data", data.DataTypeName, deviceID, len(files.Data))
	assert.NotNil(queuedData)
	assert.NoError(err)

	assert.NoError(handler.Start(e.Ctx, &messages.IngestionReceived{
		UserID:   user.ID,
		QueuedID: queuedMeta.ID,
		Verbose:  true,
	}, mc))

	assert.NoError(handler.Start(e.Ctx, &messages.IngestionReceived{
		UserID:   user.ID,
		QueuedID: queuedData.ID,
		Verbose:  true,
	}, mc))

	ir := repositories.NewIngestionRepository(e.DB)

	miAfter, err := ir.QueryQueuedByID(e.Ctx, queuedMeta.ID)
	assert.NoError(err)
	assert.NotNil(miAfter)
	assert.Equal(int64(1), *miAfter.TotalRecords)
	assert.Equal(int64(0), *miAfter.DataErrors)
	assert.Equal(int64(0), *miAfter.MetaErrors)
	assert.NotNil(miAfter.Completed)

	diAfter, err := ir.QueryQueuedByID(e.Ctx, queuedData.ID)
	assert.NoError(err)
	assert.NotNil(miAfter)
	assert.Equal(int64(16), *diAfter.TotalRecords)
	assert.Equal(int64(0), *diAfter.DataErrors)
	assert.Equal(int64(0), *diAfter.MetaErrors)
	assert.NotNil(diAfter.Completed)
}

func TestIngestionReceivedMetaAndDataWithMultipleMeta(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	assert.NoError(e.AddRandomSensors())

	fd, err := e.AddStations(1)
	assert.NoError(err)

	deviceID := fd.Stations[0].DeviceID

	files, err := e.NewFilePair(4, 16)
	assert.NoError(err)

	publisher := jobs.NewDevNullMessagePublisher()
	mc := jobs.NewMessageContext(publisher, nil)
	memoryFiles := tests.NewInMemoryArchive(map[string][]byte{
		"/meta": files.Meta,
		"/data": files.Data,
	})
	handler := NewIngestionReceivedHandler(e.DB, e.DbPool, memoryFiles, logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}), publisher, nil)

	queuedMeta, _, err := e.AddIngestion(user, "/meta", data.MetaTypeName, deviceID, len(files.Meta))
	assert.NoError(err)

	queuedData, _, err := e.AddIngestion(user, "/data", data.DataTypeName, deviceID, len(files.Data))
	assert.NoError(err)

	assert.NoError(handler.Start(e.Ctx, &messages.IngestionReceived{
		UserID:   user.ID,
		QueuedID: queuedMeta.ID,
		Verbose:  true,
	}, mc))

	assert.NoError(handler.Start(e.Ctx, &messages.IngestionReceived{
		UserID:   user.ID,
		QueuedID: queuedData.ID,
		Verbose:  true,
	}, mc))

	ir := repositories.NewIngestionRepository(e.DB)

	miAfter, err := ir.QueryQueuedByID(e.Ctx, queuedMeta.ID)
	assert.NoError(err)
	assert.Equal(int64(4), *miAfter.TotalRecords)
	assert.Equal(int64(0), *miAfter.DataErrors)
	assert.Equal(int64(0), *miAfter.MetaErrors)
	assert.NotNil(miAfter.Completed)

	diAfter, err := ir.QueryQueuedByID(e.Ctx, queuedData.ID)
	assert.NoError(err)
	assert.Equal(int64(4*16), *diAfter.TotalRecords)
	assert.Equal(int64(0), *diAfter.DataErrors)
	assert.Equal(int64(0), *diAfter.MetaErrors)
	assert.NotNil(diAfter.Completed)

	found := 0
	assert.NoError(e.DB.GetContext(e.Ctx, &found, "SELECT COUNT(*) FROM fieldkit.station_ingestion WHERE data_ingestion_id = $1", diAfter.ID))
	assert.Equal(1, found)
}

func TestIngestionReceivedMetaAndDataWithMultipleMetaAndStationAlreadyAdded(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	deviceID := fd.Stations[0].DeviceID

	files, err := e.NewFilePair(4, 16)
	assert.NoError(err)

	publisher := jobs.NewDevNullMessagePublisher()
	mc := jobs.NewMessageContext(publisher, nil)
	memoryFiles := tests.NewInMemoryArchive(map[string][]byte{
		"/meta": files.Meta,
		"/data": files.Data,
	})
	handler := NewIngestionReceivedHandler(e.DB, e.DbPool, memoryFiles, logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}), publisher, nil)

	queuedMeta, metaIngestion, err := e.AddIngestion(fd.Owner, "/meta", data.MetaTypeName, deviceID, len(files.Meta))
	assert.NoError(err)
	assert.NotNil(metaIngestion)

	queuedData, dataIngestion, err := e.AddIngestion(fd.Owner, "/data", data.DataTypeName, deviceID, len(files.Data))
	assert.NoError(err)
	assert.NotNil(dataIngestion)

	assert.NoError(handler.Start(e.Ctx, &messages.IngestionReceived{
		UserID:   fd.Owner.ID,
		QueuedID: queuedMeta.ID,
		Verbose:  true,
	}, mc))

	assert.NoError(handler.Start(e.Ctx, &messages.IngestionReceived{
		UserID:   fd.Owner.ID,
		QueuedID: queuedData.ID,
		Verbose:  true,
	}, mc))

	ir := repositories.NewIngestionRepository(e.DB)

	miAfter, err := ir.QueryQueuedByID(e.Ctx, queuedMeta.ID)
	assert.NoError(err)
	assert.Equal(int64(4), *miAfter.TotalRecords)
	assert.Equal(int64(0), *miAfter.DataErrors)
	assert.Equal(int64(0), *miAfter.MetaErrors)
	assert.Equal(int64(0), *miAfter.OtherErrors)
	assert.NotNil(miAfter.Completed)

	diAfter, err := ir.QueryQueuedByID(e.Ctx, queuedData.ID)
	assert.NoError(err)
	assert.Equal(int64(4*16), *diAfter.TotalRecords)
	assert.Equal(int64(0), *diAfter.DataErrors)
	assert.Equal(int64(0), *diAfter.MetaErrors)
	assert.Equal(int64(0), *diAfter.OtherErrors)
	assert.NotNil(diAfter.Completed)

	found := 0
	assert.NoError(e.DB.GetContext(e.Ctx, &found, "SELECT COUNT(*) FROM fieldkit.station_ingestion WHERE data_ingestion_id = $1", dataIngestion.ID))
	assert.Equal(1, found)
}
