package api

import (
	"encoding/json"
	"fmt"
	"strings"
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

	p["id"] = d.ID
	p["timestamp"] = d.Timestamp.UnixNano() / int64(time.Millisecond)

	m := map[string]string{}
	err := json.Unmarshal(d.Data, &m)
	if err == nil {
		for k, v := range m {
			// Just in case, I'm thinking of giving special values
			// an obscure prefix instead, like _ or $
			if p[k] != nil {
				k = "data" + strings.Title(k)
			}
			p[k] = v
		}
	}

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
