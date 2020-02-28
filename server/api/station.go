package api

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"

	"image"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware/security/jwt"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

func StationType(station *data.Station, ingestions []*data.Ingestion, note_media []*data.FieldNoteMediaForStation) (*app.Station, error) {
	status, err := station.GetStatus()
	if err != nil {
		return nil, err
	}

	images := make([]*app.ImageRef, 0)
	for _, row := range note_media {
		if row.URL != "" && strings.Contains(row.ContentType, "image") {
			images = append(images, &app.ImageRef{
				URL: fmt.Sprintf("/stations/%d/field-note-media/%d", station.ID, row.ID),
			})
		}
	}

	lastUploads := make([]*app.LastUpload, len(ingestions))
	for i, ingestion := range ingestions {
		lastUploads[i] = &app.LastUpload{
			ID:       int(ingestion.ID),
			Time:     ingestion.Time,
			UploadID: ingestion.UploadID,
			Size:     int(ingestion.Size),
			Type:     ingestion.Type,
			URL:      ingestion.URL,
			Blocks:   ingestion.Blocks.ToIntArray(),
		}
	}

	return &app.Station{
		ID:          int(station.ID),
		OwnerID:     int(station.OwnerID),
		DeviceID:    hex.EncodeToString(station.DeviceID),
		Name:        station.Name,
		LastUploads: lastUploads,
		StatusJSON:  status,
		Images:      images,
		Photos: &app.StationPhotos{
			Small: fmt.Sprintf("/stations/%d/photo", station.ID),
		},
	}, nil
}

func sortByStation(ingestions []*data.Ingestion) map[string][]*data.Ingestion {
	m := make(map[string][]*data.Ingestion)
	for _, i := range ingestions {
		key := hex.EncodeToString(i.DeviceID)
		if m[key] == nil {
			m[key] = make([]*data.Ingestion, 0)
		}
		m[key] = append(m[key], i)
	}
	return m
}

func sortMediaByStation(all []*data.FieldNoteMediaForStation) map[int32][]*data.FieldNoteMediaForStation {
	m := make(map[int32][]*data.FieldNoteMediaForStation)
	for _, row := range all {
		if m[row.StationID] == nil {
			m[row.StationID] = make([]*data.FieldNoteMediaForStation, 0)
		}
		m[row.StationID] = append(m[row.StationID], row)
	}
	return m
}

func StationsType(stations []*data.Station, ingestions []*data.Ingestion, note_media []*data.FieldNoteMediaForStation) (*app.Stations, error) {
	stationsCollection := make([]*app.Station, len(stations))

	noteMediaByStation := sortMediaByStation(note_media)
	ingestionsByStation := sortByStation(ingestions)
	for i, station := range stations {
		key := hex.EncodeToString(station.DeviceID)
		ingestionsByStation := ingestionsByStation[key]
		if ingestionsByStation == nil {
			ingestionsByStation = make([]*data.Ingestion, 0)
		}
		noteMediaByStation := noteMediaByStation[station.ID]
		if noteMediaByStation == nil {
			noteMediaByStation = make([]*data.FieldNoteMediaForStation, 0)
		}
		appStation, err := StationType(station, ingestionsByStation, noteMediaByStation)
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
	options *ControllerOptions
}

func NewStationController(service *goa.Service, options *ControllerOptions) *StationController {
	return &StationController{
		Controller: service.NewController("StationController"),
		options:    options,
	}
}

func (c *StationController) Add(ctx *app.AddStationContext) error {
	log := Logger(ctx).Sugar()

	p, err := NewPermissions(ctx)
	if err != nil {
		return err
	}

	deviceId, err := hex.DecodeString(ctx.Payload.DeviceID)
	if err != nil {
		return err
	}

	log.Infow("adding station", "device_id", ctx.Payload.DeviceID)

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
		svm, err := StationType(existing, make([]*data.Ingestion, 0), make([]*data.FieldNoteMediaForStation, 0))
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

	svm, err := StationType(station, make([]*data.Ingestion, 0), make([]*data.FieldNoteMediaForStation, 0))
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
		if err == sql.ErrNoRows {
			return ctx.NotFound()
		}
		return err
	}

	ingestions := []*data.Ingestion{}
	if err := c.options.Database.SelectContext(ctx, &ingestions, "SELECT * FROM fieldkit.ingestion WHERE device_id = $1 ORDER BY time DESC LIMIT 10", station.DeviceID); err != nil {
		return err
	}

	note_media := []*data.FieldNoteMediaForStation{}
	if err := c.options.Database.SelectContext(ctx, &note_media, `
		SELECT s.id AS station_id, fnm.* FROM fieldkit.station AS s JOIN fieldkit.field_note AS fn ON (fn.station_id = s.id) JOIN fieldkit.field_note_media AS fnm ON (fn.media_id = fnm.id)
		WHERE s.id = $1 ORDER BY fnm.created DESC`, ctx.StationID); err != nil {
		return err
	}

	svm, err := StationType(station, ingestions, note_media)
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

	ingestions := []*data.Ingestion{}
	if err := c.options.Database.SelectContext(ctx, &ingestions, "SELECT * FROM fieldkit.ingestion WHERE device_id = $1 ORDER BY time DESC LIMIT 10", station.DeviceID); err != nil {
		return err
	}

	note_media := []*data.FieldNoteMediaForStation{}
	if err := c.options.Database.SelectContext(ctx, &note_media, `
		SELECT s.id AS station_id, fnm.* FROM fieldkit.station AS s JOIN fieldkit.field_note AS fn ON (fn.station_id = s.id) JOIN fieldkit.field_note_media AS fnm ON (fn.media_id = fnm.id)
		WHERE s.id = $1 ORDER BY fnm.created DESC`, ctx.StationID); err != nil {
		return err
	}

	svm, err := StationType(station, ingestions, note_media)
	if err != nil {
		return err
	}
	return ctx.OK(svm)
}

func (c *StationController) ListProject(ctx *app.ListProjectStationContext) error {
	token := jwt.ContextJWT(ctx)
	if token == nil {
		return fmt.Errorf("JWT token is missing from context")
	}

	stations := []*data.Station{}
	if err := c.options.Database.SelectContext(ctx, &stations, "SELECT * FROM fieldkit.station WHERE id IN (SELECT station_id FROM fieldkit.project_station WHERE project_id = $1)", ctx.ProjectID); err != nil {
		return err
	}

	ingestions := []*data.Ingestion{}
	if err := c.options.Database.SelectContext(ctx, &ingestions, "SELECT * FROM fieldkit.ingestion WHERE device_id IN (SELECT s.device_id FROM fieldkit.station AS s JOIN fieldkit.project_station AS ps ON (s.id = ps.station_id) WHERE project_id = $1) ORDER BY time DESC", ctx.ProjectID); err != nil {
		return err
	}

	note_media := []*data.FieldNoteMediaForStation{}
	if err := c.options.Database.SelectContext(ctx, &note_media, `
		SELECT s.id AS station_id, fnm.* FROM fieldkit.station AS s JOIN fieldkit.field_note AS fn ON (fn.station_id = s.id) JOIN fieldkit.field_note_media AS fnm ON (fn.media_id = fnm.id)
		WHERE s.id IN (SELECT station_id FROM fieldkit.project_station WHERE project_id = $1) ORDER BY fnm.created DESC`, ctx.ProjectID); err != nil {
		return err
	}

	stationsWm, err := StationsType(stations, ingestions, note_media)
	if err != nil {
		return err
	}
	return ctx.OK(stationsWm)
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

	ownerId := claims["sub"]

	stations := []*data.Station{}
	if err := c.options.Database.SelectContext(ctx, &stations, "SELECT * FROM fieldkit.station WHERE owner_id = $1", ownerId); err != nil {
		return err
	}

	ingestions := []*data.Ingestion{}
	if err := c.options.Database.SelectContext(ctx, &ingestions, "SELECT * FROM fieldkit.ingestion WHERE device_id IN (SELECT device_id FROM fieldkit.station WHERE owner_id = $1) ORDER BY time DESC", ownerId); err != nil {
		return err
	}

	note_media := []*data.FieldNoteMediaForStation{}
	if err := c.options.Database.SelectContext(ctx, &note_media, `
		SELECT s.id AS station_id, fnm.* FROM fieldkit.station AS s JOIN fieldkit.field_note AS fn ON (fn.station_id = s.id) JOIN fieldkit.field_note_media AS fnm ON (fn.media_id = fnm.id)
		WHERE s.owner_id = $1 ORDER BY fnm.created DESC`, ownerId); err != nil {
		return err
	}

	stationsWm, err := StationsType(stations, ingestions, note_media)
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

func (c *StationController) Photo(ctx *app.PhotoStationContext) error {
	x := uint(124)
	y := uint(100)

	defaultPhotoContentType := "image/png"
	defaultPhoto, err := StationDefaultPicture(int64(ctx.StationID))
	if err != nil {
		// NOTE This, hopefully never happens because we've got no image to send back.
		return err
	}

	allMedia := []*data.FieldNoteMediaForStation{}
	if err := c.options.Database.SelectContext(ctx, &allMedia, `
		SELECT s.id AS station_id, fnm.* FROM fieldkit.station AS s JOIN fieldkit.field_note AS fn ON (fn.station_id = s.id) JOIN fieldkit.field_note_media AS fnm ON (fn.media_id = fnm.id)
		WHERE s.id = $1 ORDER BY fnm.created DESC`, ctx.StationID); err != nil {
		return LogErrorAndSendData(ctx, ctx.ResponseData, err, defaultPhotoContentType, defaultPhoto)
	}

	if len(allMedia) == 0 {
		return SendData(ctx.ResponseData, defaultPhotoContentType, defaultPhoto)
	}

	mr := repositories.NewMediaRepository(c.options.Session)

	lm, err := mr.LoadByURL(ctx, allMedia[0].URL)
	if err != nil {
		return LogErrorAndSendData(ctx, ctx.ResponseData, err, defaultPhotoContentType, defaultPhoto)
	}

	original, _, err := image.Decode(lm.Reader)
	if err != nil {
		return LogErrorAndSendData(ctx, ctx.ResponseData, err, defaultPhotoContentType, defaultPhoto)
	}

	cropped, err := SmartCrop(original, x, y)
	if err != nil {
		return LogErrorAndSendData(ctx, ctx.ResponseData, err, defaultPhotoContentType, defaultPhoto)
	}

	err = SendImage(ctx.ResponseData, cropped)
	if err != nil {
		return LogErrorAndSendData(ctx, ctx.ResponseData, err, defaultPhotoContentType, defaultPhoto)
	}

	return nil
}
