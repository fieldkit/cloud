package api

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/iancoleman/strcase"

	"goa.design/goa/v3/security"

	goa "goa.design/goa/v3/pkg"

	station "github.com/fieldkit/cloud/server/api/gen/station"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/common"
	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/messages"
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

	sr := repositories.NewStationRepository(c.options.Database)

	now := time.Now().UTC()

	if rawStatusPb != nil {
		log.Infow("updating station from status", "station_id", station.ID)

		if err := station.UpdateFromStatus(ctx, *rawStatusPb); err != nil {
			log.Errorw("error updating from status", "error", err)
		}

		if err := sr.UpdateStationModelFromStatus(ctx, station, *rawStatusPb); err != nil {
			log.Errorw("error updating model status", "error", err)
		}

		if station.Location != nil && c.options.locations != nil {
			lastUpdated := now.Sub(station.UpdatedAt)
			if lastUpdated > time.Hour {
				if err := c.options.Publisher.Publish(ctx, &messages.StationLocationUpdated{
					StationID: station.ID,
					Time:      now,
					Location:  station.Location.Coordinates(),
				}); err != nil {
					return err
				}
			}
		}

		station.SyncedAt = &now
	}

	station.UpdatedAt = now

	if err := sr.UpdateStation(ctx, station); err != nil {
		return err
	}

	return nil
}

func (c *StationService) Add(ctx context.Context, payload *station.AddPayload) (response *station.StationFull, err error) {
	tx, err := c.options.Database.Begin(ctx)
	if err != nil {
		return nil, err
	}

	txCtx := context.WithValue(ctx, sqlxcache.TxContextKey, tx)
	response, err = c.add(txCtx, payload)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()

	return response, err
}

func (c *StationService) add(ctx context.Context, payload *station.AddPayload) (response *station.StationFull, err error) {
	log := Logger(ctx).Sugar()

	deviceId, err := hex.DecodeString(payload.DeviceID)
	if err != nil {
		return nil, err
	}

	log.Infow("adding station", "device_id", payload.DeviceID)

	sr := repositories.NewStationRepository(c.options.Database)

	stations, err := sr.QueryStationsByDeviceID(ctx, deviceId)
	if err != nil {
		return nil, err
	}

	if len(stations) > 0 {
		existing := stations[0]

		p, err := NewPermissions(ctx, c.options).ForStation(existing)
		if err != nil {
			return nil, err
		}

		if err := p.CanModify(); err != nil {
			log.Infow("permission:denied", "device_id", deviceId, "user_id", p.UserID(), "owner_id", existing.OwnerID, "station_id", existing.ID)
			return nil, station.MakeStationOwnerConflict(errors.New("permission-denied"))
		}

		if payload.LocationName != nil {
			existing.LocationName = payload.LocationName
		}

		if err := c.updateStation(ctx, existing, payload.StatusPb); err != nil {
			return nil, err
		}

		return c.Get(ctx, &station.GetPayload{
			Auth: &payload.Auth,
			ID:   existing.ID,
		})
	}

	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	adding := &data.Station{
		OwnerID:      p.UserID(),
		Name:         payload.Name,
		DeviceID:     deviceId,
		ModelID:      data.FieldKitModelID,
		CreatedAt:    now,
		UpdatedAt:    now,
		SyncedAt:     &now,
		LocationName: payload.LocationName,
	}

	added, err := sr.AddStation(ctx, adding)
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
		Auth: &payload.Auth,
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

	sr := repositories.NewStationRepository(c.options.Database)

	sf, err := sr.QueryStationFull(ctx, payload.ID)
	if err != nil {
		return nil, err
	}

	pr := repositories.NewProjectRepository(c.options.Database)

	projects, err := pr.QueryProjectsByStationIDForPermissions(ctx, payload.ID)
	if err != nil {
		return nil, err
	}

	preciseLocation := false
	for _, project := range projects {
		if project.Privacy == data.Public {
			preciseLocation = true
		}
	}

	if !p.Anonymous() && !preciseLocation {
		if p.UserID() == sf.Owner.ID {
			preciseLocation = true
		} else {
			relationships, err := pr.QueryUserProjectRelationships(ctx, p.UserID())
			if err != nil {
				return nil, err
			}
			if row, ok := relationships[payload.ID]; ok {
				preciseLocation = row.MemberRole >= 0
			}
		}
	}

	mmr := repositories.NewModuleMetaRepository(c.options.Database)
	mm, err := mmr.FindAllModulesMeta(ctx)
	if err != nil {
		return nil, err
	}

	return transformStationFull(c.options.signer, p, sf, preciseLocation, false, mm)
}

func (c *StationService) Transfer(ctx context.Context, payload *station.TransferPayload) (err error) {
	sr := repositories.NewStationRepository(c.options.Database)

	updating, err := sr.QueryStationByID(ctx, payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return station.MakeNotFound(errors.New("station not found"))
		}
		return err
	}

	p, err := NewPermissions(ctx, c.options).ForStation(updating)
	if err != nil {
		return err
	}

	if err := p.RequireAdmin(); err != nil {
		return err
	}

	updating.UpdatedAt = time.Now()
	updating.OwnerID = payload.OwnerID

	if err := sr.UpdateOwner(ctx, updating); err != nil {
		return err
	}

	return nil
}

func (c *StationService) DefaultPhoto(ctx context.Context, payload *station.DefaultPhotoPayload) (err error) {
	sr := repositories.NewStationRepository(c.options.Database)

	updating, err := sr.QueryStationByID(ctx, payload.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return station.MakeNotFound(errors.New("station not found"))
		}
		return err
	}

	p, err := NewPermissions(ctx, c.options).ForStation(updating)
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	updating.UpdatedAt = time.Now()
	updating.PhotoID = &payload.PhotoID

	if err := sr.UpdatePhoto(ctx, updating); err != nil {
		return err
	}

	return nil
}

func (c *StationService) Update(ctx context.Context, payload *station.UpdatePayload) (response *station.StationFull, err error) {
	sr := repositories.NewStationRepository(c.options.Database)

	updating, err := sr.QueryStationByID(ctx, payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, station.MakeNotFound(errors.New("station not found"))
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
	if payload.LocationName != nil {
		updating.LocationName = payload.LocationName
	}

	if err := c.updateStation(ctx, updating, payload.StatusPb); err != nil {
		return nil, err
	}

	return c.Get(ctx, &station.GetPayload{
		Auth: &payload.Auth,
		ID:   payload.ID,
	})
}

func (c *StationService) ListMine(ctx context.Context, payload *station.ListMinePayload) (response *station.StationsFull, err error) {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	sr := repositories.NewStationRepository(c.options.Database)

	sfs, err := sr.QueryStationFullByOwnerID(ctx, p.UserID())
	if err != nil {
		return nil, err
	}

	mmr := repositories.NewModuleMetaRepository(c.options.Database)
	mm, err := mmr.FindAllModulesMeta(ctx)
	if err != nil {
		return nil, err
	}

	stations, err := transformAllStationFull(c.options.signer, p, sfs, true, false, mm, true)
	if err != nil {
		return nil, err
	}

	response = &station.StationsFull{
		Stations: stations,
	}

	return
}

func (c *StationService) ListProject(ctx context.Context, payload *station.ListProjectPayload) (response *station.StationsFull, err error) {
	pr := repositories.NewProjectRepository(c.options.Database)
	sr := repositories.NewStationRepository(c.options.Database)
	mmr := repositories.NewModuleMetaRepository(c.options.Database)

	p, err := NewPermissions(ctx, c.options).ForProjectByID(payload.ID)
	if err != nil {
		return nil, err
	}

	if err := p.CanView(); err != nil {
		return nil, err
	}

	project, err := pr.QueryByID(ctx, payload.ID)
	if err != nil {
		return nil, err
	}

	preciseLocation := project.Privacy == data.Public
	if !p.Anonymous() && !preciseLocation {
		// NOTE This may already be in Permissions.
		// NOTE Compare with GET
		relationships, err := pr.QueryUserProjectRelationships(ctx, p.UserID())
		if err != nil {
			return nil, err
		}
		if row, ok := relationships[payload.ID]; ok {
			preciseLocation = row.MemberRole >= 0
		}
	}

	sfs, err := sr.QueryStationFullByProjectID(ctx, payload.ID)
	if err != nil {
		return nil, err
	}

	mm, err := mmr.FindAllModulesMeta(ctx)
	if err != nil {
		return nil, err
	}

	disableFiltering := false
	if payload.DisableFiltering != nil {
		disableFiltering = *payload.DisableFiltering
	}

	stations, err := transformAllStationFull(c.options.signer, p, sfs, preciseLocation, false, mm, !disableFiltering)
	if err != nil {
		return nil, err
	}

	response = &station.StationsFull{
		Stations: stations,
	}

	return
}

func isForbidden(err error) bool {
	if se, ok := err.(*goa.ServiceError); ok {
		return se.Name == "forbidden"
	}
	return false
}

type associatedStationSorter struct {
	stations         []*station.AssociatedStation
	queriedStationID int32
}

func (a associatedStationSorter) Len() int { return len(a.stations) }
func (a associatedStationSorter) Less(i, j int) bool {
	return a.stations[i].Station.ID == a.queriedStationID
}
func (a associatedStationSorter) Swap(i, j int) {
	a.stations[i], a.stations[j] = a.stations[j], a.stations[i]
}

func (c *StationService) ListProjectAssociated(ctx context.Context, payload *station.ListProjectAssociatedPayload) (response *station.AssociatedStations, err error) {
	byID := make(map[int32]*station.AssociatedStation)

	stationIDs := make([]int32, 0)
	disableFiltering := true
	projectStations, err := c.ListProject(ctx, &station.ListProjectPayload{
		ID:               payload.ProjectID,
		DisableFiltering: &disableFiltering,
	})
	if err != nil {
		if isForbidden(err) {
			return nil, err
		}
		return nil, err
	} else {
		for _, fullStation := range projectStations.Stations {
			stationIDs = append(stationIDs, fullStation.ID)

			associated := &station.AssociatedStation{
				Station: fullStation,
				Hidden:  fullStation.Model.OnlyVisibleViaAssociation,
				Project: []*station.AssociatedViaProject{&station.AssociatedViaProject{
					ID: payload.ProjectID,
				},
				},
			}

			byID[fullStation.ID] = associated
		}
	}

	sr := repositories.NewStationRepository(c.options.Database)

	associations, err := sr.QueryAssociatedStations(ctx, stationIDs)
	if err != nil {
		return nil, err
	}

	associated := make([]*station.AssociatedStation, 0)

	for _, associatedStation := range byID {
		associated = append(associated, associatedStation)

		for _, association := range associations[associatedStation.Station.ID] {
			if associatedStation.Manual == nil {
				associatedStation.Manual = make([]*station.AssociatedViaManual, 0)
			}

			associatedStation.Manual = append(associatedStation.Manual, &station.AssociatedViaManual{
				OtherStationID: association.AssociatedStationID,
				Priority:       association.Priority,
			})

			associatedStation.Hidden = false
		}

		if associatedStation.Station.Location != nil && associatedStation.Station.Location.Precise != nil {
			for _, projectAssociation := range associatedStation.Project {
				location := data.NewLocation(associatedStation.Station.Location.Precise)
				if nearby, err := sr.QueryNearbyProjectStations(ctx, projectAssociation.ID, location); err != nil {
					return nil, err
				} else {
					for _, ns := range nearby {
						if associatedStation.Location == nil {
							associatedStation.Location = make([]*station.AssociatedViaLocation, 0)
						}

						associatedStation.Location = append(associatedStation.Location, &station.AssociatedViaLocation{
							StationID: ns.StationID,
							Distance:  ns.Distance,
						})
					}
				}
			}
		}
	}

	sort.Sort(associatedStationSorter{
		stations:         associated,
		queriedStationID: 0, // TODO Remove this when we remove ListAssociated (maybe)
	})

	response = &station.AssociatedStations{
		Stations: associated,
	}

	return
}

func (c *StationService) ListAssociated(ctx context.Context, payload *station.ListAssociatedPayload) (response *station.AssociatedStations, err error) {
	log := Logger(ctx).Sugar()

	p, err := NewPermissions(ctx, c.options).ForStationByID(int(payload.ID))
	if err != nil {
		return nil, err
	}

	if err := p.CanView(); err != nil {
		return nil, err
	}

	pr := repositories.NewProjectRepository(c.options.Database)

	projects, err := pr.QueryProjectsByStationIDForPermissions(ctx, payload.ID)
	if err != nil {
		return nil, err
	}

	if len(projects) > 1 {
		log.Warnw("associated:ambiguous-projects", "station_id", payload.ID, "projects", len(projects))
	}

	for _, project := range projects {
		projectPermissions, err := NewPermissions(ctx, c.options).ForProject(project)
		if err != nil {
			return nil, err
		}

		if err := projectPermissions.CanView(); err == nil {
			return c.ListProjectAssociated(ctx, &station.ListProjectAssociatedPayload{
				ProjectID: project.ID,
			})
		}
	}

	get, err := c.Get(ctx, &station.GetPayload{
		ID: payload.ID,
	})
	if err != nil {
		return nil, err
	}

	return &station.AssociatedStations{
		Stations: []*station.AssociatedStation{
			&station.AssociatedStation{
				Station: get,
			},
		},
	}, nil
}

func (c *StationService) queriedToPage(queried *repositories.QueriedEssential) (*station.PageOfStations, error) {
	stationsWm := make([]*station.EssentialStation, 0)

	for _, es := range queried.Stations {
		var lastIngestionAt *int64

		if es.LastIngestionAt != nil {
			jsEpoch := es.LastIngestionAt.Unix() * 1000
			lastIngestionAt = &jsEpoch
		}

		var location *station.StationLocation = nil
		if es.Location != nil {
			location = &station.StationLocation{
				Precise: []float64{es.Location.Longitude(), es.Location.Latitude()},
			}
		}

		var recordingStartedAt *int64
		if es.RecordingStartedAt != nil {
			unix := es.RecordingStartedAt.Unix() * 1000
			recordingStartedAt = &unix
		}

		stationsWm = append(stationsWm, &station.EssentialStation{
			ID:       es.ID,
			Name:     es.Name,
			DeviceID: hex.EncodeToString(es.DeviceID),
			Owner: &station.StationOwner{
				ID:    es.OwnerID,
				Name:  es.OwnerName,
				Email: &es.OwnerEmail,
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

	return &station.PageOfStations{
		Stations: stationsWm,
		Total:    queried.Total,
	}, nil
}

func (c *StationService) ListAll(ctx context.Context, payload *station.ListAllPayload) (*station.PageOfStations, error) {
	sr := repositories.NewStationRepository(c.options.Database)

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

	return c.queriedToPage(queried)
}

func (c *StationService) Delete(ctx context.Context, payload *station.DeletePayload) error {
	sr := repositories.NewStationRepository(c.options.Database)

	return sr.Delete(ctx, payload.StationID)
}

func findStaticImageOrGif(media []*data.FieldNoteMedia) *data.FieldNoteMedia {
	if len(media) == 0 {
		return nil
	}
	for _, m := range media {
		if m.ContentType != "image/gif" {
			return m
		}
	}
	return media[0]
}

func (c *StationService) DownloadPhoto(ctx context.Context, payload *station.DownloadPhotoPayload) (*station.DownloadedPhoto, error) {
	x := uint(124)
	y := uint(100)

	allMedia := []*data.FieldNoteMedia{}
	if err := c.options.Database.SelectContext(ctx, &allMedia, `
        SELECT id, user_id, content_type, created_at, url, key, station_id
		FROM fieldkit.notes_media WHERE id IN (SELECT photo_id FROM fieldkit.station WHERE id = $1) ORDER BY created_at DESC
        `, payload.StationID); err != nil {
		return nil, err
	}

	if len(allMedia) == 0 {
		return defaultPhoto(payload.StationID)
	}

	media := findStaticImageOrGif(allMedia)

	var resize *PhotoResizeSettings
	if payload.Size != nil {
		resize = &PhotoResizeSettings{
			Size: *payload.Size,
		}
	}

	photo, err := c.options.photoCache.Load(ctx, &ExternalMedia{
		URL:         media.URL,
		ContentType: media.ContentType,
	},
		resize,
		&PhotoCropSettings{
			X: x,
			Y: y,
		},
		payload.IfNoneMatch,
	)
	if err != nil {
		return nil, err
	}

	if len(photo.Bytes) == 0 && !photo.EtagMatch {
		log := Logger(ctx).Sugar()
		log.Warnw("empty-image", "station_id", payload.StationID)
		return defaultPhoto(payload.StationID)
	}

	return &station.DownloadedPhoto{
		Length:      int64(photo.Size),
		ContentType: photo.ContentType,
		Etag:        photo.Etag,
		Body:        photo.Bytes,
	}, nil
}

func (c *StationService) AdminSearch(ctx context.Context, payload *station.AdminSearchPayload) (*station.PageOfStations, error) {
	sr := repositories.NewStationRepository(c.options.Database)

	queried, err := sr.Search(ctx, payload.Query)
	if err != nil {
		return nil, err
	}

	return c.queriedToPage(queried)
}

func (c *StationService) Progress(ctx context.Context, payload *station.ProgressPayload) (response *station.StationProgress, err error) {
	p, err := NewPermissions(ctx, c.options).ForStationByID(int(payload.StationID))
	if err != nil {
		return nil, err
	}

	if err := p.CanView(); err != nil {
		return nil, err
	}

	jobsWm := make([]*station.StationJob, 0)

	return &station.StationProgress{
		Jobs: jobsWm,
	}, nil
}

func (s *StationService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, common.AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     func(m string) error { return station.MakeNotFound(errors.New(m)) },
		Unauthorized: func(m string) error { return station.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return station.MakeForbidden(errors.New(m)) },
	})
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

	if math.IsNaN(*s.ReadingValue) {
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

func transformConfigurations(from *data.StationFull, transformAll bool, moduleMeta *repositories.AllModuleMeta) (to []*station.StationConfiguration, err error) {
	to = make([]*station.StationConfiguration, 0)
	for _, v := range from.Configurations {
		modules, err := transformModules(from, v.ID, moduleMeta)
		if err != nil {
			return nil, err
		}

		to = append(to, &station.StationConfiguration{
			ID:      v.ID,
			Time:    v.UpdatedAt.Unix() * 1000,
			Modules: modules,
		})

		if !transformAll {
			break
		}
	}
	return
}

func convertToTypedMap(obj interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func transformModules(from *data.StationFull, configurationID int64, moduleMeta *repositories.AllModuleMeta) (to []*station.StationModule, err error) {
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

		var serializedModuleMeta map[string]interface{}

		sensorsWm := make([]*station.StationSensor, 0)
		translatedName := translateModuleName(v.Name, sensors)
		moduleKey := translateModuleKey(translatedName)

		for _, s := range sensors {
			key := strcase.ToLowerCamel(s.Name)
			fullKey := moduleKey + "." + key

			var serializedSensorMeta map[string]interface{}

			sensorMeta := moduleMeta.FindSensorByFullKey(fullKey)
			if sensorMeta != nil {
				converted, err := convertToTypedMap(sensorMeta.Sensor)
				if err != nil {
					return nil, err
				}
				serializedSensorMeta = converted

				if serializedModuleMeta == nil {
					converted, err = convertToTypedMap(sensorMeta.Module)
					if err != nil {
						return nil, err
					}
					serializedModuleMeta = converted
				}
			}

			sensorsWm = append(sensorsWm, &station.StationSensor{
				Key:           key,
				FullKey:       fullKey,
				Name:          s.Name,
				UnitOfMeasure: s.UnitOfMeasure,
				Reading:       transformReading(s),
				Meta:          serializedSensorMeta,
			})
		}

		hardwareID := hex.EncodeToString(v.HardwareID)
		hardwareIDBase64 := base64.StdEncoding.EncodeToString(v.HardwareID)

		to = append(to, &station.StationModule{
			ID:               v.ID,
			FullKey:          moduleKey,
			HardwareID:       &hardwareID,
			HardwareIDBase64: &hardwareIDBase64,
			Name:             translatedName,
			Position:         int32(v.Position),
			Flags:            int32(v.Flags),
			Internal:         v.Flags > 0 || v.Position == 255,
			Meta:             serializedModuleMeta,
			Sensors:          sensorsWm,
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

func transformLocation(sf *data.StationFull, preciseLocation bool) *station.StationLocation {
	if l := sf.Station.Location; l != nil {
		sl := &station.StationLocation{}

		if preciseLocation {
			sl.Precise = []float64{l.Longitude(), l.Latitude()}
		}

		regions := make([]*station.StationRegion, 0)
		for _, area := range sf.Areas {
			regions = append(regions, &station.StationRegion{
				Name:  area.Name,
				Shape: area.Geometry.Coordinates(),
			})
		}
		if false {
			sl.Regions = regions
		}

		return sl
	}
	return nil
}

func transformStationFull(signer *Signer, p Permissions, sf *data.StationFull, preciseLocation bool, transformAllConfigurations bool, moduleMeta *repositories.AllModuleMeta) (*station.StationFull, error) {
	readOnly := true
	if p != nil {
		sp, err := p.ForStation(sf.Station)
		if err != nil {
			return nil, err
		}

		readOnly, err = sp.IsReadOnly()
		if err != nil {
			return nil, err
		}
	}

	configurations, err := transformConfigurations(sf, transformAllConfigurations, moduleMeta)
	if err != nil {
		return nil, err
	}

	var lastReadingAt *int64

	for _, configuration := range configurations {
		for _, module := range configuration.Modules {
			for _, sensor := range module.Sensors {
				if sensor.Reading != nil && sensor.Reading.Time > 0 {
					if lastReadingAt == nil {
						lastReadingAt = &sensor.Reading.Time
					} else {
						if *lastReadingAt < sensor.Reading.Time {
							lastReadingAt = &sensor.Reading.Time
						}
					}
				}
			}
		}
	}

	dataSummary := transformDataSummary(sf.DataSummary)

	location := transformLocation(sf, preciseLocation)

	var recordingStartedAt *int64
	if sf.Station.RecordingStartedAt != nil {
		unix := sf.Station.RecordingStartedAt.Unix() * 1000
		recordingStartedAt = &unix
	}

	var photos *station.StationPhotos
	if sf.Station.PhotoID != nil {
		photos = &station.StationPhotos{
			Small: fmt.Sprintf("/stations/%d/photo", sf.Station.ID),
		}
	}

	windows := make([]*station.StationInterestingnessWindow, 0)
	for _, iness := range sf.Interestingness {
		windows = append(windows, &station.StationInterestingnessWindow{
			Seconds:         iness.WindowSeconds,
			Interestingness: float64(iness.Interestingness),
			Value:           iness.ReadingValue,
			Time:            iness.ReadingTime.Unix() * 1000,
		})
	}

	attributes := make([]*station.StationProjectAttribute, 0)
	for _, attribute := range sf.Attributes {
		attributes = append(attributes, &station.StationProjectAttribute{
			ProjectID:   attribute.ProjectID,
			AttributeID: attribute.AttributeID,
			Name:        attribute.Name,
			StringValue: attribute.StringValue,
			Priority:    attribute.Priority,
		})
	}

	stationFull := &station.StationFull{
		ID:       sf.Station.ID,
		Name:     sf.Station.Name,
		ReadOnly: readOnly,
		DeviceID: hex.EncodeToString(sf.Station.DeviceID),
		Interestingness: &station.StationInterestingness{
			Windows: windows,
		},
		Attributes: &station.StationProjectAttributes{
			Attributes: attributes,
		},
		Uploads: transformUploads(sf.Ingestions),
		Configurations: &station.StationConfigurations{
			All: configurations,
		},
		RecordingStartedAt: recordingStartedAt,
		Battery:            sf.Station.Battery,
		MemoryUsed:         sf.Station.MemoryUsed,
		MemoryAvailable:    sf.Station.MemoryAvailable,
		FirmwareNumber:     sf.Station.FirmwareNumber,
		FirmwareTime:       sf.Station.FirmwareTime,
		UpdatedAt:          sf.Station.UpdatedAt.Unix() * 1000,
		LastReadingAt:      lastReadingAt,
		SyncedAt:           optionalTime(sf.Station.SyncedAt),
		IngestionAt:        optionalTime(sf.Station.IngestionAt),
		LocationName:       sf.Station.LocationName,
		PlaceNameOther:     sf.Station.PlaceOther,
		PlaceNameNative:    sf.Station.PlaceNative,
		Location:           location,
		Data:               dataSummary,
		Hidden:             sf.Station.Hidden,
		Status:             sf.Station.Status,
		Model: &station.StationFullModel{
			Name:                      sf.Model.Name,
			OnlyVisibleViaAssociation: sf.Model.OnlyVisibleViaAssociation,
		},
		Owner: &station.StationOwner{
			ID:   sf.Owner.ID,
			Name: sf.Owner.Name,
		},
		Photos: photos,
	};

    if sf.Owner.MediaURL != nil {
        url := fmt.Sprintf("/user/%d/media", sf.Owner.ID)
        stationFull.Owner.Photo = &station.UserPhoto{
            URL: &url,
        }
    }

	return stationFull, nil
}

func optionalTime(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	value := t.Unix() * 1000
	return &value
}

func transformDataSummary(ads *data.AggregatedDataSummary) *station.StationDataSummary {
	if ads == nil {
		return nil
	}
	if ads.Start == nil || ads.End == nil || ads.NumberSamples == nil {
		return nil
	}
	return &station.StationDataSummary{
		Start:           (*ads.Start).Unix() * 1000,
		End:             (*ads.End).Unix() * 1000,
		NumberOfSamples: *ads.NumberSamples,
	}
}

func transformAllStationFull(signer *Signer, p Permissions, sfs []*data.StationFull, preciseLocation bool, transformAllConfigurations bool, moduleMeta *repositories.AllModuleMeta, filtering bool) ([]*station.StationFull, error) {
	stations := make([]*station.StationFull, 0)

	for _, sf := range sfs {
		after, err := transformStationFull(signer, p, sf, preciseLocation, transformAllConfigurations, moduleMeta)
		if err != nil {
			return nil, err
		}

		if !filtering || !after.Model.OnlyVisibleViaAssociation {
			stations = append(stations, after)
		}
	}

	return stations, nil
}

func defaultPhoto(id int32) (*station.DownloadedPhoto, error) {
	defaultPhotoContentType := "image/png"
	defaultPhoto, err := StationDefaultPicture(int64(id))
	if err != nil {
		// NOTE This, hopefully never happens because we've got no image to send back.
		return nil, err
	}

	return &station.DownloadedPhoto{
		ContentType: defaultPhotoContentType,
		Length:      int64(len(defaultPhoto)),
		Body:        defaultPhoto,
	}, nil
}
