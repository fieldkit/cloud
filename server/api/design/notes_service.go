package design

import (
	. "goa.design/goa/v3/dsl"
)

var NoteMedia = ResultType("application/vnd.app.note.media", func() {
	TypeName("NoteMedia")
	Attributes(func() {
		Attribute("id", Int64)
		Attribute("url", String)
		Required("id", "url")
	})
	View("default", func() {
		Attribute("id")
		Attribute("url")
	})
})

var FieldNote = Type("FieldNote", func() {
	Attribute("id", Int64)
	Attribute("createdAt", Int64)
	Attribute("key", String)
	Attribute("body", String)
	Attribute("mediaId", Int64)
	Required("id", "createdAt")
})

var ExistingFieldNote = Type("ExistingFieldNote", func() {
	Attribute("id", Int64)
	Attribute("key", String)
	Attribute("body", String)
	Attribute("mediaId", Int64)
	Required("id")
})

var NewFieldNote = Type("NewFieldNote", func() {
	Attribute("key", String)
	Attribute("body", String)
	Attribute("mediaId", Int64)
})

var FieldNoteUpdate = ResultType("application/vnd.app.notes.update", func() {
	TypeName("FieldNoteUpdate")
	Attributes(func() {
		Attribute("notes", ArrayOf(ExistingFieldNote))
		Attribute("creating", ArrayOf(NewFieldNote))
		Required("notes")
		Required("creating")
	})
	View("default", func() {
		Attribute("notes")
		Attribute("creating")
	})
})

var StationFieldNotes = ResultType("application/vnd.app.notes", func() {
	TypeName("FieldNotes")
	Attributes(func() {
		Attribute("notes", ArrayOf(FieldNote))
		Required("notes")
	})
	View("default", func() {
		Attribute("notes")
	})
})

var _ = Service("notes", func() {
	Method("update", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("stationId", Int32)
			Required("stationId")
			Attribute("notes", FieldNoteUpdate)
			Required("notes")
		})

		Result(StationFieldNotes)

		HTTP(func() {
			PATCH("stations/{stationId}/notes")

			httpAuthentication()
		})
	})

	Method("get", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("stationId", Int32)
			Required("stationId")
		})

		Result(StationFieldNotes)

		HTTP(func() {
			GET("stations/{stationId}/notes")

			httpAuthentication()
		})
	})

	Method("media", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("mediaId", Int32)
			Required("mediaId")
		})

		Result(func() {
			Attribute("length", Int64)
			Required("length")
			Attribute("contentType", String)
			Required("contentType")
		})

		HTTP(func() {
			GET("notes/media/{mediaId}")

			SkipResponseBodyEncodeDecode()

			Response(func() {
				Header("length:Content-Length")
				Header("contentType:Content-Type")
			})

			httpAuthenticationQueryString()
		})
	})

	Method("upload", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("contentLength", Int64)
			Required("contentLength")
			Attribute("contentType", String)
			Required("contentType")
		})

		Result(NoteMedia)

		HTTP(func() {
			POST("notes/media")

			Header("contentType:Content-Type")
			Header("contentLength:Content-Length")

			SkipRequestBodyEncodeDecode()

			httpAuthenticationQueryString()
		})
	})

	commonOptions()
})
