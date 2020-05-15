package data

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

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

func (l *Location) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Coordinates())
}

func (l *Location) UnmarshalJSON(data []byte) error {
	c := []float64{}
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
	l.point = geo.NewPointFromLatLng(c[1], c[0])
	return nil
}

func (l *Location) IsZero() bool {
	return l.point.Lng() == 0 && l.point.Lat() == 0
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
	if l == nil || l.point == nil {
		return nil, nil
	}
	return l.point.ToWKT(), nil
}

func (l *Location) String() string {
	if l.point != nil {
		return l.point.String()
	}
	return "<nil>"
}

func (l *Location) Distance(o *Location) float64 {
	return l.point.GeoDistanceFrom(o.point)
}

func (l *Location) Latitude() float64 {
	return l.point.Lat()
}

func (l *Location) Longitude() float64 {
	return l.point.Lng()
}

func (l *Location) Valid() bool {
	if l == nil {
		return false
	}

	// Simple range checks...
	if l.Latitude() > 90 || l.Latitude() < -90 || l.Longitude() > 180 || l.Longitude() < -180 {
		return false
	}

	// HACK jlewallen
	if l.Longitude() == 0.0 && l.Latitude() == 0.0 {
		return false
	}

	return true
}
