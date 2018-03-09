package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("Export", func() {
	Action("list by source", func() {
		Routing(GET("sources/:sourceId/csv"))
		Description("Export data for an source.")
		Params(func() {
			Param("sourceId", Integer)
			Required("sourceId")
		})
		Response(BadRequest)
		Response(OK)
	})
})
