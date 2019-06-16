package backend

import (
	"context"
	"time"

	"github.com/conservify/sqlxcache"
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

func (p *Pregenerator) TemporalClusters(ctx context.Context, sourceID int64) (summaries []*GeometryClusterSummary, err error) {
	Logger(ctx).Sugar().Infow("Generating temporal tracks", "source_id", sourceID)

	summaries = []*GeometryClusterSummary{}
	err = p.db.SelectContext(ctx, &summaries, `
	      SELECT
                source_id, cluster_id, updated_at, number_of_features, start_time, end_time, ST_AsBinary(envelope) AS envelope, ST_AsBinary(centroid) AS centroid, radius
              FROM fieldkit.fk_update_temporal_clusters($1)`, sourceID)
	if err != nil {
		return nil, err
	}

	return
}

func (p *Pregenerator) TemporalGeometries(ctx context.Context, sourceID int64) (summaries []*TemporalGeometry, err error) {
	Logger(ctx).Sugar().Infow("Generating temporal geometries", "source_id", sourceID)

	summaries = []*TemporalGeometry{}
	err = p.db.SelectContext(ctx, &summaries, `
               SELECT
                 source_id, cluster_id, updated_at, ST_AsBinary(geometry) AS geometry
               FROM fieldkit.fk_update_temporal_geometries($1)`, sourceID)
	if err != nil {
		return nil, err
	}

	return
}

func (p *Pregenerator) SpatialClusters(ctx context.Context, sourceID int64) (summaries []*GeometryClusterSummary, err error) {
	Logger(ctx).Sugar().Infow("Generating spatial clusters", "source_id", sourceID)

	summaries = []*GeometryClusterSummary{}
	err = p.db.SelectContext(ctx, &summaries, `
               SELECT
                 source_id, cluster_id, updated_at, number_of_features, start_time, end_time, ST_AsBinary(envelope) AS envelope, ST_AsBinary(centroid) AS centroid, radius
               FROM fieldkit.fk_update_spatial_clusters($1)`, sourceID)
	if err != nil {
		return nil, err
	}

	return
}

func (p *Pregenerator) Summaries(ctx context.Context, sourceID int64) (summaries []*FeatureSummary, err error) {
	Logger(ctx).Sugar().Infow("Generating summaries", "source_id", sourceID)

	summaries = []*FeatureSummary{}
	err = p.db.SelectContext(ctx, &summaries, `
               SELECT
                 source_id, updated_at, number_of_features, start_time, end_time, ST_AsBinary(envelope) AS envelope, ST_AsBinary(centroid) AS centroid, radius
               FROM fieldkit.fk_update_source_summary($1)`, sourceID)
	if err != nil {
		return nil, err
	}

	return
}

func (p *Pregenerator) Pregenerate(ctx context.Context, sourceID int64) (*Pregenerated, error) {
	pregenerated := &Pregenerated{}

	started := time.Now()

	if err := p.db.WithNewTransaction(ctx, func(txCtx context.Context) error {
		summaries, err := p.Summaries(txCtx, sourceID)
		if err != nil {
			return err
		}

		spatial, err := p.SpatialClusters(txCtx, sourceID)
		if err != nil {
			return err
		}

		temporal, err := p.TemporalClusters(txCtx, sourceID)
		if err != nil {
			return err
		}

		temporalGeometries, err := p.TemporalGeometries(txCtx, sourceID)
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

	elapsed := time.Since(started)

	log := Logger(ctx).Sugar()

	log.Infow("Pregeneration done", "source_id", sourceID, "elapsed", elapsed)

	return pregenerated, nil
}

type Pregenerated struct {
	Summaries          []*FeatureSummary
	Spatial            []*GeometryClusterSummary
	Temporal           []*GeometryClusterSummary
	TemporalGeometries []*TemporalGeometry
}
