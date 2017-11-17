package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
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

func createProperties(d *data.Document) map[string]string {
	p := make(map[string]string)

	p["id"] = fmt.Sprintf("%v", d.ID)
	p["timestamp"] = fmt.Sprintf("%v", d.Timestamp)

	m := map[string]string{}
	err := json.Unmarshal(d.Data, &m)
	if err == nil {
		for k, v := range m {
			// Just in case, I'm thinking of giving special values
			// an obscure prefix instead, like _ or $
			if p[k] != "" {
				k = "data" + strings.Title(k)
			}
			p[k] = v
		}
	}

	return p
}

func (c *GeoJSONController) List(ctx *app.ListGeoJSONContext) error {
	docs, err := c.options.Backend.ListDocuments(ctx, ctx.Project, ctx.Expedition)
	if err != nil {
		return err
	}

	features := make([]*app.GeoJSONFeature, 0)
	for _, d := range docs {
		c := d.Location.Coordinates()
		f := &app.GeoJSONFeature{
			Type: "Feature",
			Geometry: &app.GeoJSONGeometry{
				Type: "Point",
				Coordinates: []float64{
					c[1], c[0],
				},
			},
			Properties: createProperties(d),
		}
		features = append(features, f)
	}

	return ctx.OK(&app.GeoJSON{
		Type:     "FeatureCollection",
		Features: features,
	})
}
