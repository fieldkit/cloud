package design

import (
	. "goa.design/goa/v3/dsl"
)

var PostAuthor = Type("PostAuthor", func() {
	Attribute("id", Int32, "id", func() {
		Meta("struct:tag:json", "id")
	})
	Attribute("name", String, "name", func() {
		Meta("struct:tag:json", "name")
	})
	Attribute("mediaUrl", String, "media url", func() {
		Meta("struct:tag:json", "mediaUrl")
	})
	Required("id", "name")
})

var DiscussionSummary = ResultType("application/vnd.app.discussion.summary", func() {
	TypeName("DiscussionSummary")
	Attributes(func() {
		Attribute("total", Int32, "Total", func() {
			Meta("struct:tag:json", "total")
		})
		Required("total")
	})
	View("default", func() {
		Attribute("total")
	})
})

var Discussion = ResultType("application/vnd.app.discussion", func() {
	TypeName("Discussion")
	Attributes(func() {
		Attribute("summary", DiscussionSummary, "Summary", func() {
			Meta("struct:tag:json", "summary")
		})
		Attribute("posts", ArrayOf(ThreadedPost), "Posts", func() {
			Meta("struct:tag:json", "posts")
		})
		Required("summary")
		Required("posts")
	})
	View("default", func() {
		Attribute("summary")
		Attribute("posts")
	})
})

var ThreadedPost = ResultType("application/vnd.app.discussion.post", func() {
	TypeName("ThreadedPost")
	Attributes(func() {
		Attribute("id", Int64, "id", func() {
			Meta("struct:tag:json", "id")
		})
		Attribute("createdAt", Int64, "created at", func() {
			Meta("struct:tag:json", "createdAt")
		})
		Attribute("updatedAt", Int64, "updated at", func() {
			Meta("struct:tag:json", "updatedAt")
		})
		Attribute("author", PostAuthor, "author", func() {
			Meta("struct:tag:json", "author")
		})
		// Wish this could be 'Discussion'
		Attribute("replies", Any, "replies", func() {
			Meta("struct:tag:json", "replies")
		})
		Attribute("body", String, "body", func() {
			Meta("struct:tag:json", "body")
		})
		Attribute("bookmark", String, "bookmark", func() {
			Meta("struct:tag:json", "bookmark")
		})
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
