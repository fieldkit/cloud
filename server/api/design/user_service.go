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

	commonOptions()
})
