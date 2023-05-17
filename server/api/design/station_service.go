package design

import (
	. "goa.design/goa/v3/dsl"
)

var Owner = Type("StationOwner", func() {
	Attribute("id", Int32)
	Attribute("name", String)
	Attribute("email", String)
	Attribute("photo", UserPhoto)
	Required("id", "name")
})

var StationFullModel = Type("StationFullModel", func() {
	Attribute("name", String)
	Attribute("only_visible_via_association", Boolean)
	Required("name", "only_visible_via_association")
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
	Attribute("meta", MapOf(String, Any))
	Required("name", "unitOfMeasure", "key", "fullKey", "ranges")
})

var StationModule = Type("StationModule", func() {
	Attribute("id", Int64)
	Attribute("hardwareId", String)
	Attribute("hardwareIdBase64", String)
	Attribute("metaRecordId", Int64)
	Attribute("name", String)
	Attribute("label", String)
	Attribute("position", Int32)
	Attribute("flags", Int32)
	Attribute("internal", Boolean)
	Attribute("fullKey", String)
	Attribute("sensors", ArrayOf(StationSensor))
	Attribute("meta", MapOf(String, Any))
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

var StationRegion = Type("StationRegion", func() {
	Attribute("name", String)
	Attribute("shape", ArrayOf(ArrayOf(ArrayOf(Float64))))
	Required("name", "shape")
})

var StationLocation = ResultType("application/vnd.app.station.location", func() {
	TypeName("StationLocation")
	Attributes(func() {
		Attribute("precise", ArrayOf(Float64))
		Attribute("regions", ArrayOf(StationRegion))
		Attribute("url", String)
	})
	View("default", func() {
		Attribute("precise")
		Attribute("regions")
		Attribute("url")
	})
})

var StationDataSummary = Type("StationDataSummary", func() {
	Attribute("start", Int64)
	Attribute("end", Int64)
	Attribute("numberOfSamples", Int64)
	Required("start", "end", "numberOfSamples")
})

var StationInterestingnessWindow = Type("StationInterestingnessWindow", func() {
	Attribute("seconds", Int32)
	Attribute("interestingness", Float64)
	Attribute("value", Float64)
	Attribute("time", Int64)
	Required("seconds", "interestingness", "value", "time")
})

var StationInterestingness = Type("StationInterestingness", func() {
	Attribute("windows", ArrayOf(StationInterestingnessWindow))
	Required("windows")
})

var StationProjectAttribute = Type("StationProjectAttribute", func() {
	Attribute("projectId", Int32)
	Attribute("attributeId", Int64)
	Attribute("name", String)
	Attribute("stringValue", String)
	Attribute("priority", Int32)
	Required("projectId", "attributeId", "name", "stringValue", "priority")
})

var StationProjectAttributes = Type("StationProjectAttributes", func() {
	Attribute("attributes", ArrayOf(StationProjectAttribute))
	Required("attributes")
})

var StationFull = ResultType("application/vnd.app.station.full", func() {
	TypeName("StationFull")
	Attributes(func() {
		Attribute("id", Int32)
		Attribute("name")
		Attribute("model", StationFullModel)
		Attribute("owner", Owner)
		Attribute("deviceId", String)
		Attribute("interestingness", StationInterestingness)
		Attribute("attributes", StationProjectAttributes)
		Attribute("uploads", ArrayOf(StationUpload))
		Attribute("photos", StationFullPhotos)
		Attribute("readOnly", Boolean)
		Attribute("status", String)
		Attribute("hidden", Boolean)
		Attribute("description", String)
		Required("id", "name", "model", "owner", "deviceId", "interestingness", "attributes", "uploads", "photos", "readOnly")

		Attribute("battery", Float32)
		Attribute("recordingStartedAt", Int64)
		Attribute("memoryUsed", Int32)
		Attribute("memoryAvailable", Int32)
		Attribute("firmwareNumber", Int32)
		Attribute("firmwareTime", Int64)
		Attribute("configurations", StationConfigurations)
		Required("configurations")

		Attribute("updatedAt", Int64)
		Attribute("lastReadingAt", Int64)
		Attribute("locationName", String)
		Attribute("placeNameOther", String)
		Attribute("placeNameNative", String)
		Attribute("location", StationLocation)
		Required("updatedAt")

		Attribute("syncedAt", Int64)
		Attribute("ingestionAt", Int64)

		Attribute("data", StationDataSummary)
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("model")
		Attribute("owner")
		Attribute("deviceId")
		Attribute("interestingness")
		Attribute("attributes")
		Attribute("uploads")
		Attribute("photos")
		Attribute("readOnly")
		Attribute("hidden")
		Attribute("description")
		Attribute("status")

		Attribute("battery")
		Attribute("recordingStartedAt")
		Attribute("memoryUsed")
		Attribute("memoryAvailable")
		Attribute("firmwareNumber")
		Attribute("firmwareTime")
		Attribute("configurations")

		Attribute("updatedAt")
		Attribute("lastReadingAt")
		Attribute("location")
		Attribute("locationName")
		Attribute("placeNameOther")
		Attribute("placeNameNative")

		Attribute("syncedAt", Int64)
		Attribute("ingestionAt", Int64)

		Attribute("data")
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

var AssociatedViaProject = Type("AssociatedViaProject", func() {
	Attribute("id", Int32)
	Required("id")
})

var AssociatedViaLocation = Type("AssociatedViaLocation", func() {
	Attribute("stationID", Int32)
	Attribute("distance", Float32)
	Required("stationID", "distance")
})

var AssociatedViaManual = Type("AssociatedViaManual", func() {
	Attribute("otherStationID", Int32)
	Attribute("priority", Int32)
	Required("otherStationID", "priority")
})

var AssociatedStation = ResultType("application/vnd.app.associated.station", func() {
	TypeName("AssociatedStation")
	Attributes(func() {
		Attribute("station", StationFull)
		Required("station")
		Attribute("project", ArrayOf(AssociatedViaProject))
		Attribute("location", ArrayOf(AssociatedViaLocation))
		Attribute("manual", ArrayOf(AssociatedViaManual))
		Attribute("hidden", Boolean)
		Required("hidden")
	})
	View("default", func() {
		Attribute("station")
		Attribute("project")
		Attribute("location")
		Attribute("manual")
		Attribute("hidden")
	})
})

var AssociatedStations = ResultType("application/vnd.app.associated.stations", func() {
	TypeName("AssociatedStations")
	Attributes(func() {
		Attribute("stations", CollectionOf(AssociatedStation))
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

var StationJob = ResultType("application/vnd.app.stations.job", func() {
	TypeName("StationJob")
	Attributes(func() {
		Attribute("title", String)
		Attribute("startedAt", Int64)
		Attribute("completedAt", Int64)
		Attribute("progress", Float32)
		Required("startedAt")
		Required("progress")
		Required("title")
	})
	View("default", func() {
		Attribute("startedAt")
		Attribute("completedAt")
		Attribute("progress")
		Attribute("title")
	})
})

var StationProgress = ResultType("application/vnd.app.stations.progress", func() {
	TypeName("StationProgress")
	Attributes(func() {
		Attribute("jobs", ArrayOf(StationJob))
		Required("jobs")
	})
	View("default", func() {
		Attribute("jobs")
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
			Attribute("description", String)
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
			// Optional
		})

		Payload(func() {
			Token("auth")
			Attribute("id", Int32)
			Required("id")
		})

		Result(StationFull)

		HTTP(func() {
			GET("stations/{id}")

			httpAuthentication()
		})
	})

	Method("transfer", func() {
		Security(JWTAuth, func() {
			Scope("api:admin")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("id", Int32)
			Required("id")
			Attribute("ownerId", Int32)
			Required("ownerId")
		})

		HTTP(func() {
			POST("stations/{id}/transfer/{ownerId}")

			httpAuthentication()
		})
	})

	Method("default photo", func() {
		Security(JWTAuth, func() {
			// Optional
		})

		Payload(func() {
			Token("auth")
			Attribute("id", Int32)
			Required("id")
			Attribute("photoId", Int32)
			Required("photoId")
		})

		HTTP(func() {
			POST("stations/{id}/photo/{photoId}")

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
			Attribute("description", String)
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
			// Optional
		})

		Payload(func() {
			Token("auth")
			Attribute("id", Int32)
			Required("id")
			Attribute("disable_filtering", Boolean)
		})

		Result(StationsFull)

		HTTP(func() {
			GET("projects/{id}/stations")

			Params(func() {
				Param("disable_filtering")
			})

			httpAuthentication()
		})
	})

	Method("list associated", func() {
		Security(JWTAuth, func() {
			// Optional
		})

		Payload(func() {
			Token("auth")
			Attribute("id", Int32)
			Required("id")
		})

		Result(AssociatedStations)

		HTTP(func() {
			GET("stations/{id}/associated")

			httpAuthentication()
		})
	})

	Method("list project associated", func() {
		Security(JWTAuth, func() {
			// Optional
		})

		Payload(func() {
			Token("auth")
			Attribute("projectId", Int32)
			Required("projectId")
		})

		Result(AssociatedStations)

		HTTP(func() {
			GET("projects/{projectId}/associated")

			httpAuthentication()
		})
	})

	Method("download photo", func() {
		Security(JWTAuth, func() {
			// Optional
		})

		Payload(func() {
			Token("auth")
			Attribute("stationId", Int32)
			Required("stationId")
			Attribute("size", Int32)
			Attribute("ifNoneMatch", String)
		})

		Result(DownloadedPhoto)

		HTTP(func() {
			GET("stations/{stationId}/photo")

			Header("ifNoneMatch:If-None-Match")

			Params(func() {
				Param("size")
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

	Method("admin search", func() {
		Security(JWTAuth, func() {
			Scope("api:admin")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("query", String)
			Required("query")
		})

		Result(PageOfStations)

		HTTP(func() {
			POST("admin/stations/search")

			Params(func() {
				Param("query")
			})

			httpAuthentication()
		})
	})

	Method("progress", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("stationId", Int32)
			Required("stationId")
		})

		Result(StationProgress)

		HTTP(func() {
			GET("stations/{stationId}/progress")

			httpAuthentication()
		})
	})

	Method("update module", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("id", Int32)
			Required("id")
			Attribute("moduleId", Int32)
			Required("moduleId")
			Attribute("label", String)
			Required("label")
		})

		Result(StationFull)

		HTTP(func() {
			PATCH("stations/{id}/modules/{moduleId}")

			httpAuthentication()
		})
	})

	commonOptions()
})
