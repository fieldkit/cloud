package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var AddTwitterAccountSourcePayload = Type("AddTwitterAccountSourcePayload", func() {
	Reference(Source)
	Attribute("name")
	Required("name")
})

var UpdateTwitterAccountSourcePayload = Type("UpdateTwitterAccountSourcePayload", func() {
	Reference(Source)
	Attribute("name")
	Attribute("teamId")
	Attribute("userId")
})

var TwitterAccountSource = MediaType("application/vnd.app.twitter_account_source+json", func() {
	TypeName("TwitterAccountSource")
	Reference(Source)
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

var TwitterAccountSources = MediaType("application/vnd.app.twitter_account_intputs+json", func() {
	TypeName("TwitterAccountSources")
	Attributes(func() {
		Attribute("twitterAccountSources", CollectionOf(TwitterAccountSource))
		Required("twitterAccountSources")
	})
	View("default", func() {
		Attribute("twitterAccountSources")
	})
})

var _ = Resource("twitter", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("add", func() {
		Routing(POST("expeditions/:expeditionId/sources/twitter-accounts"))
		Description("Add a Twitter account source")
		Params(func() {
			Param("expeditionId", Integer)
			Required("expeditionId")
		})
		Payload(AddTwitterAccountSourcePayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(Location)
		})
	})

	Action("get id", func() {
		Routing(GET("sources/twitter-accounts/:sourceId"))
		Description("Get a Twitter account source")
		Params(func() {
			Param("sourceId", Integer)
			Required("sourceId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(TwitterAccountSource)
		})
	})

	Action("list", func() {
		Routing(GET("projects/@/:project/expeditions/@/:expedition/sources/twitter-accounts"))
		Description("List an expedition's Twitter account sources")
		Params(func() {
			Param("project", String, ProjectSlug)
			Param("expedition", String, ExpeditionSlug)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(TwitterAccountSources)
		})
	})

	Action("list id", func() {
		Routing(GET("expeditions/:expeditionId/sources/twitter-accounts"))
		Description("List an expedition's Twitter account sources")
		Params(func() {
			Param("expeditionId", Integer)
			Required("expeditionId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(TwitterAccountSources)
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
