package main

import (
	"context"
	"flag"
	_ "fmt"
	"log"
	"time"

	"encoding/json"
	_ "strings"

	_ "github.com/lib/pq"
	_ "github.com/paulmach/go.geo"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type options struct {
	PostgresUrl string
	DryRun      bool
}

type Row struct {
	ID        int64
	Timestamp time.Time
	Location  data.Location
	Data      []byte
}

var translations = map[string]string{
	"DO":       "do",
	"pH":       "ph",
	"ORP":      "orp",
	"Temp":     "temp",
	"Ec":       "ec",
	"Salinity": "salinity",
	"TDS":      "tds",
	"SG":       "sg",
	"td":       "tds",

	"Temp #1":         "temp_1",
	"Humidity":        "humidity",
	"Temp #2":         "temp_2",
	"Pressure":        "pressure",
	"Altitude":        "altitude",
	"Ambient IR":      "light_ir",
	"Ambient Visible": "light_visible",
	"Ambient Lux":     "light_lux",
	"IMU Cal":         "imu_cal",
	"IMU Orien X":     "imu_orien_x",
	"IMU Orien Y":     "imu_orien_y",
	"IMU Orien Z":     "imu_orien_z",

	"Wind Speed":     "wind_speed",
	"Wind Dir":       "wind_dir",
	"Hr Wind Speed":  "wind_hr_speed",
	"Hr Wind Dir":    "wind_hr_dir",
	"10m Wind Gust":  "wind_10m_gust",
	"10m Wind Dir":   "wind_10m_dir",
	"2m Wind Gust":   "wind_2m_gust",
	"2m Wind Dir":    "wind_2m_dir",
	"Prev Hrly Rain": "rain_prev_hour",
	"Hourly Rain":    "rain_this_hour",

	"Daily Rain":     "daily_rain",
	"Day Wind Speed": "day_wind_speed",
	"Day Wind Dir":   "day_wind_dir",
}

func reverseMap(m map[string]string) map[string]string {
	n := make(map[string]string)
	for k, v := range m {
		n[v] = k
	}
	return n
}

var reverseTranslations = reverseMap(translations)

func fixKeyNames(row *Row) (bool, error) {
	original := make(map[string]string)
	modified := make(map[string]string)

	err := json.Unmarshal(row.Data, &original)
	if err != nil {
		return false, err
	}

	fixed := false
	for key, value := range original {
		if reverseTranslations[key] == "" && translations[key] != "" {
			fixed = true
			modified[translations[key]] = value
		} else {
			modified[key] = value
		}
		_, _ = key, value
	}

	raw, err := json.Marshal(&modified)
	if err != nil {
		return false, err
	}

	row.Data = raw

	return fixed, nil
}

func fixAll(ctx context.Context, o *options, db *sqlxcache.DB) error {
	log.Printf("Querying...")

	updated := 0
	pageSize := 200
	page := 0
	for {
		batch := []*Row{}
		if err := db.SelectContext(ctx, &batch, `
		  SELECT d.id, d.timestamp, ST_AsBinary(d.location) AS location, d.data
		  FROM fieldkit.document AS d
		  ORDER BY timestamp
		  LIMIT $1 OFFSET $2
		`, pageSize, page*pageSize); err != nil {
			panic(err)
		}

		if len(batch) == 0 {
			break
		}

		log.Printf("Processing %d records", len(batch))
		for _, row := range batch {
			modified, err := fixKeyNames(row)
			if err != nil {
				return err
			}

			if modified {
				if o.DryRun {
					log.Printf("%v %v", row.ID, string(row.Data))
				} else {
					_, err = db.ExecContext(ctx, `UPDATE fieldkit.document SET data = $1 WHERE id = $2`, row.Data, row.ID)
					if err != nil {
						return err
					}
				}

				updated += 1
			}
		}

		page += 1
	}

	log.Printf("Fixed %d rows", updated)

	return nil
}

func main() {
	o := options{}

	flag.StringVar(&o.PostgresUrl, "postgres-url", "postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable", "url to the postgres server")
	flag.BoolVar(&o.DryRun, "dry-run", true, "dry run (verbose display and no updates)")

	flag.Parse()

	db, err := sqlxcache.Open("postgres", o.PostgresUrl)
	if err != nil {
		panic(err)
	}

	ctx := context.TODO()

	fixAll(ctx, &o, db)
}
