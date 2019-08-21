package design

import (
  . "github.com/goadesign/goa/design"
  . "github.com/goadesign/goa/design/apidsl"
)

var AddStationLogPayload = Type("AddStationLogPayload", func (){
  Attribute("ID", Integer)
  Attribute("station_id", Integer)
  Attribute("body", String)
  Attribute("timestamp", String)
  Required("ID", "station_id", "body", "timestamp")
})

var StationLog = MediaType("application/vnd.app.stationlog+json", func() {
  TypeName("StationLog")
  Reference(AddStationLogPayload)
  Attributes(func() {
    Attribute("ID", Integer)
    Attribute("station_id")
    Attribute("body")
    Attribute("timestamp")
    Required("ID", "station_id", "body", "timestamp")
  })
  View("default", func(){
    Attribute("ID")
    Attribute("station_id")
    Attribute("body")
    Attribute("timestamp")
  })
})

var AddStationLogsPayload = Type("AddStationLogsPayload", func() {
  Attribute("station_logs", ArrayOf(AddStationLogPayload))
  Required("station_logs")
})

var StationLogs = MediaType("application/vnd.app.stationlogs+json", func() {
  TypeName("StationLogs")
  Attributes(func() {
    Attribute("station_logs", CollectionOf(StationLog))
    Required("station_logs")
  })
  View("default", func(){
    Attribute("station_logs")
  })
})

var _ = Resource("stationLog", func(){
  Security(JWT, func(){
    Scope("api: access")
  })

  Action("add", func(){
    Routing(POST("stationLog"))
    Description("Add a station log")
    Payload(AddStationLogPayload)
    Response(BadRequest)
    Response(OK, func(){
      Media(StationLog)
    })
  })

  Action("addMultiple", func() {
    Routing(POST("stationLogs"))
    Description("Add multiple station logs")
    Payload(AddStationLogsPayload)
    Response(BadRequest)
    Response(OK, func(){
    })
  })

  Action("update", func() {
    Routing(PATCH("stationlogs/:stationLogId"))
    Description("Update a station log")
    Params(func() {
      Param("stationLogId", Integer)
      Required("stationLogId")
    })
    Payload(AddStationLogPayload)
    Response(BadRequest)
    Response(OK, func() {
      Media(StationLog)
    })
  })

  Action("get", func() {
    NoSecurity()
    Routing(GET("stationlogs/@/:stationLog"))
    Description("Get a station log")
    Params(func() {
      Param("stationLog", String)
    })
    Response(BadRequest)
    Response(OK, func(){
      Media(StationLog)
    })
  })
})
