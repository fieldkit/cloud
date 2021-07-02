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

	Error("unauthorized", String, "credentials are invalid")

	HTTP(func() {
		Response("unauthorized", StatusUnauthorized)
	})

	commonOptions()
})
