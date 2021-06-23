package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("oidc", func() {
	Error("user-unverified", func() {
	})

	Method("required", func() {
		Payload(func() {
			Attribute("token", String)
			Attribute("follow", Boolean)
			Attribute("after", String)
		})

		Result(func() {
			Attribute("location", String)
			Required("location")
		})

		HTTP(func() {
			GET("oidc/required")

			Header("token:Authorization", String, "authentication token", func() {
				Pattern("^Bearer [^ ]+$")
			})

			Params(func() {
				Param("after")
				Param("follow")
			})

			Response(StatusTemporaryRedirect, func() {
				Headers(func() {
					Header("location:Location")
				})
			})
		})
	})

	Method("url", func() {
		Payload(func() {
			Attribute("token", String)
			Attribute("follow", Boolean)
			Attribute("after", String)
		})

		Result(func() {
			Attribute("location", String)
			Required("location")
		})

		HTTP(func() {
			GET("oidc/url")

			Header("token:Authorization", String, "authentication token", func() {
				Pattern("^Bearer [^ ]+$")
			})

			Params(func() {
				Param("after")
				Param("follow")
			})

			Response(StatusOK)
		})
	})

	Method("authenticate", func() {
		Payload(func() {
			Attribute("state", String)
			Required("state")
			Attribute("session_state", String)
			Required("session_state")
			Attribute("code", String)
			Required("code")
		})

		Result(func() {
			Attribute("location", String)
			Required("location")
			Attribute("token", String)
			Required("token")
			Attribute("header", String)
			Required("header")
		})

		HTTP(func() {
			POST("oidc/auth")

			Params(func() {
				Param("state")
				Param("session_state")
				Param("code")
			})

			Response(StatusOK)

			Response("user-unverified", StatusForbidden)
		})
	})

	commonOptions()
})
