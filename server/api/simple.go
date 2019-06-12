package api

import (
	"context"

	"github.com/goadesign/goa"

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

func (sc *SimpleController) MyCsvData(ctx *app.MyCsvDataSimpleContext) error {
	return nil
}

func (sc *SimpleController) MyFeatures(ctx *app.MyFeaturesSimpleContext) error {
	return nil
}
