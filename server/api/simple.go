package api

import (
	"context"
	"fmt"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware/security/jwt"

	"github.com/conservify/sqlxcache"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
)

type SimpleControllerOptions struct {
	Config        *ApiConfiguration
	Session       *session.Session
	Database      *sqlxcache.DB
	Backend       *backend.Backend
	ConcatWorkers *backend.ConcatenationWorkers
}

type SimpleController struct {
	options SimpleControllerOptions
	*goa.Controller
}

func NewSimpleController(ctx context.Context, service *goa.Service, options SimpleControllerOptions) *SimpleController {
	return &SimpleController{
		options:    options,
		Controller: service.NewController("SimpleController"),
	}
}

func (sc *SimpleController) MyFeatures(ctx *app.MyFeaturesSimpleContext) error {
	log := Logger(ctx).Sugar()

	userID, err := getCurrentUserId(ctx)
	if err != nil {
		return err
	}

	log.Infow("my features", "user_id", userID)

	bb, err := ExtractBoundingBoxFromQueryString(ctx.RequestData)
	if err != nil {
		return fmt.Errorf("Error parsing bounding box: %v", err)
	}
	log.Infow("my features", "bb", bb)

	mapFeatures, err := sc.options.Backend.QueryMapFeatures(ctx, bb, userID)
	if err != nil {
		return err
	}

	features := make([]*app.GeoJSONFeature, 0)
	geoJson := &app.GeoJSON{
		Type:     "FeatureCollection",
		Features: features,
	}
	return ctx.OK(&app.MapFeatures{
		Geometries: ClusterGeometriesType(mapFeatures.TemporalGeometries),
		Spatial:    ClusterSummariesType(mapFeatures.SpatialClusters),
		Temporal:   ClusterSummariesType(mapFeatures.TemporalClusters),
		GeoJSON: &app.PagedGeoJSON{
			Geo: geoJson,
		},
	})
}

func (sc *SimpleController) MyCsvData(ctx *app.MyCsvDataSimpleContext) error {
	return nil
}

func getCurrentUserId(ctx context.Context) (id int64, err error) {
	token := jwt.ContextJWT(ctx)
	if token == nil {
		return 0, fmt.Errorf("JWT token is missing from context")
	}

	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return 0, fmt.Errorf("JWT claims error")
	}

	id = int64(claims["sub"].(float64))

	return
}
