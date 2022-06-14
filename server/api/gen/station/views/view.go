// Code generated by goa v3.2.4, DO NOT EDIT.
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

// StationsFull is the viewed result type that is projected based on a view.
type StationsFull struct {
	// Type to project
	Projected *StationsFullView
	// View to render
	View string
}

// AssociatedStations is the viewed result type that is projected based on a
// view.
type AssociatedStations struct {
	// Type to project
	Projected *AssociatedStationsView
	// View to render
	View string
}

// DownloadedPhoto is the viewed result type that is projected based on a view.
type DownloadedPhoto struct {
	// Type to project
	Projected *DownloadedPhotoView
	// View to render
	View string
}

// PageOfStations is the viewed result type that is projected based on a view.
type PageOfStations struct {
	// Type to project
	Projected *PageOfStationsView
	// View to render
	View string
}

// StationProgress is the viewed result type that is projected based on a view.
type StationProgress struct {
	// Type to project
	Projected *StationProgressView
	// View to render
	View string
}

// StationFullView is a type that runs validations on a projected type.
type StationFullView struct {
	ID                 *int32
	Name               *string
	Model              *StationFullModelView
	Owner              *StationOwnerView
	DeviceID           *string
	Interestingness    *StationInterestingnessView
	Attributes         *StationProjectAttributesView
	Uploads            []*StationUploadView
	Photos             *StationPhotosView
	ReadOnly           *bool
	Battery            *float32
	RecordingStartedAt *int64
	MemoryUsed         *int32
	MemoryAvailable    *int32
	FirmwareNumber     *int32
	FirmwareTime       *int64
	Configurations     *StationConfigurationsView
	UpdatedAt          *int64
	LocationName       *string
	PlaceNameOther     *string
	PlaceNameNative    *string
	Location           *StationLocationView
	SyncedAt           *int64
	IngestionAt        *int64
	Data               *StationDataSummaryView
}

// StationFullModelView is a type that runs validations on a projected type.
type StationFullModelView struct {
	Name                      *string
	OnlyVisibleViaAssociation *bool
}

// StationOwnerView is a type that runs validations on a projected type.
type StationOwnerView struct {
	ID   *int32
	Name *string
}

// StationInterestingnessView is a type that runs validations on a projected
// type.
type StationInterestingnessView struct {
	Windows []*StationInterestingnessWindowView
}

// StationInterestingnessWindowView is a type that runs validations on a
// projected type.
type StationInterestingnessWindowView struct {
	Seconds         *int32
	Interestingness *float64
	Value           *float64
	Time            *int64
}

// StationProjectAttributesView is a type that runs validations on a projected
// type.
type StationProjectAttributesView struct {
	Attributes []*StationProjectAttributeView
}

// StationProjectAttributeView is a type that runs validations on a projected
// type.
type StationProjectAttributeView struct {
	ProjectID   *int32
	AttributeID *int64
	Name        *string
	StringValue *string
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

// StationPhotosView is a type that runs validations on a projected type.
type StationPhotosView struct {
	Small *string
}

// StationConfigurationsView is a type that runs validations on a projected
// type.
type StationConfigurationsView struct {
	All []*StationConfigurationView
}

// StationConfigurationView is a type that runs validations on a projected type.
type StationConfigurationView struct {
	ID           *int64
	Time         *int64
	ProvisionID  *int64
	MetaRecordID *int64
	SourceID     *int32
	Modules      []*StationModuleView
}

// StationModuleView is a type that runs validations on a projected type.
type StationModuleView struct {
	ID           *int64
	HardwareID   *string
	MetaRecordID *int64
	Name         *string
	Position     *int32
	Flags        *int32
	Internal     *bool
	FullKey      *string
	Sensors      []*StationSensorView
	Meta         map[string]interface{}
}

// StationSensorView is a type that runs validations on a projected type.
type StationSensorView struct {
	Name          *string
	UnitOfMeasure *string
	Reading       *SensorReadingView
	Key           *string
	FullKey       *string
	Ranges        []*SensorRangeView
	Meta          map[string]interface{}
}

// SensorReadingView is a type that runs validations on a projected type.
type SensorReadingView struct {
	Last *float32
	Time *int64
}

// SensorRangeView is a type that runs validations on a projected type.
type SensorRangeView struct {
	Minimum *float32
	Maximum *float32
}

// StationLocationView is a type that runs validations on a projected type.
type StationLocationView struct {
	Precise []float64
	Regions []*StationRegionView
	URL     *string
}

// StationRegionView is a type that runs validations on a projected type.
type StationRegionView struct {
	Name  *string
	Shape [][][]float64
}

// StationDataSummaryView is a type that runs validations on a projected type.
type StationDataSummaryView struct {
	Start           *int64
	End             *int64
	NumberOfSamples *int64
}

// StationsFullView is a type that runs validations on a projected type.
type StationsFullView struct {
	Stations StationFullCollectionView
}

// StationFullCollectionView is a type that runs validations on a projected
// type.
type StationFullCollectionView []*StationFullView

// AssociatedStationsView is a type that runs validations on a projected type.
type AssociatedStationsView struct {
	Stations AssociatedStationCollectionView
}

// AssociatedStationCollectionView is a type that runs validations on a
// projected type.
type AssociatedStationCollectionView []*AssociatedStationView

// AssociatedStationView is a type that runs validations on a projected type.
type AssociatedStationView struct {
	Station  *StationFullView
	Project  *AssociatedViaProjectView
	Location *AssociatedViaLocationView
	Manual   *AssociatedViaManualView
}

// AssociatedViaProjectView is a type that runs validations on a projected type.
type AssociatedViaProjectView struct {
	ID *int32
}

// AssociatedViaLocationView is a type that runs validations on a projected
// type.
type AssociatedViaLocationView struct {
	Distance *float32
}

// AssociatedViaManualView is a type that runs validations on a projected type.
type AssociatedViaManualView struct {
	Priority *int32
}

// DownloadedPhotoView is a type that runs validations on a projected type.
type DownloadedPhotoView struct {
	Length      *int64
	ContentType *string
	Etag        *string
	Body        []byte
}

// PageOfStationsView is a type that runs validations on a projected type.
type PageOfStationsView struct {
	Stations []*EssentialStationView
	Total    *int32
}

// EssentialStationView is a type that runs validations on a projected type.
type EssentialStationView struct {
	ID                 *int64
	DeviceID           *string
	Name               *string
	Owner              *StationOwnerView
	CreatedAt          *int64
	UpdatedAt          *int64
	RecordingStartedAt *int64
	MemoryUsed         *int32
	MemoryAvailable    *int32
	FirmwareNumber     *int32
	FirmwareTime       *int64
	Location           *StationLocationView
	LastIngestionAt    *int64
}

// StationProgressView is a type that runs validations on a projected type.
type StationProgressView struct {
	Jobs []*StationJobView
}

// StationJobView is a type that runs validations on a projected type.
type StationJobView struct {
	Title       *string
	StartedAt   *int64
	CompletedAt *int64
	Progress    *float32
}

var (
	// StationFullMap is a map of attribute names in result type StationFull
	// indexed by view name.
	StationFullMap = map[string][]string{
		"default": []string{
			"id",
			"name",
			"model",
			"owner",
			"deviceId",
			"interestingness",
			"attributes",
			"uploads",
			"photos",
			"readOnly",
			"battery",
			"recordingStartedAt",
			"memoryUsed",
			"memoryAvailable",
			"firmwareNumber",
			"firmwareTime",
			"configurations",
			"updatedAt",
			"location",
			"locationName",
			"placeNameOther",
			"placeNameNative",
			"syncedAt",
			"ingestionAt",
			"data",
		},
	}
	// StationsFullMap is a map of attribute names in result type StationsFull
	// indexed by view name.
	StationsFullMap = map[string][]string{
		"default": []string{
			"stations",
		},
	}
	// AssociatedStationsMap is a map of attribute names in result type
	// AssociatedStations indexed by view name.
	AssociatedStationsMap = map[string][]string{
		"default": []string{
			"stations",
		},
	}
	// DownloadedPhotoMap is a map of attribute names in result type
	// DownloadedPhoto indexed by view name.
	DownloadedPhotoMap = map[string][]string{
		"default": []string{
			"length",
			"body",
			"contentType",
			"etag",
		},
	}
	// PageOfStationsMap is a map of attribute names in result type PageOfStations
	// indexed by view name.
	PageOfStationsMap = map[string][]string{
		"default": []string{
			"stations",
			"total",
		},
	}
	// StationProgressMap is a map of attribute names in result type
	// StationProgress indexed by view name.
	StationProgressMap = map[string][]string{
		"default": []string{
			"jobs",
		},
	}
	// StationLocationMap is a map of attribute names in result type
	// StationLocation indexed by view name.
	StationLocationMap = map[string][]string{
		"default": []string{
			"precise",
			"regions",
			"url",
		},
	}
	// StationFullCollectionMap is a map of attribute names in result type
	// StationFullCollection indexed by view name.
	StationFullCollectionMap = map[string][]string{
		"default": []string{
			"id",
			"name",
			"model",
			"owner",
			"deviceId",
			"interestingness",
			"attributes",
			"uploads",
			"photos",
			"readOnly",
			"battery",
			"recordingStartedAt",
			"memoryUsed",
			"memoryAvailable",
			"firmwareNumber",
			"firmwareTime",
			"configurations",
			"updatedAt",
			"location",
			"locationName",
			"placeNameOther",
			"placeNameNative",
			"syncedAt",
			"ingestionAt",
			"data",
		},
	}
	// AssociatedStationCollectionMap is a map of attribute names in result type
	// AssociatedStationCollection indexed by view name.
	AssociatedStationCollectionMap = map[string][]string{
		"default": []string{
			"station",
			"project",
			"location",
			"manual",
		},
	}
	// AssociatedStationMap is a map of attribute names in result type
	// AssociatedStation indexed by view name.
	AssociatedStationMap = map[string][]string{
		"default": []string{
			"station",
			"project",
			"location",
			"manual",
		},
	}
	// StationJobMap is a map of attribute names in result type StationJob indexed
	// by view name.
	StationJobMap = map[string][]string{
		"default": []string{
			"startedAt",
			"completedAt",
			"progress",
			"title",
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

// ValidateStationsFull runs the validations defined on the viewed result type
// StationsFull.
func ValidateStationsFull(result *StationsFull) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateStationsFullView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default"})
	}
	return
}

// ValidateAssociatedStations runs the validations defined on the viewed result
// type AssociatedStations.
func ValidateAssociatedStations(result *AssociatedStations) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateAssociatedStationsView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default"})
	}
	return
}

// ValidateDownloadedPhoto runs the validations defined on the viewed result
// type DownloadedPhoto.
func ValidateDownloadedPhoto(result *DownloadedPhoto) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateDownloadedPhotoView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default"})
	}
	return
}

// ValidatePageOfStations runs the validations defined on the viewed result
// type PageOfStations.
func ValidatePageOfStations(result *PageOfStations) (err error) {
	switch result.View {
	case "default", "":
		err = ValidatePageOfStationsView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default"})
	}
	return
}

// ValidateStationProgress runs the validations defined on the viewed result
// type StationProgress.
func ValidateStationProgress(result *StationProgress) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateStationProgressView(result.Projected)
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
	if result.Model == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("model", "result"))
	}
	if result.Owner == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("owner", "result"))
	}
	if result.DeviceID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("deviceId", "result"))
	}
	if result.Interestingness == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("interestingness", "result"))
	}
	if result.Attributes == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("attributes", "result"))
	}
	if result.Uploads == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("uploads", "result"))
	}
	if result.Photos == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("photos", "result"))
	}
	if result.ReadOnly == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("readOnly", "result"))
	}
	if result.Configurations == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("configurations", "result"))
	}
	if result.UpdatedAt == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("updatedAt", "result"))
	}
	if result.Model != nil {
		if err2 := ValidateStationFullModelView(result.Model); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if result.Owner != nil {
		if err2 := ValidateStationOwnerView(result.Owner); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if result.Interestingness != nil {
		if err2 := ValidateStationInterestingnessView(result.Interestingness); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if result.Attributes != nil {
		if err2 := ValidateStationProjectAttributesView(result.Attributes); err2 != nil {
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
	if result.Photos != nil {
		if err2 := ValidateStationPhotosView(result.Photos); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if result.Configurations != nil {
		if err2 := ValidateStationConfigurationsView(result.Configurations); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if result.Data != nil {
		if err2 := ValidateStationDataSummaryView(result.Data); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if result.Location != nil {
		if err2 := ValidateStationLocationView(result.Location); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateStationFullModelView runs the validations defined on
// StationFullModelView.
func ValidateStationFullModelView(result *StationFullModelView) (err error) {
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	if result.OnlyVisibleViaAssociation == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("only_visible_via_association", "result"))
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

// ValidateStationInterestingnessView runs the validations defined on
// StationInterestingnessView.
func ValidateStationInterestingnessView(result *StationInterestingnessView) (err error) {
	if result.Windows == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("windows", "result"))
	}
	for _, e := range result.Windows {
		if e != nil {
			if err2 := ValidateStationInterestingnessWindowView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateStationInterestingnessWindowView runs the validations defined on
// StationInterestingnessWindowView.
func ValidateStationInterestingnessWindowView(result *StationInterestingnessWindowView) (err error) {
	if result.Seconds == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("seconds", "result"))
	}
	if result.Interestingness == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("interestingness", "result"))
	}
	if result.Value == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("value", "result"))
	}
	if result.Time == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("time", "result"))
	}
	return
}

// ValidateStationProjectAttributesView runs the validations defined on
// StationProjectAttributesView.
func ValidateStationProjectAttributesView(result *StationProjectAttributesView) (err error) {
	if result.Attributes == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("attributes", "result"))
	}
	for _, e := range result.Attributes {
		if e != nil {
			if err2 := ValidateStationProjectAttributeView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateStationProjectAttributeView runs the validations defined on
// StationProjectAttributeView.
func ValidateStationProjectAttributeView(result *StationProjectAttributeView) (err error) {
	if result.ProjectID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("project_id", "result"))
	}
	if result.AttributeID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("attribute_id", "result"))
	}
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	if result.StringValue == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("string_value", "result"))
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
		err = goa.MergeErrors(err, goa.MissingFieldError("uploadId", "result"))
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

// ValidateStationPhotosView runs the validations defined on StationPhotosView.
func ValidateStationPhotosView(result *StationPhotosView) (err error) {
	if result.Small == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("small", "result"))
	}
	return
}

// ValidateStationConfigurationsView runs the validations defined on
// StationConfigurationsView.
func ValidateStationConfigurationsView(result *StationConfigurationsView) (err error) {
	if result.All == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("all", "result"))
	}
	for _, e := range result.All {
		if e != nil {
			if err2 := ValidateStationConfigurationView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateStationConfigurationView runs the validations defined on
// StationConfigurationView.
func ValidateStationConfigurationView(result *StationConfigurationView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.ProvisionID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("provisionId", "result"))
	}
	if result.Time == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("time", "result"))
	}
	if result.Modules == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("modules", "result"))
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
	if result.Flags == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("flags", "result"))
	}
	if result.Internal == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("internal", "result"))
	}
	if result.FullKey == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fullKey", "result"))
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
		err = goa.MergeErrors(err, goa.MissingFieldError("unitOfMeasure", "result"))
	}
	if result.Key == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("key", "result"))
	}
	if result.FullKey == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fullKey", "result"))
	}
	if result.Ranges == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("ranges", "result"))
	}
	if result.Reading != nil {
		if err2 := ValidateSensorReadingView(result.Reading); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	for _, e := range result.Ranges {
		if e != nil {
			if err2 := ValidateSensorRangeView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateSensorReadingView runs the validations defined on SensorReadingView.
func ValidateSensorReadingView(result *SensorReadingView) (err error) {
	if result.Last == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("last", "result"))
	}
	if result.Time == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("time", "result"))
	}
	return
}

// ValidateSensorRangeView runs the validations defined on SensorRangeView.
func ValidateSensorRangeView(result *SensorRangeView) (err error) {
	if result.Minimum == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("minimum", "result"))
	}
	if result.Maximum == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("maximum", "result"))
	}
	return
}

// ValidateStationLocationView runs the validations defined on
// StationLocationView using the "default" view.
func ValidateStationLocationView(result *StationLocationView) (err error) {
	for _, e := range result.Regions {
		if e != nil {
			if err2 := ValidateStationRegionView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateStationRegionView runs the validations defined on StationRegionView.
func ValidateStationRegionView(result *StationRegionView) (err error) {
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	if result.Shape == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("shape", "result"))
	}
	return
}

// ValidateStationDataSummaryView runs the validations defined on
// StationDataSummaryView.
func ValidateStationDataSummaryView(result *StationDataSummaryView) (err error) {
	if result.Start == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("start", "result"))
	}
	if result.End == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("end", "result"))
	}
	if result.NumberOfSamples == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("numberOfSamples", "result"))
	}
	return
}

// ValidateStationsFullView runs the validations defined on StationsFullView
// using the "default" view.
func ValidateStationsFullView(result *StationsFullView) (err error) {

	if result.Stations != nil {
		if err2 := ValidateStationFullCollectionView(result.Stations); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateStationFullCollectionView runs the validations defined on
// StationFullCollectionView using the "default" view.
func ValidateStationFullCollectionView(result StationFullCollectionView) (err error) {
	for _, item := range result {
		if err2 := ValidateStationFullView(item); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateAssociatedStationsView runs the validations defined on
// AssociatedStationsView using the "default" view.
func ValidateAssociatedStationsView(result *AssociatedStationsView) (err error) {

	if result.Stations != nil {
		if err2 := ValidateAssociatedStationCollectionView(result.Stations); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateAssociatedStationCollectionView runs the validations defined on
// AssociatedStationCollectionView using the "default" view.
func ValidateAssociatedStationCollectionView(result AssociatedStationCollectionView) (err error) {
	for _, item := range result {
		if err2 := ValidateAssociatedStationView(item); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateAssociatedStationView runs the validations defined on
// AssociatedStationView using the "default" view.
func ValidateAssociatedStationView(result *AssociatedStationView) (err error) {
	if result.Project != nil {
		if err2 := ValidateAssociatedViaProjectView(result.Project); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if result.Location != nil {
		if err2 := ValidateAssociatedViaLocationView(result.Location); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if result.Manual != nil {
		if err2 := ValidateAssociatedViaManualView(result.Manual); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if result.Station != nil {
		if err2 := ValidateStationFullView(result.Station); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateAssociatedViaProjectView runs the validations defined on
// AssociatedViaProjectView.
func ValidateAssociatedViaProjectView(result *AssociatedViaProjectView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	return
}

// ValidateAssociatedViaLocationView runs the validations defined on
// AssociatedViaLocationView.
func ValidateAssociatedViaLocationView(result *AssociatedViaLocationView) (err error) {
	if result.Distance == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("distance", "result"))
	}
	return
}

// ValidateAssociatedViaManualView runs the validations defined on
// AssociatedViaManualView.
func ValidateAssociatedViaManualView(result *AssociatedViaManualView) (err error) {
	if result.Priority == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("priority", "result"))
	}
	return
}

// ValidateDownloadedPhotoView runs the validations defined on
// DownloadedPhotoView using the "default" view.
func ValidateDownloadedPhotoView(result *DownloadedPhotoView) (err error) {
	if result.Length == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("length", "result"))
	}
	if result.ContentType == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("contentType", "result"))
	}
	if result.Etag == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("etag", "result"))
	}
	return
}

// ValidatePageOfStationsView runs the validations defined on
// PageOfStationsView using the "default" view.
func ValidatePageOfStationsView(result *PageOfStationsView) (err error) {
	if result.Stations == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("stations", "result"))
	}
	if result.Total == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("total", "result"))
	}
	for _, e := range result.Stations {
		if e != nil {
			if err2 := ValidateEssentialStationView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateEssentialStationView runs the validations defined on
// EssentialStationView.
func ValidateEssentialStationView(result *EssentialStationView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.DeviceID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("deviceId", "result"))
	}
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	if result.Owner == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("owner", "result"))
	}
	if result.CreatedAt == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("createdAt", "result"))
	}
	if result.UpdatedAt == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("updatedAt", "result"))
	}
	if result.Owner != nil {
		if err2 := ValidateStationOwnerView(result.Owner); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if result.Location != nil {
		if err2 := ValidateStationLocationView(result.Location); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateStationProgressView runs the validations defined on
// StationProgressView using the "default" view.
func ValidateStationProgressView(result *StationProgressView) (err error) {
	if result.Jobs == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("jobs", "result"))
	}
	for _, e := range result.Jobs {
		if e != nil {
			if err2 := ValidateStationJobView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateStationJobView runs the validations defined on StationJobView using
// the "default" view.
func ValidateStationJobView(result *StationJobView) (err error) {
	if result.StartedAt == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("startedAt", "result"))
	}
	if result.Progress == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("progress", "result"))
	}
	if result.Title == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("title", "result"))
	}
	return
}
