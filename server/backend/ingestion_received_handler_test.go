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
	handler := NewIngestionReceivedHandler(e.DB, tests.NewInMemoryArchive(map[string][]byte{}), logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}), publisher)

	err = handler.Handle(e.Ctx, &messages.IngestionReceived{
		QueuedID: int64(30342),
		Verbose:  true,
		UserID:   user.ID,
	})

	assert.Errorf(err, "queued ingestion missing: %d", 30342)
}

func TestIngestionReceivedCorruptedFile(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	randomData, err := e.NewRandomData(1024)
	assert.NoError(err)

	publisher := jobs.NewDevNullMessagePublisher()
	files := tests.NewInMemoryArchive(map[string][]byte{
		"/file": []byte{},
	})
	handler := NewIngestionReceivedHandler(e.DB, files, logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}), publisher)

	queued, _, err := e.AddIngestion(user, "/file", data.MetaTypeName, e.MustDeviceID(), len(randomData))
	assert.NoError(err)

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		QueuedID: queued.ID,
		UserID:   user.ID,
		Verbose:  true,
	}))
}
func TestIngestionReceivedMetaOnly(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	files, err := e.NewFilePair(1, 16)
	assert.NoError(err)

	publisher := jobs.NewDevNullMessagePublisher()
	memoryFiles := tests.NewInMemoryArchive(map[string][]byte{
		"/meta": files.Meta,
		"/data": files.Data,
	})
	handler := NewIngestionReceivedHandler(e.DB, memoryFiles, logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}), publisher)

	queued, _, err := e.AddIngestion(user, "/meta", data.MetaTypeName, e.MustDeviceID(), len(files.Meta))
	assert.NoError(err)

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		QueuedID: queued.ID,
		UserID:   user.ID,
		Verbose:  true,
	}))
}

func TestIngestionReceivedMetaAndData(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	deviceID := e.MustDeviceID()

	files, err := e.NewFilePair(1, 16)
	assert.NoError(err)

	publisher := jobs.NewDevNullMessagePublisher()
	memoryFiles := tests.NewInMemoryArchive(map[string][]byte{
		"/meta": files.Meta,
		"/data": files.Data,
	})
	handler := NewIngestionReceivedHandler(e.DB, memoryFiles, logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}), publisher)

	queuedMeta, _, err := e.AddIngestion(user, "/meta", data.MetaTypeName, deviceID, len(files.Meta))
	assert.NotNil(queuedMeta)
	assert.NoError(err)

	queuedData, _, err := e.AddIngestion(user, "/data", data.DataTypeName, deviceID, len(files.Data))
	assert.NotNil(queuedData)
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

	ir, err := repositories.NewIngestionRepository(e.DB)
	assert.NoError(err)

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

	deviceID := e.MustDeviceID()

	files, err := e.NewFilePair(4, 16)
	assert.NoError(err)

	publisher := jobs.NewDevNullMessagePublisher()
	memoryFiles := tests.NewInMemoryArchive(map[string][]byte{
		"/meta": files.Meta,
		"/data": files.Data,
	})
	handler := NewIngestionReceivedHandler(e.DB, memoryFiles, logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}), publisher)

	queuedMeta, _, err := e.AddIngestion(user, "/meta", data.MetaTypeName, deviceID, len(files.Meta))
	assert.NoError(err)

	queuedData, _, err := e.AddIngestion(user, "/data", data.DataTypeName, deviceID, len(files.Data))
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

	ir, err := repositories.NewIngestionRepository(e.DB)
	assert.NoError(err)

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
	assert.Equal(0, found)
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
	memoryFiles := tests.NewInMemoryArchive(map[string][]byte{
		"/meta": files.Meta,
		"/data": files.Data,
	})
	handler := NewIngestionReceivedHandler(e.DB, memoryFiles, logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}), publisher)

	queuedMeta, metaIngestion, err := e.AddIngestion(fd.Owner, "/meta", data.MetaTypeName, deviceID, len(files.Meta))
	assert.NoError(err)
	assert.NotNil(metaIngestion)

	queuedData, dataIngestion, err := e.AddIngestion(fd.Owner, "/data", data.DataTypeName, deviceID, len(files.Data))
	assert.NoError(err)
	assert.NotNil(dataIngestion)

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		QueuedID: queuedMeta.ID,
		UserID:   fd.Owner.ID,
		Verbose:  true,
	}))

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		QueuedID: queuedData.ID,
		UserID:   fd.Owner.ID,
		Verbose:  true,
	}))

	ir, err := repositories.NewIngestionRepository(e.DB)
	assert.NoError(err)

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
