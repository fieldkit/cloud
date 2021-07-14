package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("notifications", func() {
	Method("listen", func() {
		HTTP(func() {
			GET("notifications")
		})

		StreamingResult(MapOf(String, Any))

		StreamingPayload(MapOf(String, Any))
	})

	Method("seen", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("ids", ArrayOf(Int64))
			Required("ids")
		})

		HTTP(func() {
			POST("notifications/seen")

			httpAuthentication()
		})
	})

	Error("unauthorized", String, "credentials are invalid")

	HTTP(func() {
		Response("unauthorized", StatusUnauthorized)
	})

	commonOptions()
})
