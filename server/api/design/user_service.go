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
	Error("user-unverified", func() {
	})

	Error("user-email-registered", func() {
	})

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
		Payload(func() {
			Attribute("userId", Int32)
			Required("userId")
			Attribute("size", Int32)
			Attribute("ifNoneMatch", String)
		})

		Result(DownloadedPhoto)

		HTTP(func() {
			GET("user/{userId}/media")

			Header("ifNoneMatch:If-None-Match")

			Params(func() {
				Param("size")
			})
		})
	})

	Method("login", func() {
		Payload(func() {
			Attribute("login", LoginFields)
			Required("login")
		})

		Result(func() {
			Attribute("authorization", String)
			Required("authorization")
		})

		HTTP(func() {
			POST("login")

			Body("login")

			Response(StatusNoContent, func() {
				Headers(func() {
					Header("authorization:Authorization")
				})
			})

			Response("user-unverified", StatusForbidden)
		})
	})

	Method("recovery lookup", func() {
		Payload(func() {
			Attribute("recovery", RecoveryLookupFields)
			Required("recovery")
		})

		HTTP(func() {
			POST("user/recovery/lookup")

			Body("recovery")
		})
	})

	Method("recovery", func() {
		Payload(func() {
			Attribute("recovery", RecoveryFields)
			Required("recovery")
		})

		HTTP(func() {
			POST("user/recovery")

			Body("recovery")
		})
	})

	Method("resume", func() {
		Payload(func() {
			Attribute("token", String)
			Required("token")
		})

		Result(func() {
			Attribute("authorization", String)
			Required("authorization")
		})

		HTTP(func() {
			POST("user/resume")

			Body(func() {
				Attribute("token")
			})

			Response(StatusNoContent, func() {
				Headers(func() {
					Header("authorization:Authorization")
				})
			})
		})
	})

	Method("logout", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
		})

		HTTP(func() {
			POST("logout")

			Response(StatusNoContent)

			httpAuthentication()
		})
	})

	Method("refresh", func() {
		Payload(func() {
			Attribute("refreshToken", String)
			Required("refreshToken")
		})

		Result(func() {
			Attribute("authorization", String)
			Required("authorization")
		})

		HTTP(func() {
			POST("refresh")

			Response(StatusNoContent, func() {
				Headers(func() {
					Header("authorization:Authorization")
				})
			})
		})
	})

	Method("send validation", func() {
		Payload(func() {
			Attribute("userId", Int32)
			Required("userId")
		})

		HTTP(func() {
			POST("users/{userId}/validate-email")

			Response(StatusNoContent)
		})
	})

	Method("validate", func() {
		Payload(func() {
			Attribute("token", String)
			Required("token")
		})

		Result(func() {
			Attribute("location", String)
			Required("location")
		})

		HTTP(func() {
			GET("validate")

			Params(func() {
				Param("token")
			})

			Response(StatusFound, func() {
				Headers(func() {
					Header("location:Location")
				})
			})
		})
	})

	Method("add", func() {
		Payload(func() {
			Attribute("user", AddUserFields)
			Required("user")
		})

		Result(User)

		HTTP(func() {
			POST("users")

			Body("user")

			Response("user-email-registered", StatusBadRequest)
		})
	})

	Method("update", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("userId", Int32)
			Required("userId")
			Attribute("update", UpdateUserFields)
			Required("update")
		})

		Result(User)

		HTTP(func() {
			PATCH("users/{userId}")

			Body("update")

			httpAuthentication()
		})
	})

	Method("change password", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("userId", Int32)
			Required("userId")
			Attribute("change", UpdateUserPasswordFields)
			Required("change")
		})

		Result(User)

		HTTP(func() {
			PATCH("users/{userId}/password")

			Body("change")

			httpAuthentication()
		})
	})

	Method("accept tnc", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("userId", Int32)
			Required("userId")
			Attribute("accept", AcceptTncFields)
			Required("accept")
		})

		Result(User)

		HTTP(func() {
			PATCH("users/{userId}/accept-tnc")

			Body("accept")

			httpAuthentication()
		})
	})

	Method("get current", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
		})

		Result(User)

		HTTP(func() {
			GET("user")

			httpAuthentication()
		})
	})

	Method("list by project", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("projectId", Int32)
			Required("projectId")
		})

		Result(ProjectUsers)

		HTTP(func() {
			GET("users/project/{projectId}")

			httpAuthentication()
		})
	})

	Method("issue transmission token", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
		})

		Result(TransmissionToken)

		HTTP(func() {
			GET("user/transmission-token")

			httpAuthentication()
		})
	})

	Method("project roles", func() {
		Result(CollectionOf(ProjectRole))

		HTTP(func() {
			GET("projects/roles")
		})
	})

	Method("admin delete", func() {
		Security(JWTAuth, func() {
			Scope("api:admin")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("delete", AdminDeleteFields)
			Required("delete")
		})

		HTTP(func() {
			DELETE("admin/user")

			Body("delete")

			httpAuthentication()
		})
	})

	Method("admin search", func() {
		Security(JWTAuth, func() {
			Scope("api:admin")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("query", String)
			Required("query")
		})

		Result(func() {
			Attribute("users", CollectionOf(User))
			Required("users")
		})

		HTTP(func() {
			POST("admin/users/search")

			Params(func() {
				Param("query")
			})

			httpAuthentication()
		})
	})

	Method("mentionables", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("projectId", Int32)
			Attribute("bookmark", String)
			Attribute("query", String)
			Required("query")
		})

		Result(MentionableOptions)

		HTTP(func() {
			GET("mentionables")

			Params(func() {
				Param("projectId")
				Param("bookmark")
				Param("query")
			})

			httpAuthentication()
		})
	})

	commonOptions()
})

var AddUserFields = Type("AddUserFields", func() {
	Attribute("name", String, func() {
		Pattern(`\S`)
		MaxLength(256)
	})
	Attribute("email", String, func() {
		Format("email")
		MaxLength(40)
	})
	Attribute("password", String, func() {
		MinLength(10)
	})
	Attribute("invite_token", String)
	Attribute("tncAccept", Boolean)
	Required("name", "email", "password")
})

var UpdateUserFields = Type("UpdateUserFields", func() {
	Reference(AddUserFields)
	Attribute("name")
	Attribute("email")
	Attribute("bio")
	Required("name", "email", "bio")
})

var UpdateUserPasswordFields = Type("UpdateUserPasswordFields", func() {
	Attribute("oldPassword", String, func() {
		MinLength(10)
	})
	Attribute("newPassword", String, func() {
		MinLength(10)
	})
	Required("oldPassword", "newPassword")
})

var AcceptTncFields = Type("AcceptTncFields", func() {
	Attribute("accept", Boolean)
	Required("accept")
})

var LoginFields = Type("LoginFields", func() {
	Reference(AddUserFields)
	Attribute("email")
	Attribute("password", String, func() {
		MinLength(10)
	})
	Required("email", "password")
})

var RecoveryLookupFields = Type("RecoveryLookupFields", func() {
	Attribute("email", String)
	Required("email")
})

var RecoveryFields = Type("RecoveryFields", func() {
	Attribute("token", String)
	Attribute("password", String, func() {
		MinLength(10)
	})
	Required("token", "password")
})

var UserPhoto = Type("UserPhoto", func() {
	Attribute("url", String)
})

var User = ResultType("application/vnd.app.user+json", func() {
	TypeName("User")
	Reference(AddUserFields)
	Attributes(func() {
		Attribute("id", Int32)
		Attribute("name")
		Attribute("email")
		Attribute("bio")
		Attribute("photo", UserPhoto)
		Attribute("admin", Boolean)
		Attribute("updatedAt", Int64)
		Attribute("tncDate", Int64)
		Required("id", "name", "email", "bio", "admin", "updatedAt", "tncDate")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("email")
		Attribute("bio")
		Attribute("photo")
		Attribute("admin")
		Attribute("updatedAt")
		Attribute("tncDate")
	})
})

var ProjectRole = ResultType("application/vnd.app.project.role+json", func() {
	TypeName("ProjectRole")
	Attributes(func() {
		Attribute("id", Int32)
		Attribute("name", String)
		Required("id", "name")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
	})
})

var ProjectUser = ResultType("application/vnd.app.project.user+json", func() {
	TypeName("ProjectUser")
	Attributes(func() {
		Attribute("user", User)
		Attribute("role", String)
		Attribute("membership", String)
		Attribute("invited", Boolean)
		Attribute("accepted", Boolean)
		Attribute("rejected", Boolean)
		Required("user", "role", "membership", "invited", "accepted", "rejected")
	})
	View("default", func() {
		Attribute("user")
		Attribute("role")
		Attribute("membership")
		Attribute("invited")
		Attribute("accepted")
		Attribute("rejected")
	})
})

var ProjectUsers = ResultType("application/vnd.app.users+json", func() {
	TypeName("ProjectUsers")
	Attributes(func() {
		Attribute("users", CollectionOf(ProjectUser))
		Required("users")
	})
	View("default", func() {
		Attribute("users")
	})
})

var MentionableUser = ResultType("application/vnd.app.mentionable.user+json", func() {
	TypeName("MentionableUser")
	Attributes(func() {
		Attribute("id", Int32)
		Attribute("name", String)
		Attribute("mention", String)
		Attribute("photo", UserPhoto)
		Required("id")
		Required("name")
		Required("mention")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("mention")
		Attribute("photo")
	})
})

var MentionableOptions = ResultType("application/vnd.app.mentionables+json", func() {
	TypeName("MentionableOptions")
	Attributes(func() {
		Attribute("users", CollectionOf(MentionableUser))
		Required("users")
	})
	View("default", func() {
		Attribute("users")
	})
})

var TransmissionToken = ResultType("application/vnd.app.user.transmission.token+json", func() {
	TypeName("TransmissionToken")
	Attributes(func() {
		Attribute("token", String)
		Attribute("url", String)
		Required("token")
		Required("url")
	})
	View("default", func() {
		Attribute("token")
		Attribute("url")
	})
})

var AdminDeleteFields = ResultType("application/vnd.app.admin.user.delete+json", func() {
	TypeName("AdminDeleteFields")
	Attributes(func() {
		Attribute("email", String)
		Required("email")
		Attribute("password", String)
		Required("password")
	})
	/*
		View("default", func() {
			Attribute("email")
			Attribute("password")
		})
	*/
})
