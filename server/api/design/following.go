package design

import (
	. "goa.design/goa/v3/dsl"
)

var Avatar = Type("Avatar", func() {
	Attribute("url", String)
	Required("url")
})

var Follower = ResultType("application/vnd.app.follower", func() {
	TypeName("Follower")
	Attributes(func() {
		Attribute("id", Int64)
		Attribute("name", String)
		Attribute("avatar", Avatar)
		Required("id", "name")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("avatar")
	})
})

var FollowersPage = ResultType("application/vnd.app.followers", func() {
	TypeName("FollowersPage")
	Attributes(func() {
		Attribute("followers", CollectionOf(Follower))
		Attribute("total", Int32)
		Attribute("page", Int32)
		Required("followers", "total", "page")
	})
})

var _ = Service("following", func() {
	Method("follow", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Attribute("id", Int64)
		})

		HTTP(func() {
			POST("projects/{id}/follow")

			httpAuthentication()
		})
	})

	Method("unfollow", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Attribute("id", Int64)
		})

		HTTP(func() {
			POST("projects/{id}/unfollow")

			httpAuthentication()
		})
	})

	Method("followers", func() {
		Payload(func() {
			Attribute("id", Int64)
			Attribute("page", Int64)
		})

		Result(FollowersPage)

		HTTP(func() {
			GET("projects/{id}/followers")

			Params(func() {
				Param("page")
			})
		})
	})

	commonOptions()
})
