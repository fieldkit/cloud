package api

import (
	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
)

type DocumentControllerOptions struct {
	Backend *backend.Backend
}

// DocumentController implements the document resource.
type DocumentController struct {
	*goa.Controller
	options DocumentControllerOptions
}

func NewDocumentController(service *goa.Service, options DocumentControllerOptions) *DocumentController {
	return &DocumentController{
		Controller: service.NewController("DocumentController"),
		options:    options,
	}
}

func (c *DocumentController) List(ctx *app.ListDocumentContext) error {
	return ctx.OK(&app.Documents{})
}
