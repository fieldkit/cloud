package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"time"

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
	if err := c.options.Database.GetContext(ctx, owner, `
		SELECT * FROM fieldkit.user WHERE id = $1
		`, p.UserID()); err != nil {
		return nil, err
	}

	sr, err := repositories.NewStationRepository(c.options.Database)
	if err != nil {
		return nil, err
	}

	stations, err := sr.QueryStationsByDeviceID(ctx, deviceId)
	if err != nil {
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
		OwnerID:      p.UserID(),
		Name:         payload.Name,
		DeviceID:     deviceId,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LocationName: payload.LocationName,
	}

	if payload.StatusPb != nil {
		if err := adding.UpdateFromStatus(*payload.StatusPb); err != nil {
			log.Infow("status error", "error", err, "status", payload.StatusPb)
		}
	}

	adding.SetStatus(payload.StatusJSON)

	if _, err := sr.Add(ctx, adding); err != nil {
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

	return transformStationFull(p, sf)
}

func (c *StationService) Update(ctx context.Context, payload *station.UpdatePayload) (response *station.StationFull, err error) {
	log := Logger(ctx).Sugar()

	r, err := repositories.NewStationRepository(c.options.Database)
	if err != nil {
		return nil, err
	}

	updating, err := r.QueryStationByID(ctx, payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, station.NotFound("station not found")
		}
		return nil, err
	}

	p, err := NewPermissions(ctx, c.options).ForStation(updating)
	if err != nil {
		return nil, err
	}

	if err := p.CanModify(); err != nil {
		return nil, err
	}

	updating.Name = payload.Name
	updating.UpdatedAt = time.Now()
	if payload.LocationName != nil {
		updating.LocationName = payload.LocationName
	}

	if payload.StatusPb != nil {
		if err := updating.UpdateFromStatus(*payload.StatusPb); err != nil {
			log.Infow("status error", "error", err, "status", payload.StatusPb)
		}
	}

	updating.SetStatus(payload.StatusJSON)

	if err := r.Update(ctx, updating); err != nil {
		return nil, err
	}

	return c.Get(ctx, &station.GetPayload{
		Auth: payload.Auth,
		ID:   payload.ID,
	})
}

func (c *StationService) ListMine(ctx context.Context, payload *station.ListMinePayload) (response *station.StationsFull, err error) {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	r, err := repositories.NewStationRepository(c.options.Database)
	if err != nil {
		return nil, err
	}

	sfs, err := r.QueryStationFullByOwnerID(ctx, p.UserID())
	if err != nil {
		return nil, err
	}

	stations, err := transformAllStationFull(p, sfs)
	if err != nil {
		return nil, err
	}

	response = &station.StationsFull{
		Stations: stations,
	}

	return
}

func (c *StationService) ListProject(ctx context.Context, payload *station.ListProjectPayload) (response *station.StationsFull, err error) {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	r, err := repositories.NewStationRepository(c.options.Database)
	if err != nil {
		return nil, err
	}

	sfs, err := r.QueryStationFullByProjectID(ctx, payload.ID)
	if err != nil {
		return nil, err
	}

	stations, err := transformAllStationFull(p, sfs)
	if err != nil {
		return nil, err
	}

	response = &station.StationsFull{
		Stations: stations,
	}

	return
}

func (c *StationService) Photo(ctx context.Context, payload *station.PhotoPayload) (*station.PhotoResult, io.ReadCloser, error) {
	x := uint(124)
	y := uint(100)

	allMedia := []*data.MediaForStation{}
	if err := c.options.Database.SelectContext(ctx, &allMedia, `
		SELECT s.id AS station_id, fnm.* FROM fieldkit.station AS s JOIN fieldkit.field_note AS fn ON (fn.station_id = s.id) JOIN fieldkit.field_note_media AS fnm ON (fn.media_id = fnm.id)
		WHERE s.id = $1 ORDER BY fnm.created DESC`, payload.ID); err != nil {
		return defaultPhoto(payload.ID)
	}

	if len(allMedia) == 0 {
		return defaultPhoto(payload.ID)
	}

	mr := repositories.NewMediaRepository(c.options.Session, c.options.Buckets.Media)

	lm, err := mr.LoadByURL(ctx, allMedia[0].URL)
	if err != nil {
		return defaultPhoto(payload.ID)
	}

	original, _, err := image.Decode(lm.Reader)
	if err != nil {
		return defaultPhoto(payload.ID)
	}

	cropped, err := smartCrop(original, x, y)
	if err != nil {
		return defaultPhoto(payload.ID)
	}

	options := jpeg.Options{
		Quality: 80,
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, cropped, &options); err != nil {
		return nil, nil, err
	}

	return &station.PhotoResult{
		Length:      int64(len(buf.Bytes())),
		ContentType: "image/jpg",
	}, ioutil.NopCloser(buf), nil
}

func (s *StationService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     func(m string) error { return station.NotFound(m) },
		Unauthorized: func(m string) error { return station.Unauthorized(m) },
		Forbidden:    func(m string) error { return station.Forbidden(m) },
	})
}

func transformImages(id int32, from []*data.MediaForStation) (to []*station.ImageRef) {
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

func transformReading(s *data.ModuleSensor) *station.SensorReading {
	if s.ReadingValue == nil {
		return nil
	}

	time := int64(0)
	if s.ReadingTime != nil {
		time = s.ReadingTime.Unix() * 1000
	}

	return &station.SensorReading{
		Last: float32(*s.ReadingValue),
		Time: time,
	}
}

func transformModules(from *data.StationFull) (to []*station.StationModule) {
	to = make([]*station.StationModule, 0)
	for _, v := range from.Modules {
		sensors := make([]*station.StationSensor, 0)

		for _, s := range from.Sensors {
			if s.ModuleID == v.ID {
				sensors = append(sensors, &station.StationSensor{
					Name:          s.Name,
					UnitOfMeasure: s.UnitOfMeasure,
					Reading:       transformReading(s),
				})
			}
		}

		hardwareID := hex.EncodeToString(v.HardwareID)

		to = append(to, &station.StationModule{
			ID:           v.ID,
			HardwareID:   &hardwareID,
			MetaRecordID: v.MetaRecordID,
			Name:         translateModuleName(v.Name, sensors),
			Position:     int32(v.Position),
			Flags:        int32(v.Flags),
			Internal:     v.Flags > 0,
			Sensors:      sensors,
		})

	}
	return
}

var (
	NameMap = map[string]string{
		"distance":    "module.distance",
		"weather":     "module.weather",
		"diagnostics": "module.diagnostics",
		"ultrasonic":  "module.distance",
	}
)

// This is going to go away eventually once no modules with old names
// are coming in. Then we can do a database migration and get rid of
// them completely. I'm deciding to leave them in so we can see them
// disappear over time.
func translateModuleName(old string, sensors []*station.StationSensor) string {
	if newName, ok := NameMap[old]; ok {
		return newName
	}

	if old == "water" {
		if len(sensors) == 1 {
			return "modules.water." + sensors[0].Name
		} else {
			return "modules.water.ec"
		}
	}

	return old
}

func transformStationFull(p Permissions, sf *data.StationFull) (*station.StationFull, error) {
	sp, err := p.ForStation(sf.Station)
	if err != nil {
		return nil, err
	}

	status, err := sf.Station.GetStatus()
	if err != nil {
		return nil, err
	}

	return &station.StationFull{
		ID:       sf.Station.ID,
		Name:     sf.Station.Name,
		ReadOnly: sp.IsReadOnly(),
		Owner: &station.StationOwner{
			ID:   sf.Owner.ID,
			Name: sf.Owner.Name,
		},
		DeviceID:        hex.EncodeToString(sf.Station.DeviceID),
		Uploads:         transformUploads(sf.Ingestions),
		Images:          transformImages(sf.Station.ID, sf.Media),
		Modules:         transformModules(sf),
		StatusJSON:      status,
		Battery:         sf.Station.Battery,
		MemoryUsed:      sf.Station.MemoryUsed,
		MemoryAvailable: sf.Station.MemoryAvailable,
		FirmwareNumber:  sf.Station.FirmwareNumber,
		FirmwareTime:    sf.Station.FirmwareTime,
		Photos: &station.StationPhotos{
			Small: fmt.Sprintf("/stations/%d/photo", sf.Station.ID),
		},
	}, nil
}

func transformAllStationFull(p Permissions, sfs []*data.StationFull) ([]*station.StationFull, error) {
	stations := make([]*station.StationFull, 0)

	for _, sf := range sfs {
		after, err := transformStationFull(p, sf)
		if err != nil {
			return nil, err
		}

		stations = append(stations, after)
	}

	return stations, nil
}

func defaultPhoto(id int32) (*station.PhotoResult, io.ReadCloser, error) {
	defaultPhotoContentType := "image/png"
	defaultPhoto, err := StationDefaultPicture(int64(id))
	if err != nil {
		// NOTE This, hopefully never happens because we've got no image to send back.
		return nil, nil, err
	}

	return &station.PhotoResult{
		ContentType: defaultPhotoContentType,
		Length:      int64(len(defaultPhoto)),
	}, ioutil.NopCloser(bytes.NewBuffer(defaultPhoto)), nil
}
