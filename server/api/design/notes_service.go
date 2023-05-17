package design

import (
	. "goa.design/goa/v3/dsl"
)

var NoteMedia = ResultType("application/vnd.app.note.media", func() {
	TypeName("NoteMedia")
	Attributes(func() {
		Attribute("id", Int64)
		Attribute("url", String)
		Attribute("key", String)
		Attribute("contentType", String)
		Required("id", "url", "key", "contentType")
	})
	View("default", func() {
		Attribute("id")
		Attribute("url")
		Attribute("key")
	})
})

var FieldNoteAuthor = Type("FieldNoteAuthor", func() {
	Attribute("id", Int32)
	Attribute("name", String)
	Attribute("mediaUrl", String)
	Required("id", "name", "mediaUrl")
})

var FieldNote = Type("FieldNote", func() {
	Attribute("id", Int64)
	Attribute("createdAt", Int64)
	Attribute("updatedAt", Int64)
	Attribute("author", FieldNoteAuthor)
	Attribute("key", String)
	Attribute("title", String)
	Attribute("body", String)
	Attribute("version", Int64)
	Attribute("media", ArrayOf(NoteMedia))
	Required("id", "createdAt", "updatedAt", "author", "media", "version")
})

var ExistingFieldNote = Type("ExistingFieldNote", func() {
	Attribute("id", Int64)
	Attribute("key", String)
	Attribute("title", String)
	Attribute("body", String)
	Attribute("mediaIds", ArrayOf(Int64))
	Required("id")
})

var NewFieldNote = Type("NewFieldNote", func() {
	Attribute("key", String)
	Attribute("title", String)
	Attribute("body", String)
	Attribute("mediaIds", ArrayOf(Int64))
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

var FieldNoteStation = ResultType("application/vnd.app.notes.station", func() {
	TypeName("FieldNoteStation")
	Attributes(func() {
		Attribute("readOnly", Boolean)
		Required("readOnly")
	})
	View("default", func() {
		Attribute("readOnly")
	})
})

var StationFieldNotes = ResultType("application/vnd.app.notes", func() {
	TypeName("FieldNotes")
	Attributes(func() {
		Attribute("notes", ArrayOf(FieldNote))
		Required("notes")
		Attribute("media", ArrayOf(NoteMedia))
		Required("media")
		Attribute("station", FieldNoteStation)
		Required("station")
	})
	View("default", func() {
		Attribute("notes")
		Attribute("media")
		Attribute("station")
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
			// Optional
		})

		Payload(func() {
			Token("auth")
			Attribute("stationId", Int32)
			Required("stationId")
		})

		Result(StationFieldNotes)

		HTTP(func() {
			GET("stations/{stationId}/notes")

			httpAuthentication()
		})
	})

	Method("download media", func() {
		Security(JWTAuth, func() {
			// Optional
		})

		Payload(func() {
			Token("auth")
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

			httpAuthentication()
		})
	})

	Method("upload media", func() {
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
			Attribute("stationId", Int32)
			Required("stationId")
			Attribute("key", String)
			Required("key")
		})

		Result(NoteMedia)

		HTTP(func() {
			POST("stations/{stationId}/media")

			Header("contentType:Content-Type")
			Header("contentLength:Content-Length")

			Params(func() {
				Param("key")
			})

			SkipRequestBodyEncodeDecode()

			httpAuthentication()
		})

	})

	Method("delete media", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("mediaId", Int32)
			Required("mediaId")
		})

		HTTP(func() {
			DELETE("notes/media/{mediaId}")

			httpAuthentication()
		})
	})

	commonOptions()
})
