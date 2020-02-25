package api

import (
	"context"

	modules "github.com/fieldkit/cloud/server/api/gen/modules"
)

type ModulesService struct {
	options *ControllerOptions
}

func NewModulesService(ctx context.Context, options *ControllerOptions) *ModulesService {
	return &ModulesService{
		options: options,
	}
}

func (ms *ModulesService) Meta(ctx context.Context) (*modules.MetaResult, error) {
	v := make(map[string]interface{})

	v["value"] = 100

	return &modules.MetaResult{
		Object: v,
	}, nil
}
