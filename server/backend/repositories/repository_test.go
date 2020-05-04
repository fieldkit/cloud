package repositories

import (
	"testing"

	"github.com/fieldkit/cloud/server/tests"

	"github.com/stretchr/testify/assert"
)

func TestQueryStationByID(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	r, err := NewStationRepository(e.DB)
	assert.NoError(err)

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

	r, err := NewStationRepository(e.DB)
	assert.NoError(err)

	sfs, err := r.QueryStationFullByOwnerID(e.Ctx, fd.OwnerID)
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

	r, err := NewStationRepository(e.DB)
	assert.NoError(err)

	sfs, err := r.QueryStationFullByProjectID(e.Ctx, fd.ProjectID)
	assert.NoError(err)

	assert.NotNil(sfs)
	assert.Equal(len(sfs), len(fd.Stations))
}
