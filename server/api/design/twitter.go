package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var TwitterAccount = MediaType("application/vnd.app.twitter_account+json", func() {
	TypeName("TwitterAccount")
	Attributes(func() {
		Attribute("input_id", Integer)
		Attribute("id", Integer)
		Attribute("screen_name", String)
		Required("input_id", "id", "screen_name")
	})
	View("default", func() {
		Attribute("input_id")
		Attribute("id")
		Attribute("screen_name")
	})
})

var TwitterAccounts = MediaType("application/vnd.app.twitter_accounts+json", func() {
	TypeName("TwitterAccounts")
	Attributes(func() {
		Attribute("twitter_accounts", CollectionOf(TwitterAccount))
		Required("twitter_accounts")
	})
	View("default", func() {
		Attribute("twitter_accounts")
	})
})

var _ = Resource("twitter", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("add", func() {
		Routing(POST("expeditions/:expedition_id/inputs/twitter-account"))
		Description("Add a Twitter account input")
		Params(func() {
			Param("expedition_id", Integer)
			Required("expedition_id")
		})
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
			Media(TwitterAccount)
		})
	})

	Action("list", func() {
		Routing(GET("projects/@/:project/expeditions/@/:expedition/inputs/twitter-accounts"))
		Description("List an expedition's Twitter account inputs")
		Params(func() {
			Param("project", String, func() {
				Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
				Description("Project slug")
			})
			Param("expedition", String, func() {
				Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
				Description("Expedition slug")
			})
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(TwitterAccounts)
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
			Media(TwitterAccounts)
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
