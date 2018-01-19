package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var AddTwitterAccountInputPayload = Type("AddTwitterAccountInputPayload", func() {
	Reference(Input)
	Attribute("name")
	Required("name")
})

var UpdateTwitterAccountInputPayload = Type("UpdateTwitterAccountInputPayload", func() {
	Reference(Input)
	Attribute("name")
	Attribute("teamId")
	Attribute("userId")
})

var TwitterAccountInput = MediaType("application/vnd.app.twitter_account_input+json", func() {
	TypeName("TwitterAccountInput")
	Reference(Input)
	Attributes(func() {
		Attribute("id")
		Attribute("expeditionId")
		Attribute("name")
		Attribute("teamId")
		Attribute("userId")
		Attribute("twitterAccountId", Integer)
		Attribute("screenName", String)
		Required("id", "expeditionId", "name", "twitterAccountId", "screenName")
	})
	View("default", func() {
		Attribute("id")
		Attribute("expeditionId")
		Attribute("name")
		Attribute("teamId")
		Attribute("userId")
		Attribute("twitterAccountId")
		Attribute("screenName")
	})
})

var TwitterAccountInputs = MediaType("application/vnd.app.twitter_account_intputs+json", func() {
	TypeName("TwitterAccountInputs")
	Attributes(func() {
		Attribute("twitterAccountInputs", CollectionOf(TwitterAccountInput))
		Required("twitterAccountInputs")
	})
	View("default", func() {
		Attribute("twitterAccountInputs")
	})
})

var _ = Resource("twitter", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("add", func() {
		Routing(POST("expeditions/:expeditionId/inputs/twitter-accounts"))
		Description("Add a Twitter account input")
		Params(func() {
			Param("expeditionId", Integer)
			Required("expeditionId")
		})
		Payload(AddTwitterAccountInputPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(Location)
		})
	})

	Action("get id", func() {
		Routing(GET("inputs/twitter-accounts/:inputId"))
		Description("Get a Twitter account input")
		Params(func() {
			Param("inputId", Integer)
			Required("inputId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(TwitterAccountInput)
		})
	})

	Action("list", func() {
		Routing(GET("projects/@/:project/expeditions/@/:expedition/inputs/twitter-accounts"))
		Description("List an expedition's Twitter account inputs")
		Params(func() {
			Param("project", String, ProjectSlug)
			Param("expedition", String, ExpeditionSlug)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(TwitterAccountInputs)
		})
	})

	Action("list id", func() {
		Routing(GET("expeditions/:expeditionId/inputs/twitter-accounts"))
		Description("List an expedition's Twitter account inputs")
		Params(func() {
			Param("expeditionId", Integer)
			Required("expeditionId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(TwitterAccountInputs)
		})
	})

	Action("callback", func() {
		Routing(GET("twitter/callback"))
		Description("OAuth callback endpoint for Twitter")
		NoSecurity()
		Params(func() {
			Param("oauthToken", String)
			Param("oauthVerifier", String)
			Required("oauthToken", "oauth_verifier")
		})
		Response(BadRequest)
		Response(Found, func() {
			Headers(func() {
				Header("Location", String)
			})
		})
	})
})
