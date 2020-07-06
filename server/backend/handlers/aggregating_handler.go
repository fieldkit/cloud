package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/montanaflynn/stats"

	"github.com/conservify/sqlxcache"
	"github.com/jmoiron/sqlx"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

var (
	AggregateNames = []string{"24h", "12h", "6h", "1h", "30m", "10m", "1m"}
)

var (
	AggregateTableNames = map[string]string{
		"24h": "fieldkit.aggregated_24h",
		"12h": "fieldkit.aggregated_12h",
		"6h":  "fieldkit.aggregated_6h",
		"1h":  "fieldkit.aggregated_1h",
		"30m": "fieldkit.aggregated_30m",
		"10m": "fieldkit.aggregated_10m",
		"1m":  "fieldkit.aggregated_1m",
	}
)

type aggregation struct {
	opened   time.Time
	interval time.Duration
	table    string
	values   map[string][]float64
	location []float64
}

type Aggregated struct {
	Time     time.Time
	Location []float64
	Values   map[string]float64
}

func (a *aggregation) getTime(original time.Time) time.Time {
	return original.Truncate(a.interval)
}

func (a *aggregation) canAdd(t time.Time) bool {
	return a.opened.IsZero() || a.opened == t
}

func (a *aggregation) add(t time.Time, key string, location []float64, value float64) error {
	if !a.canAdd(t) {
		return fmt.Errorf("unable to add to aggregation")
	}

	if a.opened != t && a.opened.IsZero() {
		a.opened = t
	}

	if a.values[key] == nil {
		a.values[key] = make([]float64, 0, 1)
	}

	a.values[key] = append(a.values[key], value)
	if len(location) == 3 && location[0] != 0 && location[1] != 0 {
		a.location = location
	}

	return nil
}

func (a *aggregation) close() (*Aggregated, error) {
	agg := make(map[string]float64)

	for key, values := range a.values {
		if len(values) > 0 {
			mean, err := stats.Mean(values)
			if err != nil {
				return nil, fmt.Errorf("error taking mean: %v", err)
			}

			agg[key] = mean

			a.values[key] = a.values[key][:0]
		}
	}

	aggregated := &Aggregated{
		Time:   a.opened,
		Values: agg,
	}

	a.opened = time.Time{}

	return aggregated, nil
}

type AggregatingHandler struct {
	db           *sqlxcache.DB
	metaFactory  *repositories.MetaFactory
	aggregations []*aggregation
	sensors      map[string]*data.Sensor
	stations     map[int64]int32
}

func NewAggregatingHandler(db *sqlxcache.DB) *AggregatingHandler {
	return &AggregatingHandler{
		db:          db,
		metaFactory: repositories.NewMetaFactory(),
		stations:    make(map[int64]int32),
		aggregations: []*aggregation{
			&aggregation{
				interval: time.Minute * 1,
				table:    "fieldkit.aggregated_1m",
				values:   make(map[string][]float64),
			},
			&aggregation{
				interval: time.Minute * 10,
				table:    "fieldkit.aggregated_10m",
				values:   make(map[string][]float64),
			},
			&aggregation{
				interval: time.Minute * 30,
				table:    "fieldkit.aggregated_30m",
				values:   make(map[string][]float64),
			},
			&aggregation{
				interval: time.Hour * 1,
				table:    "fieldkit.aggregated_1h",
				values:   make(map[string][]float64),
			},
			&aggregation{
				interval: time.Hour * 6,
				table:    "fieldkit.aggregated_6h",
				values:   make(map[string][]float64),
			},
			&aggregation{
				interval: time.Hour * 12,
				table:    "fieldkit.aggregated_12h",
				values:   make(map[string][]float64),
			},
			&aggregation{
				interval: time.Hour * 24,
				table:    "fieldkit.aggregated_24h",
				values:   make(map[string][]float64),
			},
		},
	}
}

func (v *AggregatingHandler) refreshSensors(ctx context.Context) error {
	sensors := []*data.Sensor{}
	if err := v.db.SelectContext(ctx, &sensors, `SELECT * FROM fieldkit.aggregated_sensor ORDER BY key`); err != nil {
		return err
	}

	v.sensors = make(map[string]*data.Sensor)

	for _, sensor := range sensors {
		v.sensors[sensor.Key] = sensor
	}

	return nil
}

func (v *AggregatingHandler) getStationID() (int32, error) {
	for _, v := range v.stations {
		return v, nil
	}
	return 0, fmt.Errorf("no such station id")
}

func (v *AggregatingHandler) upsertAggregated(ctx context.Context, a *aggregation, d *Aggregated) error {
	var location *data.Location

	if d.Location != nil {
		location = data.NewLocation(d.Location)
	}

	if v.sensors == nil {
		if err := v.refreshSensors(ctx); err != nil {
			return err
		}
	}

	stationID, err := v.getStationID()
	if err != nil {
		return err
	}

	keeping := make([]int64, 0)

	for key, value := range d.Values {
		if v.sensors[key] == nil {
			newSensor := &data.Sensor{
				Key: key,
			}
			if err := v.db.NamedGetContext(ctx, newSensor, `
				INSERT INTO fieldkit.aggregated_sensor (key) VALUES (:key) RETURNING id
				`, newSensor); err != nil {
				return fmt.Errorf("error adding aggregated sensor: %v", err)
			}
			v.sensors[key] = newSensor
		}

		row := &data.AggregatedReading{
			StationID: stationID,
			SensorID:  v.sensors[key].ID,
			Time:      data.NumericWireTime(d.Time),
			Location:  location,
			Value:     value,
		}

		if err := v.db.NamedGetContext(ctx, row, fmt.Sprintf(`
			INSERT INTO %s (station_id, sensor_id, time, location, value)
			VALUES (:station_id, :sensor_id, :time, ST_SetSRID(ST_GeomFromText(:location), 4326), :value)
			ON CONFLICT (time, station_id, sensor_id) DO UPDATE SET value = EXCLUDED.value, location = EXCLUDED.location
			RETURNING id
			`, a.table), row); err != nil {
			return fmt.Errorf("error upserting sensor reading: %v", err)
		}

		keeping = append(keeping, row.ID)
	}

	// Delete any other rows for this station in this time.
	if len(keeping) > 0 {
		query, args, err := sqlx.In(fmt.Sprintf(`
			DELETE FROM %s WHERE station_id = ? AND time = ? AND id NOT IN (?)
		`, a.table), stationID, d.Time, keeping)
		if err != nil {
			return err
		}
		if _, err := v.db.ExecContext(ctx, v.db.Rebind(query), args...); err != nil {
			return fmt.Errorf("error deleting old sensor readings: %v", err)
		}
	} else {
		if _, err := v.db.ExecContext(ctx, fmt.Sprintf(`
			DELETE FROM %s WHERE station_id = $1 AND time = $2
		`, a.table), stationID, d.Time); err != nil {
			return fmt.Errorf("error deleting old sensor readings: %v", err)
		}
	}

	return nil
}

func (v *AggregatingHandler) OnMeta(ctx context.Context, p *data.Provision, r *pb.DataRecord, meta *data.MetaRecord) error {
	if _, ok := v.stations[p.ID]; !ok {
		sr, err := repositories.NewStationRepository(v.db)
		if err != nil {
			return err
		}

		station, err := sr.QueryStationByDeviceID(ctx, p.DeviceID)
		if err != nil || station == nil {
			// mark giving up
			return nil
		}

		v.stations[p.ID] = station.ID
	}

	_, err := v.metaFactory.Add(ctx, meta, true)
	if err != nil {
		return err
	}

	return nil
}

func (v *AggregatingHandler) OnData(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.DataRecord, meta *data.MetaRecord) error {
	if _, ok := v.stations[p.ID]; !ok {
		return nil
	}

	filtered, err := v.metaFactory.Resolve(ctx, db, false, true)
	if err != nil {
		return fmt.Errorf("error resolving: %v", err)
	}
	if filtered == nil {
		return nil
	}

	for _, aggregation := range v.aggregations {
		time := aggregation.getTime(db.Time)

		if !aggregation.canAdd(time) {
			if values, err := aggregation.close(); err != nil {
				return fmt.Errorf("error closing aggregation: %v", err)
			} else {
				if err := v.upsertAggregated(ctx, aggregation, values); err != nil {
					return fmt.Errorf("error upserting aggregated: %v", err)
				}
				if false {
					fmt.Printf("%v\n", values)
				}
			}
		}

		for key, value := range filtered.Record.Readings {
			if err := aggregation.add(time, key, filtered.Record.Location, value.Value); err != nil {
				return fmt.Errorf("error adding: %v", err)
			}
		}
	}

	return nil
}

func (v *AggregatingHandler) OnDone(ctx context.Context) error {
	if len(v.stations) == 0 {
		return nil
	}
	for _, aggregation := range v.aggregations {
		if values, err := aggregation.close(); err != nil {
			return fmt.Errorf("error closing aggregation: %v", err)
		} else {
			if err := v.upsertAggregated(ctx, aggregation, values); err != nil {
				return fmt.Errorf("error upserting aggregated: %v", err)
			}
			if false {
				fmt.Printf("%v\n", values)
			}
		}
	}
	return nil
}
