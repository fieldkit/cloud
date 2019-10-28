package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var AddFieldNotePayload = Type("AddFieldNotePayload", func() {
	Attribute("station_id", Integer)
	Attribute("created", DateTime)
	Attribute("user_id", Integer)
	Attribute("category_id", Integer)
	Attribute("note", String)
	Attribute("media_id", Integer)
	Required("station_id", "created", "user_id")
})

var FieldNoteQueryResult = MediaType("application/vnd.app.field_note_result+json", func() {
	TypeName("FieldNoteQueryResult")
	Reference(AddFieldNotePayload)
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("created")
		Attribute("user_id")
		Attribute("category_key")
		Attribute("note")
		Attribute("media_url")
		Attribute("media_content_type")
		Attribute("username")
		Required("id", "created", "user_id", "category_key", "username")
	})
	View("default", func() {
		Attribute("id")
		Attribute("created")
		Attribute("user_id")
		Attribute("category_key")
		Attribute("note")
		Attribute("media_url")
		Attribute("media_content_type")
		Attribute("username")
	})
})

var FieldNote = MediaType("application/vnd.app.field_note+json", func() {
	TypeName("FieldNote")
	Reference(AddFieldNotePayload)
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("station_id")
		Attribute("created")
		Attribute("user_id")
		Attribute("category_id")
		Attribute("note")
		Attribute("media_id")
		Required("id", "station_id", "created", "user_id", "category_id", "note")
	})
	View("default", func() {
		Attribute("id")
		Attribute("station_id")
		Attribute("created")
		Attribute("user_id")
		Attribute("category_id")
		Attribute("note")
		Attribute("media_id")
	})
})

var FieldNotes = MediaType("application/vnd.app.field_notes+json", func() {
	TypeName("FieldNotes")
	Attributes(func() {
		Attribute("notes", CollectionOf(FieldNoteQueryResult))
		Required("notes")
	})
	View("default", func() {
		Attribute("notes")
	})
})

var _ = Resource("field_note", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("add", func() {
		Routing(POST("/stations/:stationId/field-notes"))
		Description("Add a field note")
		Params(func() {
			Param("stationId", Integer)
			Required("stationId")
		})
		Payload(AddFieldNotePayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(FieldNote)
		})
	})

	Action("update", func() {
		Routing(PATCH("/stations/:stationId/field-notes/:fieldNoteId"))
		Description("Update a field note")
		Params(func() {
			Param("stationId", Integer)
			Param("fieldNoteId", Integer)
			Required("stationId", "fieldNoteId")
		})
		Payload(AddFieldNotePayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(FieldNote)
		})
	})

	Action("get", func() {
		Routing(GET("stations/:stationId/field-notes"))
		Description("Get all field notes for a station")
		Params(func() {
			Param("stationId", Integer)
			Required("stationId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(FieldNotes)
		})
	})

	Action("delete", func() {
		Routing(DELETE("/stations/:stationId/field-notes/:fieldNoteId"))
		Description("Remove a field note")
		Params(func() {
			Param("stationId", Integer)
			Param("fieldNoteId", Integer)
			Required("stationId", "fieldNoteId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Status(204)
		})
	})
})
