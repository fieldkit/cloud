package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("picture", func() {
	Action("user get id", func() {
		Routing(GET("users/:user_id/picture"))
		Description("Get a user's picture")
		Params(func() {
			Param("user_id", Integer)
			Required("user_id")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media("image/png")
		})
	})

	Action("project get id", func() {
		Routing(GET("projects/:project_id/picture"))
		Description("Get a project's picture")
		Params(func() {
			Param("project_id", Integer)
			Required("project_id")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media("image/png")
		})
	})

	Action("expedition get id", func() {
		Routing(GET("expeditions/:expedition_id/picture"))
		Description("Get a expedition's picture")
		Params(func() {
			Param("expedition_id", Integer)
			Required("expedition_id")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media("image/png")
		})
	})
})
