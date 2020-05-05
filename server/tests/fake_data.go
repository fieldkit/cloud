package tests

import (
	"crypto/sha1"
	"fmt"

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
