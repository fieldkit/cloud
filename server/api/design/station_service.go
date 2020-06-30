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
	Attribute("uploadId", String)
	Attribute("size", Int64)
	Attribute("url", String)
	Attribute("type", String)
	Attribute("blocks", ArrayOf(Int64))
	Required("id", "time", "uploadId", "size", "url", "type", "blocks")
})

var SensorReading = Type("SensorReading", func() {
	Attribute("last", Float32)
	Attribute("time", Int64)
	Required("last", "time")
})

var SensorRange = Type("SensorRange", func() {
	Attribute("minimum", Float32)
	Attribute("maximum", Float32)
	Required("minimum", "maximum")
})

var StationSensor = Type("StationSensor", func() {
	Attribute("name", String)
	Attribute("unitOfMeasure", String)
	Attribute("reading", SensorReading)
	Attribute("key", String)
	Attribute("ranges", ArrayOf(SensorRange))
	Required("name", "unitOfMeasure", "key", "ranges")

})

var StationModule = Type("StationModule", func() {
	Attribute("id", Int64)
	Attribute("hardwareId", String)
	Attribute("metaRecordId", Int64)
	Attribute("name", String)
	Attribute("position", Int32)
	Attribute("flags", Int32)
	Attribute("internal", Boolean)
	Attribute("sensors", ArrayOf(StationSensor))
	Required("id", "name", "position", "flags", "internal", "sensors")
})

var StationConfiguration = Type("StationConfiguration", func() {
	Attribute("id", Int64)
	Attribute("time", Int64)
	Attribute("provisionId", Int64)
	Attribute("metaRecordId", Int64)
	Attribute("sourceId", Int32)
	Attribute("modules", ArrayOf(StationModule))
	Required("id")
	Required("provisionId")
	Required("time")
	Required("modules")
})

var StationConfigurations = Type("StationConfigurations", func() {
	Attribute("all", ArrayOf(StationConfiguration))
	Required("all")
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
		Attribute("deviceId", String)
		Attribute("uploads", ArrayOf(StationUpload))
		Attribute("images", ArrayOf(StationFullImageRef))
		Attribute("photos", StationFullPhotos)
		Attribute("readOnly", Boolean)
		Required("id", "name", "owner", "deviceId", "uploads", "images", "photos", "readOnly")

		Attribute("battery", Float32)
		Attribute("recordingStartedAt", Int64)
		Attribute("memoryUsed", Int32)
		Attribute("memoryAvailable", Int32)
		Attribute("firmwareNumber", Int32)
		Attribute("firmwareTime", Int64)
		Attribute("configurations", StationConfigurations)
		Required("configurations")

		Attribute("updated", Int64)
		Attribute("locationName", String)
		Attribute("placeName", String)
		Attribute("nativeLandName", String)
		Attribute("location", StationLocation)
		Required("updated")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("owner")
		Attribute("deviceId")
		Attribute("uploads")
		Attribute("images")
		Attribute("photos")
		Attribute("readOnly")

		Attribute("battery")
		Attribute("recordingStartedAt")
		Attribute("memoryUsed")
		Attribute("memoryAvailable")
		Attribute("firmwareNumber")
		Attribute("firmwareTime")
		Attribute("configurations")

		Attribute("updated")
		Attribute("location")
		Attribute("locationName")
		Attribute("placeName")
		Attribute("nativeLandName")
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
			Attribute("deviceId", String)
			Attribute("locationName", String)
			Attribute("statusPb", String)
			Required("name", "deviceId")
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
			Attribute("status_pb", String)
			Required("name")
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
