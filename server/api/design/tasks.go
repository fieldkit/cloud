package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("tasks", func() {
	Method("five", func() {
		Response(OK)

		HTTP(func() {
			GET("tasks/five")
		})
	})
})
