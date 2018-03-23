package data

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/paulmach/go.geo"
)

var (
	invalidLocationError = errors.New("Invalid location")
)

type Location struct {
	point *geo.Point
}

func NewLocation(coordinates []float64) *Location {
	return &Location{
		point: geo.NewPointFromLatLng(coordinates[1], coordinates[0]),
	}
}

func (l *Location) Coordinates() []float64 {
	return []float64{l.point.Lng(), l.point.Lat()}
}

func (l *Location) Scan(data interface{}) error {
	point := &geo.Point{}
	if err := point.Scan(data); err != nil {
		return err
	}

	l.point = point
	return nil
}

func (l *Location) Value() (driver.Value, error) {
	return l.point.ToWKT(), nil
}

func (l *Location) String() string {
	return l.point.String()
}

type RecordAnalysis struct {
	RecordID         int64 `db:"record_id"`
	ManuallyExcluded bool  `db:"manually_excluded"`
	Outlier          bool  `db:"outlier"`
}

type Record struct {
	ID        int64          `db:"id,omitempty"`
	SchemaID  int32          `db:"schema_id"`
	SourceID  int32          `db:"source_id"`
	TeamID    *int32         `db:"team_id"`
	UserID    *int32         `db:"user_id"`
	Insertion time.Time      `db:"insertion"`
	Timestamp time.Time      `db:"timestamp"`
	Location  *Location      `db:"location"`
	Data      types.JSONText `db:"data"`
	Fixed     bool           `db:"fixed"`
	Visible   bool           `db:"visible"`
	Metadata  bool           `db:"metadata"`
}

type AnalysedRecord struct {
	Record
	RecordAnalysis
}

type RecordsPage struct {
	Records []*Record
}

type AnalysedRecordsPage struct {
	Records []*AnalysedRecord
}

func (d *Record) SetData(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	d.Data = jsonData
	return nil
}

func (d *Record) GetRawFields() (fields map[string]interface{}, err error) {
	err = json.Unmarshal(d.Data, &fields)
	if err != nil {
		return nil, err
	}
	return
}

func (d *Record) GetParsedFields() (fields map[string]interface{}, err error) {
	return d.GetRawFields()
}
