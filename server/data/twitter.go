package data

type TwitterOAuth struct {
	InputID       int32  `db:"input_id"`
	RequestToken  string `db:"request_token"`
	RequestSecret string `db:"request_secret"`
}

type TwitterAccount struct {
	InputID      int32  `db:"input_id"`
	ID           int64  `db:"id"`
	ScreenName   string `db:"screen_name"`
	AccessToken  string `db:"access_token"`
	AccessSecret string `db:"access_secret"`
}
