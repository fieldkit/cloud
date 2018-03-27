package backend

import (
	"context"
	"fmt"
	_ "log"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/paulmach/go.geo"

	"github.com/conservify/sqlxcache"
	"github.com/fieldkit/cloud/server/data"
)

type Backend struct {
	db *sqlxcache.DB
}

func New(url string) (*Backend, error) {
	db, err := sqlxcache.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &Backend{
		db: db,
	}, nil
}

func (b *Backend) AddSource(ctx context.Context, source *data.Source) error {
	return b.db.NamedGetContext(ctx, source, `
		INSERT INTO fieldkit.source (expedition_id, name) VALUES (:expedition_id, :name) RETURNING *
		`, source)
}

func (b *Backend) AddSourceToken(ctx context.Context, sourceToken *data.SourceToken) error {
	return b.db.NamedGetContext(ctx, sourceToken, `
		INSERT INTO fieldkit.source_token (expedition_id, token) VALUES (:expedition_id, :token) RETURNING *
		`, sourceToken)
}

func (b *Backend) CheckInviteToken(ctx context.Context, inviteToken data.Token) (bool, error) {
	count := 0
	if err := b.db.GetContext(ctx, &count, `
		SELECT COUNT(*) FROM fieldkit.invite_token WHERE token = $1
		`, inviteToken); err != nil {
		return false, err
	}

	return count > 0, nil
}

func (b *Backend) DeleteSourceToken(ctx context.Context, sourceTokenID int32) error {
	_, err := b.db.ExecContext(ctx, `
		DELETE FROM fieldkit.source_token WHERE id = $1
		`, sourceTokenID)
	return err
}

func (b *Backend) ListSourceTokens(ctx context.Context, project, expedition string) ([]*data.SourceToken, error) {
	sourceTokens := []*data.SourceToken{}
	if err := b.db.SelectContext(ctx, &sourceTokens, `
		SELECT it.*
			FROM fieldkit.source_token AS it
				JOIN fieldkit.expedition AS e ON e.id = it.expedition_id
				JOIN fieldkit.project AS p ON p.id = e.project_id
					WHERE p.slug = $1 AND e.slug = $2
		`, project, expedition); err != nil {
		return nil, err
	}

	return sourceTokens, nil
}

func (b *Backend) CheckSourceToken(ctx context.Context, sourceID int32, token data.Token) (bool, error) {
	count := 0
	if err := b.db.GetContext(ctx, &count, `
		SELECT COUNT(*)
			FROM fieldkit.source_token AS it
				JOIN fieldkit.source AS i ON i.expedition_id = it.expedition_id
					WHERE i.id = $1 AND it.token = $2
		`, sourceID, token); err != nil {
		return false, err
	}

	return count > 0, nil
}

func (b *Backend) ListSourceTokensID(ctx context.Context, expeditionID int32) ([]*data.SourceToken, error) {
	sourceTokens := []*data.SourceToken{}
	if err := b.db.SelectContext(ctx, &sourceTokens, `
		SELECT it.* FROM fieldkit.source_token AS it WHERE it.expedition_id = $1
		`, expeditionID); err != nil {
		return nil, err
	}

	return sourceTokens, nil
}

func (b *Backend) AddTwitterOAuth(ctx context.Context, twitterOAuth *data.TwitterOAuth) error {
	_, err := b.db.NamedExecContext(ctx, `
		INSERT INTO fieldkit.twitter_oauth (source_id, request_token, request_secret)
			VALUES (:source_id, :request_token, :request_secret)
			ON CONFLICT (source_id)
				DO UPDATE SET request_token = :request_token, request_secret = :request_secret
		`, twitterOAuth)
	return err
}

func (b *Backend) TwitterOAuth(ctx context.Context, requestToken string) (*data.TwitterOAuth, error) {
	twitterOAuth := &data.TwitterOAuth{}
	if err := b.db.GetContext(ctx, twitterOAuth, `
		SELECT * FROM fieldkit.twitter_oauth WHERE request_token = $1
		`, requestToken); err != nil {
		return nil, err
	}

	return twitterOAuth, nil
}

func (b *Backend) DeleteTwitterOAuth(ctx context.Context, requestToken string) error {
	_, err := b.db.ExecContext(ctx, `
		DELETE FROM fieldkit.twitter_oauth WHERE request_token = $1
		`, requestToken)
	return err
}

func (b *Backend) DeleteInviteToken(ctx context.Context, inviteToken data.Token) error {
	_, err := b.db.ExecContext(ctx, `
		DELETE FROM fieldkit.invite_token WHERE token = $1
		`, inviteToken)
	return err
}

func (b *Backend) AddTwitterAccountSource(ctx context.Context, twitterAccount *data.TwitterAccountSource) error {
	if _, err := b.db.NamedExecContext(ctx, `
		INSERT INTO fieldkit.twitter_account (id, screen_name, access_token, access_secret)
			VALUES (:twitter_account_id, :screen_name, :access_token, :access_secret)
			ON CONFLICT (id)
				DO UPDATE SET screen_name = :screen_name, access_token = :access_token, access_secret = :access_secret
		`, twitterAccount); err != nil {
		return err
	}

	if _, err := b.db.NamedExecContext(ctx, `
		INSERT INTO fieldkit.source_twitter_account (source_id, twitter_account_id)
			VALUES (:id, :twitter_account_id)
			ON CONFLICT (source_id)
				DO UPDATE SET twitter_account_id = :twitter_account_id
		`, twitterAccount); err != nil {
		return err
	}

	return nil
}

func (b *Backend) Expedition(ctx context.Context, expeditionID int32) (*data.Expedition, error) {
	expedition := &data.Expedition{}
	if err := b.db.GetContext(ctx, expedition, `
		SELECT * FROM fieldkit.expedition WHERE id = $1
		`, expeditionID); err != nil {
		return nil, err
	}

	return expedition, nil
}

func (b *Backend) Source(ctx context.Context, sourceID int32) (*data.Source, error) {
	source := &data.Source{}
	if err := b.db.GetContext(ctx, source, `
		SELECT * FROM fieldkit.source WHERE id = $1
		`, sourceID); err != nil {
		return nil, err
	}

	return source, nil
}

func (b *Backend) UpdateSource(ctx context.Context, source *data.Source) error {
	if _, err := b.db.NamedExecContext(ctx, `
		UPDATE fieldkit.source
			SET name = :name, team_id = :team_id, user_id = :user_id
				WHERE id = :id
		`, source); err != nil {
		return err
	}

	return nil
}

func (b *Backend) TwitterAccountSource(ctx context.Context, sourceID int32) (*data.TwitterAccountSource, error) {
	twitterAccount := &data.TwitterAccountSource{}
	if err := b.db.GetContext(ctx, twitterAccount, `
		SELECT i.*, ita.twitter_account_id, ta.screen_name
			FROM fieldkit.twitter_account AS ta
				JOIN fieldkit.source_twitter_account AS ita ON ita.twitter_account_id = ta.id
				JOIN fieldkit.source AS i ON i.id = ita.source_id
					WHERE i.id = $1
		`, sourceID); err != nil {
		return nil, err
	}

	return twitterAccount, nil
}

func (b *Backend) ListAllDeviceSources(ctx context.Context) ([]*data.DeviceSource, error) {
	devices := []*data.DeviceSource{}
	if err := b.db.SelectContext(ctx, &devices, `
		SELECT i.*, d.source_id, d.key, d.token
			FROM fieldkit.device AS d
				JOIN fieldkit.source AS i ON i.id = d.source_id
				JOIN fieldkit.expedition AS e ON e.id = i.expedition_id
				JOIN fieldkit.project AS p ON p.id = e.project_id
                        ORDER BY i.name
		`); err != nil {
		return nil, err
	}

	return devices, nil
}

func (b *Backend) ListDeviceSources(ctx context.Context, project, expedition string) ([]*data.DeviceSource, error) {
	devices := []*data.DeviceSource{}
	if err := b.db.SelectContext(ctx, &devices, `
		SELECT s.*, d.source_id, d.key, d.token
			FROM fieldkit.device AS d
				JOIN fieldkit.source AS s ON s.id = d.source_id
				JOIN fieldkit.expedition AS e ON e.id = s.expedition_id
				JOIN fieldkit.project AS p ON p.id = e.project_id
					WHERE p.slug = $1 AND e.slug = $2 AND s.visible
		`, project, expedition); err != nil {
		return nil, err
	}

	return devices, nil
}

func (b *Backend) ListTwitterAccountSources(ctx context.Context, project, expedition string) ([]*data.TwitterAccountSource, error) {
	twitterAccounts := []*data.TwitterAccountSource{}
	if err := b.db.SelectContext(ctx, &twitterAccounts, `
		SELECT s.*, ita.twitter_account_id, ta.screen_name, ta.access_token, ta.access_secret
			FROM fieldkit.twitter_account AS ta
				JOIN fieldkit.source_twitter_account AS ita ON ita.twitter_account_id = ta.id
				JOIN fieldkit.source AS s ON s.id = ita.source_id
				JOIN fieldkit.expedition AS e ON e.id = s.expedition_id
				JOIN fieldkit.project AS p ON p.id = e.project_id
					WHERE p.slug = $1 AND e.slug = $2 AND s.visible
		`, project, expedition); err != nil {
		return nil, err
	}

	return twitterAccounts, nil
}

func (b *Backend) ListDeviceSourcesByExpeditionID(ctx context.Context, expeditionID int32) ([]*data.DeviceSource, error) {
	devices := []*data.DeviceSource{}
	if err := b.db.SelectContext(ctx, &devices, `
		SELECT i.*, d.source_id, d.key, d.token
			FROM fieldkit.device AS d
				JOIN fieldkit.source AS i ON i.id = d.source_id
					WHERE i.expedition_id = $1
		`, expeditionID); err != nil {
		return nil, err
	}

	return devices, nil
}

func (b *Backend) GetSourceByID(ctx context.Context, id int32) (*data.Source, error) {
	sources := []*data.Source{}
	if err := b.db.SelectContext(ctx, &sources, `SELECT s.* FROM fieldkit.source AS s WHERE s.id = $1`, id); err != nil {
		return nil, err
	}

	if len(sources) != 1 {
		return nil, fmt.Errorf("No such Source")
	}

	return sources[0], nil
}

func (b *Backend) GetDeviceSourceByID(ctx context.Context, id int32) (*data.DeviceSource, error) {
	devices := []*data.DeviceSource{}
	if err := b.db.SelectContext(ctx, &devices, `
		SELECT i.*, d.source_id, d.key, d.token
			FROM fieldkit.device AS d
				JOIN fieldkit.source AS i ON i.id = d.source_id
					WHERE i.id = $1
		`, id); err != nil {
		return nil, err
	}

	if len(devices) != 1 {
		return nil, fmt.Errorf("No such Device Source")
	}

	return devices[0], nil
}

func (b *Backend) ListTwitterAccountSourcesByExpeditionID(ctx context.Context, expeditionID int32) ([]*data.TwitterAccountSource, error) {
	twitterAccountSources := []*data.TwitterAccountSource{}
	if err := b.db.SelectContext(ctx, &twitterAccountSources, `
		SELECT i.*, ita.twitter_account_id, ta.screen_name, ta.access_token, ta.access_secret
			FROM fieldkit.twitter_account AS ta
				JOIN fieldkit.source_twitter_account AS ita ON ita.twitter_account_id = ta.id
				JOIN fieldkit.source AS i ON i.id = ita.source_id
					WHERE i.expedition_id = $1
			`, expeditionID); err != nil {
		return nil, err
	}

	return twitterAccountSources, nil
}

func (b *Backend) ListTwitterAccounts(ctx context.Context) ([]*data.TwitterAccount, error) {
	twitterAccounts := []*data.TwitterAccount{}
	if err := b.db.SelectContext(ctx, &twitterAccounts, `
		SELECT id AS twitter_account_id, screen_name, access_token, access_secret FROM fieldkit.twitter_account
		`); err != nil {
		return nil, err
	}

	return twitterAccounts, nil
}

func (b *Backend) ListTwitterAccountSourcesByAccountID(ctx context.Context, accountID int64) ([]*data.TwitterAccountSource, error) {
	twitterAccountSources := []*data.TwitterAccountSource{}
	if err := b.db.SelectContext(ctx, &twitterAccountSources, `
		SELECT i.*, ita.twitter_account_id, ta.screen_name, ta.access_token, ta.access_secret
			FROM fieldkit.twitter_account AS ta
				JOIN fieldkit.source_twitter_account AS ita ON ita.twitter_account_id = ta.id
				JOIN fieldkit.source AS i ON i.id = ita.source_id
					WHERE ta.id = $1
			`, accountID); err != nil {
		return nil, err
	}

	return twitterAccountSources, nil
}

func (b *Backend) AddDevice(ctx context.Context, device *data.Device) error {
	return b.db.NamedGetContext(ctx, device, `
		INSERT INTO fieldkit.device (source_id, key, token) VALUES (:source_id, :key, :token) RETURNING *
		`, device)
}

func (b *Backend) AddRawSchema(ctx context.Context, schema *data.RawSchema) error {
	return b.db.NamedGetContext(ctx, schema, `
		INSERT INTO fieldkit.schema (project_id, json_schema) VALUES (:project_id, :json_schema) RETURNING *
		`, schema)
}

func (b *Backend) AddSchema(ctx context.Context, schema *data.Schema) error {
	return b.db.NamedGetContext(ctx, schema, `
		INSERT INTO fieldkit.schema (project_id, json_schema) VALUES (:project_id, :json_schema) RETURNING *
		`, schema)
}

func (b *Backend) UpdateSchema(ctx context.Context, schema *data.Schema) error {
	return b.db.NamedGetContext(ctx, schema, `
		UPDATE fieldkit.schema SET project_id = :project_id, json_schema = :json_schema RETURNING *
		`, schema)
}

func (b *Backend) ListSchemas(ctx context.Context, project string) ([]*data.Schema, error) {
	schemas := []*data.Schema{}
	if err := b.db.SelectContext(ctx, &schemas, `
		SELECT s.*
			FROM fieldkit.schema AS s
				JOIN fieldkit.project AS p ON p.id = s.project_id
					WHERE p.slug = $1
		`, project); err != nil {
		return nil, err
	}

	return schemas, nil
}

func (b *Backend) ListSchemasByID(ctx context.Context, projectID int32) ([]*data.Schema, error) {
	schemas := []*data.Schema{}
	if err := b.db.SelectContext(ctx, &schemas, `
		SELECT s.*
			FROM fieldkit.schema AS s
				JOIN fieldkit.project AS p ON p.id = s.project_id
					WHERE s.project_id = $1
			`, projectID); err != nil {
		return nil, err
	}

	return schemas, nil
}

func (b *Backend) AddRecord(ctx context.Context, record *data.Record) error {
	_, err := b.db.NamedExecContext(ctx, `
		INSERT INTO fieldkit.record (schema_id, source_id, team_id, user_id, timestamp, location, data, fixed, visible)
			VALUES (:schema_id, :source_id, :team_id, :user_id, :timestamp, ST_SetSRID(ST_GeomFromText(:location), 4326), :data, :fixed, :visible)
		`, record)

	return err
}

func (b *Backend) SetSchemaID(ctx context.Context, schema *data.Schema) (int32, error) {
	var schemaID int32
	if err := b.db.NamedGetContext(ctx, &schemaID, `
		INSERT INTO fieldkit.schema (project_id, json_schema)
			VALUES (:project_id, :json_schema)
			ON CONFLICT ((json_schema->'id'))
				DO UPDATE SET project_id = :project_id, json_schema = :json_schema
			RETURNING id
		`, schema); err != nil {
		return int32(0), err
	}

	return schemaID, nil
}

const DefaultPageSize = 100

func (b *Backend) ListRecordsBySource(ctx context.Context, sourceID int, descending, includeInvisible bool, token *PagingToken) (*data.AnalysedRecordsPage, *PagingToken, error) {
	if token == nil {
		after := time.Now().AddDate(-10, 0, 0)
		token = &PagingToken{
			time: after.UnixNano(),
			page: 0,
		}
	}

	after := time.Unix(0, token.time)
	pageSize := int32(DefaultPageSize)
	before := time.Now()
	records := []*data.AnalysedRecord{}
	order := "ASC"
	if descending {
		order = "DESC"
	}
	if err := b.db.SelectContext(ctx, &records, fmt.Sprintf(`
		SELECT d.id, d.schema_id, d.source_id, d.team_id, d.user_id, d.timestamp, ST_AsBinary(d.location) AS location, d.data, d.visible, d.fixed,
		       COALESCE(ra.record_id, 0) AS record_id, COALESCE(ra.outlier, false) AS outlier, COALESCE(ra.manually_excluded, false) AS manually_excluded
			FROM fieldkit.record AS d
				JOIN fieldkit.source AS i ON i.id = d.source_id
				JOIN fieldkit.expedition AS e ON e.id = i.expedition_id
				JOIN fieldkit.project AS p ON p.id = e.project_id
                                LEFT JOIN fieldkit.record_analysis AS ra ON ra.record_id = d.id
			WHERE
				(d.visible OR $6) AND i.id = $1 AND d.insertion < $2 AND (insertion >= $3) AND
				ST_X(d.location) != 0 AND ST_Y(d.location) != 0
                        ORDER BY timestamp %s
                        LIMIT $4 OFFSET $5
		`, order), sourceID, before, after, pageSize, token.page*pageSize, includeInvisible); err != nil {
		return nil, nil, err
	}

	nextToken := &PagingToken{}
	if int32(len(records)) < pageSize {
		nextToken.time = before.UnixNano()
	} else {
		nextToken.time = token.time
		nextToken.page = token.page + 1
	}

	return &data.AnalysedRecordsPage{
		Records: records,
	}, nextToken, nil
}

type ReadingSummary struct {
	Name string
}

type FeatureSummary struct {
	SourceID         int           `db:"source_id"`
	UpdatedAt        time.Time     `db:"updated_at"`
	NumberOfFeatures int           `db:"number_of_features"`
	LastFeatureID    int           `db:"last_feature_id"`
	StartTime        time.Time     `db:"start_time"`
	EndTime          time.Time     `db:"end_time"`
	Centroid         data.Location `db:"centroid"`
	Radius           float64       `db:"radius"`
}

type NotificationStatus struct {
	ID        int       `db:"id"`
	SourceID  int       `db:"source_id"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (b *Backend) FeatureSummaryBySourceID(ctx context.Context, sourceId int) (*FeatureSummary, error) {
	summaries := []*FeatureSummary{}
	if err := b.db.SelectContext(ctx, &summaries, `
		  SELECT
		    c.source_id, c.updated_at, c.number_of_features, c.last_feature_id, c.start_time, c.end_time, ST_AsBinary(c.centroid) AS centroid, radius
		  FROM
		    fieldkit.sources_summaries c
		  WHERE c.source_id = $1
	      `, sourceId); err != nil {
		return nil, err
	}
	if len(summaries) != 1 {
		return &FeatureSummary{
			UpdatedAt: time.Now(),
			StartTime: time.Now(),
			EndTime:   time.Now(),
			Centroid:  *data.NewLocation([]float64{0, 0}),
		}, nil
	}

	return summaries[0], nil
}

type GeometryClusterSummary struct {
	ClusterID        int           `db:"cluster_id"`
	SourceID         int           `db:"source_id"`
	UpdatedAt        time.Time     `db:"updated_at"`
	NumberOfFeatures int           `db:"number_of_features"`
	StartTime        time.Time     `db:"start_time"`
	EndTime          time.Time     `db:"end_time"`
	Centroid         data.Location `db:"centroid"`
	Radius           float64       `db:"radius"`
	Geometry         TemporalPath  `db:"geometry"`
}

func (b *Backend) SpatialClustersBySourceID(ctx context.Context, sourceId int) (summaries []*GeometryClusterSummary, err error) {
	summaries = []*GeometryClusterSummary{}
	if err := b.db.SelectContext(ctx, &summaries, `
		  SELECT
		    c.cluster_id, c.source_id, c.updated_at, c.number_of_features, c.start_time, c.end_time, ST_AsBinary(c.centroid) AS centroid, radius
		  FROM
		    fieldkit.sources_spatial_clusters c
		  WHERE c.source_id = $1
	      `, sourceId); err != nil {
		return nil, err
	}
	return summaries, nil
}

func (b *Backend) TemporalClustersBySourceID(ctx context.Context, sourceId int) (summaries []*GeometryClusterSummary, err error) {
	summaries = []*GeometryClusterSummary{}
	if err := b.db.SelectContext(ctx, &summaries, `
		  SELECT
		    c.cluster_id, c.source_id, c.updated_at, c.number_of_features, c.start_time, c.end_time, ST_AsBinary(c.centroid) AS centroid, c.radius, ST_AsBinary(g.geometry) AS geometry
		  FROM
		    fieldkit.sources_temporal_clusters c JOIN
                    fieldkit.sources_temporal_geometries g ON (c.source_id = g.source_id AND c.cluster_id = g.cluster_id)
		  WHERE c.source_id = $1
	      `, sourceId); err != nil {
		return nil, err
	}
	return summaries, nil
}

type BoundingBox struct {
	NorthEast *data.Location
	SouthWest *data.Location
}

type MapFeatures struct {
	TemporalClusters   []*GeometryClusterSummary
	SpatialClusters    []*GeometryClusterSummary
	TemporalGeometries []*GeometryClusterSummary
}

func (b *Backend) QueryMapFeatures(ctx context.Context, bb *BoundingBox) (mf *MapFeatures, err error) {
	ne := bb.NorthEast.Coordinates()
	sw := bb.SouthWest.Coordinates()

	spatialClusters := []*GeometryClusterSummary{}
	temporalClusters := []*GeometryClusterSummary{}

	// TODO: Note that this uses the centroid and doesn't take into consideration the radius.
	if err := b.db.SelectContext(ctx, &spatialClusters, `
		  SELECT
		    c.cluster_id, c.source_id, c.updated_at, c.number_of_features, c.start_time, c.end_time, ST_AsBinary(c.centroid) AS centroid, radius
		  FROM
		    fieldkit.sources_spatial_clusters c
                  WHERE
                    c.centroid && ST_SetSRID(ST_MakeEnvelope($1, $2, $3, $4), 4326)
		`, ne[0], ne[1], sw[0], sw[1]); err != nil {
		return nil, err
	}

	if err := b.db.SelectContext(ctx, &temporalClusters, `
		  SELECT
		    c.cluster_id, c.source_id, c.updated_at, c.number_of_features, c.start_time, c.end_time, ST_AsBinary(c.centroid) AS centroid, c.radius, ST_AsBinary(g.geometry) AS geometry
		  FROM
		    fieldkit.sources_temporal_clusters c JOIN
                    fieldkit.sources_temporal_geometries g ON (c.source_id = g.source_id AND c.cluster_id = g.cluster_id)
                  WHERE
                    g.geometry && ST_SetSRID(ST_MakeEnvelope($1, $2, $3, $4), 4326)
		`, ne[0], ne[1], sw[0], sw[1]); err != nil {
		return nil, err
	}

	mf = &MapFeatures{
		SpatialClusters:    spatialClusters,
		TemporalClusters:   temporalClusters,
		TemporalGeometries: temporalClusters,
	}

	return
}

func (b *Backend) ReadingsBySourceID(ctx context.Context, sourceId int) (summaries []*ReadingSummary, err error) {
	summaries = []*ReadingSummary{}
	if err := b.db.SelectContext(ctx, &summaries, `SELECT jsonb_object_keys(q.data) AS name FROM (SELECT data FROM fieldkit.record r WHERE r.source_id = $1 ORDER BY r.timestamp DESC LIMIT 10) q GROUP BY name`, sourceId); err != nil {
		return nil, err
	}
	return summaries, nil
}

func (b *Backend) ListRecordsByID(ctx context.Context, id int) (*data.AnalysedRecordsPage, error) {
	records := []*data.AnalysedRecord{}
	if err := b.db.SelectContext(ctx, &records, `
		SELECT d.id, d.schema_id, d.source_id, d.team_id, d.user_id, d.timestamp, ST_AsBinary(d.location) AS location, d.data, ra.*
			FROM fieldkit.record AS d
				JOIN fieldkit.source AS i ON i.id = d.source_id
				JOIN fieldkit.expedition AS e ON e.id = i.expedition_id
				JOIN fieldkit.project AS p ON p.id = e.project_id
                                LEFT JOIN fieldkit.record_analysis AS ra ON ra.record_id = d.id
			WHERE d.visible AND d.id = $1
		`, id); err != nil {
		return nil, err
	}

	return &data.AnalysedRecordsPage{
		Records: records,
	}, nil
}

func (b *Backend) TwitterListCredentialer() *TwitterListCredentialer {
	return &TwitterListCredentialer{b}
}

type TwitterListCredentialer struct {
	b *Backend
}

func (t *TwitterListCredentialer) UserList() ([]int64, error) {
	ids := []int64{}
	if err := t.b.db.Select(&ids, `
		SELECT id FROM fieldkit.twitter_account
		`); err != nil {
		return nil, err
	}

	return ids, nil
}

func (t *TwitterListCredentialer) UserCredentials(id int64) (accessToken, accessSecret string, err error) {
	twitterAccount := &data.TwitterAccount{}
	if err := t.b.db.Get(twitterAccount, `
		SELECT id AS twitter_account_id, screen_name, access_token, access_secret FROM fieldkit.twitter_account WHERE id = $1
		`, id); err != nil {
		return "", "", err
	}

	return twitterAccount.AccessToken, twitterAccount.AccessSecret, nil
}
