package api

import (
	"context"

	modules "github.com/fieldkit/cloud/server/api/gen/modules"

	"github.com/fieldkit/cloud/server/backend/repositories"
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
	r := repositories.NewModuleMetaRepository()

	mm, err := r.FindModuleMeta()
	if err != nil {
		return nil, err
	}

	return &modules.MetaResult{
		Object: mm,
	}, nil
}
