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

	Error("unauthorized", String, "credentials are invalid")
	Error("not-found", String, "not found")
	Error("bad-request", String, "bad request")

	HTTP(func() {
		Response("unauthorized", StatusUnauthorized)
		Response("not-found", StatusNotFound)
		Response("bad-request", StatusBadRequest)
	})

	commonOptions()
})
