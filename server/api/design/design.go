package design

import (
	. "github.com/goadesign/goa/design/apidsl"
)

const (
	SlugPattern = `^[\da-z]+(?:-[\da-z]+)*$`
)

var cors = func() {
	Headers("Authorization", "Content-Type")
	Expose("Authorization", "Content-Type")
	Methods("GET", "OPTIONS", "POST", "DELETE", "PATCH", "PUT")
}

var _ = API("fieldkit", func() {
	Title("Fieldkit API")
	Description("A one-click open platform for field researchers and explorers.")
	Host("api.fieldkit.org")
	Scheme("https")
	Origin("https://fieldkit.org:8080", cors)
	Origin("https://*.fieldkit.org:8080", cors)
	Origin("https://fieldkit.org", cors)
	Origin("https://*.fieldkit.org", cors)
	Origin("https://fieldkit.team", cors)
	Origin("https://*.fieldkit.team", cors)
	Origin("https://fkdev.org", cors)
	Origin("https://*.fkdev.org", cors)
	Origin("/(.+[.])?localhost:\\d+/", cors)       // Dev
	Origin("/(.+[.])?fieldkit.org:\\d+/", cors)    // Dev
	Origin("/(.+[.])?local.fkdev.org:\\d+/", cors) // Dev
	Consumes("application/json")
	Produces("application/json")
})

var JWT = JWTSecurity("jwt", func() {
	Header("Authorization")
	Scope("api:access", "API access")
	Scope("api:admin", "API admin access")
})

var Location = MediaType("application/vnd.app.location+json", func() {
	TypeName("Location")
	Attributes(func() {
		Attribute("location", func() {
			Format("uri")
		})
		Required("location")
	})
	View("default", func() {
		Attribute("location")
	})
})

var (
	ProjectSlug = func() {
		Pattern(SlugPattern)
		Description("Project slug")
		MaxLength(40)
	}
	ExpeditionSlug = func() {
		Pattern(SlugPattern)
		Description("Expedition slug")
		MaxLength(40)
	}
	TeamSlug = func() {
		Pattern(SlugPattern)
		Description("Team slug")
		MaxLength(40)
	}
)
