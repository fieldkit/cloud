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

var SensorReading = Type("SensorReading", func() {
	Attribute("last", Float32)
	Attribute("time", Int64)
	Required("last", "time")
})

var StationSensor = Type("StationSensor", func() {
	Attribute("name", String)
	Attribute("unit_of_measure", String)
	Attribute("reading", SensorReading)
	Required("name", "unit_of_measure")
})

var StationModule = Type("StationModule", func() {
	Attribute("id", Int64)
	Attribute("hardware_id", String)
	Attribute("meta_record_id", Int64)
	Attribute("name", String)
	Attribute("position", Int32)
	Attribute("flags", Int32)
	Attribute("internal", Boolean)
	Attribute("sensors", ArrayOf(StationSensor))
	Required("id", "name", "position", "flags", "internal", "sensors")
})

var StationLocation = Type("StationLocation", func() {
	Attribute("latitude", Float64)
	Attribute("longitude", Float64)
	Required("latitude", "longitude")
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
		Attribute("status_json", MapOf(String, Any))
		Required("id", "name", "owner", "device_id", "uploads", "images", "photos", "read_only", "status_json")

		Attribute("battery", Float32)
		Attribute("recording_started_at", Int64)
		Attribute("memory_used", Int32)
		Attribute("memory_available", Int32)
		Attribute("firmware_number", Int32)
		Attribute("firmware_time", Int64)
		Attribute("modules", ArrayOf(StationModule))
		Required("modules")

		Attribute("updated", Int64)
		Attribute("location_name", String)
		Attribute("location", StationLocation)
		Required("updated")
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
		Attribute("status_json")

		Attribute("battery")
		Attribute("recording_started_at")
		Attribute("memory_used")
		Attribute("memory_available")
		Attribute("firmware_number")
		Attribute("firmware_time")
		Attribute("modules")

		Attribute("updated")
		Attribute("location")
		Attribute("location_name")
	})
})

var StationsFull = ResultType("application/vnd.app.stations.full", func() {
	TypeName("StationsFull")
	Attributes(func() {
		Attribute("stations", CollectionOf(StationFull))
		Required("stations")
	})
	View("default", func() {
		Attribute("stations")
	})
})

var _ = Service("station", func() {
	Method("add", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("name", String)
			Attribute("device_id", String)
			Attribute("location_name", String)
			Attribute("status_json", MapOf(String, Any))
			Attribute("status_pb", String)
			Required("name", "device_id", "status_json")
		})

		Result(StationFull)

		HTTP(func() {
			POST("stations")

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
			Attribute("id", Int32)
			Required("id")
		})

		Result(StationFull)

		HTTP(func() {
			GET("stations/{id}")
			GET("stations/@/{id}")

			httpAuthentication()
		})
	})

	Method("update", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("id", Int32)
			Required("id")
			Attribute("name", String)
			Attribute("location_name", String)
			Attribute("status_json", MapOf(String, Any))
			Attribute("status_pb", String)
			Required("name", "status_json")
		})

		Result(StationFull)

		HTTP(func() {
			PATCH("stations/{id}")

			httpAuthentication()
		})
	})

	Method("list mine", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
		})

		Result(StationsFull)

		HTTP(func() {
			GET("stations")

			httpAuthentication()
		})
	})

	Method("list project", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("id", Int32)
			Required("id")
		})

		Result(StationsFull)

		HTTP(func() {
			GET("projects/{id}/stations")

			httpAuthentication()
		})
	})

	Method("photo", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("id", Int32)
			Required("id")
		})

		Result(func() {
			Attribute("length", Int64)
			Required("length")
			Attribute("content_type", String)
			Required("content_type")
		})

		HTTP(func() {
			GET("stations/{id}/photo")

			SkipResponseBodyEncodeDecode()

			Response(func() {
				Header("length:Content-Length")
				Header("content_type:Content-Type")
			})

			httpAuthenticationQueryString()
		})
	})

	commonOptions()
})
