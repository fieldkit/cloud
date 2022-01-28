package data

import (
	"encoding/json"
	"fmt"
)

type VizBookmark = []interface{} // [Stations, Sensors, [number, number], [[number, number], [number, number]] | [], ChartType, FastTime] | [];
type GroupBookmark = [][]VizBookmark

// {"v":1,"g":[ [[ [[159],[2],[-8640000000000000,8640000000000000],[],0,0] ]]],"s":[]}
type Bookmark struct {
	Version  int32           `json:"v"`
	Groups   []GroupBookmark `json:"g"`
	Stations []int32         `json:"s"`
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
