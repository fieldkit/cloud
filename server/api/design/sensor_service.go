package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("sensor", func() {
	Method("meta", func() {
		Result(func() {
			Attribute("object", Any)
			Required("object")
		})

		HTTP(func() {
			GET("sensors")

			Response(func() {
				Body("object")
			})
		})
	})

	Method("data", func() {
		Security(JWTAuth, func() {
			// Optional
		})

		Payload(func() {
			Token("auth")
			Attribute("start", Int64)
			Attribute("end", Int64)
			Attribute("stations", String)
			Attribute("sensors", String)
			Attribute("resolution", Int32)
			Attribute("aggregate", String)
			Attribute("complete", Boolean)
			Attribute("tail", Int32)
			Attribute("influx", Boolean)
		})

		Result(func() {
			Attribute("object", Any)
			Required("object")
		})

		HTTP(func() {
			GET("sensors/data")

			Params(func() {
				Param("start")
				Param("end")
				Param("stations")
				Param("sensors")
				Param("resolution")
				Param("aggregate")
				Param("complete")
				Param("tail")
				Param("influx")
			})

			Response(func() {
				Body("object")
			})

			httpAuthentication()
		})
	})

	Method("bookmark", func() {
		Security(JWTAuth, func() {
			// Optional
		})

		Payload(func() {
			Token("auth")
			Attribute("bookmark", String)
			Required("bookmark")
		})

		Result(SavedBookmark)

		HTTP(func() {
			POST("bookmarks/save")

			Params(func() {
				Param("bookmark")
			})

			httpAuthentication()
		})
	})

	Method("resolve", func() {
		Security(JWTAuth, func() {
			// Optional
		})

		Payload(func() {
			Token("auth")
			Attribute("v", String)
			Required("v")
		})

		Result(SavedBookmark)

		HTTP(func() {
			GET("bookmarks/resolve")

			Params(func() {
				Param("v")
			})

			httpAuthentication()
		})
	})

	commonOptions()
})

var SavedBookmark = ResultType("application/vnd.app.bookmark", func() {
	TypeName("SavedBookmark")
	Attributes(func() {
		Attribute("url", String)
		Attribute("bookmark", String)
		Attribute("token", String)
		Required("url", "bookmark", "token")
	})
	View("default", func() {
		Attribute("url")
		Attribute("bookmark")
		Attribute("token")
	})
})
