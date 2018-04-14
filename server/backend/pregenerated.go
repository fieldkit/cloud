package backend

import (
	"context"
	"log"

	"github.com/conservify/sqlxcache"
	"github.com/fieldkit/cloud/server/data"
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

func (p *Pregenerator) TemporalClusters(ctx context.Context, sourceId int64) error {
	log.Printf("Generating temporal tracks (%v)...", sourceId)

	summaries := []*GeometryClusterSummary{}
	err := p.db.SelectContext(ctx, &summaries, `SELECT source_id, cluster_id, updated_at, number_of_features, start_time, end_time, ST_AsBinary(envelope) AS envelope, ST_AsBinary(centroid) AS centroid, radius FROM fieldkit.fk_update_temporal_clusters($1)`, sourceId)
	if err != nil {
		return err
	}

	log.Printf("%v", summaries)

	return err
}

type ClusteredRow struct {
	ClusterId int64
	Location  data.Location
}

func (p *Pregenerator) TemporalGeometries(ctx context.Context, sourceId int64) error {
	log.Printf("Generating temporal geometries (%v)...", sourceId)

	summaries := []*TemporalGeometry{}
	err := p.db.SelectContext(ctx, &summaries, "SELECT source_id, cluster_id, updated_at, ST_AsBinary(geometry) AS geometry FROM fieldkit.fk_update_temporal_geometries($1)", sourceId)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pregenerator) SpatialClusters(ctx context.Context, sourceId int64) error {
	log.Printf("Generating spatial clusters (%v)...", sourceId)

	summaries := []*GeometryClusterSummary{}
	err := p.db.SelectContext(ctx, &summaries, `SELECT source_id, cluster_id, updated_at, number_of_features, start_time, end_time, ST_AsBinary(envelope) AS envelope, ST_AsBinary(centroid) AS centroid, radius FROM fieldkit.fk_update_spatial_clusters($1)`, sourceId)
	if err != nil {
		return err
	}

	log.Printf("%v", summaries)

	return err
}

func (p *Pregenerator) Summaries(ctx context.Context, sourceId int64) error {
	log.Printf("Generating summaries (%v)...", sourceId)

	summaries := []*FeatureSummary{}
	err := p.db.SelectContext(ctx, &summaries, "SELECT source_id, updated_at, number_of_features, last_feature_id, start_time, end_time, ST_AsBinary(envelope) AS envelope, ST_AsBinary(centroid) AS centroid, radius FROM fieldkit.fk_update_source_summary($1)", sourceId)
	if err != nil {
		return err
	}

	log.Printf("%v", summaries)

	return err
}

func (p *Pregenerator) Pregenerate(ctx context.Context, sourceId int64) error {
	return p.db.WithNewTransaction(ctx, func(txCtx context.Context) error {
		err := p.Summaries(txCtx, sourceId)
		if err != nil {
			return err
		}

		err = p.SpatialClusters(txCtx, sourceId)
		if err != nil {
			return err
		}

		err = p.TemporalClusters(txCtx, sourceId)
		if err != nil {
			return err
		}

		err = p.TemporalGeometries(txCtx, sourceId)
		if err != nil {
			return err
		}

		return nil
	})
}
