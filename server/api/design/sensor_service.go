package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("sensor", func() {
	Method("meta", func() {
		Result(func() {
			Attribute("object", Any)
			Required("object")
		})

		HTTP(func() {
			GET("sensors")

			Response(func() {
				Body("object")
			})
		})
	})

	Method("data", func() {
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
			Attribute("object", Any)
			Required("object")
		})

		HTTP(func() {
			GET("sensors/data")

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
				Body("object")
			})

			httpAuthentication()
		})
	})

	commonOptions()
})
