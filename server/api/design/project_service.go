package design

import (
	. "goa.design/goa/v3/dsl"
)

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

var _ = Service("project", func() {
	Method("update", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Attribute("id", Int64)
			Attribute("body", String)
			Required("auth")
			Required("id")
			Required("body")
		})

		Result(Empty)

		HTTP(func() {
			POST("projects/{id}/update")

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

	commonOptions()
})
