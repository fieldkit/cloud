package backend

import (
	"context"
	"fmt"
	"github.com/conservify/sqlxcache"
	"github.com/fieldkit/cloud/server/data"
	_ "github.com/lib/pq"
	"time"
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

func (b *Backend) AddInput(ctx context.Context, input *data.Input) error {
	return b.db.NamedGetContext(ctx, input, `
		INSERT INTO fieldkit.input (expedition_id, name) VALUES (:expedition_id, :name) RETURNING *
		`, input)
}

func (b *Backend) AddInputToken(ctx context.Context, inputToken *data.InputToken) error {
	return b.db.NamedGetContext(ctx, inputToken, `
		INSERT INTO fieldkit.input_token (expedition_id, token) VALUES (:expedition_id, :token) RETURNING *
		`, inputToken)
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

func (b *Backend) DeleteInputToken(ctx context.Context, inputTokenID int32) error {
	_, err := b.db.ExecContext(ctx, `
		DELETE FROM fieldkit.input_token WHERE id = $1
		`, inputTokenID)
	return err
}

func (b *Backend) ListInputTokens(ctx context.Context, project, expedition string) ([]*data.InputToken, error) {
	inputTokens := []*data.InputToken{}
	if err := b.db.SelectContext(ctx, &inputTokens, `
		SELECT it.*
			FROM fieldkit.input_token AS it
				JOIN fieldkit.expedition AS e ON e.id = it.expedition_id
				JOIN fieldkit.project AS p ON p.id = e.project_id
					WHERE p.slug = $1 AND e.slug = $2
		`, project, expedition); err != nil {
		return nil, err
	}

	return inputTokens, nil
}

func (b *Backend) CheckInputToken(ctx context.Context, inputID int32, token data.Token) (bool, error) {
	count := 0
	if err := b.db.GetContext(ctx, &count, `
		SELECT COUNT(*)
			FROM fieldkit.input_token AS it
				JOIN fieldkit.input AS i ON i.expedition_id = it.expedition_id
					WHERE i.id = $1 AND it.token = $2
		`, inputID, token); err != nil {
		return false, err
	}

	return count > 0, nil
}

func (b *Backend) ListInputTokensID(ctx context.Context, expeditionID int32) ([]*data.InputToken, error) {
	inputTokens := []*data.InputToken{}
	if err := b.db.SelectContext(ctx, &inputTokens, `
		SELECT it.* FROM fieldkit.input_token AS it WHERE it.expedition_id = $1
		`, expeditionID); err != nil {
		return nil, err
	}

	return inputTokens, nil
}

func (b *Backend) AddTwitterOAuth(ctx context.Context, twitterOAuth *data.TwitterOAuth) error {
	_, err := b.db.NamedExecContext(ctx, `
		INSERT INTO fieldkit.twitter_oauth (input_id, request_token, request_secret)
			VALUES (:input_id, :request_token, :request_secret)
			ON CONFLICT (input_id)
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

func (b *Backend) AddTwitterAccountInput(ctx context.Context, twitterAccount *data.TwitterAccountInput) error {
	if _, err := b.db.NamedExecContext(ctx, `
		INSERT INTO fieldkit.twitter_account (id, screen_name, access_token, access_secret)
			VALUES (:twitter_account_id, :screen_name, :access_token, :access_secret)
			ON CONFLICT (id)
				DO UPDATE SET screen_name = :screen_name, access_token = :access_token, access_secret = :access_secret
		`, twitterAccount); err != nil {
		return err
	}

	if _, err := b.db.NamedExecContext(ctx, `
		INSERT INTO fieldkit.input_twitter_account (input_id, twitter_account_id)
			VALUES (:id, :twitter_account_id)
			ON CONFLICT (input_id)
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

func (b *Backend) Input(ctx context.Context, inputID int32) (*data.Input, error) {
	input := &data.Input{}
	if err := b.db.GetContext(ctx, input, `
		SELECT * FROM fieldkit.input WHERE id = $1
		`, inputID); err != nil {
		return nil, err
	}

	return input, nil
}

func (b *Backend) UpdateInput(ctx context.Context, input *data.Input) error {
	if _, err := b.db.NamedExecContext(ctx, `
		UPDATE fieldkit.input
			SET name = :name, team_id = :team_id, user_id = :user_id
				WHERE id = :id
		`, input); err != nil {
		return err
	}

	return nil
}

func (b *Backend) TwitterAccountInput(ctx context.Context, inputID int32) (*data.TwitterAccountInput, error) {
	twitterAccount := &data.TwitterAccountInput{}
	if err := b.db.GetContext(ctx, twitterAccount, `
		SELECT i.*, ita.twitter_account_id, ta.screen_name
			FROM fieldkit.twitter_account AS ta
				JOIN fieldkit.input_twitter_account AS ita ON ita.twitter_account_id = ta.id
				JOIN fieldkit.input AS i ON i.id = ita.input_id
					WHERE i.id = $1
		`, inputID); err != nil {
		return nil, err
	}

	return twitterAccount, nil
}

func (b *Backend) ListDeviceInputs(ctx context.Context, project, expedition string) ([]*data.DeviceInput, error) {
	devices := []*data.DeviceInput{}
	if err := b.db.SelectContext(ctx, &devices, `
		SELECT i.*, d.input_id,  d.key, d.token
			FROM fieldkit.device AS d
				JOIN fieldkit.input AS i ON i.id = d.input_id
				JOIN fieldkit.expedition AS e ON e.id = i.expedition_id
				JOIN fieldkit.project AS p ON p.id = e.project_id
					WHERE p.slug = $1 AND e.slug = $2
		`, project, expedition); err != nil {
		return nil, err
	}

	return devices, nil
}

func (b *Backend) ListTwitterAccountInputs(ctx context.Context, project, expedition string) ([]*data.TwitterAccountInput, error) {
	twitterAccounts := []*data.TwitterAccountInput{}
	if err := b.db.SelectContext(ctx, &twitterAccounts, `
		SELECT i.*, ita.twitter_account_id, ta.screen_name, ta.access_token, ta.access_secret
			FROM fieldkit.twitter_account AS ta
				JOIN fieldkit.input_twitter_account AS ita ON ita.twitter_account_id = ta.id
				JOIN fieldkit.input AS i ON i.id = ita.input_id
				JOIN fieldkit.expedition AS e ON e.id = i.expedition_id
				JOIN fieldkit.project AS p ON p.id = e.project_id
					WHERE p.slug = $1 AND e.slug = $2
		`, project, expedition); err != nil {
		return nil, err
	}

	return twitterAccounts, nil
}

func (b *Backend) ListDeviceInputsByID(ctx context.Context, expeditionID int32) ([]*data.DeviceInput, error) {
	devices := []*data.DeviceInput{}
	if err := b.db.SelectContext(ctx, &devices, `
		SELECT i.*, d.input_id, d.key, d.token
			FROM fieldkit.device AS d
				JOIN fieldkit.input AS i ON i.id = d.input_id
					WHERE i.expedition_id = $1
		`, expeditionID); err != nil {
		return nil, err
	}

	return devices, nil
}

func (b *Backend) GetDeviceInputByID(ctx context.Context, id int32) (*data.DeviceInput, error) {
	devices := []*data.DeviceInput{}
	if err := b.db.SelectContext(ctx, &devices, `
		SELECT i.*, d.input_id, d.key, d.token
			FROM fieldkit.device AS d
				JOIN fieldkit.input AS i ON i.id = d.input_id
					WHERE i.id = $1
		`, id); err != nil {
		return nil, err
	}

	if len(devices) != 1 {
		return nil, fmt.Errorf("No such Device Input")
	}

	return devices[0], nil
}

func (b *Backend) ListTwitterAccountInputsByID(ctx context.Context, expeditionID int32) ([]*data.TwitterAccountInput, error) {
	twitterAccountInputs := []*data.TwitterAccountInput{}
	if err := b.db.SelectContext(ctx, &twitterAccountInputs, `
		SELECT i.*, ita.twitter_account_id, ta.screen_name, ta.access_token, ta.access_secret
			FROM fieldkit.twitter_account AS ta
				JOIN fieldkit.input_twitter_account AS ita ON ita.twitter_account_id = ta.id
				JOIN fieldkit.input AS i ON i.id = ita.input_id
					WHERE i.expedition_id = $1
			`, expeditionID); err != nil {
		return nil, err
	}

	return twitterAccountInputs, nil
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

func (b *Backend) ListTwitterAccountInputsByAccountID(ctx context.Context, accountID int64) ([]*data.TwitterAccountInput, error) {
	twitterAccountInputs := []*data.TwitterAccountInput{}
	if err := b.db.SelectContext(ctx, &twitterAccountInputs, `
		SELECT i.*, ita.twitter_account_id, ta.screen_name, ta.access_token, ta.access_secret
			FROM fieldkit.twitter_account AS ta
				JOIN fieldkit.input_twitter_account AS ita ON ita.twitter_account_id = ta.id
				JOIN fieldkit.input AS i ON i.id = ita.input_id
					WHERE ta.id = $1
			`, accountID); err != nil {
		return nil, err
	}

	return twitterAccountInputs, nil
}

func (b *Backend) AddDevice(ctx context.Context, device *data.Device) error {
	return b.db.NamedGetContext(ctx, device, `
		INSERT INTO fieldkit.device (input_id, key, token) VALUES (:input_id, :key, :token) RETURNING *
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

func (b *Backend) AddDocument(ctx context.Context, document *data.Document) error {
	_, err := b.db.NamedExecContext(ctx, `
		INSERT INTO fieldkit.document (schema_id, input_id, team_id, user_id, timestamp, location, data)
			VALUES (:schema_id, :input_id, :team_id, :user_id, :timestamp, ST_SetSRID(ST_GeomFromText(:location),4326), :data)
		`, document)
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

func (b *Backend) ListDocuments(ctx context.Context, project string, expedition string, token *PagingToken) (*data.DocumentsPage, *PagingToken, error) {
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
	documents := []*data.Document{}
	if err := b.db.SelectContext(ctx, &documents, `
		SELECT d.id, d.schema_id, d.input_id, d.team_id, d.user_id, d.timestamp, ST_AsBinary(d.location) AS location, d.data
			FROM fieldkit.document AS d
				JOIN fieldkit.input AS i ON i.id = d.input_id
				JOIN fieldkit.expedition AS e ON e.id = i.expedition_id
				JOIN fieldkit.project AS p ON p.id = e.project_id
					WHERE d.visible AND p.slug = $1 AND e.slug = $2 AND d.insertion < $3 AND (insertion >= $4) AND
                                              ST_X(d.location) != 0 AND ST_Y(d.location) != 0
                        ORDER BY timestamp
                        LIMIT $5 OFFSET $6
		`, project, expedition, before, after, pageSize, token.page*pageSize); err != nil {
		return nil, nil, err
	}

	nextToken := &PagingToken{}
	if int32(len(documents)) < pageSize {
		nextToken.time = before.UnixNano()
	} else {
		nextToken.time = token.time
		nextToken.page += 1
	}

	return &data.DocumentsPage{
		Documents: documents,
	}, nextToken, nil
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
