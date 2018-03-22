package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("Tasks", func() {
	Action("check", func() {
		Routing(GET("tasks/check"))
		Description("Run periodic checks")
		Params(func() {
		})
		Response(BadRequest)
		Response(OK)
	})
})
