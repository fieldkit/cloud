package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var AddUserPayload = Type("AddUserPayload", func() {
	Attribute("email", String, func() {
		Format("email")
	})
	Attribute("username", String, func() {
		Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
		MaxLength(40)
	})
	Attribute("password", String, func() {
		MinLength(10)
	})
	Attribute("invite_token", String)
	Required("email", "username", "password", "invite_token")
})

var User = MediaType("application/vnd.app.user+json", func() {
	TypeName("User")
	Reference(AddUserPayload)
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("username")
		Required("id", "username")
	})
	View("default", func() {
		Attribute("id")
		Attribute("username")
	})
})

var Users = MediaType("application/vnd.app.users+json", func() {
	TypeName("Users")
	Attributes(func() {
		Attribute("users", CollectionOf(User))
		Required("users")
	})
	View("default", func() {
		Attribute("users")
	})
})

var _ = Resource("user", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("login", func() {
		Routing(POST("login"))
		Description("Creates a valid JWT given login credentials.")
		NoSecurity()
		Payload(AddUserPayload)
		Response(NoContent, func() {
			Headers(func() {
				Header("Authorization", String, "Generated JWT")
			})
		})
		Response(Unauthorized)
		Response(BadRequest)
	})

	Action("logout", func() {
		Routing(POST("logout"))
		Description("Creates a valid JWT given login credentials.")
		Response(NoContent)
		Response(BadRequest)
	})

	Action("refresh", func() {
		Routing(POST("refresh"))
		Description("Creates a valid JWT given a refresh token.")
		NoSecurity()
		Payload(func() {
			Param("refresh_token", String)
			Required("refresh_token")
		})
		Response(NoContent, func() {
			Headers(func() {
				Header("Authorization", String, "Generated JWT")
			})
		})
		Response(Unauthorized)
	})

	Action("add", func() {
		Routing(POST("user"))
		Description("Add a user")
		NoSecurity()
		Payload(func() {
			Param("email", String, func() {
				Format("email")
			})
			Param("username", String, func() {
				Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
				MaxLength(40)
			})
			Param("password", String, func() {
				MinLength(10)
			})
			Param("invite_token", String)
			Required("email", "username", "password", "invite_token")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(User)
		})
	})

	Action("get current", func() {
		Routing(GET("user"))
		Description("Get the authenticated user")
		Response(OK, func() {
			Media(User)
		})
	})

	Action("get", func() {
		Routing(GET("user/:username"))
		Description("Get a user")
		Response(OK, func() {
			Media(User)
		})
	})

	Action("list", func() {
		Routing(GET("users"))
		Description("List users")
		Response(OK, func() {
			Media(Users)
		})
	})
})
