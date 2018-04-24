package backend

import (
	"context"

	"github.com/Conservify/sqlxcache"
)

type Pregenerator struct {
	db      *sqlxcache.DB
	backend *Backend
}

func NewPregenerator(backend *Backend) *Pregenerator {
	return &Pregenerator{
		db:      backend.db,
		backend: backend,
	}
}

func (p *Pregenerator) TemporalClusters(ctx context.Context, sourceId int64) (summaries []*GeometryClusterSummary, err error) {
	Logger(ctx).Sugar().Infow("Generating temporal tracks", "sourceId", sourceId)

	summaries = []*GeometryClusterSummary{}
	err = p.db.SelectContext(ctx, &summaries, `
	      SELECT
                source_id, cluster_id, updated_at, number_of_features, start_time, end_time, ST_AsBinary(envelope) AS envelope, ST_AsBinary(centroid) AS centroid, radius
              FROM fieldkit.fk_update_temporal_clusters($1)`, sourceId)
	if err != nil {
		return nil, err
	}

	return
}

func (p *Pregenerator) TemporalGeometries(ctx context.Context, sourceId int64) (summaries []*TemporalGeometry, err error) {
	Logger(ctx).Sugar().Infow("Generating temporal geometries", "sourceId", sourceId)

	summaries = []*TemporalGeometry{}
	err = p.db.SelectContext(ctx, &summaries, `
               SELECT
                 source_id, cluster_id, updated_at, ST_AsBinary(geometry) AS geometry
               FROM fieldkit.fk_update_temporal_geometries($1)`, sourceId)
	if err != nil {
		return nil, err
	}

	return
}

func (p *Pregenerator) SpatialClusters(ctx context.Context, sourceId int64) (summaries []*GeometryClusterSummary, err error) {
	Logger(ctx).Sugar().Infow("Generating spatial clusters", "sourceId", sourceId)

	summaries = []*GeometryClusterSummary{}
	err = p.db.SelectContext(ctx, &summaries, `
               SELECT
                 source_id, cluster_id, updated_at, number_of_features, start_time, end_time, ST_AsBinary(envelope) AS envelope, ST_AsBinary(centroid) AS centroid, radius
               FROM fieldkit.fk_update_spatial_clusters($1)`, sourceId)
	if err != nil {
		return nil, err
	}

	return
}

func (p *Pregenerator) Summaries(ctx context.Context, sourceId int64) (summaries []*FeatureSummary, err error) {
	Logger(ctx).Sugar().Infow("Generating summaries", "sourceId", sourceId)

	summaries = []*FeatureSummary{}
	err = p.db.SelectContext(ctx, &summaries, `
               SELECT
                 source_id, updated_at, number_of_features, last_feature_id, start_time, end_time, ST_AsBinary(envelope) AS envelope, ST_AsBinary(centroid) AS centroid, radius
               FROM fieldkit.fk_update_source_summary($1)`, sourceId)
	if err != nil {
		return nil, err
	}

	return
}

func (p *Pregenerator) Pregenerate(ctx context.Context, sourceId int64) (*Pregenerated, error) {
	pregenerated := &Pregenerated{}

	if err := p.db.WithNewTransaction(ctx, func(txCtx context.Context) error {
		summaries, err := p.Summaries(txCtx, sourceId)
		if err != nil {
			return err
		}

		spatial, err := p.SpatialClusters(txCtx, sourceId)
		if err != nil {
			return err
		}

		temporal, err := p.TemporalClusters(txCtx, sourceId)
		if err != nil {
			return err
		}

		temporalGeometries, err := p.TemporalGeometries(txCtx, sourceId)
		if err != nil {
			return err
		}

		pregenerated.Summaries = summaries
		pregenerated.Spatial = spatial
		pregenerated.Temporal = temporal
		pregenerated.TemporalGeometries = temporalGeometries

		return nil
	}); err != nil {
		return nil, err
	}

	return pregenerated, nil
}

type Pregenerated struct {
	Summaries          []*FeatureSummary
	Spatial            []*GeometryClusterSummary
	Temporal           []*GeometryClusterSummary
	TemporalGeometries []*TemporalGeometry
}
