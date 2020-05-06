package backend

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"
	_ "github.com/paulmach/go.geo"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type Backend struct {
	db  *sqlxcache.DB
	url string
}

func OpenDatabase(url string) (*sqlxcache.DB, error) {
	return sqlxcache.Open("postgres", url)
}

func New(url string) (*Backend, error) {
	db, err := sqlxcache.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &Backend{
		url: url,
		db:  db,
	}, nil
}

func (b *Backend) URL() string {
	return b.url
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

func (b *Backend) GetOrCreateDefaultProject(ctx context.Context) (id int32, err error) {
	projects := []*data.Project{}
	if err := b.db.SelectContext(ctx, &projects, `SELECT * FROM fieldkit.project WHERE name = $1`, "Default Project"); err != nil {
		return 0, fmt.Errorf("Error querying for project: %v", err)
	}

	if len(projects) == 1 {
		return projects[0].ID, nil
	}

	project := &data.Project{
		Name: "Default Project",
		Slug: "default-project",
	}

	if err := b.db.NamedGetContext(ctx, project, `INSERT INTO fieldkit.project (name, slug) VALUES (:name, :slug) RETURNING *`, project); err != nil {
		return 0, fmt.Errorf("Error inserting project: %v", err)
	}

	return project.ID, nil
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
