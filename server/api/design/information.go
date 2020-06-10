package design

import (
	. "goa.design/goa/v3/dsl"
)

var StationConfiguration = Type("StationConfiguration", func() {
	Attribute("id", Int64)
	Attribute("time", Int64)
	Attribute("provision_id", Int64)
	Attribute("meta_record_id", Int64)
	Attribute("source_id", Int32)
	Attribute("modules", ArrayOf(StationModule))
	Required("id")
	Required("provision_id")
	Required("time")
	Required("modules")
})

var DeviceLayoutResponse = ResultType("application/vnd.app.data.device.layout", func() {
	TypeName("DeviceLayoutResponse")
	Attributes(func() {
		Attribute("configurations", ArrayOf(StationConfiguration))
		Required("configurations")
	})
	View("default", func() {
		Attribute("configurations")
	})
})

var _ = Service("information", func() {
	Method("device layout", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("deviceId", String)
			Required("deviceId")
		})

		Result(DeviceLayoutResponse)

		HTTP(func() {
			GET("data/devices/{deviceId}/layout")

			httpAuthentication()
		})
	})

	commonOptions()
})
