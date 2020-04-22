package design

import (
	cors "goa.design/plugins/v3/cors/dsl"

	. "goa.design/goa/v3/dsl"
)

var JWTAuth = JWTSecurity("jwt", func() {
	Scope("api:access", "API access")
	Scope("api:admin", "API admin access")
	Scope("api:ingestion", "Ingestion access")
})

func commonOptions() {
	corsRules := func() {
		cors.Headers("Authorization", "Content-Type")
		cors.Expose("Authorization", "Content-Type")
		cors.Methods("GET", "OPTIONS", "POST", "DELETE", "PATCH", "PUT")
	}

	cors.Origin("https://fieldkit.org:8080", corsRules)
	cors.Origin("https://*.fieldkit.org:8080", corsRules)
	cors.Origin("https://fieldkit.org", corsRules)
	cors.Origin("https://*.fieldkit.org", corsRules)
	cors.Origin("https://fkdev.org", corsRules)
	cors.Origin("https://*.fkdev.org", corsRules)
	cors.Origin("/(.+[.])?localhost:\\d+/", corsRules)       // Dev
	cors.Origin("/(.+[.])?fieldkit.org:\\d+/", corsRules)    // Dev
	cors.Origin("/(.+[.])?local.fkdev.org:\\d+/", corsRules) // Dev
}

func DateTimeFormatting() {
	Format(FormatDateTime)
}
