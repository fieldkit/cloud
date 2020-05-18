package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("ingestion", func() {
	Method("process pending", func() {
		Security(JWTAuth, func() {
			Scope("api:admin")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
		})

		HTTP(func() {
			POST("data/process")

			httpAuthentication()
		})
	})

	Method("process station", func() {
		Security(JWTAuth, func() {
			Scope("api:admin")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("stationId", Int32)
			Required("stationId")
		})

		HTTP(func() {
			POST("data/stations/:stationId/process")

			httpAuthentication()
		})
	})

	Method("process ingestion", func() {
		Security(JWTAuth, func() {
			Scope("api:admin")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("ingestionId", Int64)
			Required("ingestionId")
		})

		HTTP(func() {
			POST("data/ingestions/:ingestionId/process")

			httpAuthentication()
		})
	})

	Method("delete", func() {
		Security(JWTAuth, func() {
			Scope("api:admin")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("ingestionId", Int64)
			Required("ingestionId")
		})

		HTTP(func() {
			DELETE("data/ingestions/:ingestionId")

			httpAuthentication()
		})
	})

	commonOptions()
})
