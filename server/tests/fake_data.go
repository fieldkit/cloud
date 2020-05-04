package tests

import (
	"crypto/sha1"
	"fmt"

	"github.com/bxcodec/faker/v3"

	"github.com/fieldkit/cloud/server/data"
)

type FakeStations struct {
	OwnerID   int32
	ProjectID int32
	Stations  []*data.Station
}

func (e *TestEnv) AddStations(number int) (*FakeStations, error) {
	email := faker.Email()
	owner := &data.User{
		Name:     faker.Name(),
		Username: email,
		Email:    email,
		Password: []byte("IGNORE"),
		Bio:      faker.Sentence(),
	}

	if err := e.DB.NamedGetContext(e.Ctx, owner, `
		INSERT INTO fieldkit.user (name, username, email, password, bio)
		VALUES (:name, :email, :email, :password, :bio)
		RETURNING *
		`, owner); err != nil {
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
		INSERT INTO fieldkit.project_user (project_id, user_id) VALUES ($1, $2)
		`, project.ID, owner.ID); err != nil {
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
		OwnerID:   owner.ID,
		ProjectID: project.ID,
		Stations:  stations,
	}, nil
}
