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
	JWTHMACKey    []byte
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

	userID, err := getCurrentUserIdFromContext(ctx)
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

func (sc *SimpleController) MySimpleSummary(ctx *app.MySimpleSummarySimpleContext) error {
	/*
		log := Logger(ctx).Sugar()

		userID, err := getCurrentUserIdFromContext(ctx)
		if err != nil {
			return err
		}

		log.Infow("my data (csv)", "user_id", userID)

		devices, err := sc.options.Backend.ListAllDeviceSourcesByUser(ctx, userID)
		if err != nil {
			return err
		}

		log.Infow("my data (csv)", "number_of_devices", len(devices))

		if len(devices) != 1 {
			return ctx.OK(&app.MyDataUrls{})
		}
	*/

	token := jwt.ContextJWT(ctx)
	if token == nil {
		return fmt.Errorf("JWT token is missing from context")
	}

	return ctx.OK(&app.MyDataUrls{
		Csv: sc.options.Config.MakeApiUrl("/my/simple/download?token=" + token.Raw),
	})
}

func (sc *SimpleController) Download(ctx *app.DownloadSimpleContext) error {
	log := Logger(ctx).Sugar()

	fr, err := backend.NewFileRepository(sc.options.Session, "fk-streams")
	if err != nil {
		return err
	}

	token, err := jwtgo.Parse(ctx.Token, func(token *jwtgo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return sc.options.JWTHMACKey, nil
	})
	if err != nil {
		return err
	}

	userID, err := getCurrentUserIdFromToken(token)
	if err != nil {
		return err
	}

	log.Infow("download", "user_id", userID)

	devices, err := sc.options.Backend.ListAllDeviceSourcesByUser(ctx, userID)
	if err != nil {
		return err
	}

	log.Infow("download", "number_of_devices", len(devices))

	if len(devices) != 1 {
		return fmt.Errorf("Too many devices linked to user.")
	}

	device := devices[0]

	urls := DeviceSummaryUrls(sc.options.Config, device.Key)

	data, err := fr.Info(ctx, urls.Data.ID)
	if err != nil {
		return err
	}

	log.Infow("download", "data", data)

	log.Infow("download", "files", urls)

	iterator, err := backend.LookupFile(ctx, sc.options.Session, sc.options.Database, urls.Data.ID)
	if err != nil {
		return err
	}

	exporter := backend.NewSimpleCsvExporter(ctx.ResponseData)

	return backend.ExportAllFiles(ctx, ctx.ResponseData, true, iterator, exporter)
}

func getCurrentUserIdFromToken(token *jwtgo.Token) (id int64, err error) {
	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return 0, fmt.Errorf("JWT claims error")
	}

	id = int64(claims["sub"].(float64))

	return
}

func getCurrentUserIdFromContext(ctx context.Context) (id int64, err error) {
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
