package gonaturalist

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type SimplePlace struct {
	Id            int64     `json:"id"`
	Name          string    `json:"name"`
	DisplayName   string    `json:"display_name"`
	Code          string    `json:"code"`
	PlaceType     int32     `json:"place_type"`
	PlaceTypeName string    `json:"place_type_name"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Slug          string    `json:"slug"`
	Latitude      float64   `json:"latitude,string"`
	Longitude     float64   `json:"longitude,string"`
	SwLat         float64   `json:"swlat,string"`
	SwLon         float64   `json:"swlng,string"`
	NeLat         float64   `json:"nelat,string"`
	NeLon         float64   `json:"nelng,string"`
}

func (p *SimplePlace) Rectangle() (r Rectangle, err error) {
	/*
		swLon, err := strconv.ParseFloat(p.SwLon, 64)
		if err != nil {
			return
		}
		swLat, err := strconv.ParseFloat(p.SwLat, 64)
		if err != nil {
			return
		}
		neLon, err := strconv.ParseFloat(p.NeLon, 64)
		if err != nil {
			return
		}
		neLat, err := strconv.ParseFloat(p.NeLat, 64)
		if err != nil {
			return
		}
	*/

	r = Rectangle{
		Southwest: Location{
			Longitude: p.SwLon,
			Latitude:  p.SwLat,
		},
		Northeast: Location{
			Longitude: p.NeLon,
			Latitude:  p.NeLat,
		},
	}

	return
}

type PlacesPage struct {
	Paging *PageHeaders
	Places []*SimplePlace
}

type GetPlacesOpt struct {
	Page      *int
	Longitude *float64
	Latitude  *float64
}

func (c *Client) GetPlaces(opt *GetPlacesOpt) (*PlacesPage, error) {
	var result []*SimplePlace

	u := c.buildUrl("/places.json")
	if opt != nil {
		v := url.Values{}
		if opt.Page != nil {
			v.Set("page", strconv.Itoa(*opt.Page))
		}
		if opt.Longitude != nil {
			v.Set("longitude", fmt.Sprintf("%v", *opt.Longitude))
		}
		if opt.Latitude != nil {
			v.Set("latitude", fmt.Sprintf("%v", *opt.Latitude))
		}
		if params := v.Encode(); params != "" {
			u += "?" + params
		}
	}
	p, err := c.get(u, &result)
	if err != nil {
		return nil, err
	}

	return &PlacesPage{
		Places: result,
		Paging: p,
	}, nil
}
