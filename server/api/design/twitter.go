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
	Attribute("team_id")
	Attribute("user_id")
})

var TwitterAccountInput = MediaType("application/vnd.app.twitter_account_input+json", func() {
	TypeName("TwitterAccountInput")
	Reference(Input)
	Attributes(func() {
		Attribute("id")
		Attribute("expedition_id")
		Attribute("name")
		Attribute("team_id")
		Attribute("user_id")
		Attribute("twitter_account_id", Integer)
		Attribute("screen_name", String)
		Required("id", "expedition_id", "name", "twitter_account_id", "screen_name")
	})
	View("default", func() {
		Attribute("id")
		Attribute("expedition_id")
		Attribute("name")
		Attribute("team_id")
		Attribute("user_id")
		Attribute("twitter_account_id")
		Attribute("screen_name")
	})
})

var TwitterAccountInputs = MediaType("application/vnd.app.twitter_account_intputs+json", func() {
	TypeName("TwitterAccountInputs")
	Attributes(func() {
		Attribute("twitter_account_inputs", CollectionOf(TwitterAccountInput))
		Required("twitter_account_inputs")
	})
	View("default", func() {
		Attribute("twitter_account_inputs")
	})
})

var _ = Resource("twitter", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("add", func() {
		Routing(POST("expeditions/:expedition_id/inputs/twitter-accounts"))
		Description("Add a Twitter account input")
		Params(func() {
			Param("expedition_id", Integer)
			Required("expedition_id")
		})
		Payload(AddTwitterAccountInputPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(Location)
		})
	})

	Action("get id", func() {
		Routing(GET("inputs/twitter-accounts/:input_id"))
		Description("Get a Twitter account input")
		Params(func() {
			Param("input_id", Integer)
			Required("input_id")
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
		Routing(GET("expeditions/:expedition_id/inputs/twitter-accounts"))
		Description("List an expedition's Twitter account inputs")
		Params(func() {
			Param("expedition_id", Integer)
			Required("expedition_id")
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
			Param("oauth_token", String)
			Param("oauth_verifier", String)
			Required("oauth_token", "oauth_verifier")
		})
		Response(BadRequest)
		Response(Found, func() {
			Headers(func() {
				Header("Location", String)
			})
		})
	})
})
