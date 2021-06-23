package data

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/jmoiron/sqlx/types"

	"github.com/lib/pq"

	"github.com/golang/protobuf/proto"

	pb "github.com/fieldkit/data-protocol"
)

const (
	MetaTypeName = "meta"
	DataTypeName = "data"
)

type QueuedIngestion struct {
	ID           int64      `db:"id"`
	IngestionID  int64      `db:"ingestion_id"`
	Queued       time.Time  `db:"queued"`
	Attempted    *time.Time `db:"attempted"`
	Completed    *time.Time `db:"completed"`
	TotalRecords *int64     `db:"total_records"`
	OtherErrors  *int64     `db:"other_errors"`
	MetaErrors   *int64     `db:"meta_errors"`
	DataErrors   *int64     `db:"data_errors"`
}

type Ingestion struct {
	ID           int64         `db:"id"`
	Time         time.Time     `db:"time"`
	UploadID     string        `db:"upload_id"`
	UserID       int32         `db:"user_id"`
	DeviceID     []byte        `db:"device_id"`
	GenerationID []byte        `db:"generation"`
	Size         int64         `db:"size"`
	URL          string        `db:"url"`
	Type         string        `db:"type"`
	Blocks       Int64Range    `db:"blocks"`
	Flags        pq.Int64Array `db:"flags"`
}

type Provision struct {
	ID           int64     `db:"id"`
	Created      time.Time `db:"created"`
	Updated      time.Time `db:"updated"`
	DeviceID     []byte    `db:"device_id"`
	GenerationID []byte    `db:"generation"`
}

type Int64Range []int64

// Scan implements the sql.Scanner interface.
func (a *Int64Range) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.parseString(string(src))
	case string:
		return a.parseString(src)
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pq: cannot convert %T to Int64Range", src)
}

func (a *Int64Range) ToInt64Array() []int64 {
	rv := make([]int64, len(*a))
	for i, value := range *a {
		rv[i] = int64(value)
	}
	return rv
}

func (a *Int64Range) ToIntArray() []int {
	rv := make([]int, len(*a))
	for i, value := range *a {
		rv[i] = int(value)
	}
	return rv
}

func (a *Int64Range) parseString(s string) error {
	if s[0] != '[' || s[len(s)-1] != ')' {
		return fmt.Errorf("Unexpected range boundaries. I was lazy.")
	}

	values := s[1 : len(s)-1]
	b, err := ParseBlocks(values)
	if err != nil {
		return err
	}

	b[1] = b[1] - 1

	*a = b

	return nil
}

// Value implements the driver.Valuer interface.
func (a Int64Range) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '['

		b = strconv.AppendInt(b, a[0], 10)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendInt(b, a[i], 10)
		}

		return string(append(b, ']')), nil
	}

	return "{}", nil
}

func ParseBlocks(s string) ([]int64, error) {
	parts := strings.Split(s, ",")

	if len(parts) != 2 {
		return nil, fmt.Errorf("malformed block range")
	}

	blocks := make([]int64, 2)
	for i, p := range parts {
		b, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil, err
		}
		blocks[i] = int64(b)
	}

	return blocks, nil
}

type DataRecord struct {
	ID           int64          `db:"id" json:"id"`
	ProvisionID  int64          `db:"provision_id" json:"provision_id"`
	Time         time.Time      `db:"time" json:"time"`
	Number       int64          `db:"number" json:"number"`
	MetaRecordID int64          `db:"meta_record_id" json:"meta_record_id"`
	Location     *Location      `db:"location" json:"location"`
	Data         types.JSONText `db:"raw" json:"raw"`
	PB           []byte         `db:"pb" json:"pb"`
}

func (d *DataRecord) SetData(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	d.Data = jsonData
	return nil
}

func (d *DataRecord) GetData() (fields map[string]interface{}, err error) {
	err = json.Unmarshal(d.Data, &fields)
	if err != nil {
		return nil, err
	}
	return
}

func (d *DataRecord) Unmarshal(r *pb.DataRecord) error {
	if d.PB != nil {
		return r.Unmarshal(d.PB)
	}
	err := json.Unmarshal(d.Data, r)
	if err != nil {
		return fmt.Errorf("error parsing data record json: %v", err)
	}
	return nil
}

type MetaRecord struct {
	ID          int64          `db:"id" json:"id"`
	ProvisionID int64          `db:"provision_id" json:"provision_id"`
	Time        time.Time      `db:"time" json:"time"`
	Number      int64          `db:"number" json:"number"`
	Data        types.JSONText `db:"raw" json:"raw"`
	PB          []byte         `db:"pb" json:"pb"`
}

func (d *MetaRecord) SetData(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	d.Data = jsonData
	return nil
}

func (d *MetaRecord) GetData() (fields map[string]interface{}, err error) {
	err = json.Unmarshal(d.Data, &fields)
	if err != nil {
		return nil, err
	}
	return
}

func (d *MetaRecord) Unmarshal(r *pb.DataRecord) error {
	if d.PB != nil {
		buffer := proto.NewBuffer(d.PB)
		if err := buffer.Unmarshal(r); err != nil {
			return fmt.Errorf("error parsing meta record pb: %v", err)
		}
		return nil
	}
	err := json.Unmarshal(d.Data, r)
	if err != nil {
		return fmt.Errorf("error parsing meta record json: %v", err)
	}
	return nil
}

func DecodeBinaryString(s string) ([]byte, error) {
	if s[8] == '-' {
		uid, err := uuid.Parse(s)
		if err == nil {
			return uid[:], nil
		}
	}

	bytes, err := hex.DecodeString(s)
	if err == nil {
		return bytes, nil
	}

	bytes, err = base64.StdEncoding.DecodeString(s)
	if err == nil {
		return bytes, nil
	}

	return nil, fmt.Errorf("unable to decode binary string: %s", s)
}

type VisibleStationConfiguration struct {
	VisibleConfiguration
	StationConfiguration
}

type VisibleConfiguration struct {
	StationID       int32 `db:"station_id"`
	ConfigurationID int64 `db:"configuration_id"`
}

type StationConfiguration struct {
	ID           int64     `db:"id" json:"id"`
	ProvisionID  int64     `db:"provision_id" json:"provision_id"`
	MetaRecordID *int64    `db:"meta_record_id" json:"meta_record_id"`
	SourceID     *int32    `db:"source_id" json:"source_id"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type StationModule struct {
	ID              int64  `db:"id" json:"id"`
	ConfigurationID int64  `db:"configuration_id" json:"configuration_id"`
	HardwareID      []byte `db:"hardware_id" json:"hardware_id"`
	Index           uint32 `db:"module_index" json:"module_index"`
	Position        uint32 `db:"position" json:"position"`
	Name            string `db:"name" json:"name"`
	Manufacturer    uint32 `db:"manufacturer" json:"manufacturer"`
	Kind            uint32 `db:"kind" json:"kind"`
	Version         uint32 `db:"version" json:"version"`
	Flags           uint32 `db:"flags" json:"flags"`
}

type ModuleSensor struct {
	ID              int64      `db:"id" json:"id"`
	ModuleID        int64      `db:"module_id" json:"module_id"`
	ConfigurationID int64      `db:"configuration_id" json:"configuration_id"`
	Index           uint32     `db:"sensor_index" json:"sensor_index"`
	UnitOfMeasure   string     `db:"unit_of_measure" json:"unit_of_measure"`
	Name            string     `db:"name" json:"name"`
	ReadingValue    *float64   `db:"reading_last" json:"reading_last"`
	ReadingTime     *time.Time `db:"reading_time" json:"reading_time"`
}

type RecordRangeFlagEnum int

const (
	NoneRecordRangeFlag   RecordRangeFlagEnum = 0
	HiddenRecordRangeFlag RecordRangeFlagEnum = 1
)

type RecordRangeMeta struct {
	ID        int64     `db:"id" json:"id"`
	StationID int32     `db:"station_id" json:"station_id"`
	StartTime time.Time `db:"start_time" json:"start_time"`
	EndTime   time.Time `db:"end_time" json:"end_time"`
	Flags     uint32    `db:"flags" json:"flags"`
}

func sanitizeProtobufData(bytes []byte) ([]byte, types.JSONText, error) {
	if bytes == nil {
		return nil, nil, nil
	}

	buffer := proto.NewBuffer(bytes)
	parsed := &pb.DataRecord{}
	if err := buffer.Unmarshal(parsed); err != nil {
		return nil, nil, err
	}

	if parsed.Network != nil && parsed.Network.Networks != nil {
		for _, n := range parsed.Network.Networks {
			if n.Ssid != "" || n.Password != "" {
				n.Ssid = "<sensitive>"
				n.Password = "<sensitive>"
			}
		}
	}

	if parsed.Lora != nil {
		if parsed.Lora.AppKey != nil {
			parsed.Lora.AppKey = []byte("<sensitive>")
		}
		if parsed.Lora.NetworkSessionKey != nil {
			parsed.Lora.NetworkSessionKey = []byte("<sensitive>")
		}
		if parsed.Lora.AppSessionKey != nil {
			parsed.Lora.AppSessionKey = []byte("<sensitive>")
		}
	}

	if parsed.Transmission != nil && parsed.Transmission.Wifi != nil {
		if parsed.Transmission.Wifi.Token != "" {
			parsed.Transmission.Wifi.Token = "<sensitive>"
		}
	}

	newBytes := proto.NewBuffer(make([]byte, 0))
	newBytes.EncodeMessage(parsed)

	newJson, err := json.Marshal(parsed)
	if err != nil {
		return nil, nil, err
	}

	return newBytes.Bytes(), newJson, nil
}

func SanitizeMetaRecord(r *MetaRecord) *MetaRecord {
	newBytes, newJson, err := sanitizeProtobufData(r.PB)
	if err != nil {
		return nil
	}

	return &MetaRecord{
		ID:          r.ID,
		ProvisionID: r.ProvisionID,
		Time:        r.Time,
		Number:      r.Number,
		Data:        newJson,
		PB:          newBytes,
	}
}

func SanitizeDataRecord(r *DataRecord) *DataRecord {
	newBytes, newJson, err := sanitizeProtobufData(r.PB)
	if err != nil {
		return nil
	}

	return &DataRecord{
		ID:           r.ID,
		ProvisionID:  r.ProvisionID,
		Time:         r.Time,
		Number:       r.Number,
		MetaRecordID: r.MetaRecordID,
		Location:     nil,
		Data:         newJson,
		PB:           newBytes,
	}
}
