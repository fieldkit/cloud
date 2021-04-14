package backend

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/tests"
)

func TestExportCsv(t *testing.T) {
	ctx := context.Background()

	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	ids := make([]int32, 0)
	for _, s := range fd.Stations {
		ids = append(ids, s.ID)
		_, err := e.AddMetaAndData(s, fd.Owner, 100)
		assert.NoError(err)
	}

	formatter := NewCsvFormatter()
	exporter, err := NewExporter(e.DB)
	assert.Nil(err)
	assert.Nil(exporter.Export(ctx, &ExportCriteria{
		Stations: ids,
		Sensors:  []ModuleAndSensor{},
		Start:    time.Now().Add(-24 * time.Hour),
		End:      time.Now().Add(24 * time.Hour),
	}, formatter, ExportProgressNoop, ioutil.Discard))
}
