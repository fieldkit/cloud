package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("modules", func() {
	Method("meta", func() {
		Result(func() {
			Attribute("object", Any)
			Required("object")
		})

		HTTP(func() {
			GET("modules/meta")

			Response(func() {
				Body("object")
			})
		})
	})

	commonOptions()
})
