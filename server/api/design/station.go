package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var AddStationPayload = Type("AddStationPayload", func(){
  Attribute("name", String)
	Attribute("user_id", Integer)
	Required("name", "user_id")
})

var Station = MediaType("application/vnd.app.station+json", func() {
  TypeName("Station")
  Reference(AddStationPayload)
  Attributes(func() {
		Attribute("ID", Integer)
    Attribute("name")
		Attribute("user_id", Integer)
    Required("ID", "name", "user_id")
  })
  View("default", func(){
		Attribute("ID")
    Attribute("name")
		Attribute("user_id")
  })
})

var Stations = MediaType("application/vnd.app.stations+json", func() {
  TypeName("Stations")
  Attributes(func(){
    Attribute("stations", CollectionOf(Station))
    Required("stations")
  })
  View("default", func(){
    Attribute("stations")
  })
})

var _ = Resource("station", func(){
  Security(JWT, func(){
    Scope("api:access")
  })

  Action("add", func(){
    Routing(POST("stations"))
    Description("Add a station")
    Payload(AddStationPayload)
    Response(BadRequest)
    Response(OK, func(){
      Media(Station)
    })
  })

  Action("update", func(){
    Routing(PATCH("stations/:stationId"))
    Description("Update a station")
    Params(func() {
      Param("stationId", Integer)
      Required("stationId")
    })
    Payload(AddStationPayload)
    Response(BadRequest)
    Response(OK, func() {
      Media(Station)
    })
  })

  Action("get", func() {
    NoSecurity()
    Routing(GET("stations/@/:station"))
    Description("Get a station")
    Params(func(){
      Param("station", String)
    })
    Response(BadRequest)
    Response(OK, func(){
      Media(Station)
    })
  })

	Action("list", func() {
		Routing(GET("stations"))
		Description("List stations")
		Response(BadRequest)
		Response(OK, func() {
			Media(Stations)
		})
	})
})
