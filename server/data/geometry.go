package data

import (
	"database/sql/driver"

	"github.com/paulmach/go.geo"
)

type TemporalPath struct {
	path *geo.Path
}

func (l *TemporalPath) Coordinates() [][]float64 {
	if l.path == nil {
		return make([][]float64, 0)
	}
	ps := l.path.Points()
	c := make([][]float64, len(ps))
	for i, p := range ps {
		c[i] = []float64{p.Lng(), p.Lat()}
	}
	return c
}

func (l *TemporalPath) Scan(data interface{}) error {
	path := &geo.Path{}
	if err := path.Scan(data); err != nil {
		return err
	}

	l.path = path
	return nil
}

func (l *TemporalPath) Value() (driver.Value, error) {
	return l.path.ToWKT(), nil
}

type Envelope struct {
	geometry *geo.PointSet
}

func (l *Envelope) Coordinates() [][]float64 {
	if l.geometry == nil {
		return make([][]float64, 0)
	}
	c := make([][]float64, l.geometry.Length())
	for i := 0; i < l.geometry.Length(); i += 1 {
		p := l.geometry.GetAt(i)
		c[i] = []float64{p.Lng(), p.Lat()}
	}
	return c
}

func (l *Envelope) Scan(data interface{}) error {
	geometry := &geo.PointSet{}
	if err := geometry.Scan(data); err != nil {
		p := &geo.Point{}
		if err := p.Scan(data); err != nil {
			return err
		}
		geometry.Push(p)
	}

	l.geometry = geometry
	return nil
}

func (l *Envelope) Value() (driver.Value, error) {
	return l.geometry.ToWKT(), nil
}
