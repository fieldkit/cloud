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
		Attribute("format", String)
		Attribute("progress", Float32)
		Attribute("message", String)
		Attribute("statusUrl", String)
		Attribute("downloadUrl", String)
		Attribute("size", Int32)
		Attribute("args", Any)
		Required("id", "token", "createdAt", "format", "statusUrl", "progress", "args")
	})
	View("default", func() {
		Attribute("id")
		Attribute("token")
		Attribute("createdAt")
		Attribute("completedAt")
		Attribute("format")
		Attribute("progress")
		Attribute("message")
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
			Attribute("contentDisposition", String)
			Required("contentDisposition")
		})

		HTTP(func() {
			GET("export/{id}/download")

			Params(func() {
				Param("auth")
			})

			Response(func() {
				Header("length:Content-Length")
				Header("contentType:Content-Type")
				Header("contentDisposition:Content-Disposition")
			})

			SkipResponseBodyEncodeDecode()
		})
	})

	exportPayload := func() {
		Token("auth")
		Required("auth")
		Attribute("start", Int64)
		Attribute("end", Int64)
		Attribute("stations", String)
		Attribute("sensors", String)
		Attribute("resolution", Int32)
		Attribute("aggregate", String)
		Attribute("complete", Boolean)
		Attribute("tail", Int32)
	}

	exportParams := func() {
		Param("start")
		Param("end")
		Param("stations")
		Param("sensors")
		Param("resolution")
		Param("aggregate")
		Param("complete")
		Param("tail")
	}

	Method("csv", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(exportPayload)

		Result(func() {
			Attribute("location", String)
			Required("location")
		})

		HTTP(func() {
			POST("export/csv")

			Params(exportParams)

			Response(StatusFound, func() {
				Headers(func() {
					Header("location:Location")
				})
			})

			httpAuthentication()
		})
	})

	Method("json lines", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(exportPayload)

		Result(func() {
			Attribute("location", String)
			Required("location")
		})

		HTTP(func() {
			POST("export/json-lines")

			Params(exportParams)

			Response(StatusFound, func() {
				Headers(func() {
					Header("location:Location")
				})
			})

			httpAuthentication()
		})
	})

	commonOptions()
})
