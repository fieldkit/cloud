package backend

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

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
	Sensors  []ModuleAndSensor
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
	// sensorsByID map[int64]*data.Sensor
	errors int32
}

type stationInfo struct {
	provision *data.Provision
	meta      *data.MetaRecord
	station   *data.Station
}

func NewExporter(db *sqlxcache.DB) (*Exporter, error) {
	return &Exporter{
		db:          db,
		metaFactory: repositories.NewMetaFactory(db),
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
	mr := repositories.NewModuleMetaRepository(e.db)

	sr := repositories.NewStationRepository(e.db)

	allSensors := []*data.Sensor{}
	if err := e.db.SelectContext(ctx, &allSensors, `SELECT * FROM fieldkit.aggregated_sensor ORDER BY key`); err != nil {
		return err
	}

	allSensorsByID := make(map[int64]*data.Sensor)
	for _, sensor := range allSensors {
		allSensorsByID[sensor.ID] = sensor
	}

	log := Logger(ctx).Sugar()

	e.sensors = make([]*repositories.SensorAndModuleMeta, 0)
	for _, moduleAndSensor := range criteria.Sensors {
		if sensor, ok := allSensorsByID[moduleAndSensor.SensorID]; !ok {
			return fmt.Errorf("invalid sensor")
		} else {
			if sensor.ID != moduleAndSensor.SensorID {
				panic("wtf")
			}

			sam, err := mr.FindByFullKey(ctx, sensor.Key)
			if err != nil {
				return err
			}

			log.Infow("exporting", "lookup_key", sensor.Key, "sensor_key", sam.Sensor.FullKey, "sensor_id", moduleAndSensor.SensorID, "module_id", moduleAndSensor.ModuleID)
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
		sr := repositories.NewStationRepository(e.db)

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
		e.errors += 1
		return nil
	}

	sensors := make([]*repositories.ReadingValue, 0)
	for _, sam := range e.sensors {
		appended := false
		for key, rv := range filtered.Record.Readings {
			if key.SensorKey == sam.Sensor.FullKey {
				sensors = append(sensors, rv)
				appended = true
				break
			}
		}
		if !appended {
			sensors = append(sensors, nil)
		}
	}

	if err := e.writer.Row(ctx, info.station, sensors, filtered); err != nil {
		return err
	}

	return nil
}

func (e *Exporter) OnDone(ctx context.Context) error {
	if e.errors > 0 {
		log := Logger(ctx).Sugar()
		log.Warnw("exporting", "errors", e.errors)
	}
	return nil
}
