package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("csv", func() {
	Method("export", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("start", Int64)
			Attribute("end", Int64)
			Attribute("stations", String)
			Attribute("sensors", String)
			Attribute("resolution", Int32)
			Attribute("aggregate", String)
			Attribute("complete", Boolean)
			Attribute("tail", Int32)
		})

		Result(func() {
			Attribute("location", String)
			Required("location")
		})

		HTTP(func() {
			POST("export/csv")

			Params(func() {
				Param("start")
				Param("end")
				Param("stations")
				Param("sensors")
				Param("resolution")
				Param("aggregate")
				Param("complete")
				Param("tail")
			})

			Response(StatusFound, func() {
				Headers(func() {
					Header("location:Location")
				})
			})

			httpAuthentication()
		})
	})

	commonOptions()
})
