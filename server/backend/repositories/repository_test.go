package repositories

import (
	"testing"

	"github.com/fieldkit/cloud/server/tests"

	"github.com/stretchr/testify/assert"
)

func TestQueryStationByID(t *testing.T) {
	assert := assert.New(t)
	_, err := tests.NewTestEnv()
	assert.NoError(err)
}
