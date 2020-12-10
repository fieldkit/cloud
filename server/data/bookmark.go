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
				if len(viz) != 6 {
					return nil, fmt.Errorf("viz bookmark malformed: len(%v) != 6", viz)
				}
				if stationPart, ok := viz[0].([]interface{}); !ok {
					return nil, fmt.Errorf("viz bookmark stations malformed: %v", viz[0])
				} else {
					for _, id := range stationPart {
						if _, ok := id.(float64); !ok {
							return nil, fmt.Errorf("viz bookmark stations malformed: %v", id)
						}
					}
				}
			}
		}
	}
	return &b, nil
}

func (b *Bookmark) StationIDs() []int32 {
	ids := make([]int32, 0)

	for _, gg := range b.Groups {
		for _, g := range gg {
			for _, viz := range g {
				for _, id := range viz[0].([]interface{}) {
					ids = append(ids, int32(id.(float64)))
				}
			}
		}
	}

	return ids
}
