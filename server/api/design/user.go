package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var AddUserPayload = Type("AddUserPayload", func() {
	Attribute("name", String, func() {
		Pattern(`\S`)
		MaxLength(256)
	})
	Attribute("email", String, func() {
		Format("email")
	})
	Attribute("password", String, func() {
		MinLength(10)
	})
	Attribute("invite_token", String)
	Required("name", "email", "password")
})

var UpdateUserPayload = Type("UpdateUserPayload", func() {
	Reference(AddUserPayload)
	Attribute("name")
	Attribute("email")
	Attribute("bio")
	Required("name", "email", "bio")
})

var UpdateUserPasswordPayload = Type("UpdateUserPasswordPayload", func() {
	Attribute("oldPassword", String, func() {
		MinLength(10)
	})
	Attribute("newPassword", String, func() {
		MinLength(10)
	})
	Required("oldPassword", "newPassword")
})

var LoginPayload = Type("LoginPayload", func() {
	Reference(AddUserPayload)
	Attribute("email")
	Attribute("password", String, func() {
		MinLength(10)
	})
	Required("email", "password")
})

var RecoveryLookupPayload = Type("RecoveryLookupPayload", func() {
	Attribute("email", String)
	Required("email")
})

var RecoveryPayload = Type("RecoveryPayload", func() {
	Attribute("token", String)
	Attribute("password", String, func() {
		MinLength(10)
	})
	Required("token", "password")
})

var UserPhoto = Type("UserPhoto", func() {
	Attribute("url", String)
})

var User = MediaType("application/vnd.app.user+json", func() {
	TypeName("User")
	Reference(AddUserPayload)
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("name")
		Attribute("email")
		Attribute("bio")
		Attribute("photo", UserPhoto)
		Attribute("admin", Boolean)
		Required("id", "name", "email", "bio", "admin")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("email")
		Attribute("bio")
		Attribute("photo")
		Attribute("admin")
	})
})

var ProjectRole = MediaType("application/vnd.app.project.role+json", func() {
	TypeName("ProjectRole")
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("name", String)
		Required("id", "name")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
	})
})

var ProjectUser = MediaType("application/vnd.app.project.user+json", func() {
	TypeName("ProjectUser")
	Attributes(func() {
		Attribute("user", User)
		Attribute("role", String)
		Attribute("membership", String)
		Required("user", "role", "membership")
	})
	View("default", func() {
		Attribute("user")
		Attribute("role")
		Attribute("membership")
	})
})

var ProjectUsers = MediaType("application/vnd.app.users+json", func() {
	TypeName("ProjectUsers")
	Attributes(func() {
		Attribute("users", CollectionOf(ProjectUser))
		Required("users")
	})
	View("default", func() {
		Attribute("users")
	})
})

var TransmissionToken = MediaType("application/vnd.app.user.transmission.token+json", func() {
	TypeName("TransmissionToken")
	Attributes(func() {
		Attribute("token", String)
		Required("token")
	})
	View("default", func() {
		Attribute("token")
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

	Action("recovery lookup", func() {
		Routing(POST("user/recovery/lookup"))
		NoSecurity()
		Payload(RecoveryLookupPayload)
		Response(Unauthorized)
		Response(OK, func() {
		})
	})

	Action("recovery", func() {
		Routing(POST("user/recovery"))
		NoSecurity()
		Payload(RecoveryPayload)
		Response(Unauthorized)
		Response(OK, func() {
		})
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

	Action("send validation", func() {
		Routing(POST("users/:userId/validate-email"))
		NoSecurity()
		Params(func() {
			Param("userId", Integer)
			Required("userId")
		})
		Response(Unauthorized)
		Response(NoContent)
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
		Payload(AddUserPayload)
		Response(Unauthorized)
		Response(BadRequest)
		Response(OK, func() {
			Media(User)
		})
	})

	Action("update", func() {
		Routing(PATCH("users/:userId"))
		Description("Update a user")
		Params(func() {
			Param("userId", Integer)
			Required("userId")
		})
		Payload(UpdateUserPayload)
		Response(OK, func() {
			Media(User)
		})
	})

	Action("change password", func() {
		Routing(PATCH("users/:userId/password"))
		Description("Update a user password")
		Params(func() {
			Param("userId", Integer)
			Required("userId")
		})
		Payload(UpdateUserPasswordPayload)
		Response(BadRequest)
		Response(Forbidden)
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

	Action("get id", func() {
		Routing(GET("users/:userId"))
		Description("Get a user")
		Params(func() {
			Param("userId", Integer)
			Required("userId")
		})
		Response(OK, func() {
			Media(User)
		})
	})

	Action("list by project", func() {
		Routing(GET("users/project/:projectId"))
		Description("List users by project")
		Response(OK, func() {
			Media(ProjectUsers)
		})
	})

	Action("save current user image", func() {
		Routing(POST("/user/media"))
		Description("Save the authenticated user's image")
		Response(BadRequest)
		Response(OK, func() {
			Media(User)
		})
	})

	Action("get current user image", func() {
		Routing(GET("/user/media"))
		Description("Get the authenticated user's image")
		Response(BadRequest)
		Response(OK, func() {
			Media("image/png")
		})
	})

	Action("get user image", func() {
		NoSecurity()
		Routing(GET("/user/:userId/media"))
		Description("Get a user image")
		Params(func() {
			Param("userId", Integer)
			Required("userId")
		})
		Response(OK, func() {
			Media("image/png")
		})
	})

	Action("transmission token", func() {
		Routing(GET("/user/transmission-token"))
		Response(OK, func() {
			Media(TransmissionToken)
		})
	})

	Action("project roles", func() {
		NoSecurity()
		Routing(GET("/projects/roles"))
		Response(OK, func() {
			Media(CollectionOf(ProjectRole))
		})
	})
})
