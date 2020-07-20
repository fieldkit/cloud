package api

import (
	"context"
	"net/http"
)

func CreateApi(ctx context.Context, controllerOptions *ControllerOptions) (http.Handler, error) {
	goaV3Handler, err := CreateGoaV3Handler(ctx, controllerOptions)
	if err != nil {
		return nil, err
	}

	return goaV3Handler, nil
}
