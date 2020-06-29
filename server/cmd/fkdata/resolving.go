package main

import (
	"context"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

func NewResolvingVisitor() *ResolvingVisitor {
	return &ResolvingVisitor{
		metaFactory: repositories.NewMetaFactory(),
	}
}

type ResolvingVisitor struct {
	metaFactory *repositories.MetaFactory
}

func (v *ResolvingVisitor) VisitMeta(ctx context.Context, meta *data.MetaRecord) error {
	_, err := v.metaFactory.Add(ctx, meta, true)
	if err != nil {
		return err
	}

	return nil
}

func (v *ResolvingVisitor) VisitData(ctx context.Context, meta *data.MetaRecord, data *data.DataRecord) error {
	row, err := v.metaFactory.Resolve(ctx, data, false, true)
	if err != nil {
		return err
	}

	_ = row

	return nil
}
