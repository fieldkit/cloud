package api

import (
	"encoding/hex"
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

func StationType(station *data.Station) (*app.Station, error) {
	status, err := station.GetStatus()
	if err != nil {
		return nil, err
	}

	return &app.Station{
		ID:         int(station.ID),
		OwnerID:    int(station.OwnerID),
		DeviceID:   hex.EncodeToString(station.DeviceID),
		Name:       station.Name,
		StatusJSON: status,
	}, nil
}

func StationsType(stations []*data.Station) (*app.Stations, error) {
	stationsCollection := make([]*app.Station, len(stations))

	for i, station := range stations {
		appStation, err := StationType(station)
		if err != nil {
			return nil, err
		}
		stationsCollection[i] = appStation
	}

	return &app.Stations{
		Stations: stationsCollection,
	}, nil
}

type StationController struct {
	*goa.Controller
	options StationControllerOptions
}

func NewStationController(service *goa.Service, options StationControllerOptions) *StationController {
	return &StationController{
		Controller: service.NewController("StationController"),
		options:    options,
	}
}

func (c *StationController) Add(ctx *app.AddStationContext) error {
	p, err := NewPermissions(ctx)
	if err != nil {
		return err
	}

	deviceId, err := hex.DecodeString(ctx.Payload.DeviceID)
	if err != nil {
		return err
	}

	stations := []*data.Station{}
	if err := c.options.Database.SelectContext(ctx, &stations, "SELECT * FROM fieldkit.station WHERE device_id = $1", deviceId); err != nil {
		return err
	}

	if len(stations) > 0 {
		existing := stations[0]
		if existing.OwnerID != p.UserID {
			return ctx.BadRequest(&app.BadRequestResponse{
				Key:     "stationAlreadyRegistered",
				Message: "This station is already registered.",
			})
		}
		svm, err := StationType(existing)
		if err != nil {
			return err
		}
		return ctx.OK(svm)
	}

	station := &data.Station{
		Name:     ctx.Payload.Name,
		OwnerID:  p.UserID,
		DeviceID: deviceId,
	}

	station.SetStatus(ctx.Payload.StatusJSON)

	if err := c.options.Database.NamedGetContext(ctx, station, "INSERT INTO fieldkit.station (name, device_id, owner_id, status_json) VALUES (:name, :device_id, :owner_id, :status_json) RETURNING *", station); err != nil {
		return err
	}

	svm, err := StationType(station)
	if err != nil {
		return err
	}
	return ctx.OK(svm)
}

func (c *StationController) Update(ctx *app.UpdateStationContext) error {
	p, err := NewPermissions(ctx)
	if err != nil {
		return err
	}

	err = p.CanModifyStationByStationID(int32(ctx.StationID))
	if err != nil {
		return err
	}

	station := &data.Station{
		ID:   int32(ctx.StationID),
		Name: ctx.Payload.Name,
	}

	station.SetStatus(ctx.Payload.StatusJSON)

	if err := c.options.Database.NamedGetContext(ctx, station, "UPDATE fieldkit.station SET name = :name, status_json = :status_json WHERE id = :id RETURNING *", station); err != nil {
		return err
	}

	svm, err := StationType(station)
	if err != nil {
		return err
	}
	return ctx.OK(svm)
}

func (c *StationController) Get(ctx *app.GetStationContext) error {
	p, err := NewPermissions(ctx)
	if err != nil {
		return err
	}

	err = p.CanViewStationByStationID(int32(ctx.StationID))
	if err != nil {
		return err
	}

	station := &data.Station{}
	if err := c.options.Database.GetContext(ctx, station, "SELECT * FROM fieldkit.station WHERE id = $1", ctx.StationID); err != nil {
		return err
	}

	svm, err := StationType(station)
	if err != nil {
		return err
	}
	return ctx.OK(svm)
}

func (c *StationController) List(ctx *app.ListStationContext) error {
	token := jwt.ContextJWT(ctx)
	if token == nil {
		return fmt.Errorf("JWT token is missing from context")
	}

	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return fmt.Errorf("JWT claims error")
	}

	stations := []*data.Station{}
	if err := c.options.Database.SelectContext(ctx, &stations, "SELECT * FROM fieldkit.station WHERE owner_id = $1", claims["sub"]); err != nil {
		return err
	}

	stationsWm, err := StationsType(stations)
	if err != nil {
		return err
	}
	return ctx.OK(stationsWm)
}

func (c *StationController) Delete(ctx *app.DeleteStationContext) error {
	p, err := NewPermissions(ctx)
	if err != nil {
		return err
	}

	err = p.CanModifyStationByStationID(int32(ctx.StationID))
	if err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, "DELETE FROM fieldkit.station WHERE id = $1", ctx.StationID); err != nil {
		return err
	}

	return ctx.OK()
}
