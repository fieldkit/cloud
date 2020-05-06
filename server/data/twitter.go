package data

type TwitterOAuth struct {
	SourceID      int32  `db:"source_id"`
	RequestToken  string `db:"request_token"`
	RequestSecret string `db:"request_secret"`
}

type TwitterAccount struct {
	TwitterAccountID int64  `db:"twitter_account_id"`
	ScreenName       string `db:"screen_name"`
	AccessToken      string `db:"access_token"`
	AccessSecret     string `db:"access_secret"`
}
