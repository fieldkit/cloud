package api

import (
	"github.com/conservify/sqlxcache"
	"github.com/goadesign/goa"

	"github.com/O-C-R/fieldkit/server/api/app"
	"github.com/O-C-R/fieldkit/server/data"
)

type MemberControllerOptions struct {
	Database *sqlxcache.DB
}

func TeamMemberType(member *data.Member) *app.TeamMember {
	return &app.TeamMember{
		TeamID: int(member.TeamID),
		UserID: int(member.UserID),
		Role:   member.Role,
	}
}

func TeamMembersType(members []*data.Member) *app.TeamMembers {
	membersCollection := make([]*app.TeamMember, len(members))
	for i, member := range members {
		membersCollection[i] = TeamMemberType(member)
	}

	return &app.TeamMembers{
		Members: membersCollection,
	}
}

// MemberController implements the user resource.
type MemberController struct {
	*goa.Controller
	options MemberControllerOptions
}

func NewMemberController(service *goa.Service, options MemberControllerOptions) *MemberController {
	return &MemberController{
		Controller: service.NewController("MemberController"),
		options:    options,
	}
}

func (c *MemberController) Add(ctx *app.AddMemberContext) error {
	member := &data.Member{
		TeamID: int32(ctx.TeamID),
		UserID: int32(ctx.Payload.UserID),
		Role:   ctx.Payload.Role,
	}

	if err := c.options.Database.NamedGetContext(ctx, member, "INSERT INTO fieldkit.team_user (team_id, user_id, role) VALUES (:team_id, :user_id, :role) RETURNING *", member); err != nil {
		return err
	}

	return ctx.OK(TeamMemberType(member))
}

func (c *MemberController) Update(ctx *app.UpdateMemberContext) error {
	member := &data.Member{
		TeamID: int32(ctx.TeamID),
		UserID: int32(ctx.UserID),
		Role:   ctx.Payload.Role,
	}

	if err := c.options.Database.NamedGetContext(ctx, member, "UPDATE fieldkit.team_user SET role = :role WHERE team_id = :team_id AND user_id = :user_id RETURNING *", member); err != nil {
		return err
	}

	return ctx.OK(TeamMemberType(member))
}

func (c *MemberController) Delete(ctx *app.DeleteMemberContext) error {
	member := &data.Member{}
	if err := c.options.Database.GetContext(ctx, member, "DELETE FROM fieldkit.team_user WHERE team_id = $1 AND user_id = $2 RETURNING *", int32(ctx.TeamID), int32(ctx.UserID)); err != nil {
		return err
	}

	return ctx.OK(TeamMemberType(member))
}

func (c *MemberController) Get(ctx *app.GetMemberContext) error {
	member := &data.Member{}
	if err := c.options.Database.GetContext(ctx, member, "SELECT tu.* FROM fieldkit.team_user AS tu JOIN fieldkit.team AS t ON t.id = tu.team_id JOIN fieldkit.user AS u ON u.id = tu.user_id JOIN fieldkit.expedition AS e ON e.id = t.expedition_id JOIN fieldkit.project AS p ON p.id = e.project_id WHERE p.slug = $1 AND e.slug = $2 AND t.slug = $3 AND u.username = $4", ctx.Project, ctx.Expedition, ctx.Team, ctx.Username); err != nil {
		return err
	}

	return ctx.OK(TeamMemberType(member))
}

func (c *MemberController) GetID(ctx *app.GetIDMemberContext) error {
	member := &data.Member{}
	if err := c.options.Database.GetContext(ctx, member, "SELECT * FROM fieldkit.team_user WHERE team_id = $1 AND user_id = $2", ctx.TeamID, ctx.UserID); err != nil {
		return err
	}

	return ctx.OK(TeamMemberType(member))
}

func (c *MemberController) List(ctx *app.ListMemberContext) error {
	members := []*data.Member{}
	if err := c.options.Database.SelectContext(ctx, &members, "SELECT tu.* FROM fieldkit.team_user AS tu JOIN fieldkit.team AS t ON t.id = tu.team_id JOIN fieldkit.expedition AS e ON e.id = t.expedition_id JOIN fieldkit.project AS p ON p.id = e.project_id WHERE p.slug = $1 AND e.slug = $2 AND t.slug = $3", ctx.Project, ctx.Expedition, ctx.Team); err != nil {
		return err
	}

	return ctx.OK(TeamMembersType(members))
}

func (c *MemberController) ListID(ctx *app.ListIDMemberContext) error {
	members := []*data.Member{}
	if err := c.options.Database.SelectContext(ctx, &members, "SELECT * FROM fieldkit.team_user WHERE team_id = $1", ctx.TeamID); err != nil {
		return err
	}

	return ctx.OK(TeamMembersType(members))
}
