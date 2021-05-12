package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("records", func() {
	Security(JWTAuth, func() {
		Scope("api:admin")
	})

	Method("data", func() {
		Payload(func() {
			Token("auth")
			Attribute("recordId", Int64)
		})

		Result(func() {
			Attribute("object", Any)
			Required("object")
		})

		HTTP(func() {
			GET("records/data/{recordId}")

			Response(StatusOK, func() {
				Body("object")
			})
		})
	})

	Method("meta", func() {
		Payload(func() {
			Token("auth")
			Attribute("recordId", Int64)
		})

		Result(func() {
			Attribute("object", Any)
			Required("object")
		})

		HTTP(func() {
			GET("records/meta/{recordId}")

			Response(StatusOK, func() {
				Body("object")
			})
		})
	})

	Method("resolved", func() {
		Payload(func() {
			Token("auth")
			Attribute("recordId", Int64)
		})

		Result(func() {
			Attribute("object", Any)
			Required("object")
		})

		HTTP(func() {
			GET("records/data/{recordId}/resolved")

			Response(StatusOK, func() {
				Body("object")
			})
		})
	})

	commonOptions()
})
