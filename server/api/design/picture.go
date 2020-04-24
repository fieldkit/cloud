package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var MediaReferenceResponse = MediaType("application/vnd.app.media+json", func() {
	TypeName("MediaReferenceResponse")
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("url", String)
		Required("id")
		Required("url")
	})
	View("default", func() {
		Attribute("id")
		Attribute("url")
	})
})

var _ = Resource("picture", func() {
	Action("user save id", func() {
		Routing(POST("users/:userId/picture"))
		Description("Save a user's picture")
		Params(func() {
			Param("userId", Integer)
			Required("userId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(MediaReferenceResponse)
		})
	})

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
})
