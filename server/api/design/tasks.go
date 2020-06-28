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

	Error("unauthorized", String, "credentials are invalid")

	HTTP(func() {
		Response("unauthorized", StatusUnauthorized)
	})

	commonOptions()
})
