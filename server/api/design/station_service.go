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
	Attribute("fullKey", String)
	Attribute("ranges", ArrayOf(SensorRange))
	Required("name", "unitOfMeasure", "key", "fullKey", "ranges")
})

var StationModule = Type("StationModule", func() {
	Attribute("id", Int64)
	Attribute("hardwareId", String)
	Attribute("metaRecordId", Int64)
	Attribute("name", String)
	Attribute("position", Int32)
	Attribute("flags", Int32)
	Attribute("internal", Boolean)
	Attribute("fullKey", String)
	Attribute("sensors", ArrayOf(StationSensor))
	Required("id", "name", "position", "flags", "internal", "fullKey", "sensors")
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
		Attribute("placeNameOther", String)
		Attribute("placeNameNative", String)
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
		Attribute("placeNameOther")
		Attribute("placeNameNative")
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

var EssentialStation = Type("EssentialStation", func() {
	Attribute("id", Int64)
	Attribute("deviceId")
	Attribute("name")
	Attribute("owner", Owner)
	Attribute("createdAt", Int64)
	Attribute("updatedAt", Int64)

	Required("id", "deviceId", "name", "owner", "createdAt", "updatedAt")

	Attribute("recordingStartedAt", Int64)
	Attribute("memoryUsed", Int32)
	Attribute("memoryAvailable", Int32)
	Attribute("firmwareNumber", Int32)
	Attribute("firmwareTime", Int64)
	Attribute("location", StationLocation)
	Attribute("lastIngestionAt", Int64)
})

var PageOfStations = ResultType("application/vnd.app.stations.essential.page", func() {
	TypeName("PageOfStations")
	Attributes(func() {
		Attribute("stations", ArrayOf(EssentialStation))
		Required("stations")
		Attribute("total", Int32)
		Required("total")
	})
	View("default", func() {
		Attribute("stations")
		Attribute("total")
	})
})

var _ = Service("station", func() {
	Error("station-owner-conflict", func() {
	})

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

			Response("station-owner-conflict", StatusBadRequest)

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
			Attribute("locationName", String)
			Attribute("statusPb", String)
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
			GET("user/stations")

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

	Method("download photo", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("stationId", Int32)
			Required("stationId")
			Attribute("size", Int32)
		})

		Result(func() {
			Attribute("length", Int64)
			Required("length")
			Attribute("contentType", String)
			Required("contentType")
		})

		HTTP(func() {
			GET("stations/{stationId}/photo")

			SkipResponseBodyEncodeDecode()

			Params(func() {
				Param("size")
			})

			Response(func() {
				Header("length:Content-Length")
				Header("contentType:Content-Type")
			})

			httpAuthentication()
		})
	})

	Method("list all", func() {
		Security(JWTAuth, func() {
			Scope("api:admin")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("page", Int32)
			Attribute("pageSize", Int32)
			Attribute("ownerId", Int32)
			Attribute("query", String)
			Attribute("sortBy", String)
		})

		Result(PageOfStations)

		HTTP(func() {
			GET("admin/stations")

			Params(func() {
				Param("page")
				Param("pageSize")
				Param("ownerId", Int32)
				Param("query", String)
				Param("sortBy", String)
			})

			httpAuthentication()
		})
	})

	Method("delete", func() {
		Security(JWTAuth, func() {
			Scope("api:admin")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("stationId", Int32)
			Required("stationId")
		})

		HTTP(func() {
			DELETE("admin/stations/{stationId}")

			httpAuthentication()
		})
	})

	commonOptions()
})
