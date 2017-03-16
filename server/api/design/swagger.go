package design

import (
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("swagger", func() {
	Origin("*", func() {
		Methods("GET", "OPTIONS")
	})
	Files("/swagger.json", "api/public/swagger/swagger.json")
	Files("/swagger.yaml", "api/public/swagger/swagger.yaml")
})
