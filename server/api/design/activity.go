package design

import (
	. "goa.design/goa/v3/dsl"
)

var StationSummary = Type("StationSummary", func() {
	Attribute("id", Int64)
	Attribute("name", String)
	Required("id", "name")
})

var ProjectSummary = Type("ProjectSummary", func() {
	Attribute("id", Int64)
	Attribute("name", String)
	Required("id", "name")
})

var ActivityEntry = ResultType("application/vnd.app.activity.entry", func() {
	TypeName("ActivityEntry")
	Attributes(func() {
		Attribute("id", Int64)
		Attribute("key", String)
		Attribute("project", ProjectSummary)
		Attribute("station", StationSummary)
		Attribute("created_at", Int64)
		Attribute("type", String)
		Attribute("meta", Any)
		Required("id", "key", "project", "station", "created_at", "type", "meta")
	})
	View("default", func() {
		Attribute("id")
		Attribute("key")
		Attribute("project")
		Attribute("station")
		Attribute("created_at")
		Attribute("type")
		Attribute("meta")
	})
})

var StationActivityPage = ResultType("application/vnd.app.station.activity.page", func() {
	TypeName("StationActivityPage")
	Attributes(func() {
		Attribute("activities", CollectionOf(ActivityEntry))
		Attribute("total", Int32)
		Attribute("page", Int32)
		Required("activities", "total", "page")
	})
})

var ProjectActivityPage = ResultType("application/vnd.app.project.activity.page", func() {
	TypeName("ProjectActivityPage")
	Attributes(func() {
		Attribute("activities", CollectionOf(ActivityEntry))
		Attribute("total", Int32)
		Attribute("page", Int32)
		Required("activities", "total", "page")
	})
})

var _ = Service("activity", func() {
	Method("station", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Attribute("id", Int64)
			Required("id")
			Attribute("page", Int64)
		})

		Result(StationActivityPage)

		HTTP(func() {
			GET("stations/{id}/activity")

			Params(func() {
				Param("page")
			})

			httpAuthentication()
		})
	})

	Method("project", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Attribute("id", Int64)
			Required("id")
			Attribute("page", Int64)
		})

		Result(ProjectActivityPage)

		HTTP(func() {
			GET("projects/{id}/activity")

			Params(func() {
				Param("page")
			})

			httpAuthentication()
		})
	})

	commonOptions()
})
