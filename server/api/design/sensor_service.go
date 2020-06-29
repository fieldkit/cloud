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
		})

		Result(func() {
			Attribute("object", Any)
			Required("object")
		})

		HTTP(func() {
			GET("sensors/data")

			Response(func() {
				Body("object")
			})

			httpAuthentication()
		})
	})

	commonOptions()
})
