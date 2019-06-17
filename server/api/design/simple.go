package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var MyDataUrls = Type("MyDataUrls", func() {
	Attribute("csv", String)
	Attribute("fkpb", String)
	Required("csv", "fkpb")
})

var MySimpleSummary = MediaType("application/vnd.app.simple_summary+json", func() {
	TypeName("MySimpleSummary")
	Attributes(func() {
		Attribute("urls", MyDataUrls)
		Attribute("center", ArrayOf(Number))
		Required("urls", "center")
	})
	View("default", func() {
		Attribute("urls")
		Attribute("center")
	})
})

var _ = Resource("simple", func() {
	Security(JWT, func() {
		Scope("api:access")
	})

	Action("my simple summary", func() {
		Routing(GET("my/simple"))
		Response(NotFound)
		Response(OK, func() {
			Media(MySimpleSummary)
		})
	})

	Action("my features", func() {
		Routing(GET("my/simple/features"))
		Response(NotFound)
		Response(OK, func() {
			Media(MapFeatures)
		})
	})

	Action("download", func() {
		NoSecurity()
		Routing(GET("my/simple/download"))
		Params(func() {
			Param("token", String)
			Required("token")
		})
		Response(NotFound)
		Response(OK)
	})
})
