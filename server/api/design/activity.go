package design

import (
	. "goa.design/goa/v3/dsl"
)

var StationSummary = Type("StationSummary", func() {
	Attribute("id", Int64)
	Attribute("name", String)
	Required("id", "name")
})

var StationActivity = ResultType("application/vnd.app.station.activity", func() {
	TypeName("StationActivity")
	Attributes(func() {
		Attribute("id", Int64)
		Attribute("station", StationSummary)
		Attribute("created_at", Int64)
		Attribute("type", String)
		Attribute("meta", Any)
		Required("id", "station", "created_at", "type", "meta")
	})
	View("default", func() {
		Attribute("id")
		Attribute("station")
		Attribute("created_at")
		Attribute("type")
		Attribute("meta")
	})
})

var StationActivityPage = ResultType("application/vnd.app.station.activity.page", func() {
	TypeName("StationActivityPage")
	Attributes(func() {
		Attribute("activities", CollectionOf(StationActivity))
		Attribute("total", Int32)
		Attribute("page", Int32)
		Required("activities", "total", "page")
	})
})

var ProjectSummary = Type("ProjectSummary", func() {
	Attribute("id", Int64)
	Attribute("name", String)
	Required("id", "name")
})

var ProjectActivity = ResultType("application/vnd.app.project.activity", func() {
	TypeName("ProjectActivity")
	Attributes(func() {
		Attribute("id", Int64)
		Attribute("project", ProjectSummary)
		Attribute("created_at", Int64)
		Attribute("type", String)
		Attribute("meta", Any)
		Required("id", "project", "created_at", "type", "meta")
	})
	View("default", func() {
		Attribute("id")
		Attribute("project")
		Attribute("created_at")
		Attribute("type")
		Attribute("meta")
	})
})

var ProjectActivityPage = ResultType("application/vnd.app.project.activity.page", func() {
	TypeName("ProjectActivityPage")
	Attributes(func() {
		Attribute("activities", CollectionOf(ProjectActivity))
		Attribute("total", Int32)
		Attribute("page", Int32)
		Required("activities", "total", "page")
	})
})

var _ = Service("activity", func() {
	commonOptions()

	Method("station", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Attribute("id", Int64)
			Attribute("page", Int64)
		})

		Result(StationActivityPage)

		HTTP(func() {
			GET("stations/{id}/activity")

			Params(func() {
				Param("page")
			})

			Header("auth:Authorization", String, "authentication token", func() {
				Pattern("^Bearer [^ ]+$")
			})
		})
	})

	Method("project", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Attribute("id", Int64)
			Attribute("page", Int64)
		})

		Result(ProjectActivityPage)

		HTTP(func() {
			GET("projects/{id}/activity")

			Params(func() {
				Param("page")
			})

			Header("auth:Authorization", String, "authentication token", func() {
				Pattern("^Bearer [^ ]+$")
			})
		})
	})

	HTTP(func() {
	})
})
