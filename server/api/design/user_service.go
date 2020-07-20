package design

import (
	. "goa.design/goa/v3/dsl"
)

var AvailableRole = Type("AvailableRole", func() {
	Attribute("id", Int32)
	Attribute("name", String)
	Required("id", "name")
})

var AvailableRoles = ResultType("application/vnd.app.roles.available", func() {
	TypeName("AvailableRoles")
	Attributes(func() {
		Attribute("roles", ArrayOf(AvailableRole))
		Required("roles")
	})
	View("default", func() {
		Attribute("roles")
	})
})

var _ = Service("user", func() {
	Method("roles", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
		})

		Result(AvailableRoles)

		HTTP(func() {
			GET("roles")

			httpAuthentication()
		})
	})

	Method("delete", func() {
		Security(JWTAuth, func() {
			Scope("api:admin")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("userId", Int32)
			Required("userId")
		})

		HTTP(func() {
			DELETE("admin/users/{userId}")

			httpAuthentication()
		})
	})

	Method("upload photo", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("contentLength", Int64)
			Required("contentLength")
			Attribute("contentType", String)
			Required("contentType")
		})

		HTTP(func() {
			POST("user/media")

			Header("contentType:Content-Type")
			Header("contentLength:Content-Length")

			SkipRequestBodyEncodeDecode()

			httpAuthentication()
		})
	})

	Method("download photo", func() {
		/*
			Security(JWTAuth, func() {
				Scope("api:access")
			})
		*/

		Payload(func() {
			// Token("auth")
			// Required("auth")
			Attribute("userId", Int32)
			Required("userId")
		})

		Result(func() {
			Attribute("length", Int64)
			Required("length")
			Attribute("contentType", String)
			Required("contentType")
		})

		HTTP(func() {
			GET("user/{userId}/media")

			SkipResponseBodyEncodeDecode()

			Response(func() {
				Header("length:Content-Length")
				Header("contentType:Content-Type")
			})

			// httpAuthenticationQueryString()
		})
	})

	commonOptions()
})
