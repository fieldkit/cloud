package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("notifications", func() {
	Method("listen", func() {
		HTTP(func() {
			GET("notifications")
		})

		StreamingResult(Notification)
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
