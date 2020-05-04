package api

import (
	"context"
	"net/http"

	_ "github.com/fieldkit/cloud/server/logging"
)

func CreateApi(ctx context.Context, controllerOptions *ControllerOptions) (http.Handler, error) {
	goaV3Handler, err := CreateGoaV3Handler(ctx, controllerOptions)
	if err != nil {
		return nil, err
	}

	goaV2Handler, err := CreateGoaV2Handler(ctx, controllerOptions, goaV3Handler)
	if err != nil {
		return nil, err
	}

	return goaV2Handler, nil
}
