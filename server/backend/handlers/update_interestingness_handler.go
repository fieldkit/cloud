package handlers

import (
	"context"
	"fmt"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

type InterestingnessHandler struct {
	repository   *repositories.InterestingnessRepository
	querySensors *repositories.SensorsRepository
	sensors      map[string]*data.Sensor
}

func NewInterestingnessHandler(db *sqlxcache.DB) *InterestingnessHandler {
	return &InterestingnessHandler{
		repository:   repositories.NewInterestingnessRepository(db),
		querySensors: repositories.NewSensorsRepository(db),
	}
}

func (ih *InterestingnessHandler) ConsiderReading(ctx context.Context, r *data.IncomingReading) error {
	if ih.sensors == nil {
		sensorsMap, err := ih.querySensors.QueryAllSensors(ctx)
		if err != nil {
			return err
		}
		ih.sensors = sensorsMap
	}

	sensor := ih.sensors[r.SensorKey]

	if sensor == nil {
		return fmt.Errorf("unknown sensor: %v", r.SensorKey)
	}

	if sensor.InterestingnessPriority == nil {
		return nil
	}

	fn := data.NewMaximumInterestingessFunction()

	interestingness := fn.Calculate(r.Value)

	existing, err := ih.repository.QueryByStationID(ctx, r.StationID)
	if err != nil {
		return err
	}

	for _, window := range data.Windows {
		if !window.Includes(r.Time) {
			// Windows are sorted in descending order so that when we reach one that
			// doesn't include a reading we can be sure that none of the remaining
			// windows include the reading either.
			break
		}

		log := Logger(ctx).Sugar().With("station_id", r.StationID, "sensor_id", sensor.ID, "reading_time", r.Time, "window_duration", window.Duration, "iness", interestingness)

		rowForThis := &data.StationInterestingness{
			StationID:       r.StationID,
			WindowSeconds:   int32(window.Duration.Seconds()),
			Interestingness: interestingness,
			ReadingSensorID: sensor.ID,
			ReadingModuleID: r.ModuleID,
			ReadingValue:    r.Value,
			ReadingTime:     r.Time,
		}

		// Do we have an existing interestingness for this duration?
		if e, ok := existing[window.Duration]; ok {
			// Is the incoming reading is more interesting than the one we have?
			if fn.MoreThan(interestingness, e.Interestingness) {
				// This reading is now the most interesting.
				log.Infow("iness:updating", "old_iness", e.Interestingness)
				if _, err := ih.repository.UpsertInterestingness(ctx, rowForThis); err != nil {
					return err
				}
			}
		} else {
			// This reading is the most interesting for this window by default.
			log.Infow("iness:inserting")
			if _, err := ih.repository.UpsertInterestingness(ctx, rowForThis); err != nil {
				return err
			}
		}

	}

	return nil
}
