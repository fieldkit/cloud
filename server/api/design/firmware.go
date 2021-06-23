package design

import (
	. "goa.design/goa/v3/dsl"
)

var AddFirmwarePayload = Type("AddFirmwarePayload", func() {
	Attribute("etag")
	Required("etag")
	Attribute("module")
	Required("module")
	Attribute("profile")
	Required("profile")
	Attribute("version")
	Required("version")
	Attribute("url")
	Required("url")
	Attribute("meta")
	Required("meta")
	Attribute("logicalAddress", Int64)
})

var FirmwareSummary = ResultType("application/vnd.app.firmware+json", func() {
	TypeName("FirmwareSummary")
	Attributes(func() {
		Attribute("id", Int32)
		Attribute("time", String)
		Attribute("etag", String)
		Attribute("module", String)
		Attribute("profile", String)
		Attribute("version", String)
		Attribute("url", String)
		Attribute("meta", MapOf(String, Any))
		Attribute("buildNumber", Int32)
		Attribute("buildTime", Int64)
		Attribute("logicalAddress", Int64)
		Required("id")
		Required("time")
		Required("etag")
		Required("module")
		Required("profile")
		Required("url")
		Required("meta")
		Required("buildNumber")
		Required("buildTime")
	})
	View("default", func() {
		Attribute("id")
		Attribute("time")
		Attribute("etag")
		Attribute("module")
		Attribute("profile")
		Attribute("version")
		Attribute("url")
		Attribute("meta")
		Attribute("buildNumber")
		Attribute("buildTime")
		Attribute("logicalAddress")
	})
})

var Firmwares = ResultType("application/vnd.app.firmwares+json", func() {
	TypeName("Firmwares")
	Attributes(func() {
		Attribute("firmwares", CollectionOf(FirmwareSummary))
		Required("firmwares")
	})
	View("default", func() {
		Attribute("firmwares")
	})
})

var _ = Service("firmware", func() {
	Method("download", func() {
		Payload(func() {
			Attribute("firmwareId", Int32)
			Required("firmwareId")
		})

		Result(func() {
			Attribute("length", Int64)
			Required("length")
			Attribute("contentType", String)
			Required("contentType")
		})

		HTTP(func() {
			GET("firmware/{firmwareId}/download")

			SkipResponseBodyEncodeDecode()

			Response(func() {
				Header("length:Content-Length")
				Header("contentType:Content-Type")
			})
		})
	})

	Method("add", func() {
		Security(JWTAuth, func() {
			Scope("api:admin")
		})

		Payload(func() {
			Token("auth")
			Attribute("firmware", AddFirmwarePayload)
			Required("firmware")
		})

		HTTP(func() {
			POST("firmware")

			Body("firmware")
		})
	})

	Method("list", func() {
		Security(JWTAuth, func() {
		})

		Payload(func() {
			Token("auth")
			Attribute("module", String)
			Attribute("profile", String)
			Attribute("pageSize", Int32)
			Attribute("page", Int32)
		})

		Result(Firmwares)

		HTTP(func() {
			GET("firmware")

			Params(func() {
				Param("module")
				Param("profile")
				Param("pageSize")
				Param("page")
			})
		})
	})

	Method("delete", func() {
		Security(JWTAuth, func() {
			Scope("api:admin")
		})

		Payload(func() {
			Token("auth")
			Attribute("firmwareId", Int32)
			Required("firmwareId")
		})

		HTTP(func() {
			DELETE("firmware/{firmwareId}")

			Response(StatusOK)
		})
	})

	commonOptions()
})
