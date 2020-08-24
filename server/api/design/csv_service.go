package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("csv", func() {
	Method("noop", func() {
		HTTP(func() {
			POST("csv/noop")
		})
	})

	commonOptions()
})
