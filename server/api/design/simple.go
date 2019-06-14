package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var MyDataUrls = MediaType("application/vnd.app.my_data_urls+json", func() {
	TypeName("MyDataUrls")
	Attributes(func() {
		Attribute("csv", String)
		Attribute("fkpb", String)
		Required("csv", "fkpb")
	})
	View("default", func() {
		Attribute("csv")
		Attribute("fkpb")
	})
})

var _ = Resource("simple", func() {
	Security(JWT, func() {
		Scope("api:access")
	})

	Action("my features", func() {
		Routing(GET("my/simple/features"))
		Response(NotFound)
		Response(OK, func() {
			Media(MapFeatures)
		})
	})

	Action("my csv data", func() {
		Routing(GET("my/simple/data/csv"))
		Response(NotFound)
		Response(OK, func() {
			Media(MyDataUrls)
		})
	})
})
