package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("notifications", func() {
	Method("listen", func() {
		/*
			Security(JWTAuth, func() {
				Scope("api:access")
			})

			Payload(func() {
				Token("token", String, func() {
					Description("JWT used for authentication")
				})

				Required("token")
			})
		*/

		HTTP(func() {
			GET("notifications")
		})

		StreamingResult(MapOf(String, Any))

		StreamingPayload(MapOf(String, Any))
	})

	Error("unauthorized", String, "credentials are invalid")

	HTTP(func() {
		Response("unauthorized", StatusUnauthorized)
	})

	commonOptions()
})

var Notification = ResultType("application/vnd.app.notification", func() {
	TypeName("Notification")
	Attributes(func() {
		Attribute("id", Int64)
		Required("id")
	})
	View("default", func() {
		Attribute("id")
	})
})
