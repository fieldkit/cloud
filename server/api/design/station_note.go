package design

import (
	. "goa.design/goa/v3/dsl"
)

var StationNoteAuthorPhoto = Type("StationNoteAuthorPhoto", func() {
	Attribute("url", String)
	Required("url")
})

var StationNoteAuthor = Type("StationNoteAuthor", func() {
	Attribute("id", Int32)
	Attribute("name", String)
	Attribute("photo", StationNoteAuthorPhoto)
	Required("id", "name")
})

var StationNote = ResultType("application/vnd.app.station.note", func() {
	TypeName("StationNote")
	Attributes(func() {
		Attribute("id", Int32)
		Attribute("createdAt", Int64)
		Attribute("updatedAt", Int64)
		Attribute("author", StationNoteAuthor)
		Attribute("body", String)
		Required("id", "createdAt", "updatedAt", "author", "body")
	})
	View("default", func() {
		Attribute("id")
		Attribute("createdAt")
		Attribute("updatedAt")
		Attribute("author")
		Attribute("body")
	})
})

var StationNotes = ResultType("application/vnd.app.station.notes", func() {
	TypeName("StationNotes")
	Attributes(func() {
		Attribute("notes", ArrayOf(StationNote))
		Required("notes")
	})
	View("default", func() {
		Attribute("notes")
	})
})

var _ = Service("station_note", func() {
	Method("station", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Attribute("stationId", Int32)
			Required("stationId")
		})

		Result(StationNotes)

		HTTP(func() {
			GET("station/{stationId}/station-notes")

			httpAuthentication()
		})
	})

	Method("add note", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("userId", Int32)
			Required("userId")
			Attribute("body", String)
			Required("body")
			Attribute("stationId", Int32)
			Required("stationId")
		})

		Result(StationNote)

		HTTP(func() {
			POST("station/{stationId}/station-note")

			httpAuthentication()
		})
	})

	Method("update note", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("stationId", Int32)
			Required("stationId")
			Attribute("stationNoteId", Int32)
			Required("stationNoteId")
			Attribute("body", String)
			Required("body")
		})

		Result(StationNote)

		HTTP(func() {
			POST("station/{stationId}/station-note/{stationNoteId}")

			httpAuthentication()
		})
	})

	Method("delete note", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("stationId", Int32)
			Required("stationId")
			Attribute("stationNoteId", Int32)
			Required("stationNoteId")
		})

		HTTP(func() {
			DELETE("station/{stationId}/station-note/{stationNoteId}")

			httpAuthentication()
		})
	})

	commonOptions()
})
