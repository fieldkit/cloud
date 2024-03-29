// Code generated by goa v3.2.4, DO NOT EDIT.
//
// sensor client
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package sensor

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "sensor" service client.
type Client struct {
	MetaEndpoint        goa.Endpoint
	StationMetaEndpoint goa.Endpoint
	SensorMetaEndpoint  goa.Endpoint
	DataEndpoint        goa.Endpoint
	TailEndpoint        goa.Endpoint
	RecentlyEndpoint    goa.Endpoint
	BookmarkEndpoint    goa.Endpoint
	ResolveEndpoint     goa.Endpoint
}

// NewClient initializes a "sensor" service client given the endpoints.
func NewClient(meta, stationMeta, sensorMeta, data, tail, recently, bookmark, resolve goa.Endpoint) *Client {
	return &Client{
		MetaEndpoint:        meta,
		StationMetaEndpoint: stationMeta,
		SensorMetaEndpoint:  sensorMeta,
		DataEndpoint:        data,
		TailEndpoint:        tail,
		RecentlyEndpoint:    recently,
		BookmarkEndpoint:    bookmark,
		ResolveEndpoint:     resolve,
	}
}

// Meta calls the "meta" endpoint of the "sensor" service.
func (c *Client) Meta(ctx context.Context) (res *MetaResult, err error) {
	var ires interface{}
	ires, err = c.MetaEndpoint(ctx, nil)
	if err != nil {
		return
	}
	return ires.(*MetaResult), nil
}

// StationMeta calls the "station meta" endpoint of the "sensor" service.
func (c *Client) StationMeta(ctx context.Context, p *StationMetaPayload) (res *StationMetaResult, err error) {
	var ires interface{}
	ires, err = c.StationMetaEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*StationMetaResult), nil
}

// SensorMeta calls the "sensor meta" endpoint of the "sensor" service.
func (c *Client) SensorMeta(ctx context.Context) (res *SensorMetaResult, err error) {
	var ires interface{}
	ires, err = c.SensorMetaEndpoint(ctx, nil)
	if err != nil {
		return
	}
	return ires.(*SensorMetaResult), nil
}

// Data calls the "data" endpoint of the "sensor" service.
func (c *Client) Data(ctx context.Context, p *DataPayload) (res *DataResult, err error) {
	var ires interface{}
	ires, err = c.DataEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*DataResult), nil
}

// Tail calls the "tail" endpoint of the "sensor" service.
func (c *Client) Tail(ctx context.Context, p *TailPayload) (res *TailResult, err error) {
	var ires interface{}
	ires, err = c.TailEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*TailResult), nil
}

// Recently calls the "recently" endpoint of the "sensor" service.
func (c *Client) Recently(ctx context.Context, p *RecentlyPayload) (res *RecentlyResult, err error) {
	var ires interface{}
	ires, err = c.RecentlyEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*RecentlyResult), nil
}

// Bookmark calls the "bookmark" endpoint of the "sensor" service.
func (c *Client) Bookmark(ctx context.Context, p *BookmarkPayload) (res *SavedBookmark, err error) {
	var ires interface{}
	ires, err = c.BookmarkEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*SavedBookmark), nil
}

// Resolve calls the "resolve" endpoint of the "sensor" service.
func (c *Client) Resolve(ctx context.Context, p *ResolvePayload) (res *SavedBookmark, err error) {
	var ires interface{}
	ires, err = c.ResolveEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*SavedBookmark), nil
}
