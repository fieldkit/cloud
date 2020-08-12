package design

import (
	. "goa.design/goa/v3/dsl"
)

var ExportStatus = ResultType("application/vnd.app.export.status", func() {
	TypeName("ExportStatus")
	Attributes(func() {
		Attribute("id", Int64)
		Attribute("createdAt", Int64)
		Attribute("completedAt", Int64)
		Attribute("progress", Float32)
		Attribute("url", String)
		Attribute("kind", String)
		Attribute("args", Any)
		Required("id", "createdAt", "progress", "kind", "args")
	})
	View("default", func() {
		Attribute("id")
		Attribute("createdAt")
		Attribute("completedAt")
		Attribute("kind")
		Attribute("progress")
		Attribute("url")
		Attribute("args")
	})
})

var _ = Service("csv", func() {
	Method("export", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
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
		})

		Result(func() {
			Attribute("location", String)
			Required("location")
		})

		HTTP(func() {
			POST("export/csv")

			Params(func() {
				Param("start")
				Param("end")
				Param("stations")
				Param("sensors")
				Param("resolution")
				Param("aggregate")
				Param("complete")
				Param("tail")
			})

			Response(func() {
				Body(func() {
					Attribute("location")
				})
			})

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
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("id", String)
			Required("id")
		})

		Result(func() {
			Attribute("length", Int64)
			Required("length")
			Attribute("contentType", String)
			Required("contentType")
		})

		HTTP(func() {
			GET("export/{id}/download")

			Response(func() {
				Header("length:Content-Length")
				Header("contentType:Content-Type")
			})

			SkipResponseBodyEncodeDecode()

			httpAuthentication()
		})
	})

	commonOptions()
})
