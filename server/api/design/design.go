package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var cors = func() {
	Credentials()
	Methods("GET", "OPTIONS", "POST")
}

var _ = API("fieldkit", func() {
	Title("Fieldkit API")
	Description("A one-click open platform for field researchers and explorers.")
	Host("api.data.fieldkit.org")
	Scheme("https")
	Origin("https://fieldkit.org", cors)
	Origin("https://*.fieldkit.org", cors)
	Origin("http://localhost:3000", cors)
	Consumes("application/json")
	Produces("application/json")
})

var JWT = JWTSecurity("jwt", func() {
	Header("Authorization")
	Scope("api:access", "API access") // Define "api:access" scope
})
