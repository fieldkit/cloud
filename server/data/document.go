package data

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/paulmach/go.geo"
)

var (
	invalidLocationError = errors.New("invalid location")
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

type Document struct {
	ID        int64          `db:"id,omitempty"`
	SchemaID  int32          `db:"schema_id"`
	InputID   int32          `db:"input_id"`
	TeamID    *int32         `db:"team_id"`
	UserID    *int32         `db:"user_id"`
	Insertion time.Time      `db:"insertion"`
	Timestamp time.Time      `db:"timestamp"`
	Location  *Location      `db:"location"`
	Data      types.JSONText `db:"data"`
	Fixed     bool           `db:"fixed"`
	Visible   bool           `db:"visible"`
}

type DocumentsPage struct {
	Documents []*Document
}

func (d *Document) SetData(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	d.Data = jsonData
	return nil
}

func (d *Document) GetRawFields() (fields map[string]string, err error) {
	err = json.Unmarshal(d.Data, &fields)
	if err != nil {
		return nil, err
	}
	return
}

func (d *Document) GetParsedFields() (fields map[string]interface{}, err error) {
	raw, err := d.GetRawFields()
	if err != nil {
		return nil, err
	}

	fields = make(map[string]interface{})
	for key, value := range raw {
		value, err := strconv.ParseFloat(value, 32)
		if err == nil {
			if !math.IsNaN(value) {
				fields[key] = value
			}
		}
	}

	return
}
