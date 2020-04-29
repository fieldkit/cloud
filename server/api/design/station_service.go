package design

import (
	. "goa.design/goa/v3/dsl"
)

var Owner = Type("StationOwner", func() {
	Attribute("id", Int32)
	Attribute("name", String)
	Required("id", "name")
})

var StationFullImageRef = Type("ImageRef", func() {
	Attribute("url", String)
	Required("url")
})

var StationFullPhotos = Type("StationPhotos", func() {
	Attribute("small", String)
	Required("small")
})

var StationUpload = Type("StationUpload", func() {
	Attribute("id", Int64)
	Attribute("time", Int64)
	Attribute("upload_id", String)
	Attribute("size", Int64)
	Attribute("url", String)
	Attribute("type", String)
	Attribute("blocks", ArrayOf(Int64))
	Required("id", "time", "upload_id", "size", "url", "type", "blocks")
})

var StationFull = ResultType("application/vnd.app.station.full", func() {
	TypeName("StationFull")
	Attributes(func() {
		Attribute("id", Int32)
		Attribute("name")
		Attribute("owner", Owner)
		Attribute("device_id", String)
		Attribute("uploads", ArrayOf(StationUpload))
		Attribute("images", ArrayOf(StationFullImageRef))
		Attribute("photos", StationFullPhotos)
		Attribute("read_only", Boolean)
		Required("id", "name", "owner", "device_id", "uploads", "images", "photos", "read_only")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("owner")
		Attribute("device_id")
		Attribute("uploads")
		Attribute("images")
		Attribute("photos")
		Attribute("read_only")
	})
})

var _ = Service("station", func() {
	Method("station", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("id", Int32)
			Required("id")
		})

		Result(StationFull)

		HTTP(func() {
			GET("stations/@/{id}/details")

			httpAuthentication()
		})
	})

	Error("unauthorized", String, "credentials are invalid")
	Error("not-found", String, "not found")

	HTTP(func() {
		Response("unauthorized", StatusUnauthorized)
		Response("not-found", StatusNotFound)
	})

	commonOptions()
})
