package backend

import (
	"testing"
	"time"

	"github.com/go-test/deep"

	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/backend/repositories"
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

func TestGetVisibilityNoProjectsNoCollections(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user1, err := e.AddUser()
	assert.NoError(err)
	s1, err := e.AddStationOwnedBy(user1)
	assert.NoError(err)

	vs := NewVisibilitySlicer(e.DB)

	dvs, err := vs.Slice(e.Ctx, s1.ID)
	assert.NoError(err)
	assert.NotNil(dvs)
	assert.NotEmpty(dvs)

	expected := []*data.DataVisibility{
		{
			StationID: s1.ID,
			StartTime: data.MinimumTime,
			EndTime:   data.MaximumTime,
			UserID:    &user1.ID,
		},
	}

	if diff := deep.Equal(dvs, expected); diff != nil {
		t.Error(diff)
	}
}

func TestGetVisibilityPrivateProject(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user1, err := e.AddUser()
	assert.NoError(err)
	s1, err := e.AddStationOwnedBy(user1)
	assert.NoError(err)
	p1, err := e.SaveProject(&data.Project{
		Name:      "Private Project",
		StartTime: MustParseTime("2008/01/01 00:00:00"),
		EndTime:   MustParseTime("2009/01/01 00:00:00"),
		Privacy:   data.Private,
	})
	assert.NoError(err)
	assert.NoError(e.AddStationToProject(s1, p1, nil))

	vs := NewVisibilitySlicer(e.DB)

	dvs, err := vs.Slice(e.Ctx, s1.ID)
	assert.NoError(err)
	assert.NotNil(dvs)
	assert.NotEmpty(dvs)

	expected := []*data.DataVisibility{
		{
			StationID: s1.ID,
			StartTime: data.MinimumTime,
			EndTime:   data.MaximumTime,
			UserID:    &user1.ID,
		},
		{
			StationID: s1.ID,
			StartTime: *p1.StartTime,
			EndTime:   *p1.EndTime,
			ProjectID: &p1.ID,
		},
	}

	if diff := deep.Equal(dvs, expected); diff != nil {
		t.Error(diff)
	}
}

func TestGetVisibilityPublicProject(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user1, err := e.AddUser()
	assert.NoError(err)
	s1, err := e.AddStationOwnedBy(user1)
	assert.NoError(err)
	p1, err := e.SaveProject(&data.Project{
		Name:      "Public Project",
		StartTime: MustParseTime("2008/01/01 00:00:00"),
		EndTime:   MustParseTime("2009/01/01 00:00:00"),
		Privacy:   data.Public,
	})
	assert.NoError(err)
	assert.NoError(e.AddStationToProject(s1, p1, nil))

	vs := NewVisibilitySlicer(e.DB)

	dvs, err := vs.Slice(e.Ctx, s1.ID)
	assert.NoError(err)
	assert.NotNil(dvs)
	assert.NotEmpty(dvs)

	expected := []*data.DataVisibility{
		{
			StationID: s1.ID,
			StartTime: data.MinimumTime,
			EndTime:   data.MaximumTime,
			UserID:    &user1.ID,
		},
		{
			StationID: s1.ID,
			StartTime: *p1.StartTime,
			EndTime:   *p1.EndTime,
		},
	}

	if diff := deep.Equal(dvs, expected); diff != nil {
		t.Error(diff)
	}
}

func TestGetVisibilityStationPassesThroughTwoProjects(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user1, err := e.AddUser()
	assert.NoError(err)
	s1, err := e.AddStationOwnedBy(user1)
	assert.NoError(err)
	p1, err := e.SaveProject(&data.Project{
		Name:      "Public Project",
		StartTime: MustParseTime("2008/01/01 00:00:00"),
		EndTime:   MustParseTime("2009/01/01 00:00:00"),
		Privacy:   data.Public,
	})
	assert.NoError(err)
	assert.NoError(e.AddStationToProject(s1, p1, nil))

	p2, err := e.SaveProject(&data.Project{
		Name:      "Private Project",
		StartTime: MustParseTime("2008/02/01 00:00:00"),
		EndTime:   MustParseTime("2008/04/01 00:00:00"),
		Privacy:   data.Private,
	})
	assert.NoError(err)
	assert.NoError(e.AddStationToProject(s1, p2, nil))

	vs := NewVisibilitySlicer(e.DB)

	dvs, err := vs.Slice(e.Ctx, s1.ID)
	assert.NoError(err)
	assert.NotNil(dvs)
	assert.NotEmpty(dvs)

	expected := []*data.DataVisibility{
		{
			StationID: s1.ID,
			StartTime: data.MinimumTime,
			EndTime:   data.MaximumTime,
			UserID:    &user1.ID,
		},
		{
			StationID: s1.ID,
			StartTime: *p1.StartTime,
			EndTime:   *p2.StartTime,
		},
		{
			StationID: s1.ID,
			StartTime: *p2.StartTime,
			EndTime:   *p2.EndTime,
			ProjectID: &p2.ID,
		},
		{
			StationID: s1.ID,
			StartTime: *p2.EndTime,
			EndTime:   *p1.EndTime,
		},
	}

	if diff := deep.Equal(dvs, expected); diff != nil {
		t.Error(diff)
	}
}

func TestVisibilityMergeBrandNew(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user1, err := e.AddUser()
	assert.NoError(err)
	s1, err := e.AddStationOwnedBy(user1)
	assert.NoError(err)
	p1, err := e.SaveProject(&data.Project{
		Name:      "Public Project",
		StartTime: MustParseTime("2008/01/01 00:00:00"),
		EndTime:   MustParseTime("2009/01/01 00:00:00"),
		Privacy:   data.Public,
	})
	assert.NoError(err)
	assert.NoError(e.AddStationToProject(s1, p1, nil))

	p2, err := e.SaveProject(&data.Project{
		Name:      "Private Project",
		StartTime: MustParseTime("2008/02/01 00:00:00"),
		EndTime:   MustParseTime("2008/04/01 00:00:00"),
		Privacy:   data.Private,
	})
	assert.NoError(err)
	assert.NoError(e.AddStationToProject(s1, p2, nil))

	vs := NewVisibilitySlicer(e.DB)

	dvs, err := vs.Slice(e.Ctx, s1.ID)
	assert.NoError(err)
	assert.NotEmpty(dvs)

	r := repositories.NewVisibilityRepository(e.DB)

	_, err = r.Merge(e.Ctx, s1.ID, dvs)
	assert.NoError(err)
}

func TestVisibilityMergeTwiceIsNoop(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user1, err := e.AddUser()
	assert.NoError(err)
	s1, err := e.AddStationOwnedBy(user1)
	assert.NoError(err)
	p1, err := e.SaveProject(&data.Project{
		Name:      "Public Project",
		StartTime: MustParseTime("2008/01/01 00:00:00"),
		EndTime:   MustParseTime("2009/01/01 00:00:00"),
		Privacy:   data.Public,
	})
	assert.NoError(err)
	assert.NoError(e.AddStationToProject(s1, p1, nil))

	p2, err := e.SaveProject(&data.Project{
		Name:      "Private Project",
		StartTime: MustParseTime("2008/02/01 00:00:00"),
		EndTime:   MustParseTime("2008/04/01 00:00:00"),
		Privacy:   data.Private,
	})
	assert.NoError(err)
	assert.NoError(e.AddStationToProject(s1, p2, nil))

	vs := NewVisibilitySlicer(e.DB)

	dvs, err := vs.Slice(e.Ctx, s1.ID)
	assert.NoError(err)
	assert.NotEmpty(dvs)

	r := repositories.NewVisibilityRepository(e.DB)

	_, err = r.Merge(e.Ctx, s1.ID, dvs)
	assert.NoError(err)

	db1, err := r.QueryByStationID(e.Ctx, s1.ID)
	assert.NoError(err)

	_, err = r.Merge(e.Ctx, s1.ID, dvs)
	assert.NoError(err)

	db2, err := r.QueryByStationID(e.Ctx, s1.ID)
	assert.NoError(err)

	assert.Equal(len(db1), len(db2))
}
