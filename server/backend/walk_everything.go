package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/profile"

	"github.com/golang/protobuf/proto"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/data"

	"github.com/govau/que-go"
)

func walkEverything(ctx context.Context, j *que.Job, services *BackgroundServices, tm *jobs.TransportMessage) error {
	log := Logger(ctx).Sugar()

	if false {
		defer profile.Start().Stop()
	}

	stationIDs := make([]int32, 0)
	if err := services.database.SelectContext(ctx, &stationIDs, `SELECT id FROM fieldkit.station`); err != nil {
		return err
	}

	for _, id := range stationIDs {
		log.Infow("station", "station_id", id)

		walkParams := &WalkParameters{
			Start:      time.Time{},
			End:        time.Now().Add(1 * time.Hour),
			StationIDs: []int32{id},
			PageSize:   1000,
			Page:       0,
		}

		walker := NewRecordWalker(services.database)

		handler := &fixingHandler{
			services: services,
		}

		if err := walker.WalkStation(ctx, handler, walkParams); err != nil {
			return err
		}
	}

	return nil
}

type fixingHandler struct {
	services *BackgroundServices
}

func (h *fixingHandler) OnMeta(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.MetaRecord) error {
	if db.PB == nil {
		r := &pb.DataRecord{}
		if err := db.Unmarshal(r); err != nil {
			return err
		}

		data := proto.NewBuffer(make([]byte, 0))
		data.Marshal(r)

		if _, err := h.services.database.ExecContext(ctx, `
			UPDATE fieldkit.meta_record SET pb = $1 WHERE id = $2
		`, data.Bytes(), db.ID); err != nil {
			return fmt.Errorf("error updating pb: %v", err)
		}
	}
	return nil
}

func (h *fixingHandler) OnData(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.DataRecord, meta *data.MetaRecord) error {
	if db.PB == nil {
		r := &pb.DataRecord{}
		if err := db.Unmarshal(r); err != nil {
			return err
		}

		data := proto.NewBuffer(make([]byte, 0))
		data.Marshal(r)

		if _, err := h.services.database.ExecContext(ctx, `
			UPDATE fieldkit.data_record SET pb = $1 WHERE id = $2
		`, data.Bytes(), db.ID); err != nil {
			return fmt.Errorf("error updating pb: %v", err)
		}
	}
	return nil
}

func (h *fixingHandler) OnDone(ctx context.Context) error {
	return nil
}
