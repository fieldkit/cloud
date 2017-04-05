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
	invalidLocationError = errors.New("invalid location")
)

type Location struct {
	point *geo.Point
}

func NewLocation(longitude, latitude float64) *Location {
	return &Location{
		point: geo.NewPoint(longitude, latitude),
	}
}

func (l *Location) Coordinates() []float64 {
	coordinates := l.point.ToArray()
	return coordinates[:]
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

type Document struct {
	ID        int64          `db:"id,omitempty"`
	SchemaID  int32          `db:"schema_id"`
	InputID   int32          `db:"input_id"`
	TeamID    *int32         `db:"team_id"`
	UserID    *int32         `db:"user_id"`
	Timestamp time.Time      `db:"timestamp"`
	Location  *Location      `db:"location"`
	Data      types.JSONText `db:"data"`
}

func (d *Document) SetData(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	d.Data = jsonData
	return nil
}
