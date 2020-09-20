package backend

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/tests"
)

func TestQueryStationByID(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	r := repositories.NewStationRepository(e.DB)

	sf, err := r.QueryStationFull(e.Ctx, fd.Stations[0].ID)
	assert.NoError(err)

	assert.NotNil(sf)
	assert.NotNil(fd)
}

func TestQueryStationsByOwnerID(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	_, err = e.AddStations(5)
	assert.NoError(err)

	fd, err := e.AddStations(5)
	assert.NoError(err)

	r := repositories.NewStationRepository(e.DB)

	sfs, err := r.QueryStationFullByOwnerID(e.Ctx, fd.Owner.ID)
	assert.NoError(err)

	assert.NotNil(sfs)
	assert.Equal(len(sfs), len(fd.Stations))
}

func TestQueryStationsByProjectID(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	_, err = e.AddStations(5)
	assert.NoError(err)

	fd, err := e.AddStations(5)
	assert.NoError(err)

	r := repositories.NewStationRepository(e.DB)

	sfs, err := r.QueryStationFullByProjectID(e.Ctx, fd.Project.ID)
	assert.NoError(err)

	assert.NotNil(sfs)
	assert.Equal(len(sfs), len(fd.Stations))
}

func TestQueryStationByDeviceID(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	r := repositories.NewStationRepository(e.DB)

	sf, err := r.QueryStationByDeviceID(e.Ctx, fd.Stations[0].DeviceID)
	assert.NoError(err)

	assert.NotNil(sf)
	assert.NotNil(fd)
}

func TestQueryStationsByDeviceID(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	r := repositories.NewStationRepository(e.DB)

	sf, err := r.QueryStationsByDeviceID(e.Ctx, fd.Stations[0].DeviceID)
	assert.NoError(err)

	assert.NotNil(sf)
	assert.NotNil(fd)
}

func TestTryQueryStationByDeviceID(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	r := repositories.NewStationRepository(e.DB)

	sf, err := r.TryQueryStationByDeviceID(e.Ctx, fd.Stations[0].DeviceID)
	assert.NoError(err)

	assert.NotNil(sf)
	assert.NotNil(fd)
}

func TestIngestionRepositoryQueryPending(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	r, err := repositories.NewIngestionRepository(e.DB)
	assert.NoError(err)

	sf, err := r.QueryPending(e.Ctx)
	assert.NoError(err)

	assert.NotNil(sf)
	assert.NotNil(fd)
}

func TestIngestionRepositoryQueryByStationID(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	r, err := repositories.NewIngestionRepository(e.DB)
	assert.NoError(err)

	sf, err := r.QueryByStationID(e.Ctx, int64(fd.Stations[0].ID))
	assert.NoError(err)

	assert.NotNil(sf)
	assert.NotNil(fd)
}
