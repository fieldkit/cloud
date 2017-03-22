package backend

import (
	"context"

	"github.com/O-C-R/sqlxcache"
	_ "github.com/lib/pq"

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

func (b *Backend) AddInputID(ctx context.Context, expeditionID int32) (int32, error) {
	var inputID int32
	if err := b.db.GetContext(ctx, &inputID, `
		INSERT INTO fieldkit.input (expedition_id) VALUES ($1) RETURNING id
		`, expeditionID); err != nil {
		return int32(0), err
	}

	return inputID, nil
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

func (b *Backend) AddTwitterAccount(ctx context.Context, twitterAccount *data.TwitterAccount) error {
	if _, err := b.db.NamedExecContext(ctx, `
		BEGIN;
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

func (b *Backend) TwitterAccount(ctx context.Context, inputID int32) (*data.TwitterAccount, error) {
	twitterAccount := &data.TwitterAccount{}
	if err := b.db.GetContext(ctx, twitterAccount, `
		SELECT i.id, i.expedition_id, ita.twitter_account_id, ta.screen_name
			FROM fieldkit.twitter_account AS ta
				JOIN fieldkit.input_twitter_account AS ita ON ita.twitter_account_id = ta.id
				JOIN fieldkit.input AS i ON i.id = ita.input_id
					WHERE i.id = $1
		`, inputID); err != nil {
		return nil, err
	}

	return twitterAccount, nil
}

func (b *Backend) ListTwitterAccounts(ctx context.Context, project, expedition string) ([]*data.TwitterAccount, error) {
	twitterAccounts := []*data.TwitterAccount{}
	if err := b.db.SelectContext(ctx, &twitterAccounts, `
		SELECT i.id, i.expedition_id, ita.twitter_account_id, ta.screen_name, ta.access_token, ta.access_secret 
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

func (b *Backend) ListTwitterAccountsByID(ctx context.Context, expeditionID int32) ([]*data.TwitterAccount, error) {
	twitterAccounts := []*data.TwitterAccount{}
	if err := b.db.SelectContext(ctx, &twitterAccounts, `
		SELECT i.id, i.expedition_id, ita.twitter_account_id, ta.screen_name, ta.access_token, ta.access_secret 
			FROM fieldkit.twitter_account AS ta 
				JOIN fieldkit.input_twitter_account AS ita ON ita.twitter_account_id = ta.id
				JOIN fieldkit.input AS i ON i.id = ita.input_id
					WHERE i.expedition_id = $1
			`, expeditionID); err != nil {
		return nil, err
	}

	return twitterAccounts, nil
}
