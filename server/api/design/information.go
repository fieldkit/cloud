package design

import (
	. "goa.design/goa/v3/dsl"
)

var DeviceLayoutResponse = ResultType("application/vnd.app.data.device.layout", func() {
	TypeName("DeviceLayoutResponse")
	Attributes(func() {
		Attribute("configurations", ArrayOf(StationConfiguration))
		Attribute("sensors", MapOf(String, ArrayOf(StationSensor)))
		Required("configurations")
		Required("sensors")
	})
	View("default", func() {
		Attribute("configurations")
		Attribute("sensors")
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

	Method("firmware statistics", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
		})

		Result(func() {
			Attribute("object", Any)
			Required("object")
		})

		HTTP(func() {
			GET("devices/firmware/statistics")

			Response(func() {
				Body("object")
			})

			httpAuthentication()
		})
	})

	commonOptions()
})
