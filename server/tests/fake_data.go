package tests

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/bxcodec/faker/v3"

	"github.com/fieldkit/cloud/server/data"
)

type FakeStations struct {
	Owner    *data.User
	Project  *data.Project
	Stations []*data.Station
}

func (e *TestEnv) AddUser(pw string) (*data.User, error) {
	email := faker.Email()
	user := &data.User{
		Name:     faker.Name(),
		Username: email,
		Email:    email,
		Bio:      faker.Sentence(),
	}

	user.SetPassword(pw)

	if err := e.DB.NamedGetContext(e.Ctx, user, `
		INSERT INTO fieldkit.user (name, username, email, password, bio)
		VALUES (:name, :email, :email, :password, :bio)
		RETURNING *
		`, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (e *TestEnv) AddStations(number int) (*FakeStations, error) {
	owner, err := e.AddUser("passwordpassword")
	if err != nil {
		return nil, err
	}

	name := faker.Name()

	project := &data.Project{
		Name: name + " Project",
		Slug: name,
	}

	if err := e.DB.NamedGetContext(e.Ctx, project, `
		INSERT INTO fieldkit.project (name, slug)
		VALUES (:name, :slug)
		RETURNING *
		`, project); err != nil {
		return nil, err
	}

	if _, err := e.DB.ExecContext(e.Ctx, `
		INSERT INTO fieldkit.project_user (project_id, user_id, role) VALUES ($1, $2, $3)
		`, project.ID, owner.ID, data.AdministratorRole.ID); err != nil {
		return nil, err
	}

	stations := []*data.Station{}

	for i := 0; i < number; i += 1 {
		name := fmt.Sprintf("%s #%d", owner.Name, i)

		hasher := sha1.New()
		hasher.Write([]byte(name))
		deviceID := hasher.Sum(nil)

		station := &data.Station{
			OwnerID:  owner.ID,
			Name:     name,
			DeviceID: deviceID,
		}

		if err := e.DB.NamedGetContext(e.Ctx, station, `
			INSERT INTO fieldkit.station (name, device_id, owner_id, status_json)
			VALUES (:name, :device_id, :owner_id, :status_json)
			RETURNING *
		`, station); err != nil {
			return nil, err
		}

		if _, err := e.DB.ExecContext(e.Ctx, `
			INSERT INTO fieldkit.project_station (project_id, station_id) VALUES ($1, $2)
			`, project.ID, station.ID); err != nil {
			return nil, err
		}

		stations = append(stations, station)
	}

	return &FakeStations{
		Owner:    owner,
		Project:  project,
		Stations: stations,
	}, nil
}

func (e *TestEnv) NewStation(owner *data.User) *data.Station {
	name := faker.Name()

	hasher := sha1.New()
	hasher.Write([]byte(name))
	deviceID := hasher.Sum(nil)

	station := &data.Station{
		OwnerID:  owner.ID,
		DeviceID: deviceID,
		Name:     name,
	}

	return station
}

func (e *TestEnv) AddStationActivity(station *data.Station, user *data.User) error {
	location := data.NewLocation([]float64{0, 0})

	depoyedActivity := &data.StationDeployed{
		StationActivity: data.StationActivity{
			CreatedAt: time.Now(),
			StationID: int64(station.ID),
		},
		DeployedAt: time.Now(),
		Location:   location,
	}

	if _, err := e.DB.NamedExecContext(e.Ctx, `
		INSERT INTO fieldkit.station_deployed (created_at, station_id, deployed_at, location) VALUES (:created_at, :station_id, :deployed_at, ST_SetSRID(ST_GeomFromText(:location), 4326))
		`, depoyedActivity); err != nil {
		return err
	}

	ingestion := &data.Ingestion{
		URL:          "file:///dev/nul",
		UploadID:     "",
		UserID:       user.ID,
		DeviceID:     station.DeviceID,
		GenerationID: []byte{0x00, 0x01},
		Type:         "data",
		Size:         int64(1024),
		Blocks:       data.Int64Range([]int64{1, 100}),
		Flags:        pq.Int64Array([]int64{}),
	}

	if err := e.DB.NamedGetContext(e.Ctx, ingestion, `
			INSERT INTO fieldkit.ingestion (time, upload_id, user_id, device_id, generation, type, size, url, blocks, flags)
			VALUES (NOW(), :upload_id, :user_id, :device_id, :generation, :type, :size, :url, :blocks, :flags)
			RETURNING *
			`, ingestion); err != nil {
		return err
	}

	activity := &data.StationIngestion{
		StationActivity: data.StationActivity{
			CreatedAt: time.Now(),
			StationID: int64(station.ID),
		},
		UploaderID:      int64(user.ID),
		DataIngestionID: ingestion.ID,
		DataRecords:     1,
		Errors:          false,
	}

	if err := e.DB.NamedGetContext(e.Ctx, activity, `
		INSERT INTO fieldkit.station_ingestion (created_at, station_id, uploader_id, data_ingestion_id, data_records, errors)
		VALUES (:created_at, :station_id, :uploader_id, :data_ingestion_id, :data_records, :errors)
		ON CONFLICT (data_ingestion_id) DO NOTHING
		RETURNING *
		`, activity); err != nil {
		return err
	}

	return nil
}

func (e *TestEnv) AddProjectActivity(project *data.Project, station *data.Station, user *data.User) error {
	if err := e.AddStationActivity(station, user); err != nil {
		return err
	}

	projectUpdate := &data.ProjectUpdate{
		ProjectActivity: data.ProjectActivity{
			CreatedAt: time.Now(),
			ProjectID: int64(project.ID),
		},
		AuthorID: int64(user.ID),
		Body:     "Project update",
	}

	if _, err := e.DB.NamedExecContext(e.Ctx, `
		INSERT INTO fieldkit.project_update (created_at, project_id, author_id, body) VALUES (:created_at, :project_id, :author_id, :body)
		`, projectUpdate); err != nil {
		return err
	}

	return nil
}
