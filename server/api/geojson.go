package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/api/client"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"
)

type GeoJSONControllerOptions struct {
	Backend *backend.Backend
}

type GeoJSONController struct {
	*goa.Controller
	options GeoJSONControllerOptions
}

func NewGeoJSONController(service *goa.Service, options GeoJSONControllerOptions) *GeoJSONController {
	return &GeoJSONController{
		Controller: service.NewController("GeoJSONController"),
		options:    options,
	}
}

func createProperties(d *data.Record) map[string]interface{} {
	p := make(map[string]interface{})

	timestamp := d.Timestamp.UnixNano() / int64(time.Millisecond)

	data := make(map[string]interface{})
	m := make(map[string]interface{})
	err := json.Unmarshal(d.Data, &m)
	if err == nil {
		for k, v := range m {
			data[k] = v
		}
	}

	source := make(map[string]interface{})
	source["id"] = d.ID
	source["timestamp"] = timestamp
	source["sourceId"] = d.SourceID
	source["teamId"] = d.TeamID
	source["userId"] = d.UserID

	p["data"] = data
	p["id"] = d.ID
	p["timestamp"] = timestamp
	p["source"] = source

	return p
}

func MakeGeoJSON(docs *data.RecordsPage) *app.GeoJSON {
	features := make([]*app.GeoJSONFeature, 0)
	for _, d := range docs.Records {
		c := d.Location.Coordinates()
		f := &app.GeoJSONFeature{
			Type: "Feature",
			Geometry: &app.GeoJSONGeometry{
				Type: "Point",
				Coordinates: []float64{
					c[0], c[1],
				},
			},
			Properties: createProperties(d),
		}
		features = append(features, f)
	}

	return &app.GeoJSON{
		Type:     "FeatureCollection",
		Features: features,
	}
}

func (c *GeoJSONController) List(ctx *app.ListGeoJSONContext) error {
	token := backend.NewPagingTokenFromString(ctx.RequestData.Params.Get("token"))
	docs, nextToken, err := c.options.Backend.ListRecords(ctx, ctx.Project, ctx.Expedition, token)
	if err != nil {
		return err
	}

	geoJson := MakeGeoJSON(docs)

	return ctx.OK(&app.PagedGeoJSON{
		NextURL: client.ListGeoJSONPath(ctx.Project, ctx.Expedition) + fmt.Sprintf("?token=%s", nextToken.String()),
		Geo:     geoJson,
		HasMore: len(geoJson.Features) >= backend.DefaultPageSize,
	})
}

func (c *GeoJSONController) ListBySource(ctx *app.ListBySourceGeoJSONContext) error {
	token := backend.NewPagingTokenFromString(ctx.RequestData.Params.Get("token"))
	descending := false
	if ctx.Descending != nil {
		descending = *ctx.Descending
	}
	docs, nextToken, err := c.options.Backend.ListRecordsBySource(ctx, ctx.SourceID, descending, token)
	if err != nil {
		return err
	}

	geoJson := MakeGeoJSON(docs)

	return ctx.OK(&app.PagedGeoJSON{
		NextURL: client.ListBySourceGeoJSONPath(ctx.SourceID) + fmt.Sprintf("?token=%s&descending=%v", nextToken.String(), descending),
		Geo:     geoJson,
		HasMore: len(geoJson.Features) >= backend.DefaultPageSize,
	})
}

func (c *GeoJSONController) ListByID(ctx *app.ListByIDGeoJSONContext) error {
	docs, err := c.options.Backend.ListRecordsByID(ctx, ctx.FeatureID)
	if err != nil {
		return err
	}

	geoJson := MakeGeoJSON(docs)

	return ctx.OK(&app.PagedGeoJSON{
		Geo:     geoJson,
		HasMore: len(geoJson.Features) >= backend.DefaultPageSize,
	})
}
