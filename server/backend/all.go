package backend

import (
	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/storage"

	"github.com/fieldkit/cloud/server/backend/handlers"
)

func NewAllHandlers(db *sqlxcache.DB, tsConfig *storage.TimeScaleDBConfig) (RecordHandler, error) {
	return handlers.NewStationModelRecordHandler(db), nil
}
