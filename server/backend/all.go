package backend

import (
	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/handlers"
)

func NewAllHandlers(db *sqlxcache.DB) (RecordHandler, error) {
	stationModel := handlers.NewStationModelRecordHandler(db)
	aggregating := handlers.NewAggregatingHandler(db)

	handlers := []RecordHandler{
		stationModel,
		aggregating,
	}

	return NewMuxRecordHandler(handlers), nil
}
