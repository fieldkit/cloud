package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("following", func() {
	commonOptions()

	Method("follow", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Attribute("id", Int64)
		})

		HTTP(func() {
			POST("projects/{id}/follow")
		})
	})

	Method("unfollow", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Attribute("id", Int64)
		})

		HTTP(func() {
			POST("projects/{id}/unfollow")
		})
	})

	Error("unauthorized", String, "credentials are invalid")

	HTTP(func() {
		Response("unauthorized", StatusUnauthorized)

		Header("auth:Authorization", String, "authentication token", func() {
			Pattern("^Bearer [^ ]+$")
		})
	})
})
