package api

import (
	"github.com/goadesign/goa"

	"github.com/O-C-R/fieldkit/server/api/app"
	"github.com/O-C-R/fieldkit/server/backend"
	"github.com/O-C-R/fieldkit/server/data"
)

func InputTokenType(inputToken *data.InputToken) *app.InputToken {
	return &app.InputToken{
		ID:           int(inputToken.ID),
		ExpeditionID: int(inputToken.ExpeditionID),
		Token:        inputToken.Token.String(),
	}
}

func InputTokensType(inputTokens []*data.InputToken) *app.InputTokens {
	inputTokensCollection := make([]*app.InputToken, len(inputTokens))
	for i, inputToken := range inputTokens {
		inputTokensCollection[i] = InputTokenType(inputToken)
	}

	return &app.InputTokens{
		InputTokens: inputTokensCollection,
	}
}

type InputTokenControllerOptions struct {
	Backend *backend.Backend
}

// InputTokenController implements the inputToken resource.
type InputTokenController struct {
	*goa.Controller
	options InputTokenControllerOptions
}

func NewInputTokenController(service *goa.Service, options InputTokenControllerOptions) *InputTokenController {
	return &InputTokenController{
		Controller: service.NewController("InputTokenController"),
		options:    options,
	}
}

func (c *InputTokenController) Add(ctx *app.AddInputTokenContext) error {
	token, err := data.NewToken(40)
	if err != nil {
		return err
	}

	inputToken := &data.InputToken{
		ExpeditionID: int32(ctx.ExpeditionID),
		Token:        token,
	}

	if err := c.options.Backend.AddInputToken(ctx, inputToken); err != nil {
		return err
	}

	return ctx.OK(InputTokenType(inputToken))
}

func (c *InputTokenController) Delete(ctx *app.DeleteInputTokenContext) error {
	if err := c.options.Backend.DeleteInputToken(ctx, int32(ctx.InputTokenID)); err != nil {
		return err
	}

	return ctx.NoContent()
}

func (c *InputTokenController) List(ctx *app.ListInputTokenContext) error {
	inputTokens, err := c.options.Backend.ListInputTokens(ctx, ctx.Project, ctx.Expedition)
	if err != nil {
		return err
	}

	return ctx.OK(InputTokensType(inputTokens))
}

func (c *InputTokenController) ListID(ctx *app.ListIDInputTokenContext) error {
	inputTokens, err := c.options.Backend.ListInputTokensID(ctx, int32(ctx.ExpeditionID))
	if err != nil {
		return err
	}

	return ctx.OK(InputTokensType(inputTokens))
}
