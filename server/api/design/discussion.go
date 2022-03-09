package design

import (
	. "goa.design/goa/v3/dsl"
)

var AuthorPhoto = Type("AuthorPhoto", func() {
	Attribute("url", String)
	Required("url")
})

var PostAuthor = Type("PostAuthor", func() {
	Attribute("id", Int32)
	Attribute("name", String)
	Attribute("photo", AuthorPhoto)
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
		    // Optional
		})

		Payload(func() {
			Token("auth")
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

	Method("update message", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("postId", Int64)
			Required("postId")
			Attribute("body", String)
			Required("body")
		})

		Result(func() {
			Attribute("post", ThreadedPost)
			Required("post")
		})

		HTTP(func() {
			POST("discussion/{postId}")

			httpAuthentication()
		})
	})

	Method("delete message", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("postId", Int64)
			Required("postId")
		})

		HTTP(func() {
			DELETE("discussion/{postId}")

			httpAuthentication()
		})
	})

	commonOptions()
})
