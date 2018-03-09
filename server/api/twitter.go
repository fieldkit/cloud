package api

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	oauth1Twitter "github.com/dghubble/oauth1/twitter"
	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"
)

func TwitterAccountSourceType(twitterAccount *data.TwitterAccountSource) *app.TwitterAccountSource {
	twitterAccountType := &app.TwitterAccountSource{
		ID:               int(twitterAccount.ID),
		ExpeditionID:     int(twitterAccount.ExpeditionID),
		Name:             twitterAccount.Name,
		TwitterAccountID: int(twitterAccount.TwitterAccountID),
		ScreenName:       twitterAccount.ScreenName,
	}

	if twitterAccount.TeamID != nil {
		teamID := int(*twitterAccount.TeamID)
		twitterAccountType.TeamID = &teamID
	}

	if twitterAccount.UserID != nil {
		userID := int(*twitterAccount.UserID)
		twitterAccountType.UserID = &userID
	}

	return twitterAccountType
}

func TwitterAccountSourcesType(twitterAccounts []*data.TwitterAccountSource) *app.TwitterAccountSources {
	twitterAccountsCollection := make([]*app.TwitterAccountSource, len(twitterAccounts))
	for i, twitterAccount := range twitterAccounts {
		twitterAccountsCollection[i] = TwitterAccountSourceType(twitterAccount)
	}

	return &app.TwitterAccountSources{
		TwitterAccountSources: twitterAccountsCollection,
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
	source := &data.Source{}
	source.ExpeditionID = int32(ctx.ExpeditionID)
	source.Name = ctx.Payload.Name
	if err := c.options.Backend.AddSource(ctx, source); err != nil {
		return err
	}

	twitterOAuth := &data.TwitterOAuth{
		SourceID: source.ID,
	}

	var err error
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
	twitterAccount, err := c.options.Backend.TwitterAccountSource(ctx, int32(ctx.SourceID))
	if err != nil {
		return err
	}

	return ctx.OK(TwitterAccountSourceType(twitterAccount))
}

func (c *TwitterController) ListID(ctx *app.ListIDTwitterContext) error {
	twitterAccounts, err := c.options.Backend.ListTwitterAccountSourcesByExpeditionID(ctx, int32(ctx.ExpeditionID))
	if err != nil {
		return err
	}

	return ctx.OK(TwitterAccountSourcesType(twitterAccounts))
}

func (c *TwitterController) List(ctx *app.ListTwitterContext) error {
	twitterAccounts, err := c.options.Backend.ListTwitterAccountSources(ctx, ctx.Project, ctx.Expedition)
	if err != nil {
		return err
	}

	return ctx.OK(TwitterAccountSourcesType(twitterAccounts))
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

	twitterAccountSource := &data.TwitterAccountSource{}
	twitterAccountSource.ID = twitterOAuth.SourceID
	twitterAccountSource.AccessToken, twitterAccountSource.AccessSecret, err = c.config.AccessToken(requestToken, twitterOAuth.RequestSecret, verifier)
	if err != nil {
		return err
	}

	client := twitter.NewClient(c.config.Client(ctx, oauth1.NewToken(twitterAccountSource.AccessToken, twitterAccountSource.AccessSecret)))
	user, _, err := client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})
	if err != nil {
		return err
	}

	twitterAccountSource.TwitterAccountID = user.ID
	twitterAccountSource.ScreenName = user.ScreenName
	if err := c.options.Backend.AddTwitterAccountSource(ctx, twitterAccountSource); err != nil {
		return err
	}

	if err := c.options.Backend.DeleteTwitterOAuth(ctx, requestToken); err != nil {
		return err
	}

	ctx.ResponseData.Header().Set("Location", "https://fieldkit.org/admin")
	return ctx.Found()
}
