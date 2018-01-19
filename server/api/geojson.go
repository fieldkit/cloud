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

func createProperties(d *data.Document) map[string]interface{} {
	p := make(map[string]interface{})

	timestamp := d.Timestamp.UnixNano() / int64(time.Millisecond)

	data := make(map[string]interface{})
	m := map[string]string{}
	err := json.Unmarshal(d.Data, &m)
	if err == nil {
		for k, v := range m {
			data[k] = v
		}
	}

	source := make(map[string]interface{})
	source["id"] = d.ID
	source["timestamp"] = timestamp
	source["inputId"] = d.InputID
	source["teamId"] = d.TeamID
	source["userId"] = d.UserID

	p["data"] = data
	p["id"] = d.ID
	p["timestamp"] = timestamp
	p["source"] = source

	return p
}

func (c *GeoJSONController) List(ctx *app.ListGeoJSONContext) error {
	token := backend.NewPagingTokenFromString(ctx.RequestData.Params.Get("token"))
	docs, nextToken, err := c.options.Backend.ListDocuments(ctx, ctx.Project, ctx.Expedition, token)
	if err != nil {
		return err
	}

	features := make([]*app.GeoJSONFeature, 0)
	for _, d := range docs.Documents {
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

	geoJson := &app.GeoJSON{
		Type:     "FeatureCollection",
		Features: features,
	}

	return ctx.OK(&app.PagedGeoJSON{
		NextURL: client.ListGeoJSONPath(ctx.Project, ctx.Expedition) + fmt.Sprintf("?token=%s", nextToken.String()),
		Geo:     geoJson,
		HasMore: len(features) >= backend.DefaultPageSize,
	})
}

func (c *GeoJSONController) ListByInput(ctx *app.ListByInputGeoJSONContext) error {
	token := backend.NewPagingTokenFromString(ctx.RequestData.Params.Get("token"))
	docs, nextToken, err := c.options.Backend.ListDocumentsByInput(ctx, ctx.InputID, token)
	if err != nil {
		return err
	}

	features := make([]*app.GeoJSONFeature, 0)
	for _, d := range docs.Documents {
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

	geoJson := &app.GeoJSON{
		Type:     "FeatureCollection",
		Features: features,
	}

	return ctx.OK(&app.PagedGeoJSON{
		NextURL: client.ListByInputGeoJSONPath(ctx.InputID) + fmt.Sprintf("?token=%s", nextToken.String()),
		Geo:     geoJson,
		HasMore: len(features) >= backend.DefaultPageSize,
	})
}
