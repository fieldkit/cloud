package backend

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/tests"
)

type ExampleSaga struct {
	Name string `json:"name"`
}

func TestQueryAndInsertSaga(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	sr := jobs.NewSagaRepository(e.DbPool)

	id := jobs.NewSagaID()

	_, err = sr.FindByID(e.Ctx, id)
	assert.NoError(err)

	saga := jobs.NewSaga(jobs.WithID(id))

	if err := saga.SetBody(&ExampleSaga{Name: "Carla"}); err != nil {
		assert.NoError(err)
	}

	err = sr.Upsert(e.Ctx, saga)
	assert.NoError(err)
	assert.Equal(saga.Version, 1)

	f1, err := sr.FindByID(e.Ctx, id)
	assert.NoError(err)
	assert.NotNil(f1)
	assert.Equal(f1.ID, saga.ID)
	assert.Nil(f1.ScheduledAt)
	assert.Equal(f1.Version, 1)

	f1.Schedule(3 * time.Hour)

	err = sr.Upsert(e.Ctx, f1)
	assert.NoError(err)
	assert.Equal(f1.Version, 2)

	f2, err := sr.FindByID(e.Ctx, id)
	assert.NoError(err)
	assert.NotNil(f2)
	assert.NotNil(f2.ScheduledAt)
	assert.Equal(f2.Version, 2)

	example := ExampleSaga{}
	if err := saga.GetBody(&example); err != nil {
		assert.NoError(err)
	} else {
		assert.Equal(example.Name, "Carla")
	}

	assert.NoError(sr.Delete(e.Ctx, f2))

	f3, err := sr.FindByID(e.Ctx, id)
	assert.NoError(err)
	assert.Nil(f3)
}
