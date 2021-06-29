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

var Origins = []string{
	"https://fieldkit.org:8080",
	"https://*.fieldkit.org:8080",
	"https://fieldkit.org",
	"https://*.fieldkit.org",
	"https://fkdev.org",
	"https://*.fkdev.org",
	"/(.+[.])?192.168.\\d+.\\d+:\\d+/", // Dev
	"/(.+[.])?127.0.0.1:\\d+/",         // Dev
	"/(.+[.])?localhost:\\d+/",         // Dev
	"/(.+[.])?fieldkit.org:\\d+/",      // Dev
	"/(.+[.])?local.fkdev.org:\\d+/",   // Dev
}

func commonOptions() {
	corsRules := func() {
		cors.Headers("Authorization", "Content-Type")
		cors.Expose("Authorization", "Content-Type")
		cors.Methods("GET", "OPTIONS", "POST", "DELETE", "PATCH", "PUT")
	}

	for _, origin := range Origins {
		cors.Origin(origin, corsRules)
	}

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
