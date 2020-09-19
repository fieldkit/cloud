package design

import (
	. "goa.design/goa/v3/dsl"
)

var AddCollectionFields = Type("AddCollectionFields", func() {
	Attribute("name", String)
	Attribute("description", String)
	Attribute("tags", String)
	Attribute("privacy", Int32)
	Required("name", "description")
})

var Collection = ResultType("application/vnd.app.collection+json", func() {
	TypeName("Collection")
	Reference(AddCollectionFields)
	Attributes(func() {
		Attribute("id", Int32)
		Attribute("name")
		Attribute("description")
		Attribute("tags")
		Attribute("privacy", Int32)
		Required("id", "name", "description", "privacy", "tags")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("description")
		Attribute("tags")
		Attribute("privacy")
	})
})

var Collections = ResultType("application/vnd.app.collections+json", func() {
	TypeName("Collections")
	Attributes(func() {
		Attribute("collections", CollectionOf(Collection))
		Required("collections")
	})
	View("default", func() {
		Attribute("collections")
	})
})

var _ = Service("collection", func() {
	Method("add", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("collection", AddCollectionFields)
			Required("collection")
		})

		Result(Collection)

		HTTP(func() {
			POST("collections")

			Body("collection")

			httpAuthentication()
		})
	})

	Method("update", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("collectionId", Int32)
			Required("collectionId")
			Attribute("collection", AddCollectionFields)
			Required("collection")
		})

		Result(Collection)

		HTTP(func() {
			PATCH("collections/{collectionId}")

			Body("collection")

			httpAuthentication()
		})
	})

	Method("get", func() {
		Security(JWTAuth, func() {
			// Optional
		})

		Payload(func() {
			Token("auth")
			Attribute("collectionId", Int32)
			Required("collectionId")
		})

		Result(Collection)

		HTTP(func() {
			GET("collections/{collectionId}")

			httpAuthentication()
		})
	})

	Method("list mine", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
		})

		Result(Collections)

		HTTP(func() {
			GET("user/collections")

			httpAuthentication()
		})
	})

	Method("add station", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("collectionId", Int32)
			Required("collectionId")
			Attribute("stationId", Int32)
			Required("stationId")
		})

		HTTP(func() {
			POST("collections/{collectionId}/stations/{stationId}")

			httpAuthentication()
		})
	})

	Method("remove station", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("collectionId", Int32)
			Required("collectionId")
			Attribute("stationId", Int32)
			Required("stationId")
		})

		HTTP(func() {
			DELETE("collections/{collectionId}/stations/{stationId}")

			httpAuthentication()
		})
	})

	Method("delete", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("collectionId", Int32)
			Required("collectionId")
		})

		HTTP(func() {
			DELETE("collections/{collectionId}")

			httpAuthentication()
		})
	})

	commonOptions()
})
