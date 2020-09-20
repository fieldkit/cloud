package backend

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/tests"
)

func MustParseTime(value string) *time.Time {
	t, err := time.Parse("2006/01/02 15:04:05", value)
	if err != nil {
		panic(err)
	}
	return &t
}

func TestGetStationVisibilityNoProjectsNoCollections(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user1, err := e.AddUser()
	assert.NoError(err)
	s1, err := e.AddStation(user1)
	assert.NoError(err)

	vs := NewVisibilitySlicer(e.DB)

	dvs, err := vs.Slice(e.Ctx, s1.ID)
	assert.NoError(err)
	assert.NotNil(dvs)

	// TODO Should be visible only to the owner.
}

func TestGetStationVisibilityPrivateProject(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user1, err := e.AddUser()
	assert.NoError(err)
	s1, err := e.AddStation(user1)
	assert.NoError(err)
	p1, err := e.SaveProject(&data.Project{
		Name:      "Private Project",
		StartTime: MustParseTime("2008/01/01 00:00:00"),
		EndTime:   MustParseTime("2009/01/01 00:00:00"),
		Privacy:   data.Private,
	})
	assert.NoError(err)
	assert.NoError(e.AddStationToProject(s1, p1))

	vs := NewVisibilitySlicer(e.DB)

	dvs, err := vs.Slice(e.Ctx, s1.ID)
	assert.NoError(err)
	assert.NotNil(dvs)

	// TODO Should be visible only to project members.
}

func TestGetStationVisibilityPublicProject(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user1, err := e.AddUser()
	assert.NoError(err)
	s1, err := e.AddStation(user1)
	assert.NoError(err)
	p1, err := e.SaveProject(&data.Project{
		Name:      "Public Project",
		StartTime: MustParseTime("2008/01/01 00:00:00"),
		EndTime:   MustParseTime("2009/01/01 00:00:00"),
		Privacy:   data.Public,
	})
	assert.NoError(err)
	assert.NoError(e.AddStationToProject(s1, p1))

	vs := NewVisibilitySlicer(e.DB)

	dvs, err := vs.Slice(e.Ctx, s1.ID)
	assert.NoError(err)
	assert.NotNil(dvs)

	// TODO Should be visible to anybody.
}
