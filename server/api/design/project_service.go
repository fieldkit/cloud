package design

import (
	. "goa.design/goa/v3/dsl"
)

var AddProjectFields = Type("AddProjectFields", func() {
	Attribute("name", String)
	Attribute("description", String)
	Attribute("goal", String)
	Attribute("location", String)
	Attribute("tags", String)
	Attribute("private", Boolean)
	Attribute("startTime", String)
	Attribute("endTime", String)
	Required("name", "description")
})

var InviteUserFields = Type("InviteUserFields", func() {
	Attribute("email", String)
	Required("email")
	Attribute("role", Int32)
	Required("role")
})

var RemoveUserFields = Type("RemoveUserFields", func() {
	Attribute("email", String)
	Required("email")
})

var Project = ResultType("application/vnd.app.project+json", func() {
	TypeName("Project")
	Reference(AddProjectFields)
	Attributes(func() {
		Attribute("id", Int32)
		Attribute("name")
		Attribute("description")
		Attribute("goal")
		Attribute("location")
		Attribute("tags")
		Attribute("private", Boolean)
		Attribute("startTime", String)
		Attribute("endTime", String)
		Attribute("photo")
		Attribute("readOnly", Boolean)
		Attribute("numberOfFollowers", Int32)
		Required("id", "name", "description", "goal", "location", "private", "tags", "readOnly", "numberOfFollowers")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("description")
		Attribute("goal")
		Attribute("location")
		Attribute("tags")
		Attribute("private")
		Attribute("startTime")
		Attribute("endTime")
		Attribute("photo")
		Attribute("readOnly")
		Attribute("numberOfFollowers")
	})
})

var Projects = ResultType("application/vnd.app.projects+json", func() {
	TypeName("Projects")
	Attributes(func() {
		Attribute("projects", CollectionOf(Project))
		Required("projects")
	})
	View("default", func() {
		Attribute("projects")
	})
})

var PendingInvite = Type("PendingInvite", func() {
	Attribute("id", Int64)
	Attribute("project", ProjectSummary)
	Attribute("time", Int64)
	Attribute("role", Int32)
	Required("id")
	Required("project")
	Required("time")
	Required("role")
})

var PendingInvites = ResultType("application/vnd.app.invites.pending", func() {
	TypeName("PendingInvites")
	Attributes(func() {
		Attribute("pending", ArrayOf(PendingInvite))
		Required("pending")
	})
	View("default", func() {
		Attribute("pending")
	})
})

var ProjectUpdate = ResultType("application/vnd.app.project.update", func() {
	TypeName("ProjectUpdate")
	Attributes(func() {
		Attribute("id", Int64)
		Attribute("body", String)
		Attribute("createdAt", Int64)
		Required("id")
		Required("body")
		Required("createdAt")
	})
	View("default", func() {
		Attribute("id")
		Attribute("body")
	})
})

var _ = Service("project", func() {
	Method("add update", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("projectId", Int32)
			Attribute("body", String)
			Required("projectId")
			Required("body")
		})

		Result(ProjectUpdate)

		HTTP(func() {
			POST("projects/{projectId}/updates")

			httpAuthentication()
		})
	})

	Method("delete update", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("projectId", Int32)
			Attribute("updateId", Int64)
			Required("projectId")
			Required("updateId")
		})

		Result(Empty)

		HTTP(func() {
			DELETE("projects/{projectId}/updates/{updateId}")

			httpAuthentication()
		})
	})

	Method("modify update", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("projectId", Int32)
			Attribute("updateId", Int64)
			Attribute("body", String)
			Required("projectId")
			Required("updateId")
			Required("body")
		})

		Result(ProjectUpdate)

		HTTP(func() {
			POST("projects/{projectId}/updates/{updateId}")

			httpAuthentication()
		})
	})

	Method("invites", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
		})

		Result(PendingInvites)

		HTTP(func() {
			GET("projects/invites/pending")

			httpAuthentication()
		})
	})

	Method("lookup invite", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Token("token")
			Required("token")
		})

		Result(PendingInvites)

		HTTP(func() {
			GET("projects/invites/{token}")

			httpAuthentication()
		})
	})

	Method("accept invite", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("id", Int64)
			Required("id")
			Attribute("token", String)
		})

		HTTP(func() {
			POST("projects/invites/{id}/accept")

			Params(func() {
				Param("token")
			})

			httpAuthentication()
		})
	})

	Method("reject invite", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("id", Int64)
			Required("id")
			Attribute("token", String)
		})

		HTTP(func() {
			POST("projects/invites/{id}/reject")

			Params(func() {
				Param("token")
			})

			httpAuthentication()
		})
	})

	Method("add", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("project", AddProjectFields)
			Required("project")
		})

		Result(Project)

		HTTP(func() {
			POST("projects")

			Body("project")

			httpAuthentication()
		})
	})

	Method("update", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("projectId", Int32)
			Required("projectId")
			Attribute("project", AddProjectFields)
			Required("project")
		})

		Result(Project)

		HTTP(func() {
			PATCH("projects/{projectId}")

			Body("project")

			httpAuthentication()
		})
	})

	Method("get", func() {
		Security(JWTAuth, func() {
			// Optional
		})

		Payload(func() {
			Token("auth")
			Attribute("projectId", Int32)
			Required("projectId")
		})

		Result(Project)

		HTTP(func() {
			GET("projects/{projectId}")

			httpAuthentication()
		})
	})

	Method("list community", func() {
		Security(JWTAuth, func() {
			// Optional
		})

		Payload(func() {
			Token("auth")
		})

		Result(Projects)

		HTTP(func() {
			GET("projects")

			httpAuthentication()
		})
	})

	Method("list mine", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
		})

		Result(Projects)

		HTTP(func() {
			GET("user/projects")

			httpAuthentication()
		})
	})

	Method("invite", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("projectId", Int32)
			Required("projectId")
			Attribute("invite", InviteUserFields)
			Required("invite")
		})

		HTTP(func() {
			POST("projects/{projectId}/invite")

			Body("invite")
		})
	})

	Method("remove user", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("projectId", Int32)
			Required("projectId")
			Attribute("remove", RemoveUserFields)
			Required("remove")
		})

		HTTP(func() {
			DELETE("projects/{projectId}/members")

			Body("remove")
		})
	})

	Method("add station", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("projectId", Int32)
			Required("projectId")
			Attribute("stationId", Int32)
			Required("stationId")
		})

		HTTP(func() {
			POST("projects/{projectId}/stations/{stationId}")

			httpAuthentication()
		})
	})

	Method("remove station", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("projectId", Int32)
			Required("projectId")
			Attribute("stationId", Int32)
			Required("stationId")
		})

		HTTP(func() {
			DELETE("projects/{projectId}/stations/{stationId}")

			httpAuthentication()
		})
	})

	Method("delete", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("projectId", Int32)
			Required("projectId")
		})

		HTTP(func() {
			DELETE("projects/{projectId}")

			httpAuthentication()
		})
	})

	Method("upload media", func() {
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
			Attribute("projectId", Int32)
			Required("projectId")
		})

		HTTP(func() {
			POST("projects/{projectId}/media")

			Header("contentType:Content-Type")
			Header("contentLength:Content-Length")

			SkipRequestBodyEncodeDecode()

			httpAuthentication()
		})
	})

	Method("download media", func() {
		/*
			Security(JWTAuth, func() {
				Scope("api:access")
			})
		*/

		Payload(func() {
			// Token("auth")
			// Required("auth")
			Attribute("projectId", Int32)
			Required("projectId")
		})

		Result(func() {
			Attribute("length", Int64)
			Required("length")
			Attribute("contentType", String)
			Required("contentType")
		})

		HTTP(func() {
			GET("projects/{projectId}/media")

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
