package backend

import (
	"context"

	"github.com/vgarvardt/gue/v4"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/messages"
)

type DescribeStationLocationHandler struct {
	db        *sqlxcache.DB
	metrics   *logging.Metrics
	publisher jobs.MessagePublisher
	locations *data.DescribeLocations
}

func NewDescribeStationLocationHandler(db *sqlxcache.DB, metrics *logging.Metrics, publisher jobs.MessagePublisher, locations *data.DescribeLocations) *DescribeStationLocationHandler {
	return &DescribeStationLocationHandler{
		db:        db,
		metrics:   metrics,
		publisher: publisher,
		locations: locations,
	}
}

func (h *DescribeStationLocationHandler) Handle(ctx context.Context, m *messages.StationLocationUpdated, j *gue.Job) error {
	log := Logger(ctx).Sugar().With("station_id", m.StationID)

	if !h.locations.IsEnabled() {
		log.Infow("describing-location:disabled")
		return nil
	}

	log.Infow("describing-location")

	location := data.NewLocation(m.Location)

	names, err := h.locations.Describe(ctx, location)
	if err != nil {
		return err
	} else if names != nil {
		stations := repositories.NewStationRepository(h.db)

		station, err := stations.QueryStationByID(ctx, m.StationID)
		if err != nil {
			return err
		}

		station.PlaceOther = names.OtherLandName
		station.PlaceNative = names.NativeLandName

		if err := stations.UpdateStation(ctx, station); err != nil {
			return err
		}
	}

	return nil
}
