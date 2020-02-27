package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("test", func() {
	Method("get", func() {
		Payload(func() {
			Attribute("id", Int64)
		})

		HTTP(func() {
			GET("test/{id}")
		})
	})

	Method("error", func() {
		HTTP(func() {
			GET("test/error")
		})
	})

	Method("email", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
		})

		Result(Empty)

		HTTP(func() {
			GET("test/email")

			Header("auth:Authorization", String, "authentication token", func() {
				Pattern("^Bearer [^ ]+$")
			})
		})
	})

	Error("unauthorized", String, "credentials are invalid")

	HTTP(func() {
		Response("unauthorized", StatusUnauthorized)
	})
})
