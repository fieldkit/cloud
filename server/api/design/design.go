package design

import (
	cors "goa.design/plugins/v3/cors/dsl"

	. "goa.design/goa/v3/dsl"
)

var ServiceError = Type("ServiceError", func() {
	Attribute("code", String)
	Attribute("status", Int32)
})

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
	cors.Origin("http://192.168.\\d+.\\d+:\\d+/", corsRules) // Dev
	cors.Origin("/(.+[.])?127.0.0.1:\\d+/", corsRules)       // Dev
	cors.Origin("/(.+[.])?localhost:\\d+/", corsRules)       // Dev
	cors.Origin("/(.+[.])?fieldkit.org:\\d+/", corsRules)    // Dev
	cors.Origin("/(.+[.])?local.fkdev.org:\\d+/", corsRules) // Dev

	Error("unauthorized", func() {})
	Error("forbidden", func() {})
	Error("not-found", func() {})
	Error("bad-request", func() {})

	HTTP(func() {
		Response("unauthorized", StatusUnauthorized)
		Response("forbidden", StatusForbidden)
		Response("not-found", StatusNotFound)
		Response("bad-request", StatusBadRequest)
	})
}

func httpAuthentication() {
	Header("auth:Authorization", String, "authentication token", func() {
		Pattern("^Bearer [^ ]+$")
	})
}

func httpAuthenticationQueryString() {
	Param("auth:token", String, "authentication token")
}

func DateTimeFormatting() {
	Format(FormatDateTime)
}

var DownloadedPhoto = ResultType("application/vnd.app.photo+json", func() {
	TypeName("DownloadedPhoto")
	Attributes(func() {
		Attribute("length", Int64)
		Required("length")
		Attribute("contentType", String)
		Required("contentType")
		Attribute("etag", String)
		Required("etag")
		Attribute("body", Bytes)
	})
	View("default", func() {
		Attribute("length")
		Attribute("body")
		Attribute("contentType")
		Attribute("etag")
	})
})
