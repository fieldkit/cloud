package api

import (
  "fmt"

  jwtgo "github.com/dgrijalva/jwt-go"
  "github.com/goadesign/goa"
  "github.com/goadesign/goa/middleware/security/jwt"

  "github.com/conservify/sqlxcache"

  "github.com/fieldkit/cloud/server/api/app"
  "github.com/fieldkit/cloud/server/data"
)

type StationLogControllerOptions struct {
  Database *sqlxcache.DB
}

func StationLogType(stationlog *data.StationLog) *app.StationLog {
  return &app.StationLog{
    ID:           int(stationlog.ID),
    Body:         stationlog.Body,
    Timestamp:    stationlog.Timestamp,
    StationID:    int(stationlog.StationId),
  }
}

func StationLogsType(stationlogs []*data.StationLog) *app.StationLogs {
  stationLogsCollection := make([]*app.StationLog, len(stationlogs))
  for i, stationlog := range stationlogs {
    stationLogsCollection[i] = StationLogType(stationlog)
  }

  return &app.StationLogs{
    StationLogs: stationLogsCollection,
  }
}

//StationLogController
type StationLogController struct {
  *goa.Controller
  options StationLogControllerOptions
}

func NewStationLogController(service *goa.Service, options StationLogControllerOptions) * StationLogController {
  return &StationLogController{
    Controller: service.NewController("StationLogController"),
    options:    options,
  }
}

func (c *StationLogController) Add(ctx *app.AddStationLogContext) error{
  token := jwt.ContextJWT(ctx)
  if token == nil {
    return fmt.Errorf("JWT token is missing from context")
  }

  claims, ok := token.Claims.(jwtgo.MapClaims)
  if !ok {
    return fmt.Errorf("JWT claims error")
  }

  stationLog := &data.StationLog{
    ID:         int32(ctx.Payload.ID),
    StationId:  int32(ctx.Payload.StationID),
    Body:       ctx.Payload.Body,
    Timestamp:  ctx.Payload.Timestamp,
  }

  if err := c.options.Database.NamedGetContext(ctx, stationLog, "INSERT INTO fieldkit.stationlog (name, body, timestamp, station_id) VALUES (:name, :body, :timestamp, :station_id) RETURNING *", stationLog); err != nil {
    return err
  }

  if _, err := c.options.Database.ExecContext(ctx, "INSERT INTO fieldkit.stationlog_user (stationlog_id, user_id) VALUES ($1, $2)", stationLog.ID, claims["sub"]); err != nil{
    return err
  }

  return ctx.OK(StationLogType(stationLog))
}

func (c *StationLogController) AddMultiple(ctx *app.AddMultipleStationLogContext) error {

  for _, add_log := range ctx.Payload.StationLogs {
    stationLog := &data.StationLog{
      StationId:  int32(add_log.StationID),
      Body:       add_log.Body,
      Timestamp:  add_log.Timestamp,
    }

    if err := c.options.Database.NamedGetContext(ctx, stationLog, "INSERT INTO fieldkit.stationLog (ID, stationId, Body, Timestamp) VALUES (:ID, :stationId, :Body, :Timestamp) RETURNING *", stationLog); err != nil {
      return err
    }

  }

  //return ctx.OK(StationLogsType(stationlogs))
  /*
  if _, err := options.Database.ExecContext(ctx, "INSERT INTO fieldkit.stationlog_user ") {
    return err
  }
  */
  return nil
}

func (c * StationLogController) Update(ctx *app.UpdateStationLogContext) error {
  stationlog := &data.StationLog{
    ID:   int32(ctx.StationLogID),
    StationId:  int32(ctx.Payload.StationID),
    Body:       ctx.Payload.Body,
    Timestamp:  ctx.Payload.Timestamp,
  }
  if err := c.options.Database.NamedGetContext(ctx, stationlog, "UPDATE fieldkit.stationlog SET name = :name, station_id = :station_id, body = :body, timestamp = :timestamp, WHERE id = :id RETURNING *", stationlog); err != nil {
    return err
  }

  return ctx.OK(StationLogType(stationlog))
}

func (c *StationLogController) Get(ctx *app.GetStationLogContext) error {
  stationlog := &data.StationLog{}
  if err := c.options.Database.GetContext(ctx, stationlog, "SELECT * FROM fieldkit.stationlog WHERE name = $1", ctx.StationLog); err != nil {
    return err
  }

  return ctx.OK(StationLogType(stationlog))
}
