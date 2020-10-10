package design

import (
	. "goa.design/goa/v3/dsl"
)

var PostAuthor = Type("PostAuthor", func() {
	Attribute("id", Int32)
	Attribute("name", String)
	Attribute("mediaUrl", String)
	Required("id", "name")
})

var ThreadedPost = ResultType("application/vnd.app.discussion.post", func() {
	TypeName("ThreadedPost")
	Attributes(func() {
		Attribute("id", Int64)
		Attribute("createdAt", Int64)
		Attribute("updatedAt", Int64)
		Attribute("author", PostAuthor)
		Attribute("replies", ArrayOf("ThreadedPost"))
		Attribute("body", String)
		Attribute("bookmark", String)
		Required("id", "createdAt", "updatedAt", "author", "replies", "body")
	})
	View("default", func() {
		Attribute("id")
		Attribute("createdAt")
		Attribute("updatedAt")
		Attribute("author")
		Attribute("replies")
		Attribute("body")
		Attribute("bookmark")
	})
})

var Discussion = ResultType("application/vnd.app.discussion", func() {
	TypeName("Discussion")
	Attributes(func() {
		Attribute("posts", ArrayOf(ThreadedPost))
		Required("posts")
	})
	View("default", func() {
		Attribute("posts")
	})
})

var NewPost = ResultType("application/vnd.app.discussion.post.new", func() {
	TypeName("NewPost")
	Attributes(func() {
		Attribute("threadId", Int64)
		Attribute("body", String)
		Attribute("projectId", Int32)
		Attribute("bookmark", String)
		Required("body")
	})
	View("default", func() {
		Attribute("threadId")
		Attribute("body")
		Attribute("projectId")
		Attribute("bookmark")
	})
})

var _ = Service("discussion", func() {
	Method("project", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("projectID", Int32)
			Required("projectID")
		})

		Result(Discussion)

		HTTP(func() {
			GET("discussion/projects/{projectID}")

			httpAuthentication()
		})
	})

	Method("data", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("bookmark", String)
			Required("bookmark")
		})

		Result(Discussion)

		HTTP(func() {
			GET("discussion")

			Params(func() {
				Param("bookmark")
			})

			httpAuthentication()
		})
	})

	Method("post message", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("post", NewPost)
			Required("post")
		})

		Result(func() {
			Attribute("post", ThreadedPost)
			Required("post")
		})

		HTTP(func() {
			POST("discussion")

			httpAuthentication()
		})
	})

	commonOptions()
})
