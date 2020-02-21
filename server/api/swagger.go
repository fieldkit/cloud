package api

import (
	"github.com/goadesign/goa"
)

type SwaggerController struct {
	*goa.Controller
}

func NewSwaggerController(service *goa.Service) *SwaggerController {
	return &SwaggerController{Controller: service.NewController("SwaggerController")}
}
