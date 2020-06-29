package main

import (
	"context"

	"github.com/fieldkit/cloud/server/data"
)

type RecordVisitor interface {
	VisitMeta(ctx context.Context, meta *data.MetaRecord) error
	VisitData(ctx context.Context, meta *data.MetaRecord, data *data.DataRecord) error
	VisitEnd(ctx context.Context) error
}

type noopVisitor struct {
}

func NewNoopVisitor() RecordVisitor {
	return &noopVisitor{}
}

func (v *noopVisitor) VisitMeta(ctx context.Context, meta *data.MetaRecord) error {
	return nil
}

func (v *noopVisitor) VisitData(ctx context.Context, meta *data.MetaRecord, data *data.DataRecord) error {
	return nil
}

func (v *noopVisitor) VisitEnd(ctx context.Context) error {
	return nil
}
