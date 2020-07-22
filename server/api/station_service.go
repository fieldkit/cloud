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
	"strings"
	"time"

	"github.com/iancoleman/strcase"

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

func (c *StationService) updateStation(ctx context.Context, station *data.Station, rawStatusPb *string) error {
	log := Logger(ctx).Sugar()

	sr, err := repositories.NewStationRepository(c.options.Database)
	if err != nil {
		return err
	}

	if rawStatusPb != nil {
		log.Infow("updating station from status", "station_id", station.ID)

		if err := station.UpdateFromStatus(ctx, *rawStatusPb); err != nil {
			log.Errorw("error updating from status", "error", err)
		}

		if err := sr.UpdateStationModelFromStatus(ctx, station, *rawStatusPb); err != nil {
			log.Errorw("error updating model status", "error", err)
		}

		if station.Location != nil && c.options.locations != nil {
			log.Infow("describing location", "station_id", station.ID)
			names, err := c.options.locations.Describe(ctx, station.Location)
			if err != nil {
				log.Errorw("error updating from location", "error", err)
			} else if names != nil {
				station.PlaceOther = names.OtherLandName
				station.PlaceNative = names.NativeLandName
			}
		}
	}

	if err := sr.Update(ctx, station); err != nil {
		return err
	}

	return nil
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
			log.Infow("station conflict", "device_id", deviceId, "user_id", p.UserID(), "owner_id", existing.OwnerID, "station_id", existing.ID)
			return nil, station.BadRequest("station already registered to another user")
		}

		if err := c.updateStation(ctx, existing, payload.StatusPb); err != nil {
			return nil, err
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

	added, err := sr.Add(ctx, adding)
	if err != nil {
		return nil, err
	}

	if err := c.updateStation(ctx, added, payload.StatusPb); err != nil {
		return nil, err
	}

	pr := repositories.NewProjectRepository(c.options.Database)
	if err != nil {
		return nil, err
	}

	if err := pr.AddStationToDefaultProjectMaybe(ctx, adding); err != nil {
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
	sr, err := repositories.NewStationRepository(c.options.Database)
	if err != nil {
		return nil, err
	}

	updating, err := sr.QueryStationByID(ctx, payload.ID)
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

	if err := c.updateStation(ctx, updating, payload.StatusPb); err != nil {
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

func (c *StationService) ListAll(ctx context.Context, payload *station.ListAllPayload) (*station.PageOfStations, error) {
	sr, err := repositories.NewStationRepository(c.options.Database)
	if err != nil {
		return nil, err
	}

	qp := &repositories.EssentialQueryParams{
		Page:     0,
		PageSize: 20,
	}

	if payload.Page != nil {
		qp.Page = *payload.Page
	}
	if payload.PageSize != nil {
		qp.PageSize = *payload.PageSize
	}

	queried, err := sr.QueryEssentialStations(ctx, qp)
	if err != nil {
		return nil, err
	}

	stationsWm := make([]*station.EssentialStation, 0)

	for _, es := range queried.Stations {
		recordingStartedAt := es.RecordingStartedAt
		var lastIngestionAt *int64

		if es.LastIngestionAt != nil {
			jsEpoch := es.LastIngestionAt.Unix() * 1000
			lastIngestionAt = &jsEpoch
		}

		var location *station.StationLocation = nil
		if es.Location != nil {
			location = &station.StationLocation{
				Latitude:  es.Location.Latitude(),
				Longitude: es.Location.Longitude(),
			}
		}

		stationsWm = append(stationsWm, &station.EssentialStation{
			ID:       es.ID,
			Name:     es.Name,
			DeviceID: hex.EncodeToString(es.DeviceID),
			Owner: &station.StationOwner{
				ID:   es.OwnerID,
				Name: es.OwnerName,
			},
			CreatedAt:          es.CreatedAt.Unix() * 1000,
			UpdatedAt:          es.UpdatedAt.Unix() * 1000,
			RecordingStartedAt: recordingStartedAt,
			MemoryUsed:         es.MemoryUsed,
			MemoryAvailable:    es.MemoryAvailable,
			FirmwareNumber:     es.FirmwareNumber,
			FirmwareTime:       es.FirmwareTime,
			Location:           location,
			LastIngestionAt:    lastIngestionAt,
		})
	}

	wm := &station.PageOfStations{
		Stations: stationsWm,
		Total:    queried.Total,
	}

	return wm, nil
}

func (c *StationService) Photo(ctx context.Context, payload *station.PhotoPayload) (*station.PhotoResult, io.ReadCloser, error) {
	x := uint(124)
	y := uint(100)

	allMedia := []*data.FieldNoteMedia{}
	if err := c.options.Database.SelectContext(ctx, &allMedia, `
		SELECT * FROM fieldkit.notes_media WHERE station_id = $1 ORDER BY created_at DESC
		`, payload.ID); err != nil {
		return defaultPhoto(payload.ID)
	}

	if len(allMedia) == 0 {
		return defaultPhoto(payload.ID)
	}

	mr := repositories.NewMediaRepository(c.options.MediaFiles)

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

func transformImages(id int32, from []*data.FieldNoteMedia) (to []*station.ImageRef) {
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

func transformConfigurations(from *data.StationFull) (to []*station.StationConfiguration) {
	to = make([]*station.StationConfiguration, 0)
	for _, v := range from.Configurations {
		to = append(to, &station.StationConfiguration{
			ID:      v.ID,
			Time:    v.UpdatedAt.Unix() * 1000,
			Modules: transformModules(from, v.ID),
		})
	}
	return
}

func transformModules(from *data.StationFull, configurationID int64) (to []*station.StationModule) {
	to = make([]*station.StationModule, 0)
	for _, v := range from.Modules {
		if v.ConfigurationID != configurationID {
			continue
		}

		sensors := make([]*data.ModuleSensor, 0)
		for _, s := range from.Sensors {
			if s.ModuleID == v.ID {
				sensors = append(sensors, s)
			}
		}

		sensorsWm := make([]*station.StationSensor, 0)
		translatedName := translateModuleName(v.Name, sensors)
		moduleKey := translateModuleKey(translatedName)

		for _, s := range sensors {
			key := strcase.ToLowerCamel(s.Name)

			sensorsWm = append(sensorsWm, &station.StationSensor{
				Key:           key,
				FullKey:       moduleKey + "." + key,
				Name:          s.Name,
				UnitOfMeasure: s.UnitOfMeasure,
				Reading:       transformReading(s),
			})
		}

		hardwareID := hex.EncodeToString(v.HardwareID)

		to = append(to, &station.StationModule{
			ID:         v.ID,
			FullKey:    moduleKey,
			HardwareID: &hardwareID,
			Name:       translatedName,
			Position:   int32(v.Position),
			Flags:      int32(v.Flags),
			Internal:   v.Flags > 0 || v.Position == 255,
			Sensors:    sensorsWm,
		})

	}
	return
}

var (
	NameMap = map[string]string{
		"distance":    "modules.distance",
		"weather":     "modules.weather",
		"diagnostics": "modules.diagnostics",
		"ultrasonic":  "modules.distance",
	}
)

// This is going to go away eventually once no modules with old names
// are coming in. Then we can do a database migration and get rid of
// them completely. I'm deciding to leave them in so we can see them
// disappear over time.
func translateModuleName(old string, sensors []*data.ModuleSensor) string {
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

func translateModuleKey(name string) string {
	return strings.Replace(name, "modules.", "fk.", 1)
}

func transformLocation(sf *data.StationFull) *station.StationLocation {
	if l := sf.Station.Location; l != nil {
		return &station.StationLocation{
			Latitude:  l.Latitude(),
			Longitude: l.Longitude(),
		}
	}
	return nil
}

func transformStationFull(p Permissions, sf *data.StationFull) (*station.StationFull, error) {
	sp, err := p.ForStation(sf.Station)
	if err != nil {
		return nil, err
	}

	configurations := transformConfigurations(sf)

	return &station.StationFull{
		ID:       sf.Station.ID,
		Name:     sf.Station.Name,
		ReadOnly: sp.IsReadOnly(),
		DeviceID: hex.EncodeToString(sf.Station.DeviceID),
		Uploads:  transformUploads(sf.Ingestions),
		Images:   transformImages(sf.Station.ID, sf.Media),
		Configurations: &station.StationConfigurations{
			All: configurations,
		},
		Battery:         sf.Station.Battery,
		MemoryUsed:      sf.Station.MemoryUsed,
		MemoryAvailable: sf.Station.MemoryAvailable,
		FirmwareNumber:  sf.Station.FirmwareNumber,
		FirmwareTime:    sf.Station.FirmwareTime,
		Updated:         sf.Station.UpdatedAt.Unix() * 1000,
		LocationName:    sf.Station.LocationName,
		PlaceNameOther:  sf.Station.PlaceOther,
		PlaceNameNative: sf.Station.PlaceNative,
		Location:        transformLocation(sf),
		Owner: &station.StationOwner{
			ID:   sf.Owner.ID,
			Name: sf.Owner.Name,
		},
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
