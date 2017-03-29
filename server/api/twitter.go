package api

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	oauth1Twitter "github.com/dghubble/oauth1/twitter"
	"github.com/goadesign/goa"

	"github.com/O-C-R/fieldkit/server/api/app"
	"github.com/O-C-R/fieldkit/server/backend"
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
	Backend        *backend.Backend
	ConsumerKey    string
	ConsumerSecret string
	Domain         string
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
			CallbackURL:    "https://api." + options.Domain + "/twitter/callback",
			Endpoint:       oauth1Twitter.AuthorizeEndpoint,
		},
	}
}

func (c *TwitterController) Add(ctx *app.AddTwitterContext) error {
	var err error
	twitterOAuth := &data.TwitterOAuth{}
	twitterOAuth.InputID, err = c.options.Backend.AddInputID(ctx, int32(ctx.ExpeditionID))
	if err != nil {
		return err
	}

	twitterOAuth.RequestToken, twitterOAuth.RequestSecret, err = c.config.RequestToken()
	if err != nil {
		return err
	}

	if err := c.options.Backend.AddTwitterOAuth(ctx, twitterOAuth); err != nil {
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
	twitterAccount, err := c.options.Backend.TwitterAccount(ctx, int32(ctx.InputID))
	if err != nil {
		return err
	}

	return ctx.OK(TwitterAccountType(twitterAccount))
}

func (c *TwitterController) ListID(ctx *app.ListIDTwitterContext) error {
	twitterAccounts, err := c.options.Backend.ListTwitterAccountsByID(ctx, int32(ctx.ExpeditionID))
	if err != nil {
		return err
	}

	return ctx.OK(TwitterAccountsType(twitterAccounts))
}

func (c *TwitterController) List(ctx *app.ListTwitterContext) error {
	twitterAccounts, err := c.options.Backend.ListTwitterAccounts(ctx, ctx.Project, ctx.Expedition)
	if err != nil {
		return err
	}

	return ctx.OK(TwitterAccountsType(twitterAccounts))
}

func (c *TwitterController) Callback(ctx *app.CallbackTwitterContext) error {
	requestToken, verifier, err := oauth1.ParseAuthorizationCallback(ctx.RequestData.Request)
	if err != nil {
		return err
	}

	twitterOAuth, err := c.options.Backend.TwitterOAuth(ctx, requestToken)
	if err != nil {
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

	twitterAccount := &data.TwitterAccount{
		TwitterAccountID: user.ID,
		ScreenName:       user.ScreenName,
		AccessToken:      accessToken,
		AccessSecret:     accessSecret,
	}

	if err := c.options.Backend.AddTwitterAccount(ctx, twitterAccount); err != nil {
		return err
	}

	if err := c.options.Backend.DeleteTwitterOAuth(ctx, requestToken); err != nil {
		return err
	}

	ctx.ResponseData.Header().Set("Location", "https://fieldkit.org/admin")
	return ctx.Found()
}
