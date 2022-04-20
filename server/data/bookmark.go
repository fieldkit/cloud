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

func (b *Bookmark) Vizes() ([]*BookmarkViz, error) {
	vizes := make([]*BookmarkViz, 0)

	MaxTime := int64(8640000000000000)
	MinTime := -MaxTime

	for _, g1 := range b.Groups {
		for _, g2 := range g1 {
			for _, viz := range g2 {
				if len(viz) == 5 {
					timeRange, ok := viz[1].([]interface{})
					if !ok {
						return nil, fmt.Errorf("bookmark:malformed")
					}

					startRaw, ok := timeRange[0].(float64)
					if !ok {
						return nil, fmt.Errorf("bookmark:malformed")
					}

					endRaw, ok := timeRange[1].(float64)
					if !ok {
						return nil, fmt.Errorf("bookmark:malformed")
					}

					bookmarkViz := &BookmarkViz{
						Start:       time.Unix(0, int64(startRaw)*int64(time.Millisecond)),
						End:         time.Unix(0, int64(endRaw)*int64(time.Millisecond)),
						ExtremeTime: int64(startRaw) == MinTime || int64(endRaw) == MaxTime,
						Sensors:     make([]*BookmarkSensor, 0),
					}

					for _, s := range viz[0].([]interface{}) {
						values, ok := s.([]interface{})
						if !ok {
							return nil, fmt.Errorf("bookmark:malformed")
						}
						stationID, ok := values[0].(float64)
						if !ok {
							return nil, fmt.Errorf("bookmark:malformed")
						}

						sensorAndModule, ok := values[1].([]interface{})
						if !ok || len(sensorAndModule) != 2 {
							return nil, fmt.Errorf("bookmark:malformed")
						}

						if moduleID, ok := sensorAndModule[0].(string); !ok {
							return nil, fmt.Errorf("bookmark:malformed")
						} else {
							if sensorID, ok := sensorAndModule[1].(float64); !ok {
								return nil, fmt.Errorf("bookmark:malformed")
							} else {
								bookmarkViz.Sensors = append(bookmarkViz.Sensors, &BookmarkSensor{
									StationID: int32(stationID),
									SensorID:  int64(sensorID),
									ModuleID:  moduleID,
								})
							}
						}
					}

					vizes = append(vizes, bookmarkViz)
				}
			}
		}
	}

	return vizes, nil
}

func (b *Bookmark) StationIDs() ([]int32, error) {
	ids := make([]int32, 0)

	for _, g1 := range b.Groups {
		for _, g2 := range g1 {
			for _, viz := range g2 {
				if len(viz) == 6 {
					stationIDs, ok := viz[0].([]interface{})
					if !ok {
						return nil, fmt.Errorf("bookmark:malformed")
					}
					for _, id := range stationIDs {
						ids = append(ids, int32(id.(float64)))
					}
				}
				if len(viz) == 5 {
					stationIDs, ok := viz[0].([]interface{})
					if !ok {
						return nil, fmt.Errorf("bookmark:malformed")
					}
					for _, s := range stationIDs {
						values := s.([]interface{})
						ids = append(ids, int32(values[0].(float64)))
					}
				}
			}
		}
	}

	return ids, nil
}
