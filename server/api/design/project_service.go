package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("project", func() {
	commonOptions()

	Method("update", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Attribute("id", Int64)
			Attribute("body", String)
			Required("id", "body")
		})

		Result(Empty)

		HTTP(func() {
			POST("projects/{id}/update")

			Header("auth:Authorization", String, "authentication token", func() {
				Pattern("^Bearer [^ ]+$")
			})
		})
	})

	HTTP(func() {
	})
})
