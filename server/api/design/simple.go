package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("simple", func() {
	Action("my features", func() {
		Routing(GET("my/simple/features"))
		Response(NotFound)
		Response(OK, func() {
			Media(GeoJSON)
		})
	})

	Action("my csv data", func() {
		Routing(GET("my/simple/data/csv"))
		Response(NotFound)
		Response(OK, func() {
			Status(200)
		})
	})
})
