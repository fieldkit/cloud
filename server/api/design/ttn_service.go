package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("ttn", func() {
	Method("webhook", func() {
		Security(JWTAuth, func() {
		})

		Payload(func() {
			Token("auth")
			Attribute("contentLength", Int64)
			Required("contentLength")
			Attribute("contentType", String)
			Required("contentType")
			Attribute("token", String)
		})

		HTTP(func() {
			POST("ttn/webhook")

			Header("contentType:Content-Type")
			Header("contentLength:Content-Length")

			Params(func() {
				Param("token")
			})

			SkipRequestBodyEncodeDecode()

			httpAuthentication()
		})
	})

	Error("unauthorized", String, "credentials are invalid")

	HTTP(func() {
		Response("unauthorized", StatusUnauthorized)
	})

	commonOptions()
})
