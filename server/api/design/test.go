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
			GET("error")
		})
	})

	Method("json", func() {
		Payload(func() {
			Attribute("id", Int64)
		})

		Result(func() {
			Attribute("json", String)
			Required("json")
		})

		HTTP(func() {
			GET("test/json/{id}")

			Response(func() {
				Body("json")
			})
		})
	})
})
