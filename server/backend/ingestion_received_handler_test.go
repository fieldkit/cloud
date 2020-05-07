package backend

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/backend/repositories"
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
		UserID:  int64(user.ID),
	})

	assert.Errorf(err, "ingestion missing: %d", 30342)
}

func TestIngestionReceivedCorruptedFile(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser("")
	assert.NoError(err)

	data, err := e.NewRandomData(1024)
	assert.NoError(err)

	handler := &IngestionReceivedHandler{
		Database: e.DB,
		Metrics:  logging.NewMetrics(e.Ctx, &logging.MetricsSettings{}),
		Files: tests.NewInMemoryArchive(map[string][]byte{
			"/file": []byte{},
		}),
	}

	ingestion, err := e.AddIngestion(user, "/file", e.MustDeviceID(), len(data))
	assert.NoError(err)

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		Time:    time.Now(),
		ID:      ingestion.ID,
		UserID:  int64(user.ID),
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

	ingestion, err := e.AddIngestion(user, "/meta", e.MustDeviceID(), len(files.Meta))
	assert.NoError(err)

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		Time:    time.Now(),
		ID:      ingestion.ID,
		UserID:  int64(user.ID),
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

	metaIngestion, err := e.AddIngestion(user, "/meta", deviceID, len(files.Meta))
	assert.NoError(err)

	dataIngestion, err := e.AddIngestion(user, "/data", deviceID, len(files.Data))
	assert.NoError(err)

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		Time:    time.Now(),
		ID:      metaIngestion.ID,
		UserID:  int64(user.ID),
		Verbose: true,
	}))

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		Time:    time.Now(),
		ID:      dataIngestion.ID,
		UserID:  int64(user.ID),
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

	metaIngestion, err := e.AddIngestion(user, "/meta", deviceID, len(files.Meta))
	assert.NoError(err)

	dataIngestion, err := e.AddIngestion(user, "/data", deviceID, len(files.Data))
	assert.NoError(err)

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		Time:    time.Now(),
		ID:      metaIngestion.ID,
		UserID:  int64(user.ID),
		Verbose: true,
	}))

	assert.NoError(handler.Handle(e.Ctx, &messages.IngestionReceived{
		Time:    time.Now(),
		ID:      dataIngestion.ID,
		UserID:  int64(user.ID),
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
}
