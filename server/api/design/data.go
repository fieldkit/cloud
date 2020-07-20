package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("data", func() {
	Security(JWTAuth, func() {
		Scope("api:access")
	})

	Method("device summary", func() {
		Result(DeviceDataSummaryResponse)

		Payload(func() {
			Token("auth")
			Attribute("deviceId", String)
			Required("deviceId")
		})

		HTTP(func() {
			GET("data/devices/{deviceId}/summary")
		})
	})

	commonOptions()
})

var DeviceStreamSummary = Type("DeviceStreamSummary", func() {
	Attribute("records", Int64)
	Attribute("size", Int64)
})

var DeviceDataStreamsSummary = Type("DeviceDataStreamsSummary", func() {
	Attribute("meta", DeviceStreamSummary)
	Attribute("data", DeviceStreamSummary)
})

var DeviceProvisionSummary = ResultType("application/vnd.app.device.provision.summary+json", func() {
	TypeName("DeviceProvisionSummary")
	Attributes(func() {
		Attribute("generation", String)
		Attribute("created", Int64)
		Attribute("updated", Int64)
		Attribute("meta", DeviceMetaSummary)
		Attribute("data", DeviceDataSummary)

		Required("generation")
		Required("created")
		Required("updated")
		Required("meta")
		Required("data")
	})
	View("default", func() {
		Attribute("generation")
		Attribute("created")
		Attribute("updated")
		Attribute("meta")
		Attribute("data")
	})
})

var DeviceMetaSummary = ResultType("application/vnd.app.device.meta.summary+json", func() {
	TypeName("DeviceMetaSummary")
	Attributes(func() {
		Attribute("size", Int64)
		Attribute("first", Int64)
		Attribute("last", Int64)
		Required("size")
		Required("first")
		Required("last")
	})
	View("default", func() {
		Attribute("size")
		Attribute("first")
		Attribute("last")
	})
})

var DeviceDataSummary = ResultType("application/vnd.app.device.data.summary+json", func() {
	TypeName("DeviceDataSummary")
	Attributes(func() {
		Attribute("size", Int64)
		Attribute("first", Int64)
		Attribute("last", Int64)
		Required("size")
		Required("first")
		Required("last")
	})
	View("default", func() {
		Attribute("size")
		Attribute("first")
		Attribute("last")
	})
})

var DeviceDataSummaryResponse = ResultType("application/vnd.app.device.summary+json", func() {
	TypeName("DeviceDataSummaryResponse")
	Attributes(func() {
		Attribute("provisions", CollectionOf(DeviceProvisionSummary))
		Required("provisions")
	})
	View("default", func() {
		Attribute("provisions")
	})
})
