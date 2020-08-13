package design

import (
	. "goa.design/goa/v3/dsl"
)

var ExportStatus = ResultType("application/vnd.app.export.status", func() {
	TypeName("ExportStatus")
	Attributes(func() {
		Attribute("id", Int64)
		Attribute("token", String)
		Attribute("createdAt", Int64)
		Attribute("completedAt", Int64)
		Attribute("kind", String)
		Attribute("progress", Float32)
		Attribute("statusUrl", String)
		Attribute("downloadUrl", String)
		Attribute("size", Int32)
		Attribute("args", Any)
		Required("id", "token", "createdAt", "kind", "statusUrl", "progress", "args")
	})
	View("default", func() {
		Attribute("id")
		Attribute("token")
		Attribute("createdAt")
		Attribute("completedAt")
		Attribute("kind")
		Attribute("progress")
		Attribute("statusUrl")
		Attribute("downloadUrl")
		Attribute("size")
		Attribute("args")
	})
})

var UserExports = ResultType("application/vnd.app.exports", func() {
	TypeName("UserExports")
	Attributes(func() {
		Attribute("exports", ArrayOf(ExportStatus))
		Required("exports")
	})
	View("default", func() {
		Attribute("exports")
	})
})

var _ = Service("export", func() {
	Method("list mine", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
		})

		Result(UserExports)

		HTTP(func() {
			GET("export")

			httpAuthentication()
		})
	})

	Method("status", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("id", String)
			Required("id")
		})

		Result(ExportStatus)

		HTTP(func() {
			GET("export/{id}")

			httpAuthentication()
		})
	})

	Method("download", func() {
		Payload(func() {
			Attribute("id", String)
			Required("id")
			Attribute("auth", String)
			Required("auth")
		})

		Result(func() {
			Attribute("length", Int64)
			Required("length")
			Attribute("contentType", String)
			Required("contentType")
		})

		HTTP(func() {
			GET("export/{id}/download")

			Params(func() {
				Param("auth")
			})

			Response(func() {
				Header("length:Content-Length")
				Header("contentType:Content-Type")
			})

			SkipResponseBodyEncodeDecode()
		})
	})

	commonOptions()
})
