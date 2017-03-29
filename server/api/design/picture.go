package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("picture", func() {
	Action("get id", func() {
		Routing(GET("users/:user_id/picture"))
		Description("List a project's inputs")
		Params(func() {
			Param("user_id", Integer)
			Required("user_id")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media("image/png")
		})
	})
})
