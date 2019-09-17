package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("data", func() {
	Action("process", func() {
		Routing(GET("data/process"))
		Description("Process data")
		Response("Busy", func() {
			Status(503)
		})
		Response(OK, func() {
			Status(200)
		})
	})
})
