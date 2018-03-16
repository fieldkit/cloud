package api

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
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

func createProperties(d *data.AnalysedRecord) map[string]interface{} {
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

func MakeGeoJSON(docs *data.AnalysedRecordsPage) *app.GeoJSON {
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

func (c *GeoJSONController) ListBySource(ctx *app.ListBySourceGeoJSONContext) error {
	token := backend.NewPagingTokenFromString(ctx.RequestData.Params.Get("token"))
	descending := false
	if ctx.Descending != nil {
		descending = *ctx.Descending
	}
	records, nextToken, err := c.options.Backend.ListRecordsBySource(ctx, ctx.SourceID, descending, false, token)
	if err != nil {
		return err
	}

	geoJson := MakeGeoJSON(records)

	return ctx.OK(&app.PagedGeoJSON{
		NextURL: client.ListBySourceGeoJSONPath(ctx.SourceID) + fmt.Sprintf("?token=%s&descending=%v", nextToken.String(), descending),
		Geo:     geoJson,
		HasMore: len(geoJson.Features) >= backend.DefaultPageSize,
	})
}

func (c *GeoJSONController) ListByID(ctx *app.ListByIDGeoJSONContext) error {
	records, err := c.options.Backend.ListRecordsByID(ctx, ctx.FeatureID)
	if err != nil {
		return err
	}

	geoJson := MakeGeoJSON(records)

	return ctx.OK(&app.PagedGeoJSON{
		Geo:     geoJson,
		HasMore: len(geoJson.Features) >= backend.DefaultPageSize,
	})
}

type BoundingBox struct {
	NorthEast *data.Location
	SouthWest *data.Location
}

func ExtractLocationFromQueryString(key string, ctx *app.GeographicalQueryGeoJSONContext) (l *data.Location, err error) {
	strs := strings.Split(ctx.RequestData.Params.Get(key), ",")
	if len(strs) != 2 {
		return nil, fmt.Errorf("Expected two values for data.Location")
	}

	lng, err := strconv.ParseFloat(strs[0], 64)
	if err != nil {
		return nil, err
	}
	lat, err := strconv.ParseFloat(strs[0], 64)
	if err != nil {
		return nil, err
	}

	return data.NewLocation([]float64{lng, lat}), nil
}

func ExtractBoundingBoxFromQueryString(ctx *app.GeographicalQueryGeoJSONContext) (bb *BoundingBox, err error) {
	ne, err := ExtractLocationFromQueryString("ne", ctx)
	if err != nil {
		return nil, err
	}

	sw, err := ExtractLocationFromQueryString("sw", ctx)
	if err != nil {
		return nil, err
	}

	bb = &BoundingBox{
		NorthEast: ne,
		SouthWest: sw,
	}
	return
}

func (c *GeoJSONController) GeographicalQuery(ctx *app.GeographicalQueryGeoJSONContext) error {
	bb, err := ExtractBoundingBoxFromQueryString(ctx)
	if err != nil {
		return fmt.Errorf("Error parsing bounding box: %v", err)
	}
	log.Printf("Querying over %+v", bb)

	features := make([]*app.GeoJSONFeature, 0)
	geoJson := &app.GeoJSON{
		Type:     "FeatureCollection",
		Features: features,
	}
	return ctx.OK(&app.PagedGeoJSON{
		Geo:     geoJson,
		HasMore: false,
	})
}
