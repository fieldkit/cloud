// Code generated by goa v3.1.2, DO NOT EDIT.
//
// station views
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package views

import (
	goa "goa.design/goa/v3/pkg"
)

// StationFull is the viewed result type that is projected based on a view.
type StationFull struct {
	// Type to project
	Projected *StationFullView
	// View to render
	View string
}

// StationFullView is a type that runs validations on a projected type.
type StationFullView struct {
	ID                 *int32
	Name               *string
	Owner              *StationOwnerView
	DeviceID           *string
	Uploads            []*StationUploadView
	Images             []*ImageRefView
	Photos             *StationPhotosView
	ReadOnly           *bool
	Battery            *float32
	RecordingStartedAt *int64
	MemoryUsed         *int32
	MemoryAvailable    *int32
	FirmwareNumber     *int32
	FirmwareTime       *int32
	Modules            []*StationModuleView
}

// StationOwnerView is a type that runs validations on a projected type.
type StationOwnerView struct {
	ID   *int32
	Name *string
}

// StationUploadView is a type that runs validations on a projected type.
type StationUploadView struct {
	ID       *int64
	Time     *int64
	UploadID *string
	Size     *int64
	URL      *string
	Type     *string
	Blocks   []int64
}

// ImageRefView is a type that runs validations on a projected type.
type ImageRefView struct {
	URL *string
}

// StationPhotosView is a type that runs validations on a projected type.
type StationPhotosView struct {
	Small *string
}

// StationModuleView is a type that runs validations on a projected type.
type StationModuleView struct {
	ID       *string
	Name     *string
	Position *int32
	Sensors  []*StationSensorView
}

// StationSensorView is a type that runs validations on a projected type.
type StationSensorView struct {
	Name          *string
	UnitOfMeasure *string
}

var (
	// StationFullMap is a map of attribute names in result type StationFull
	// indexed by view name.
	StationFullMap = map[string][]string{
		"default": []string{
			"id",
			"name",
			"owner",
			"device_id",
			"uploads",
			"images",
			"photos",
			"read_only",
			"battery",
			"recording_started_at",
			"memory_used",
			"memory_available",
			"firmware_number",
			"firmware_time",
			"modules",
		},
	}
)

// ValidateStationFull runs the validations defined on the viewed result type
// StationFull.
func ValidateStationFull(result *StationFull) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateStationFullView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default"})
	}
	return
}

// ValidateStationFullView runs the validations defined on StationFullView
// using the "default" view.
func ValidateStationFullView(result *StationFullView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	if result.Owner == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("owner", "result"))
	}
	if result.DeviceID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("device_id", "result"))
	}
	if result.Uploads == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("uploads", "result"))
	}
	if result.Images == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("images", "result"))
	}
	if result.Photos == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("photos", "result"))
	}
	if result.ReadOnly == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("read_only", "result"))
	}
	if result.Battery == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("battery", "result"))
	}
	if result.RecordingStartedAt == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("recording_started_at", "result"))
	}
	if result.MemoryUsed == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("memory_used", "result"))
	}
	if result.MemoryAvailable == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("memory_available", "result"))
	}
	if result.FirmwareNumber == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("firmware_number", "result"))
	}
	if result.FirmwareTime == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("firmware_time", "result"))
	}
	if result.Modules == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("modules", "result"))
	}
	if result.Owner != nil {
		if err2 := ValidateStationOwnerView(result.Owner); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	for _, e := range result.Uploads {
		if e != nil {
			if err2 := ValidateStationUploadView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	for _, e := range result.Images {
		if e != nil {
			if err2 := ValidateImageRefView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	if result.Photos != nil {
		if err2 := ValidateStationPhotosView(result.Photos); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	for _, e := range result.Modules {
		if e != nil {
			if err2 := ValidateStationModuleView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateStationOwnerView runs the validations defined on StationOwnerView.
func ValidateStationOwnerView(result *StationOwnerView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	return
}

// ValidateStationUploadView runs the validations defined on StationUploadView.
func ValidateStationUploadView(result *StationUploadView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Time == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("time", "result"))
	}
	if result.UploadID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("upload_id", "result"))
	}
	if result.Size == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("size", "result"))
	}
	if result.URL == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("url", "result"))
	}
	if result.Type == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("type", "result"))
	}
	if result.Blocks == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("blocks", "result"))
	}
	return
}

// ValidateImageRefView runs the validations defined on ImageRefView.
func ValidateImageRefView(result *ImageRefView) (err error) {
	if result.URL == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("url", "result"))
	}
	return
}

// ValidateStationPhotosView runs the validations defined on StationPhotosView.
func ValidateStationPhotosView(result *StationPhotosView) (err error) {
	if result.Small == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("small", "result"))
	}
	return
}

// ValidateStationModuleView runs the validations defined on StationModuleView.
func ValidateStationModuleView(result *StationModuleView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	if result.Position == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("position", "result"))
	}
	if result.Sensors == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("sensors", "result"))
	}
	for _, e := range result.Sensors {
		if e != nil {
			if err2 := ValidateStationSensorView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateStationSensorView runs the validations defined on StationSensorView.
func ValidateStationSensorView(result *StationSensorView) (err error) {
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	if result.UnitOfMeasure == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("unit_of_measure", "result"))
	}
	return
}
