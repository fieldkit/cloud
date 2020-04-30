package api

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"

	"goa.design/goa/v3/security"

	station "github.com/fieldkit/cloud/server/api/gen/station"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

type StationService struct {
	options *ControllerOptions
}

func NewStationService(ctx context.Context, options *ControllerOptions) *StationService {
	return &StationService{
		options: options,
	}
}

func (c *StationService) Add(ctx context.Context, payload *station.AddPayload) (response *station.StationFull, err error) {
	log := Logger(ctx).Sugar()

	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	deviceId, err := hex.DecodeString(payload.DeviceID)
	if err != nil {
		return nil, err
	}

	log.Infow("adding station", "device_id", payload.DeviceID)

	owner := &data.User{}
	if err := c.options.Database.GetContext(ctx, owner, `SELECT * FROM fieldkit.user WHERE id = $1`, p.UserID()); err != nil {
		return nil, err
	}

	stations := []*data.Station{}
	if err := c.options.Database.SelectContext(ctx, &stations, `SELECT * FROM fieldkit.station WHERE device_id = $1`, deviceId); err != nil {
		return nil, err
	}

	if len(stations) > 0 {
		existing := stations[0]

		if existing.OwnerID != p.UserID() {
			return nil, station.BadRequest("station already registered to another user")
		}

		return c.Get(ctx, &station.GetPayload{
			Auth: payload.Auth,
			ID:   existing.ID,
		})
	}

	adding := &data.Station{
		OwnerID:  p.UserID(),
		Name:     payload.Name,
		DeviceID: deviceId,
	}

	adding.SetStatus(payload.StatusJSON)

	if err := c.options.Database.NamedGetContext(ctx, adding, `INSERT INTO fieldkit.station (name, device_id, owner_id, status_json) VALUES (:name, :device_id, :owner_id, :status_json) RETURNING *`, adding); err != nil {
		return nil, err
	}

	return c.Get(ctx, &station.GetPayload{
		Auth: payload.Auth,
		ID:   adding.ID,
	})
}

func (c *StationService) Get(ctx context.Context, payload *station.GetPayload) (response *station.StationFull, err error) {
	p, err := NewPermissions(ctx, c.options).ForStationByID(int(payload.ID))
	if err != nil {
		return nil, err
	}

	if err := p.CanView(); err != nil {
		return nil, err
	}

	r, err := repositories.NewStationRepository(c.options.Database)
	if err != nil {
		return nil, err
	}

	sf, err := r.QueryStationFull(ctx, payload.ID)
	if err != nil {
		return nil, err
	}

	response = &station.StationFull{
		ID:       sf.Station.ID,
		Name:     sf.Station.Name,
		ReadOnly: p.IsReadOnly(),
		Owner: &station.StationOwner{
			ID:   sf.Owner.ID,
			Name: sf.Owner.Name,
		},
		DeviceID: hex.EncodeToString(sf.Station.DeviceID),
		Uploads:  transformUploads(sf.Ingestions),
		Images:   transformImages(sf.Station.ID, sf.Media),
		Photos: &station.StationPhotos{
			Small: fmt.Sprintf("/stations/%d/photo", sf.Station.ID),
		},
	}

	return
}

func (c *StationService) Update(ctx context.Context, payload *station.UpdatePayload) (response *station.StationFull, err error) {
	p, err := NewPermissions(ctx, c.options).ForStationByID(int(payload.ID))
	if err != nil {
		return nil, err
	}

	if err := p.CanModify(); err != nil {
		return nil, err
	}

	updated := &data.Station{
		ID:   payload.ID,
		Name: payload.Name,
	}

	updated.SetStatus(payload.StatusJSON)

	if err := c.options.Database.NamedGetContext(ctx, updated, "UPDATE fieldkit.station SET name = :name, status_json = :status_json WHERE id = :id RETURNING *", updated); err != nil {
		if err == sql.ErrNoRows {
			return nil, station.NotFound("station not found")
		}
		return nil, err
	}

	return c.Get(ctx, &station.GetPayload{
		Auth: payload.Auth,
		ID:   payload.ID,
	})
}

func (s *StationService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:         token,
		Scheme:        scheme,
		Key:           s.options.JWTHMACKey,
		InvalidToken:  station.Unauthorized("invalid token"),
		InvalidScopes: station.Unauthorized("invalid scopes in token"),
	})
}

func transformImages(id int32, from []*data.FieldNoteMediaForStation) (to []*station.ImageRef) {
	to = make([]*station.ImageRef, 0, len(from))
	for _, v := range from {
		to = append(to, &station.ImageRef{
			URL: fmt.Sprintf("/stations/%d/field-note-media/%d", id, v.ID),
		})
	}
	return to
}

func transformUploads(from []*data.Ingestion) (to []*station.StationUpload) {
	to = make([]*station.StationUpload, 0, len(from))
	for _, v := range from {
		to = append(to, &station.StationUpload{
			ID:       v.ID,
			Time:     v.Time.Unix() * 1000,
			UploadID: v.UploadID,
			Size:     v.Size,
			Type:     v.Type,
			URL:      v.URL,
			Blocks:   v.Blocks.ToInt64Array(),
		})
	}
	return to
}
