package backend

import (
	"context"

	"github.com/O-C-R/sqlxcache"
	"github.com/lib/pq"

	"github.com/O-C-R/fieldkit/server/data"
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

func (b *Backend) ListTwitterAccountInputsByID(ctx context.Context, expeditionID int32) ([]*data.TwitterAccountInput, error) {
	TwitterAccountInputs := []*data.TwitterAccountInput{}
	if err := b.db.SelectContext(ctx, &TwitterAccountInputs, `
		SELECT i.*, ita.twitter_account_id, ta.screen_name, ta.access_token, ta.access_secret 
			FROM fieldkit.twitter_account AS ta 
				JOIN fieldkit.input_twitter_account AS ita ON ita.twitter_account_id = ta.id
				JOIN fieldkit.input AS i ON i.id = ita.input_id
					WHERE i.expedition_id = $1
			`, expeditionID); err != nil {
		return nil, err
	}

	return TwitterAccountInputs, nil
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
	TwitterAccountInputs := []*data.TwitterAccountInput{}
	if err := b.db.SelectContext(ctx, &TwitterAccountInputs, `
		SELECT i.*, ita.twitter_account_id, ta.screen_name, ta.access_token, ta.access_secret 
			FROM fieldkit.twitter_account AS ta 
				JOIN fieldkit.input_twitter_account AS ita ON ita.twitter_account_id = ta.id
				JOIN fieldkit.input AS i ON i.id = ita.input_id
					WHERE ta.id = $1
			`, accountID); err != nil {
		return nil, err
	}

	return TwitterAccountInputs, nil
}

func (b *Backend) AddFieldkitInput(ctx context.Context, fieldkitInput *data.FieldkitInput) error {
	return b.db.NamedGetContext(ctx, fieldkitInput, `
		INSERT INTO fieldkit.input_fieldkit (input_id) VALUES (:id) RETURNING input_id AS id
		`, fieldkitInput)
}

func (b *Backend) SetFieldkitInputBinary(ctx context.Context, fieldkitBinary *data.FieldkitBinary) error {
	_, err := b.db.ExecContext(ctx, `
		INSERT INTO fieldkit.fieldkit_binary (id, input_id, fields) VALUES ($1, $2, $3)
			ON CONFLICT (id, input_id)
				DO UPDATE SET fields = $3
		`, fieldkitBinary.ID, fieldkitBinary.InputID, pq.Array(fieldkitBinary.Fields))
	return err
}

func (b *Backend) FieldkitInput(ctx context.Context, inputID int32) (*data.FieldkitInput, error) {
	fieldkitInput := &data.FieldkitInput{}
	if err := b.db.GetContext(ctx, fieldkitInput, `
		SELECT i.*
			FROM fieldkit.input_fieldkit AS if
				JOIN fieldkit.input AS i ON i.id = if.input_id
					WHERE i.id = $1
		`, inputID); err != nil {
		return nil, err
	}

	return fieldkitInput, nil
}

func (b *Backend) ListFieldkitInputs(ctx context.Context, project, expedition string) ([]*data.FieldkitInput, error) {
	fieldkitInput := []*data.FieldkitInput{}
	if err := b.db.SelectContext(ctx, &fieldkitInput, `
		SELECT i.*
			FROM fieldkit.input_fieldkit AS if
				JOIN fieldkit.input AS i ON i.id = if.input_id
				JOIN fieldkit.expedition AS e ON e.id = i.expedition_id
				JOIN fieldkit.project AS p ON p.id = e.project_id
					WHERE p.slug = $1 AND e.slug = $2
		`, project, expedition); err != nil {
		return nil, err
	}

	return fieldkitInput, nil
}

func (b *Backend) ListFieldkitInputsByID(ctx context.Context, expeditionID int32) ([]*data.FieldkitInput, error) {
	fieldkitInputs := []*data.FieldkitInput{}
	if err := b.db.SelectContext(ctx, &fieldkitInputs, `
		SELECT i.*
			FROM fieldkit.input_fieldkit AS if
				JOIN fieldkit.input AS i ON i.id = if.input_id
					WHERE i.expedition_id = $1
			`, expeditionID); err != nil {
		return nil, err
	}

	return fieldkitInputs, nil
}

func (b *Backend) AddSchema(ctx context.Context, fieldkitInput *data.Schema) error {
	return b.db.NamedGetContext(ctx, fieldkitInput, `
		INSERT INTO fieldkit.schema (project_id, json_schema) VALUES (:project_id, :json_schema) RETURNING *
		`, fieldkitInput)
}

func (b *Backend) UpdateSchema(ctx context.Context, fieldkitInput *data.Schema) error {
	return b.db.NamedGetContext(ctx, fieldkitInput, `
		UPDATE fieldkit.schema SET project_id = :project_id, json_schema = :json_schema RETURNING *
		`, fieldkitInput)
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

func (b *Backend) ListDocuments(ctx context.Context, project, expedition string) ([]*data.Document, error) {
	documents := []*data.Document{}
	if err := b.db.SelectContext(ctx, &documents, `
		SELECT d.id, d.schema_id, d.input_id, d.team_id, d.user_id, d.timestamp, ST_AsBinary(d.location), d.data
			FROM fieldkit.document AS d
				JOIN fieldkit.input AS i ON i.id = d.input_id
				JOIN fieldkit.expedition AS e ON e.id = i.expedition_id
				JOIN fieldkit.project AS p ON p.id = e.project_id
					WHERE p.slug = $1 AND e.slug = $2
		`, project, expedition); err != nil {
		return nil, err
	}

	return documents, nil
}

func (b *Backend) ListDocumentsByID(ctx context.Context, expeditionID int32) ([]*data.Document, error) {
	documents := []*data.Document{}
	if err := b.db.SelectContext(ctx, &documents, `
		SELECT d.id, d.schema_id, d.input_id, d.team_id, d.user_id, d.timestamp, ST_AsBinary(d.location), d.data
			FROM fieldkit.document AS d
				JOIN fieldkit.input AS i ON i.id = d.input_id
					WHERE i.expedition_id = $1
		`, expeditionID); err != nil {
		return nil, err
	}

	return documents, nil
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
