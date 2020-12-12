package tests

import (
	crand "crypto/rand"
	"crypto/sha1"
	"fmt"
	mrand "math/rand"
	"time"

	"github.com/lib/pq"

	"github.com/bxcodec/faker/v3"
	"golang.org/x/crypto/blake2b"

	"github.com/golang/protobuf/proto"

	pbapp "github.com/fieldkit/app-protocol"
	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

const (
	GoodPassword = "goodgoodgood"
	BadPassword  = "badbadbadbad"
)

type FakeStations struct {
	Owner    *data.User
	Project  *data.Project
	Stations []*data.Station
}

func (e *TestEnv) AddInvalidUser() (*data.User, error) {
	email := faker.Email()
	user := &data.User{
		Name:     faker.Name(),
		Username: email,
		Email:    email,
		Bio:      faker.Sentence(),
		Valid:    false,
	}

	user.SetPassword(GoodPassword)

	if err := e.DB.NamedGetContext(e.Ctx, user, `
		INSERT INTO fieldkit.user (name, username, email, password, bio)
		VALUES (:name, :email, :email, :password, :bio) RETURNING *
		`, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (e *TestEnv) AddUser() (*data.User, error) {
	email := faker.Email()
	user := &data.User{
		Name:     faker.Name(),
		Username: email,
		Email:    email,
		Bio:      faker.Sentence(),
		Valid:    true,
	}

	user.SetPassword(GoodPassword)

	if err := e.DB.NamedGetContext(e.Ctx, user, `
		INSERT INTO fieldkit.user (name, username, email, password, bio, valid)
		VALUES (:name, :email, :email, :password, :bio, :valid) RETURNING *
		`, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (e *TestEnv) AddAdminUser() (*data.User, error) {
	email := faker.Email()
	user := &data.User{
		Name:     faker.Name(),
		Username: email,
		Email:    email,
		Bio:      faker.Sentence(),
		Admin:    true,
	}

	user.SetPassword(GoodPassword)

	if err := e.DB.NamedGetContext(e.Ctx, user, `
		INSERT INTO fieldkit.user (name, username, email, password, bio, admin)
		VALUES (:name, :email, :email, :password, :bio, :admin) RETURNING *
		`, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (e *TestEnv) AddProject() (*data.Project, error) {
	name := faker.Name()

	project := &data.Project{
		Name:    name + " Project",
		Privacy: data.Private,
	}

	if err := e.DB.NamedGetContext(e.Ctx, project, `
		INSERT INTO fieldkit.project (name, privacy)
		VALUES (:name, :privacy)
		RETURNING *
		`, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (e *TestEnv) AddProjectUser(p *data.Project, u *data.User, r *data.Role) error {
	if _, err := e.DB.ExecContext(e.Ctx, `
		INSERT INTO fieldkit.project_user (project_id, user_id, role) VALUES ($1, $2, $3)
		`, p.ID, u.ID, r.ID); err != nil {
		return err
	}

	return nil
}

func (e *TestEnv) AddStations(number int) (*FakeStations, error) {
	owner, err := e.AddUser()
	if err != nil {
		return nil, err
	}

	project, err := e.AddProject()
	if err != nil {
		return nil, err
	}

	if err := e.AddProjectUser(project, owner, data.AdministratorRole); err != nil {
		return nil, err
	}

	stations := []*data.Station{}

	for i := 0; i < number; i += 1 {
		name := fmt.Sprintf("%s #%d", owner.Name, i)

		hasher := sha1.New()
		hasher.Write([]byte(name))
		deviceID := hasher.Sum(nil)

		station := &data.Station{
			OwnerID:   owner.ID,
			Name:      name,
			DeviceID:  deviceID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Location:  data.NewLocation([]float64{36.9053064, -117.5297165}),
		}

		if err := e.DB.NamedGetContext(e.Ctx, station, `
			INSERT INTO fieldkit.station (name, device_id, owner_id, created_at, updated_at, location)
			VALUES (:name, :device_id, :owner_id, :created_at, :updated_at, ST_SetSRID(ST_GeomFromText(:location), 4326))
			RETURNING id
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
		OwnerID:   owner.ID,
		DeviceID:  deviceID,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return station
}

func (e *TestEnv) AddProvision(deviceID, generationID []byte) (*data.Provision, error) {
	provision := &data.Provision{
		Created:      time.Now(),
		Updated:      time.Now(),
		DeviceID:     deviceID,
		GenerationID: generationID,
	}

	if err := e.DB.NamedGetContext(e.Ctx, provision, `
			INSERT INTO fieldkit.provision (device_id, generation, created, updated)
			VALUES (:device_id, :generation, :created, :updated) ON CONFLICT (device_id, generation)
			DO UPDATE SET updated = NOW() RETURNING id
			`, provision); err != nil {
		return nil, err
	}

	return provision, nil
}

func (e *TestEnv) AddIngestion(user *data.User, url, typeName string, deviceID []byte, length int) (*data.QueuedIngestion, *data.Ingestion, error) {
	ingestion := &data.Ingestion{
		URL:          url,
		UserID:       user.ID,
		DeviceID:     deviceID,
		GenerationID: deviceID,
		Type:         typeName,
		Size:         int64(length),
		Blocks:       data.Int64Range([]int64{1, 100}),
		Flags:        pq.Int64Array([]int64{}),
	}

	if err := e.DB.NamedGetContext(e.Ctx, ingestion, `
			INSERT INTO fieldkit.ingestion (time, upload_id, user_id, device_id, generation, type, size, url, blocks, flags)
			VALUES (NOW(), :upload_id, :user_id, :device_id, :generation, :type, :size, :url, :blocks, :flags)
			RETURNING id
			`, ingestion); err != nil {
		return nil, nil, err
	}

	queued := &data.QueuedIngestion{
		Queued:      time.Now(),
		IngestionID: ingestion.ID,
	}

	if err := e.DB.NamedGetContext(e.Ctx, queued, `
			INSERT INTO fieldkit.ingestion_queue (queued, ingestion_id)
			VALUES (:queued, :ingestion_id)
			RETURNING id
			`, queued); err != nil {
		return nil, nil, err
	}

	return queued, ingestion, nil
}

func (e *TestEnv) AddStationActivity(station *data.Station, user *data.User) error {
	location := data.NewLocation([]float64{0, 0})

	depoyedActivity := &data.StationDeployed{
		StationActivity: data.StationActivity{
			CreatedAt: time.Now(),
			StationID: station.ID,
		},
		DeployedAt: time.Now(),
		Location:   location,
	}

	if _, err := e.DB.NamedExecContext(e.Ctx, `
		INSERT INTO fieldkit.station_deployed (created_at, station_id, deployed_at, location) VALUES (:created_at, :station_id, :deployed_at, ST_SetSRID(ST_GeomFromText(:location), 4326))
		`, depoyedActivity); err != nil {
		return err
	}

	_, ingestion, err := e.AddIngestion(user, "file:///dev/null", data.DataTypeName, station.DeviceID, 0)
	if err != nil {
		return err
	}

	activity := &data.StationIngestion{
		StationActivity: data.StationActivity{
			CreatedAt: time.Now(),
			StationID: station.ID,
		},
		UploaderID:      user.ID,
		DataIngestionID: ingestion.ID,
		DataRecords:     1,
		Errors:          false,
	}

	if err := e.DB.NamedGetContext(e.Ctx, activity, `
		INSERT INTO fieldkit.station_ingestion (created_at, station_id, uploader_id, data_ingestion_id, data_records, errors)
		VALUES (:created_at, :station_id, :uploader_id, :data_ingestion_id, :data_records, :errors)
		ON CONFLICT (data_ingestion_id) DO NOTHING
		RETURNING id
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
			ProjectID: project.ID,
		},
		AuthorID: user.ID,
		Body:     "Project update",
	}

	if _, err := e.DB.NamedExecContext(e.Ctx, `
		INSERT INTO fieldkit.project_update (created_at, project_id, author_id, body) VALUES (:created_at, :project_id, :author_id, :body)
		`, projectUpdate); err != nil {
		return err
	}

	return nil
}

func (e *TestEnv) NewRandomData(n int) ([]byte, error) {
	data := make([]byte, n)
	actual, err := crand.Read(data)
	if actual != n {
		return nil, fmt.Errorf("unexpected random byte read")
	}
	return data, err
}

type SignedRecordAndData struct {
	Signed *pb.SignedRecord
	Data   *pb.DataRecord
	Bytes  []byte
}

func hashString(seed string) []byte {
	hasher := sha1.New()
	hasher.Write([]byte(seed))
	return hasher.Sum(nil)
}

func (e *TestEnv) NewHttpStatusReply(s *data.Station) *pbapp.HttpReply {
	now := time.Now()
	used := uint32(1024 * 10)
	installed := uint32(512 * 1024 * 1024)
	recording := uint64(0)
	deviceID := s.DeviceID
	generation := hashString(fmt.Sprintf("%s-%d-%d", s.DeviceID, e.Seed, 0))

	return &pbapp.HttpReply{
		Type: pbapp.ReplyType_REPLY_STATUS,
		Status: &pbapp.Status{
			Version: 1,
			Uptime:  1,
			Time:    uint64(now.Unix()),
			Identity: &pbapp.Identity{
				Name:       s.Name,
				DeviceId:   deviceID,
				Generation: generation,
			},
			Recording: &pbapp.Recording{
				Enabled:     recording > 0,
				StartedTime: recording,
			},
			Memory: &pbapp.MemoryStatus{
				SramAvailable:           128 * 1024,
				ProgramFlashAvailable:   600 * 1024,
				ExtendedMemoryAvailable: 0,
				DataMemoryInstalled:     installed,
				DataMemoryUsed:          used,
				DataMemoryConsumption:   float32(used) / float32(installed) * 100.0,
			},
			Gps: &pbapp.GpsStatus{
				Fix:        1,
				Time:       uint64(now.Unix()),
				Satellites: 5,
				Longitude:  -118.2709223,
				Latitude:   34.0318047,
				Altitude:   65,
			},
			Power: &pbapp.PowerStatus{
				Battery: &pbapp.BatteryStatus{
					Voltage:    3420.0,
					Percentage: 70.0,
				},
			},
		},
		NetworkSettings: &pbapp.NetworkSettings{},
		Streams: []*pbapp.DataStream{
			&pbapp.DataStream{
				Id:      0,
				Time:    0,
				Size:    0,
				Version: 0,
				Block:   0,
				Name:    "data.fkpb",
				Path:    "/fk/v1/download/data",
			},
			&pbapp.DataStream{
				Id:      1,
				Time:    0,
				Size:    0,
				Version: 0,
				Block:   0,
				Name:    "meta.fkpb",
				Path:    "/fk/v1/download/meta",
			},
		},
		Modules: []*pbapp.ModuleCapabilities{
			&pbapp.ModuleCapabilities{
				Position: 0,
				Flags:    1,
				Name:     "modules.diagnostics",
				Id:       hashString(fmt.Sprintf("diagnostics-%v-%v", s.DeviceID, e.Seed)),
				Header: &pbapp.ModuleHeader{
					Manufacturer: repositories.ManufacturerConservify,
					Kind:         repositories.ConservifyDiagnostics,
					Version:      1,
				},
				Sensors: []*pbapp.SensorCapabilities{
					&pbapp.SensorCapabilities{
						Number:        0,
						Name:          "memory",
						UnitOfMeasure: "bytes",
						Frequency:     60,
					},
				},
			},
			&pbapp.ModuleCapabilities{
				Position: 0,
				Name:     "modules.water.ph",
				Id:       hashString(fmt.Sprintf("ph-%v-%v", s.DeviceID, e.Seed)),
				Header: &pbapp.ModuleHeader{
					Manufacturer: repositories.ManufacturerConservify,
					Kind:         repositories.ConservifyWaterPh,
					Version:      1,
				},
				Sensors: []*pbapp.SensorCapabilities{
					&pbapp.SensorCapabilities{
						Number:        0,
						Name:          "ph",
						UnitOfMeasure: "",
						Frequency:     60,
					},
				},
			},
			&pbapp.ModuleCapabilities{
				Position: 1,
				Name:     "modules.water.do",
				Id:       hashString(fmt.Sprintf("do-%v-%v", s.DeviceID, e.Seed)),
				Header: &pbapp.ModuleHeader{
					Manufacturer: repositories.ManufacturerConservify,
					Kind:         repositories.ConservifyWaterDo,
					Version:      1,
				},
				Sensors: []*pbapp.SensorCapabilities{
					&pbapp.SensorCapabilities{
						Number:        0,
						Name:          "do",
						UnitOfMeasure: "",
						Frequency:     60,
					},
				},
			},
			&pbapp.ModuleCapabilities{
				Position: 2,
				Name:     "modules.water.ec",
				Id:       hashString(fmt.Sprintf("ec-%v-%v", s.DeviceID, e.Seed)),
				Header: &pbapp.ModuleHeader{
					Manufacturer: repositories.ManufacturerConservify,
					Kind:         repositories.ConservifyWaterEc,
					Version:      1,
				},
				Sensors: []*pbapp.SensorCapabilities{
					&pbapp.SensorCapabilities{
						Number:        0,
						Name:          "ec",
						UnitOfMeasure: "µS/cm",
						Frequency:     60,
					},
					&pbapp.SensorCapabilities{
						Number:        1,
						Name:          "temperature",
						UnitOfMeasure: "C",
						Frequency:     60,
					},
					&pbapp.SensorCapabilities{
						Number:        2,
						Name:          "depth",
						UnitOfMeasure: "m",
						Frequency:     60,
					},
					&pbapp.SensorCapabilities{
						Number:        2,
						Name:          "depth (mv)",
						UnitOfMeasure: "mv",
						Frequency:     60,
						Flags:         1,
					},
				},
			},
		},
	}
}

func (e *TestEnv) NewLiveReadingsReply(s *data.Station) *pbapp.HttpReply {
	status := e.NewHttpStatusReply(s)

	moduleLiveReadings := make([]*pbapp.LiveModuleReadings, 0, len(status.Modules))

	for _, module := range status.Modules {
		sensors := make([]*pbapp.LiveSensorReading, 0, len(module.Sensors))

		for _, sensor := range module.Sensors {
			sensors = append(sensors, &pbapp.LiveSensorReading{
				Sensor: sensor,
				Value:  mrand.Float32(),
			})
		}

		moduleLiveReadings = append(moduleLiveReadings, &pbapp.LiveModuleReadings{
			Module:   module,
			Readings: sensors,
		})
	}

	status.LiveReadings = []*pbapp.LiveReadings{
		&pbapp.LiveReadings{
			Time:    uint64(0),
			Modules: moduleLiveReadings,
		},
	}

	return status
}

func (e *TestEnv) NewMetaLayout(record uint64) *SignedRecordAndData {
	cfg := &pb.DataRecord{
		Metadata: &pb.Metadata{
			DeviceId: []byte{},
			Firmware: &pb.Firmware{},
		},
		Identity: &pb.Identity{
			Name: "",
		},
		Modules: []*pb.ModuleInfo{
			&pb.ModuleInfo{
				Position: 0,
				Name:     "random-module-1",
				Id:       hashString(fmt.Sprintf("random-module-1-%d", e.Seed)),
				Header: &pb.ModuleHeader{
					Manufacturer: repositories.ManufacturerConservify,
					Kind:         repositories.ConservifyRandom,
					Version:      0x1,
				},
				Firmware: &pb.Firmware{},
				Sensors: []*pb.SensorInfo{
					&pb.SensorInfo{
						Name:          "random_0",
						UnitOfMeasure: "C",
					},
					&pb.SensorInfo{
						Name:          "random_1",
						UnitOfMeasure: "C",
					},
					&pb.SensorInfo{
						Name:          "random_2",
						UnitOfMeasure: "C",
					},
					&pb.SensorInfo{
						Name:          "random_3",
						UnitOfMeasure: "C",
					},
					&pb.SensorInfo{
						Name:          "random_4",
						UnitOfMeasure: "C",
					},
				},
			},
			&pb.ModuleInfo{
				Position: 1,
				Name:     "random-module-2",
				Id:       hashString(fmt.Sprintf("random-module-2-%d", e.Seed)),
				Header: &pb.ModuleHeader{
					Manufacturer: repositories.ManufacturerConservify,
					Kind:         repositories.ConservifyRandom,
					Version:      0x1,
				},
				Firmware: &pb.Firmware{},
				Sensors: []*pb.SensorInfo{
					&pb.SensorInfo{
						Name:          "random_0",
						UnitOfMeasure: "C",
					},
					&pb.SensorInfo{
						Name:          "random_1",
						UnitOfMeasure: "C",
					},
					&pb.SensorInfo{
						Name:          "random_2",
						UnitOfMeasure: "C",
					},
					&pb.SensorInfo{
						Name:          "random_3",
						UnitOfMeasure: "C",
					},
					&pb.SensorInfo{
						Name:          "random_4",
						UnitOfMeasure: "C",
					},
					&pb.SensorInfo{
						Name:          "random_5",
						UnitOfMeasure: "C",
					},
					&pb.SensorInfo{
						Name:          "random_6",
						UnitOfMeasure: "C",
					},
					&pb.SensorInfo{
						Name:          "random_7",
						UnitOfMeasure: "C",
					},
					&pb.SensorInfo{
						Name:          "random_8",
						UnitOfMeasure: "C",
					},
					&pb.SensorInfo{
						Name:          "random_9",
						UnitOfMeasure: "C",
					},
				},
			},
		},
	}

	plain := proto.NewBuffer(make([]byte, 0))
	plain.Marshal(cfg)

	delimited := proto.NewBuffer(make([]byte, 0))
	delimited.EncodeMessage(cfg)

	hash := blake2b.Sum256(delimited.Bytes())

	return &SignedRecordAndData{
		Signed: &pb.SignedRecord{
			Kind:   1, /* Modules */
			Time:   0,
			Data:   delimited.Bytes(),
			Hash:   hash[:],
			Record: record,
		},
		Bytes: plain.Bytes(),
		Data:  cfg,
	}
}

func (e *TestEnv) NewDataReading(meta, reading uint64) *pb.DataRecord {
	now := time.Now()

	return &pb.DataRecord{
		Readings: &pb.Readings{
			Time:    int64(now.Unix()),
			Reading: uint64(reading),
			Meta:    uint64(meta),
			Flags:   0,
			Location: &pb.DeviceLocation{
				Fix:        1,
				Time:       int64(now.Unix()),
				Longitude:  -118.2709223,
				Latitude:   34.0318047,
				Altitude:   mrand.Float32(),
				Satellites: 6,
			},
			SensorGroups: []*pb.SensorGroup{
				&pb.SensorGroup{
					Module: 0,
					Readings: []*pb.SensorAndValue{
						&pb.SensorAndValue{
							Sensor: 0,
							Value:  mrand.Float32(),
						},
						&pb.SensorAndValue{
							Sensor: 1,
							Value:  mrand.Float32(),
						},
						&pb.SensorAndValue{
							Sensor: 2,
							Value:  mrand.Float32(),
						},
						&pb.SensorAndValue{
							Sensor: 3,
							Value:  mrand.Float32(),
						},
						&pb.SensorAndValue{
							Sensor: 4,
							Value:  mrand.Float32(),
						},
					},
				},
				&pb.SensorGroup{
					Module: 1,
					Readings: []*pb.SensorAndValue{
						&pb.SensorAndValue{
							Sensor: 0,
							Value:  mrand.Float32(),
						},
						&pb.SensorAndValue{
							Sensor: 1,
							Value:  mrand.Float32(),
						},
						&pb.SensorAndValue{
							Sensor: 2,
							Value:  mrand.Float32(),
						},
						&pb.SensorAndValue{
							Sensor: 3,
							Value:  mrand.Float32(),
						},
						&pb.SensorAndValue{
							Sensor: 4,
							Value:  mrand.Float32(),
						},
						&pb.SensorAndValue{
							Sensor: 5,
							Value:  mrand.Float32(),
						},
						&pb.SensorAndValue{
							Sensor: 6,
							Value:  mrand.Float32(),
						},
						&pb.SensorAndValue{
							Sensor: 7,
							Value:  mrand.Float32(),
						},
						&pb.SensorAndValue{
							Sensor: 8,
							Value:  mrand.Float32(),
						},
						&pb.SensorAndValue{
							Sensor: 9,
							Value:  mrand.Float32(),
						},
					},
				},
			},
		},
	}
}

type AddedDataRecords struct {
	Provision *data.Provision
	Ingestion *data.Ingestion
	Records   []*data.DataRecord
}

type AddedMetaRecords struct {
	Provision *data.Provision
	Ingestion *data.Ingestion
	Records   []*data.MetaRecord
}

type MetaAndData struct {
	Meta *AddedMetaRecords
	Data *AddedDataRecords
}

func (e *TestEnv) AddMetaAndData(station *data.Station, user *data.User, numberData int) (*MetaAndData, error) {
	recordRepository, err := repositories.NewRecordRepository(e.DB)
	if err != nil {
		return nil, err
	}

	_, di, err := e.AddIngestion(user, "url", data.DataTypeName, station.DeviceID, 0)
	if err != nil {
		return nil, err
	}

	_, mi, err := e.AddIngestion(user, "url", data.MetaTypeName, station.DeviceID, 0)
	if err != nil {
		return nil, err
	}

	p, err := e.AddProvision(station.DeviceID, station.DeviceID)
	if err != nil {
		return nil, err
	}

	metaRecords := make([]*data.MetaRecord, 0)
	dataRecords := make([]*data.DataRecord, 0)

	metaNumber := uint64(1)
	dataNumber := uint64(1)

	for m := 0; m < 1; m += 1 {
		meta := e.NewMetaLayout(metaNumber)
		metaRecord, err := recordRepository.AddMetaRecord(e.Ctx, p, mi, meta.Signed, meta.Data, meta.Bytes)
		if err != nil {
			return nil, err
		}

		metaRecords = append(metaRecords, metaRecord)

		for d := 0; d < numberData; d += 1 {
			data := e.NewDataReading(meta.Signed.Record, dataNumber)

			buffer := proto.NewBuffer(make([]byte, 0))
			if err := buffer.Marshal(data); err != nil {
				return nil, err
			}

			dataRecord, _, err := recordRepository.AddDataRecord(e.Ctx, p, di, data, buffer.Bytes())
			if err != nil {
				return nil, err
			}

			dataRecords = append(dataRecords, dataRecord)

			dataNumber += 1
		}
		metaNumber += 1
	}

	return &MetaAndData{
		Meta: &AddedMetaRecords{
			Provision: p,
			Ingestion: mi,
			Records:   metaRecords,
		},
		Data: &AddedDataRecords{
			Provision: p,
			Ingestion: di,
			Records:   dataRecords,
		},
	}, nil
}

type FilePair struct {
	Meta       []byte
	Data       []byte
	MetaBlocks []uint64
	DataBlocks []uint64
}

func (e *TestEnv) NewFilePair(nmeta, ndata int) (*FilePair, error) {
	dataFile := proto.NewBuffer(make([]byte, 0))
	metaFile := proto.NewBuffer(make([]byte, 0))

	metaNumber := uint64(1)
	dataNumber := uint64(1)

	for m := 0; m < nmeta; m += 1 {
		meta := e.NewMetaLayout(metaNumber)
		if err := metaFile.EncodeMessage(meta.Signed); err != nil {
			return nil, err
		}

		metaNumber += 1

		for i := 0; i < ndata; i += 1 {
			dataRecord := e.NewDataReading(meta.Signed.Record, dataNumber)
			if err := dataFile.EncodeMessage(dataRecord); err != nil {
				return nil, err
			}

			dataNumber += 1
		}
	}

	return &FilePair{
		Meta:       metaFile.Bytes(),
		Data:       dataFile.Bytes(),
		MetaBlocks: []uint64{1, metaNumber},
		DataBlocks: []uint64{1, dataNumber},
	}, nil
}

func (e *TestEnv) MustDeviceID() []byte {
	id, err := e.NewRandomData(16)
	if err != nil {
		panic(err)
	}
	return id
}
