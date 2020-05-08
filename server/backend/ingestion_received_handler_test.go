package backend

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/logging"
	"github.com/fieldkit/cloud/server/messages"
	"github.com/fieldkit/cloud/server/tests"
)

func TestIngestionReceivedNoSuchIngestion(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser("")
	assert.NoError(err)

	handler := &IngestionReceivedHandler{
		Database: e.DB,
		Metrics:  logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}),
		Files:    tests.NewInMemoryArchive(map[string][]byte{}),
	}

	err = handler.Handle(e.Ctx, &messages.IngestionReceived{
		Time:    time.Now(),
		ID:      int64(30342),
		Verbose: true,
		UserID:  user.ID,
	})

	assert.Errorf(err, "ingestion missing: %d", 30342)
}

func TestIngestionReceivedCorruptedFile(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser("")
	assert.NoError(err)

	randomData, err := e.NewRandomData(1024)
	assert.NoError(err)

	handler := &IngestionReceivedHandler{
		Database: e.DB,
		Metrics:  logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}),
		Files: tests.NewInMemoryArchive(map[string][]byte{
			"/file": []byte{},
		}),
	}

	ingestion, err := e.AddIngestion(user, "/file", data.MetaTypeName, e.MustDeviceID(), len(randomData))
	assert.NoError(err)

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		Time:    time.Now(),
		ID:      ingestion.ID,
		UserID:  user.ID,
		Verbose: true,
	}))
}
func TestIngestionReceivedMetaOnly(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser("")
	assert.NoError(err)

	files, err := e.NewFilePair(1, 16)
	assert.NoError(err)

	handler := &IngestionReceivedHandler{
		Database: e.DB,
		Metrics:  logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}),
		Files: tests.NewInMemoryArchive(map[string][]byte{
			"/meta": files.Meta,
			"/data": files.Data,
		}),
	}

	ingestion, err := e.AddIngestion(user, "/meta", data.MetaTypeName, e.MustDeviceID(), len(files.Meta))
	assert.NoError(err)

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		Time:    time.Now(),
		ID:      ingestion.ID,
		UserID:  user.ID,
		Verbose: true,
	}))
}

func TestIngestionReceivedMetaAndData(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser("")
	assert.NoError(err)

	deviceID := e.MustDeviceID()

	files, err := e.NewFilePair(1, 16)
	assert.NoError(err)

	handler := &IngestionReceivedHandler{
		Database: e.DB,
		Metrics:  logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}),
		Files: tests.NewInMemoryArchive(map[string][]byte{
			"/meta": files.Meta,
			"/data": files.Data,
		}),
	}

	metaIngestion, err := e.AddIngestion(user, "/meta", data.MetaTypeName, deviceID, len(files.Meta))
	assert.NoError(err)

	dataIngestion, err := e.AddIngestion(user, "/data", data.DataTypeName, deviceID, len(files.Data))
	assert.NoError(err)

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		Time:    time.Now(),
		ID:      metaIngestion.ID,
		UserID:  user.ID,
		Verbose: true,
	}))

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		Time:    time.Now(),
		ID:      dataIngestion.ID,
		UserID:  user.ID,
		Verbose: true,
	}))

	ir, err := repositories.NewIngestionRepository(e.DB)
	assert.NoError(err)

	miAfter, err := ir.QueryByID(e.Ctx, metaIngestion.ID)
	assert.NoError(err)
	assert.Equal(int64(1), *miAfter.TotalRecords)
	assert.Equal(int64(0), *miAfter.DataErrors)
	assert.Equal(int64(0), *miAfter.MetaErrors)
	assert.NotNil(miAfter.Completed)

	diAfter, err := ir.QueryByID(e.Ctx, dataIngestion.ID)
	assert.NoError(err)
	assert.Equal(int64(16), *diAfter.TotalRecords)
	assert.Equal(int64(0), *diAfter.DataErrors)
	assert.Equal(int64(0), *diAfter.MetaErrors)
	assert.NotNil(diAfter.Completed)
}

func TestIngestionReceivedMetaAndDataWithMultipleMeta(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser("")
	assert.NoError(err)

	deviceID := e.MustDeviceID()

	files, err := e.NewFilePair(4, 16)
	assert.NoError(err)

	handler := &IngestionReceivedHandler{
		Database: e.DB,
		Metrics:  logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}),
		Files: tests.NewInMemoryArchive(map[string][]byte{
			"/meta": files.Meta,
			"/data": files.Data,
		}),
	}

	metaIngestion, err := e.AddIngestion(user, "/meta", data.MetaTypeName, deviceID, len(files.Meta))
	assert.NoError(err)

	dataIngestion, err := e.AddIngestion(user, "/data", data.DataTypeName, deviceID, len(files.Data))
	assert.NoError(err)

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		Time:    time.Now(),
		ID:      metaIngestion.ID,
		UserID:  user.ID,
		Verbose: true,
	}))

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		Time:    time.Now(),
		ID:      dataIngestion.ID,
		UserID:  user.ID,
		Verbose: true,
	}))

	ir, err := repositories.NewIngestionRepository(e.DB)
	assert.NoError(err)

	miAfter, err := ir.QueryByID(e.Ctx, metaIngestion.ID)
	assert.NoError(err)
	assert.Equal(int64(4), *miAfter.TotalRecords)
	assert.Equal(int64(0), *miAfter.DataErrors)
	assert.Equal(int64(0), *miAfter.MetaErrors)
	assert.NotNil(miAfter.Completed)

	diAfter, err := ir.QueryByID(e.Ctx, dataIngestion.ID)
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

	handler := &IngestionReceivedHandler{
		Database: e.DB,
		Metrics:  logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}),
		Files: tests.NewInMemoryArchive(map[string][]byte{
			"/meta": files.Meta,
			"/data": files.Data,
		}),
	}

	metaIngestion, err := e.AddIngestion(fd.Owner, "/meta", data.MetaTypeName, deviceID, len(files.Meta))
	assert.NoError(err)

	dataIngestion, err := e.AddIngestion(fd.Owner, "/data", data.DataTypeName, deviceID, len(files.Data))
	assert.NoError(err)

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		Time:    time.Now(),
		ID:      metaIngestion.ID,
		UserID:  fd.Owner.ID,
		Verbose: true,
	}))

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		Time:    time.Now(),
		ID:      dataIngestion.ID,
		UserID:  fd.Owner.ID,
		Verbose: true,
	}))

	ir, err := repositories.NewIngestionRepository(e.DB)
	assert.NoError(err)

	miAfter, err := ir.QueryByID(e.Ctx, metaIngestion.ID)
	assert.NoError(err)
	assert.Equal(int64(4), *miAfter.TotalRecords)
	assert.Equal(int64(0), *miAfter.DataErrors)
	assert.Equal(int64(0), *miAfter.MetaErrors)
	assert.Equal(int64(0), *miAfter.OtherErrors)
	assert.NotNil(miAfter.Completed)

	diAfter, err := ir.QueryByID(e.Ctx, dataIngestion.ID)
	assert.NoError(err)
	assert.Equal(int64(4*16), *diAfter.TotalRecords)
	assert.Equal(int64(0), *diAfter.DataErrors)
	assert.Equal(int64(0), *diAfter.MetaErrors)
	assert.Equal(int64(0), *diAfter.OtherErrors)
	assert.NotNil(diAfter.Completed)

	found := 0
	assert.NoError(e.DB.GetContext(e.Ctx, &found, "SELECT COUNT(*) FROM fieldkit.station_ingestion WHERE data_ingestion_id = $1", diAfter.ID))
	assert.Equal(1, found)
}
