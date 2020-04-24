package api

import (
	"github.com/dghubble/oauth1"
	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/data"

	"github.com/dghubble/go-twitter/twitter"
	oauth1Twitter "github.com/dghubble/oauth1/twitter"
)

type TwitterController struct {
	*goa.Controller
	options *ControllerOptions
	config  *oauth1.Config
}

func NewTwitterController(service *goa.Service, options *ControllerOptions) *TwitterController {
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
