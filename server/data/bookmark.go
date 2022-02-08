package data

import (
	"encoding/json"
	"fmt"
	"time"
)

type VizBookmark = []interface{} // [Stations, Sensors, [number, number], [[number, number], [number, number]] | [], ChartType, FastTime] | [];
type GroupBookmark = [][]VizBookmark

// {"v":1,"g":[ [[ [[159],[2],[-8640000000000000,8640000000000000],[],0,0] ]]],"s":[]}
type Bookmark struct {
	Version  int32           `json:"v"`
	Groups   []GroupBookmark `json:"g"`
	Stations []int32         `json:"s"`
}

type BookmarkSensor struct {
	StationID int32  `json:"station_id"`
	ModuleID  string `json:"module_id"`
	SensorID  int64  `json:"sensor_id"`
}

type BookmarkViz struct {
	Start       time.Time         `json:"start"`
	End         time.Time         `json:"end"`
	ExtremeTime bool              `json:"extreme_time"`
	Sensors     []*BookmarkSensor `json:"sensors"`
}

func ParseBookmark(s string) (*Bookmark, error) {
	b := Bookmark{}
	err := json.Unmarshal([]byte(s), &b)
	if err != nil {
		return nil, err
	}
	for _, gg := range b.Groups {
		for _, g := range gg {
			for _, viz := range g {
				if len(viz) == 6 {
					if stationPart, ok := viz[0].([]interface{}); !ok {
						return nil, fmt.Errorf("viz bookmark (6) stations malformed: %v", viz[0])
					} else {
						for _, id := range stationPart {
							if _, ok := id.(float64); !ok {
								return nil, fmt.Errorf("viz bookmark (6) sensors malformed: %v", id)
							}
						}
					}
				} else if len(viz) == 5 {
					if sensorsPart, ok := viz[0].([]interface{}); !ok {
						return nil, fmt.Errorf("viz bookmark (5) sensors malformed: %v", viz[0])
					} else {
						_ = sensorsPart
					}
				} else {
					return nil, fmt.Errorf("viz bookmark malformed: len(%v) != 6 or 5", viz)
				}
			}
		}
	}
	return &b, nil
}

func (b *Bookmark) Vizes() []*BookmarkViz {
	vizes := make([]*BookmarkViz, 0)

	MaxTime := int64(8640000000000000)
	MinTime := -MaxTime

	for _, g1 := range b.Groups {
		for _, g2 := range g1 {
			for _, viz := range g2 {
				if len(viz) == 5 {
					timeRange := viz[1].([]interface{})

					startRaw := int64(timeRange[0].(float64))
					endRaw := int64(timeRange[1].(float64))

					bookmarkViz := &BookmarkViz{
						Start:       time.Unix(0, startRaw*int64(time.Millisecond)),
						End:         time.Unix(0, endRaw*int64(time.Millisecond)),
						ExtremeTime: startRaw == MinTime || endRaw == MaxTime,
						Sensors:     make([]*BookmarkSensor, 0),
					}

					for _, s := range viz[0].([]interface{}) {
						values := s.([]interface{})
						stationID := int32(values[0].(float64))
						sensorAndModule := values[1].([]interface{})
						moduleID := sensorAndModule[0].(string)
						sensorID := int64(sensorAndModule[1].(float64))

						bookmarkViz.Sensors = append(bookmarkViz.Sensors, &BookmarkSensor{
							StationID: stationID,
							ModuleID:  moduleID,
							SensorID:  sensorID,
						})
					}

					vizes = append(vizes, bookmarkViz)
				}
			}
		}
	}

	return vizes
}

func (b *Bookmark) StationIDs() []int32 {
	ids := make([]int32, 0)

	for _, g1 := range b.Groups {
		for _, g2 := range g1 {
			for _, viz := range g2 {
				if len(viz) == 6 {
					for _, id := range viz[0].([]interface{}) {
						ids = append(ids, int32(id.(float64)))
					}
				}
				if len(viz) == 5 {
					for _, s := range viz[0].([]interface{}) {
						values := s.([]interface{})
						ids = append(ids, int32(values[0].(float64)))
					}
				}
			}
		}
	}

	return ids
}
