package design

import (
	. "goa.design/goa/v3/dsl"
)

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
			GET("sensors/data/export/csv")

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

			Response(StatusFound, func() {
				Headers(func() {
					Header("location:Location")
				})
			})

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

		Error("busy", func() {})

		HTTP(func() {
			GET("sensors/data/export/csv/{id}")

			Response(func() {
				Header("length:Content-Length")
				Header("contentType:Content-Type")
			})

			Response("busy", StatusNotFound)

			SkipResponseBodyEncodeDecode()

			httpAuthentication()
		})
	})

	commonOptions()
})
