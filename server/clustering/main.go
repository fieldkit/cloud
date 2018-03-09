package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/paulmach/go.geo"

	"github.com/hongshibao/go-kdtree"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

const (
	RadToDeg             = 180 / math.Pi
	DegToRad             = math.Pi / 180
	EarthsRadiusInMeters = 6371e3
)

type Sample struct {
	ID        int64
	Timestamp time.Time
	Location  data.Location
}

func Havershine(l1, l2 data.Location) float64 {
	x := l1.Coordinates()
	y := l2.Coordinates()

	a1 := x[1] * DegToRad
	a2 := y[1] * DegToRad

	b1 := (y[1] - x[1]) * DegToRad
	b2 := (y[0] - x[0]) * DegToRad

	a := math.Sin(b1/2)*math.Sin(b1/2) + math.Cos(a1)*math.Cos(a2)*math.Sin(b2/2)*math.Sin(b2/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthsRadiusInMeters * c
}

func (s Sample) Distance(other kdtree.Point) float64 {
	x := s.Location
	y := (other.(*Sample)).Location
	return Havershine(x, y)
}

func (s Sample) GetID() string {
	return fmt.Sprint(s.ID)
}

func (p Sample) Dim() int {
	return len(p.Location.Coordinates())
}

func (p Sample) GetValue(dimension int) float64 {
	return p.Location.Coordinates()[dimension]
}

func (p Sample) PlaneDistance(value float64, dimension int) float64 {
	t := p.GetValue(dimension) - value
	return t * t
}

type options struct {
	PostgresUrl string
}

func main() {
	o := options{}

	flag.StringVar(&o.PostgresUrl, "postgres-url", "postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable", "url to the postgres server")

	flag.Parse()

	db, err := sqlxcache.Open("postgres", o.PostgresUrl)
	if err != nil {
		panic(err)
	}

	ctx := context.TODO()

	log.Printf("Querying...")

	started := time.Now()
	clusterables := make([]Clusterable, 0)
	kps := make([]kdtree.Point, 0)

	pageSize := 200
	page := 0
	for {
		batch := []*Sample{}
		if err := db.SelectContext(ctx, &batch, `
		  SELECT d.id, d.timestamp, ST_AsBinary(d.location) AS location
		  FROM fieldkit.record AS d
		  WHERE d.visible AND d.source_id = $1 AND ST_X(d.location) != 0 AND ST_Y(d.location) != 0
		  ORDER BY timestamp
		  LIMIT $2 OFFSET $3
		`, 4, pageSize, page*pageSize); err != nil {
			panic(err)
		}

		if len(batch) == 0 {
			break
		}

		for _, c := range batch {
			// clusterables = append(clusterables, c)
			kps = append(kps, c)
		}

		page += 1
	}

	log.Printf("Query done in %v, building tree...", time.Now().Sub(started))

	tree := kdtree.NewKDTree(kps)

	log.Printf("%v", tree)

	log.Printf("Tree ready in %v, clustering...", time.Now().Sub(started))

	eps := 50.0
	minPts := 10
	clusters := Clusterize(clusterables, tree, minPts, eps)

	for _, c := range clusters {
		log.Printf("%v", c)
	}

	log.Printf("Done in %v", time.Now().Sub(started))
}
