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

type StationControllerOptions struct {
  Database *sqlxcache.DB
}

func StationType(station *data.Station) *app.Station {
  return &app.Station{
    ID:         int(station.ID),
    Name:       station.Name,
    UserID:     int(station.UserID),
  }
}

func StationsType(stations []*data.Station) *app.Stations {
  stationsCollection := make([]*app.Station, len(stations))
  for i, station := range stations {
    stationsCollection[i] = StationType(station)
  }

  return &app.Stations{
    Stations: stationsCollection,
  }
}

//StationController
type StationController struct {
  *goa.Controller
  options StationControllerOptions
}

func NewStationController(service *goa.Service, options StationControllerOptions) * StationController {
  return &StationController{
    Controller: service.NewController("StationController"),
    options:    options,
  }
}

func (c *StationController) Add(ctx *app.AddStationContext) error{
  token := jwt.ContextJWT(ctx)
  if token == nil {
    return fmt.Errorf("JWT token is missing from context")
  }

  claims, ok := token.Claims.(jwtgo.MapClaims)
  if !ok {
    return fmt.Errorf("JWT claims error")
  }

  station := &data.Station{
    Name:    ctx.Payload.Name,
    UserID:  int32(ctx.Payload.UserID),
  }

  if err := c.options.Database.NamedGetContext(ctx, station, "INSERT INTO fieldkit.station (name, user_id) VALUES (:name, :user_id) RETURNING *", station); err != nil {
    return err
  }

  if _, err := c.options.Database.ExecContext(ctx, "INSERT INTO fieldkit.station_user (station_id, user_id) VALUES ($1, $2)", station.ID, claims["sub"]); err != nil {
    return err
  }

  return ctx.OK(StationType(station))
}

func (c *StationController) Update(ctx *app.UpdateStationContext) error {
  station := &data.Station{
    ID:         int32(ctx.StationID),
    Name:       ctx.Payload.Name,
  }

  if err := c.options.Database.NamedGetContext(ctx, station, "UPDATE fieldkit.station SET name = :name, WHERE id = :id RETURNING *", station); err != nil {
		return err
	}

	return ctx.OK(StationType(station))
}

//what does $1 mean?
func (c *StationController) Get(ctx *app.GetStationContext) error {
  station := &data.Station{}
  if err := c.options.Database.GetContext(ctx, station, "SELECT * FROM fieldkit.station WHERE name = $1", ctx.Station); err != nil {
    return err
  }

  return ctx.OK(StationType(station))
}
