package design

import (
	. "goa.design/goa/v3/dsl"
)

var EventAuthorPhoto = Type("EventAuthorPhoto", func() {
	Attribute("url", String)
	Required("url")
})

var EventAuthor = Type("EventAuthor", func() {
	Attribute("id", Int32)
	Attribute("name", String)
	Attribute("photo", AuthorPhoto)
	Required("id", "name")
})

var DataEvent = ResultType("application/vnd.app.events.data.event", func() {
	TypeName("DataEvent")
	Attributes(func() {
		Attribute("id", Int64)
		Attribute("createdAt", Int64)
		Attribute("updatedAt", Int64)
		Attribute("author", PostAuthor)
		Attribute("title", String)
		Attribute("description", String)
		Attribute("bookmark", String)
		Attribute("start", Int64)
		Attribute("end", Int64)
		Required("id", "createdAt", "updatedAt", "author", "title", "description", "start", "end")
	})
	View("default", func() {
		Attribute("id")
		Attribute("createdAt")
		Attribute("updatedAt")
		Attribute("author")
		Attribute("title")
		Attribute("description")
		Attribute("bookmark")
		Attribute("start")
		Attribute("end")
	})
})

var DataEvents = ResultType("application/vnd.app.events.data", func() {
	TypeName("DataEvents")
	Attributes(func() {
		Attribute("events", ArrayOf(DataEvent))
		Required("events")
	})
	View("default", func() {
		Attribute("events")
	})
})

var NewDataEvent = ResultType("application/vnd.app.events.data.new", func() {
	TypeName("NewDataEvent")
	Attributes(func() {
		Attribute("projectId", Int32)
		Attribute("bookmark", String)
		Attribute("title", String)
		Attribute("description", String)
		Attribute("start", Int64)
		Attribute("end", Int64)
		Required("title", "description", "start", "end")
	})
	View("default", func() {
		Attribute("projectId")
		Attribute("bookmark")
		Attribute("title")
		Attribute("description")
		Attribute("start")
		Attribute("end")
	})
})

var _ = Service("data events", func() {
	Method("data events", func() {
		Security(JWTAuth, func() {
			// Optional
		})

		Payload(func() {
			Token("auth")
			Attribute("bookmark", String)
			Required("bookmark")
		})

		Result(DataEvents)

		HTTP(func() {
			GET("data-events")

			Params(func() {
				Param("bookmark")
			})

			httpAuthentication()
		})
	})

	Method("add data event", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("event", NewDataEvent)
			Required("event")
		})

		Result(func() {
			Attribute("event", DataEvent)
			Required("event")
		})

		HTTP(func() {
			POST("data-events")

			httpAuthentication()
		})
	})

	Method("update data event", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("eventId", Int64)
			Required("eventId")
			Attribute("title", String)
			Required("title")
			Attribute("description", String)
			Required("description")
			Attribute("start", Int64)
			Required("start")
			Attribute("end", Int64)
			Required("end")
		})

		Result(func() {
			Attribute("event", DataEvent)
			Required("event")
		})

		HTTP(func() {
			POST("data-events/{eventId}")

			httpAuthentication()
		})
	})

	Method("delete data event", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("eventId", Int64)
			Required("eventId")
		})

		HTTP(func() {
			DELETE("data-events/{eventId}")

			httpAuthentication()
		})
	})

	commonOptions()
})
