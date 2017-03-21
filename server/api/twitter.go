package api

import (
	"github.com/O-C-R/sqlxcache"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	oauth1Twitter "github.com/dghubble/oauth1/twitter"
	"github.com/goadesign/goa"

	"github.com/O-C-R/fieldkit/server/api/app"
	"github.com/O-C-R/fieldkit/server/data"
)

func TwitterAccountType(twitterAccount *data.TwitterAccount) *app.TwitterAccount {
	return &app.TwitterAccount{
		ID:               int(twitterAccount.ID),
		ExpeditionID:     int(twitterAccount.ExpeditionID),
		TwitterAccountID: int(twitterAccount.TwitterAccountID),
		ScreenName:       twitterAccount.ScreenName,
	}
}

func TwitterAccountsType(twitterAccounts []*data.TwitterAccount) *app.TwitterAccounts {
	twitterAccountsCollection := make([]*app.TwitterAccount, len(twitterAccounts))
	for i, twitterAccount := range twitterAccounts {
		twitterAccountsCollection[i] = TwitterAccountType(twitterAccount)
	}

	return &app.TwitterAccounts{
		TwitterAccounts: twitterAccountsCollection,
	}
}

type TwitterControllerOptions struct {
	Database       *sqlxcache.DB
	ConsumerKey    string
	ConsumerSecret string
}

// TwitterController implements the twitter resource.
type TwitterController struct {
	*goa.Controller
	options TwitterControllerOptions
	config  *oauth1.Config
}

func NewTwitterController(service *goa.Service, options TwitterControllerOptions) *TwitterController {
	return &TwitterController{
		Controller: service.NewController("TwitterController"),
		options:    options,
		config: &oauth1.Config{
			ConsumerKey:    options.ConsumerKey,
			ConsumerSecret: options.ConsumerSecret,
			CallbackURL:    "https://api.data.fieldkit.org/twitter/callback",
			Endpoint:       oauth1Twitter.AuthorizeEndpoint,
		},
	}
}

func (c *TwitterController) Add(ctx *app.AddTwitterContext) error {
	twitterOAuth := &data.TwitterOAuth{}
	if err := c.options.Database.GetContext(ctx, &twitterOAuth.InputID, "INSERT INTO fieldkit.input (expedition_id) VALUES ($1) RETURNING id", ctx.ExpeditionID); err != nil {
		return err
	}

	var err error
	twitterOAuth.RequestToken, twitterOAuth.RequestSecret, err = c.config.RequestToken()
	if err != nil {
		return err
	}

	if _, err := c.options.Database.NamedExecContext(ctx, "INSERT INTO fieldkit.twitter_oauth (input_id, request_token, request_secret) VALUES (:input_id, :request_token, :request_secret) ON CONFLICT (input_id) DO UPDATE SET request_token = :request_token, request_secret = :request_secret", twitterOAuth); err != nil {
		return err
	}

	authorizationURL, err := c.config.AuthorizationURL(twitterOAuth.RequestToken)
	if err != nil {
		return err
	}

	return ctx.OK(&app.Location{
		Location: authorizationURL.String(),
	})
}

func (c *TwitterController) GetID(ctx *app.GetIDTwitterContext) error {
	twitterAccount := &data.TwitterAccount{}
	if err := c.options.Database.GetContext(ctx, twitterAccount, "SELECT i.id, i.expedition_id, ita.twitter_account_id, ta.screen_name FROM fieldkit.twitter_account AS ta JOIN fieldkit.input_twitter_account AS ita ON ita.twitter_account_id = ta.id JOIN fieldkit.input AS i ON i.id = ita.input_id WHERE i.id = $1", ctx.InputID); err != nil {
		return err
	}

	return ctx.OK(TwitterAccountType(twitterAccount))
}

func (c *TwitterController) ListID(ctx *app.ListIDTwitterContext) error {
	twitterAccounts := []*data.TwitterAccount{}
	if err := c.options.Database.SelectContext(ctx, &twitterAccounts, "SELECT i.id, i.expedition_id, ita.twitter_account_id, ta.screen_name, ta.access_token, ta.access_secret FROM fieldkit.twitter_account AS ta JOIN fieldkit.input_twitter_account AS ita ON ita.twitter_account_id = ta.id JOIN fieldkit.input AS i ON i.id = ita.input_id WHERE i.expedition_id = $1", ctx.ExpeditionID); err != nil {
		return err
	}

	return ctx.OK(TwitterAccountsType(twitterAccounts))
}

func (c *TwitterController) List(ctx *app.ListTwitterContext) error {
	twitterAccounts := []*data.TwitterAccount{}
	if err := c.options.Database.SelectContext(ctx, &twitterAccounts, "SELECT i.id, i.expedition_id, ita.twitter_account_id, ta.screen_name, ta.access_token, ta.access_secret FROM fieldkit.twitter_account AS ta JOIN fieldkit.input_twitter_account AS ita ON ita.twitter_account_id = ta.id JOIN fieldkit.input AS i ON i.id = ita.input_id JOIN fieldkit.expedition AS e ON e.id = i.expedition_id JOIN fieldkit.project AS p ON p.id = e.project_id WHERE p.slug = $1 AND e.slug = $2", ctx.Project, ctx.Expedition); err != nil {
		return err
	}

	return ctx.OK(TwitterAccountsType(twitterAccounts))
}

func (c *TwitterController) Callback(ctx *app.CallbackTwitterContext) error {
	requestToken, verifier, err := oauth1.ParseAuthorizationCallback(ctx.RequestData.Request)
	if err != nil {
		return err
	}

	twitterOAuth := &data.TwitterOAuth{}
	if err := c.options.Database.GetContext(ctx, twitterOAuth, "SELECT * FROM fieldkit.twitter_oauth WHERE request_token = $1", requestToken); err != nil {
		return err
	}

	accessToken, accessSecret, err := c.config.AccessToken(requestToken, twitterOAuth.RequestSecret, verifier)
	if err != nil {
		return err
	}

	client := twitter.NewClient(c.config.Client(ctx, oauth1.NewToken(accessToken, accessSecret)))
	user, _, err := client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})
	if err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, "INSERT INTO fieldkit.twitter_account (id, screen_name, access_token, access_secret) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET screen_name = $2, access_token = $3, access_secret = $4", user.ID, user.ScreenName, accessToken, accessSecret); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, "INSERT INTO fieldkit.input_twitter_account (input_id, twitter_account_id) VALUES ($1, $2) ON CONFLICT (input_id) DO UPDATE SET twitter_account_id = $2", twitterOAuth.InputID, user.ID); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, "DELETE FROM fieldkit.twitter_oauth WHERE request_token = $1", requestToken); err != nil {
		return err
	}

	ctx.ResponseData.Header().Set("Location", "https://fieldkit.org/admin")
	return ctx.Found()
}
