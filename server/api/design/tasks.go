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

	Action("five", func() {
		Routing(GET("tasks/five"))
		Description("Run periodic checks")
		Params(func() {
		})
		Response(BadRequest)
		Response(OK)
	})

	Action("streams/process", func() {
		Routing(GET("tasks/streams/:id/process"))
		Description("Process an uploaded stream")
		Params(func() {
			Param("id", String)
		})
		Response(BadRequest)
		Response(NotFound)
		Response(OK)
	})
})
