package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("tasks", func() {
	Method("five", func() {
		HTTP(func() {
			GET("tasks/five")
		})
	})

	Method("refresh device", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Attribute("deviceId", String)
			Required("auth", "deviceId")
		})

		Result(Empty)

		HTTP(func() {
			POST("tasks/devices/{deviceId}/refresh")

			Header("auth:Authorization", String, "authentication token", func() {
				Pattern("^Bearer [^ ]+$")
			})
		})
	})

	Error("unauthorized", String, "credentials are invalid")
	Error("invalid-scopes", String, "token scopes are invalid")

	HTTP(func() {
		Response("unauthorized", StatusUnauthorized)
	})

	commonOptions()
})
