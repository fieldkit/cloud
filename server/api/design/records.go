package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("records", func() {
	Security(JWT, func() {
		Scope("api:access")
	})

	Action("data", func() {
		Routing(GET("records/data/:recordId"))
		Params(func() {
			Param("recordId", Integer)
		})
		Response(NotFound)
		Response(OK)
	})

	Action("meta", func() {
		Routing(GET("records/meta/:recordId"))
		Params(func() {
			Param("recordId", Integer)
		})
		Response(NotFound)
		Response(OK)
	})

	Action("resolved", func() {
		Routing(GET("records/data/:recordId/resolved"))
		Params(func() {
			Param("recordId", Integer)
		})
		Response(NotFound)
		Response(OK)
	})

	Action("filtered", func() {
		Routing(GET("records/data/:recordId/filtered"))
		Params(func() {
			Param("recordId", Integer)
		})
		Response(NotFound)
		Response(OK)
	})
})
