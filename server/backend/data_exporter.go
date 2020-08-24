package backend

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/conservify/sqlxcache"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

type ExportProgressFunc func(ctx context.Context, progress float64, message string) error

func ExportProgressNoop(ctx context.Context, progress float64, message string) error {
	return nil
}

type ExportWriter interface {
	Begin(ctx context.Context, stations []*data.Station, sensors []*repositories.SensorAndModuleMeta) error
	Row(ctx context.Context, station *data.Station, sensors []*repositories.ReadingValue, row *repositories.FilteredRecord) error
	End(ctx context.Context) error
}

type ExportFormat interface {
	MimeType() string
	CreateWriter(ctx context.Context, writer io.Writer) (ExportWriter, error)
}

type ExportCriteria struct {
	Start    time.Time
	End      time.Time
	Stations []int32
	Sensors  []int64
}

type Exporter struct {
	db          *sqlxcache.DB
	format      ExportFormat
	writer      ExportWriter
	walker      *RecordWalker
	metaFactory *repositories.MetaFactory
	byProvision map[int64]*stationInfo
	stations    []*data.Station
	sensors     []*repositories.SensorAndModuleMeta
	sensorsByID map[int64]*data.Sensor
}

type stationInfo struct {
	provision *data.Provision
	meta      *data.MetaRecord
	station   *data.Station
}

func NewExporter(db *sqlxcache.DB) (*Exporter, error) {
	return &Exporter{
		db:          db,
		metaFactory: repositories.NewMetaFactory(),
		walker:      NewRecordWalker(db),
		byProvision: make(map[int64]*stationInfo),
	}, nil
}

func contains(ids []int64, value int64) bool {
	for _, v := range ids {
		if v == value {
			return true
		}
	}
	return false
}

func (e *Exporter) Export(ctx context.Context, criteria *ExportCriteria, format ExportFormat, progressFunc ExportProgressFunc, writer io.Writer) error {
	mr := repositories.NewModuleMetaRepository()

	sr, err := repositories.NewStationRepository(e.db)
	if err != nil {
		return err
	}

	allSensors := []*data.Sensor{}
	if err := e.db.SelectContext(ctx, &allSensors, `SELECT * FROM fieldkit.aggregated_sensor ORDER BY key`); err != nil {
		return err
	}

	e.sensors = make([]*repositories.SensorAndModuleMeta, 0)
	e.sensorsByID = make(map[int64]*data.Sensor)
	for _, sensor := range allSensors {
		if len(criteria.Sensors) == 0 || contains(criteria.Sensors, sensor.ID) {
			sam, err := mr.FindByFullKey(sensor.Key)
			if err != nil {
				return err
			}
			e.sensorsByID[sensor.ID] = sensor
			e.sensors = append(e.sensors, sam)
		}
	}

	e.stations = make([]*data.Station, 0)
	for _, id := range criteria.Stations {
		if station, err := sr.QueryStationByID(ctx, id); err != nil {
			return err
		} else {
			e.stations = append(e.stations, station)
		}
	}

	formatWriter, err := format.CreateWriter(ctx, writer)
	if err != nil {
		return err
	}

	e.format = format
	e.writer = formatWriter

	if err := e.writer.Begin(ctx, e.stations, e.sensors); err != nil {
		return err
	}

	walkParams := &WalkParameters{
		Start:      criteria.Start,
		End:        criteria.End,
		StationIDs: criteria.Stations,
		PageSize:   1000,
		Page:       0,
	}

	walkerProgress := func(ctx context.Context, progress float64) error {
		return progressFunc(ctx, progress, "")
	}

	if err := e.walker.WalkStation(ctx, e, walkerProgress, walkParams); err != nil {
		return err
	}

	if err := e.writer.End(ctx); err != nil {
		return err
	}

	return nil
}

func (e *Exporter) OnMeta(ctx context.Context, p *data.Provision, r *pb.DataRecord, meta *data.MetaRecord) error {
	if _, ok := e.byProvision[p.ID]; !ok {
		sr, err := repositories.NewStationRepository(e.db)
		if err != nil {
			return err
		}

		station, err := sr.QueryStationByDeviceID(ctx, p.DeviceID)
		if err != nil {
			return err
		}

		if station == nil {
			return fmt.Errorf("no station")
		}

		e.byProvision[p.ID] = &stationInfo{
			meta:      meta,
			provision: p,
			station:   station,
		}
	}

	_, err := e.metaFactory.Add(ctx, meta, true)
	if err != nil {
		return err
	}

	return nil
}

func (e *Exporter) OnData(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.DataRecord, meta *data.MetaRecord) error {
	info := e.byProvision[p.ID]
	if info == nil {
		// This should never happen, we should error in OnMeta if
		// we're unable to find the station.
		return fmt.Errorf("unknown provision: %d", p.ID)
	}

	filtered, err := e.metaFactory.Resolve(ctx, db, false, true)
	if err != nil {
		return fmt.Errorf("resolving: %v", err)
	}
	if filtered == nil {
		return fmt.Errorf("resolving: unknown")
	}

	sensors := make([]*repositories.ReadingValue, 0)
	for _, sam := range e.sensors {
		key := sam.Sensor.FullKey
		if rv, ok := filtered.Record.Readings[key]; !ok {
			sensors = append(sensors, nil)
		} else {
			sensors = append(sensors, rv)
		}
	}

	if err := e.writer.Row(ctx, info.station, sensors, filtered); err != nil {
		return err
	}

	return nil
}

func (e *Exporter) OnDone(ctx context.Context) error {
	return nil
}
