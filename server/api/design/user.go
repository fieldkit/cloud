package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var AddUserPayload = Type("AddUserPayload", func() {
	Attribute("name", String, func() {
		MaxLength(256)
	})
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
	Attribute("bio", String)
	Attribute("invite_token", String)
	Required("name", "email", "username", "password", "bio", "invite_token")
})

var UpdateUserPayload = Type("UpdateUserPayload", func() {
	Attribute("name")
	Attribute("email")
	Attribute("username")
	Attribute("bio")
	Required("name", "email", "username", "bio")
})

var LoginPayload = Type("LoginPayload", func() {
	Reference(AddUserPayload)
	Attribute("username")
	Attribute("password")
	Required("username", "password")
})

var User = MediaType("application/vnd.app.user+json", func() {
	TypeName("User")
	Reference(AddUserPayload)
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("name")
		Attribute("username")
		Attribute("email")
		Attribute("bio")
		Required("id", "name", "username", "email", "bio")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("username")
		Attribute("email")
		Attribute("bio")
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
		Payload(LoginPayload)
		Response(NoContent, func() {
			Headers(func() {
				Header("Authorization", String, "Generated JWT")
			})
		})
		Response(Unauthorized, ErrorMedia)
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

	Action("validate", func() {
		Routing(GET("validate"))
		Description("Validate a user's email address.")
		NoSecurity()
		Params(func() {
			Param("token", String)
			Required("token")
		})
		Response(Found, func() {
			Headers(func() {
				Header("Location", String)
			})
		})
		Response(Unauthorized)
	})

	Action("add", func() {
		Routing(POST("users"))
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

	Action("update", func() {
		Routing(PATCH("users/:user_id"))
		Description("Get a user")
		Params(func() {
			Param("user_id", Integer)
			Required("user_id")
		})
		Payload(UpdateUserPayload)
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
		Routing(GET("users/@/:username"))
		Description("Get a user")
		Response(OK, func() {
			Media(User)
		})
	})

	Action("get id", func() {
		Routing(GET("users/:user_id"))
		Description("Get a user")
		Params(func() {
			Param("user_id", Integer)
			Required("user_id")
		})
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
