package design

import (
	. "github.com/goadesign/goa/design/apidsl"
)

var DeviceInput = MediaType("application/vnd.app.device_input+json", func() {
	TypeName("DeviceInput")
	Reference(Input)
	Attributes(func() {
		Attribute("id")
		Attribute("token")
		Attribute("key")
	})
	View("default", func() {
		Attribute("id")
		Attribute("token")
		Attribute("key")
	})
})

var DeviceInputs = MediaType("application/vnd.app.device_inputs+json", func() {
	TypeName("DeviceInputs")
	Attributes(func() {
		Attribute("device_inputs", CollectionOf(DeviceInput))
		Required("device_inputs")
	})
	View("default", func() {
		Attribute("device_inputs")
	})
})
