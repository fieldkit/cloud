package backend

import (
	"context"
	"database/sql/driver"
	"github.com/conservify/sqlxcache"
	"github.com/fieldkit/cloud/server/data"
	"github.com/paulmach/go.geo"
	"log"
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
	_, err := p.db.ExecContext(ctx, "DELETE FROM fieldkit.sources_temporal_clusters WHERE (source_id = $1)", sourceId)
	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, `
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

type ClusteredRow struct {
	ClusterId int64
	Location  data.Location
}

func (p *Pregenerator) GenerateTemporalGeometries(ctx context.Context, sourceId int64) error {
	log.Printf("Generating temporal cluster geometries...")

	_, err := p.db.ExecContext(ctx, "DELETE FROM fieldkit.sources_temporal_geometries WHERE (source_id = $1)", sourceId)
	if err != nil {
		return err
	}

	locations := []*ClusteredRow{}
	if err := p.db.SelectContext(ctx, &locations, "SELECT temporal_cluster_id AS clusterId, ST_AsBinary(location) AS location FROM fieldkit.fk_clustered_docs($1) WHERE spatial_cluster_id IS NULL ORDER BY timestamp", sourceId); err != nil {
		return err
	}

	geometries := make(map[int64][]geo.Point)
	for _, row := range locations {
		if geometries[row.ClusterId] == nil {
			geometries[row.ClusterId] = make([]geo.Point, 0)
		}
		c := row.Location.Coordinates()
		p := geo.NewPoint(c[0], c[1])
		geometries[row.ClusterId] = append(geometries[row.ClusterId], *p)
	}

	for clusterId, geometry := range geometries {
		if len(geometry) >= 2 {
			path := geo.NewPath()
			path.SetPoints(geometry)
			temporal := &TemporalPath{
				path: path,
			}
			_, err := p.db.ExecContext(ctx, `
			    INSERT INTO fieldkit.sources_temporal_geometries (source_id, cluster_id, updated_at, geometry) VALUES ($1, $2, NOW(), ST_SetSRID(ST_GeomFromText($3), 4326))
			    ON CONFLICT (source_id, cluster_id) DO UPDATE SET updated_at = excluded.updated_at, geometry = excluded.geometry
			`, sourceId, clusterId, temporal)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Pregenerator) GenerateSpatialClusters(ctx context.Context, sourceId int64) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM fieldkit.sources_spatial_clusters WHERE (source_id = $1)", sourceId)
	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, `
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

	err = p.GenerateTemporalGeometries(ctx, sourceId)
	if err != nil {
		return err
	}

	return nil
}

type TemporalPath struct {
	path *geo.Path
}

func (l *TemporalPath) Coordinates() [][]float64 {
	if l.path == nil {
		return make([][]float64, 0)
	}
	ps := l.path.Points()
	c := make([][]float64, len(ps))
	for i, p := range ps {
		c[i] = []float64{p.Lng(), p.Lat()}
	}
	return c
}

func (l *TemporalPath) Scan(data interface{}) error {
	path := &geo.Path{}
	if err := path.Scan(data); err != nil {
		return err
	}

	l.path = path
	return nil
}

func (l *TemporalPath) Value() (driver.Value, error) {
	return l.path.ToWKT(), nil
}
