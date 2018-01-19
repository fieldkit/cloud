package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("picture", func() {
	Action("user get id", func() {
		Routing(GET("users/:userId/picture"))
		Description("Get a user's picture")
		Params(func() {
			Param("userId", Integer)
			Required("userId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media("image/png")
		})
	})

	Action("project get id", func() {
		Routing(GET("projects/:projectId/picture"))
		Description("Get a project's picture")
		Params(func() {
			Param("projectId", Integer)
			Required("projectId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media("image/png")
		})
	})

	Action("expedition get id", func() {
		Routing(GET("expeditions/:expeditionId/picture"))
		Description("Get a expedition's picture")
		Params(func() {
			Param("expeditionId", Integer)
			Required("expeditionId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media("image/png")
		})
	})
})
