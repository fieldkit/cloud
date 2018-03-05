package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("Export", func() {
	Action("list by input", func() {
		Routing(GET("inputs/:inputId/csv"))
		Description("Export data for an input.")
		Params(func() {
			Param("inputId", Integer)
			Required("inputId")
		})
		Response(BadRequest)
		Response(OK)
	})
})
