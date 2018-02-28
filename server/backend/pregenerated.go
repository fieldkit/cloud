package backend

import (
	"context"
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

func (p *Pregenerator) GenerateTemporalClusters(ctx context.Context, sourceId int64) error {
	_, err := p.db.ExecContext(ctx, `
	      INSERT INTO fieldkit.sources_temporal_clusters
	      SELECT
		      input_id,
		      temporal_cluster_id,
                      NOW() AS updated_at,
		      COUNT(id) AS number_of_features,
		      MIN(timestamp) AS start_time,
		      MAX(timestamp) AS end_time,
		      ST_Centroid(ST_Collect(location)) AS centroid,
		      SQRT(ST_Area(ST_MinimumBoundingCircle(ST_Collect(ST_Transform(location, 2950)))) / pi()) AS radius
	      FROM
		      fieldkit.fk_clustered_docs($1)
	      WHERE spatial_cluster_id IS NULL
	      GROUP BY temporal_cluster_id, input_id
	      ORDER BY temporal_cluster_id
	      ON CONFLICT (source_id, cluster_id) DO UPDATE SET
                      updated_at = excluded.updated_at,
		      number_of_features = excluded.number_of_features,
		      start_time = excluded.start_time,
		      end_time = excluded.end_time,
		      centroid = excluded.centroid,
		      radius = excluded.radius
	`, sourceId)
	return err
}

func (p *Pregenerator) GenerateSpatialClusters(ctx context.Context, sourceId int64) error {
	_, err := p.db.ExecContext(ctx, `
	      INSERT INTO fieldkit.sources_spatial_clusters
	      SELECT
		      input_id,
		      spatial_cluster_id,
                      NOW() AS updated_at,
		      COUNT(id) AS number_of_features,
		      MIN(timestamp) AS start_time,
		      MAX(timestamp) AS end_time,
		      ST_Centroid(ST_Collect(location)) AS centroid,
		      SQRT(ST_Area(ST_MinimumBoundingCircle(ST_Collect(ST_Transform(location, 2950)))) / pi()) AS radius
	      FROM
		      fieldkit.fk_clustered_docs($1)
	      WHERE spatial_cluster_id IS NOT NULL
	      GROUP BY spatial_cluster_id, input_id
	      ORDER BY spatial_cluster_id
	      ON CONFLICT (source_id, cluster_id) DO UPDATE SET
                      updated_at = excluded.updated_at,
		      number_of_features = excluded.number_of_features,
		      start_time = excluded.start_time,
		      end_time = excluded.end_time,
		      centroid = excluded.centroid,
		      radius = excluded.radius
	`, sourceId)
	return err
}

func (p *Pregenerator) Pregenerate(ctx context.Context, sourceId int64) error {
	err := p.GenerateSpatialClusters(ctx, sourceId)
	if err != nil {
		return err
	}

	err = p.GenerateTemporalClusters(ctx, sourceId)
	if err != nil {
		return err
	}
	return nil
}
