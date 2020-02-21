package api

import (
	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/data"
)

func SourceTokenType(sourceToken *data.SourceToken) *app.SourceToken {
	return &app.SourceToken{
		ID:           int(sourceToken.ID),
		ExpeditionID: int(sourceToken.ExpeditionID),
		Token:        sourceToken.Token.String(),
	}
}

func SourceTokensType(sourceTokens []*data.SourceToken) *app.SourceTokens {
	sourceTokensCollection := make([]*app.SourceToken, len(sourceTokens))
	for i, sourceToken := range sourceTokens {
		sourceTokensCollection[i] = SourceTokenType(sourceToken)
	}

	return &app.SourceTokens{
		SourceTokens: sourceTokensCollection,
	}
}

type SourceTokenController struct {
	*goa.Controller
	options *ControllerOptions
}

func NewSourceTokenController(service *goa.Service, options *ControllerOptions) *SourceTokenController {
	return &SourceTokenController{
		Controller: service.NewController("SourceTokenController"),
		options:    options,
	}
}

func (c *SourceTokenController) Add(ctx *app.AddSourceTokenContext) error {
	token, err := data.NewToken(40)
	if err != nil {
		return err
	}

	sourceToken := &data.SourceToken{
		ExpeditionID: int32(ctx.ExpeditionID),
		Token:        token,
	}

	if err := c.options.Backend.AddSourceToken(ctx, sourceToken); err != nil {
		return err
	}

	return ctx.OK(SourceTokenType(sourceToken))
}

func (c *SourceTokenController) Delete(ctx *app.DeleteSourceTokenContext) error {
	if err := c.options.Backend.DeleteSourceToken(ctx, int32(ctx.SourceTokenID)); err != nil {
		return err
	}

	return ctx.NoContent()
}

func (c *SourceTokenController) List(ctx *app.ListSourceTokenContext) error {
	sourceTokens, err := c.options.Backend.ListSourceTokens(ctx, ctx.Project, ctx.Expedition)
	if err != nil {
		return err
	}

	return ctx.OK(SourceTokensType(sourceTokens))
}

func (c *SourceTokenController) ListID(ctx *app.ListIDSourceTokenContext) error {
	sourceTokens, err := c.options.Backend.ListSourceTokensID(ctx, int32(ctx.ExpeditionID))
	if err != nil {
		return err
	}

	return ctx.OK(SourceTokensType(sourceTokens))
}
