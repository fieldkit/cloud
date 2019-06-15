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

	Action("refresh", func() {
		Routing(GET("tasks/refresh/:deviceId/:fileTypeId"))
		Description("Refresh a device by ID")
		Params(func() {
			Param("deviceId", String)
			Param("fileTypeId", String)
			Required("deviceId", "fileTypeId")
		})
		Response(BadRequest)
		Response(NotFound)
		Response(OK)
	})
})
