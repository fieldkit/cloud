// Code generated by goagen v1.3.0, DO NOT EDIT.
//
// API "fieldkit": Application Media Types
//
// Command:
// $ main

package app

import (
	"github.com/goadesign/goa"
	"time"
	"unicode/utf8"
)

// ProjectAdministrator media type (default view)
//
// Identifier: application/vnd.app.administrator+json; view=default
type ProjectAdministrator struct {
	ProjectID int `form:"projectId" json:"projectId" xml:"projectId"`
	UserID    int `form:"userId" json:"userId" xml:"userId"`
}

// Validate validates the ProjectAdministrator media type instance.
func (mt *ProjectAdministrator) Validate() (err error) {

	return
}

// ProjectAdministratorCollection is the media type for an array of ProjectAdministrator (default view)
//
// Identifier: application/vnd.app.administrator+json; type=collection; view=default
type ProjectAdministratorCollection []*ProjectAdministrator

// Validate validates the ProjectAdministratorCollection media type instance.
func (mt ProjectAdministratorCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ProjectAdministrators media type (default view)
//
// Identifier: application/vnd.app.administrators+json; view=default
type ProjectAdministrators struct {
	Administrators ProjectAdministratorCollection `form:"administrators" json:"administrators" xml:"administrators"`
}

// Validate validates the ProjectAdministrators media type instance.
func (mt *ProjectAdministrators) Validate() (err error) {
	if mt.Administrators == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "administrators"))
	}
	return
}

// ClusterGeometrySummary media type (default view)
//
// Identifier: application/vnd.app.cluster_geometry_summary+json; view=default
type ClusterGeometrySummary struct {
	Geometry [][]float64 `form:"geometry" json:"geometry" xml:"geometry"`
	ID       int         `form:"id" json:"id" xml:"id"`
	SourceID int         `form:"sourceId" json:"sourceId" xml:"sourceId"`
}

// Validate validates the ClusterGeometrySummary media type instance.
func (mt *ClusterGeometrySummary) Validate() (err error) {

	if mt.Geometry == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "geometry"))
	}
	return
}

// ClusterGeometrySummaryCollection is the media type for an array of ClusterGeometrySummary (default view)
//
// Identifier: application/vnd.app.cluster_geometry_summary+json; type=collection; view=default
type ClusterGeometrySummaryCollection []*ClusterGeometrySummary

// Validate validates the ClusterGeometrySummaryCollection media type instance.
func (mt ClusterGeometrySummaryCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// DeviceSchema media type (default view)
//
// Identifier: application/vnd.app.device_schema+json; view=default
type DeviceSchema struct {
	Active     bool   `form:"active" json:"active" xml:"active"`
	DeviceID   int    `form:"deviceId" json:"deviceId" xml:"deviceId"`
	JSONSchema string `form:"jsonSchema" json:"jsonSchema" xml:"jsonSchema"`
	Key        string `form:"key" json:"key" xml:"key"`
	ProjectID  int    `form:"projectId" json:"projectId" xml:"projectId"`
	SchemaID   int    `form:"schemaId" json:"schemaId" xml:"schemaId"`
}

// Validate validates the DeviceSchema media type instance.
func (mt *DeviceSchema) Validate() (err error) {

	if mt.JSONSchema == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "jsonSchema"))
	}
	if mt.Key == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "key"))
	}
	return
}

// DeviceSchemaCollection is the media type for an array of DeviceSchema (default view)
//
// Identifier: application/vnd.app.device_schema+json; type=collection; view=default
type DeviceSchemaCollection []*DeviceSchema

// Validate validates the DeviceSchemaCollection media type instance.
func (mt DeviceSchemaCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// DeviceSchemas media type (default view)
//
// Identifier: application/vnd.app.device_schemas+json; view=default
type DeviceSchemas struct {
	Schemas DeviceSchemaCollection `form:"schemas" json:"schemas" xml:"schemas"`
}

// Validate validates the DeviceSchemas media type instance.
func (mt *DeviceSchemas) Validate() (err error) {
	if mt.Schemas == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "schemas"))
	}
	if err2 := mt.Schemas.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// DeviceSource media type (default view)
//
// Identifier: application/vnd.app.device_source+json; view=default
type DeviceSource struct {
	Active       bool   `form:"active" json:"active" xml:"active"`
	ExpeditionID int    `form:"expeditionId" json:"expeditionId" xml:"expeditionId"`
	ID           int    `form:"id" json:"id" xml:"id"`
	Key          string `form:"key" json:"key" xml:"key"`
	Name         string `form:"name" json:"name" xml:"name"`
	TeamID       *int   `form:"teamId,omitempty" json:"teamId,omitempty" xml:"teamId,omitempty"`
	Token        string `form:"token" json:"token" xml:"token"`
	UserID       *int   `form:"userId,omitempty" json:"userId,omitempty" xml:"userId,omitempty"`
}

// Validate validates the DeviceSource media type instance.
func (mt *DeviceSource) Validate() (err error) {

	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if mt.Token == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "token"))
	}
	if mt.Key == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "key"))
	}
	return
}

// DeviceSource media type (public view)
//
// Identifier: application/vnd.app.device_source+json; view=public
type DeviceSourcePublic struct {
	Active           bool       `form:"active" json:"active" xml:"active"`
	Centroid         []float64  `form:"centroid,omitempty" json:"centroid,omitempty" xml:"centroid,omitempty"`
	EndTime          *time.Time `form:"endTime,omitempty" json:"endTime,omitempty" xml:"endTime,omitempty"`
	ExpeditionID     int        `form:"expeditionId" json:"expeditionId" xml:"expeditionId"`
	ID               int        `form:"id" json:"id" xml:"id"`
	LastFeatureID    *int       `form:"lastFeatureId,omitempty" json:"lastFeatureId,omitempty" xml:"lastFeatureId,omitempty"`
	Name             string     `form:"name" json:"name" xml:"name"`
	NumberOfFeatures *int       `form:"numberOfFeatures,omitempty" json:"numberOfFeatures,omitempty" xml:"numberOfFeatures,omitempty"`
	Radius           *float64   `form:"radius,omitempty" json:"radius,omitempty" xml:"radius,omitempty"`
	StartTime        *time.Time `form:"startTime,omitempty" json:"startTime,omitempty" xml:"startTime,omitempty"`
	TeamID           *int       `form:"teamId,omitempty" json:"teamId,omitempty" xml:"teamId,omitempty"`
	UserID           *int       `form:"userId,omitempty" json:"userId,omitempty" xml:"userId,omitempty"`
}

// Validate validates the DeviceSourcePublic media type instance.
func (mt *DeviceSourcePublic) Validate() (err error) {

	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	return
}

// DeviceSourceCollection is the media type for an array of DeviceSource (default view)
//
// Identifier: application/vnd.app.device_source+json; type=collection; view=default
type DeviceSourceCollection []*DeviceSource

// Validate validates the DeviceSourceCollection media type instance.
func (mt DeviceSourceCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// DeviceSourceCollection is the media type for an array of DeviceSource (public view)
//
// Identifier: application/vnd.app.device_source+json; type=collection; view=public
type DeviceSourcePublicCollection []*DeviceSourcePublic

// Validate validates the DeviceSourcePublicCollection media type instance.
func (mt DeviceSourcePublicCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// DeviceSources media type (default view)
//
// Identifier: application/vnd.app.device_sources+json; view=default
type DeviceSources struct {
	DeviceSources DeviceSourceCollection `form:"deviceSources" json:"deviceSources" xml:"deviceSources"`
}

// Validate validates the DeviceSources media type instance.
func (mt *DeviceSources) Validate() (err error) {
	if mt.DeviceSources == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "deviceSources"))
	}
	if err2 := mt.DeviceSources.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// Expedition media type (default view)
//
// Identifier: application/vnd.app.expedition+json; view=default
type Expedition struct {
	Description string `form:"description" json:"description" xml:"description"`
	ID          int    `form:"id" json:"id" xml:"id"`
	Name        string `form:"name" json:"name" xml:"name"`
	Slug        string `form:"slug" json:"slug" xml:"slug"`
}

// Validate validates the Expedition media type instance.
func (mt *Expedition) Validate() (err error) {

	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if mt.Slug == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "slug"))
	}
	if mt.Description == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "description"))
	}
	if ok := goa.ValidatePattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`, mt.Slug); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.slug`, mt.Slug, `^[[:alnum:]]+(-[[:alnum:]]+)*$`))
	}
	if utf8.RuneCountInString(mt.Slug) > 40 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.slug`, mt.Slug, utf8.RuneCountInString(mt.Slug), 40, false))
	}
	return
}

// Expedition media type (detailed view)
//
// Identifier: application/vnd.app.expedition+json; view=detailed
type ExpeditionDetailed struct {
	Centroid         []float64  `form:"centroid,omitempty" json:"centroid,omitempty" xml:"centroid,omitempty"`
	Description      string     `form:"description" json:"description" xml:"description"`
	EndTime          *time.Time `form:"endTime,omitempty" json:"endTime,omitempty" xml:"endTime,omitempty"`
	ID               int        `form:"id" json:"id" xml:"id"`
	LastFeatureID    *int       `form:"lastFeatureId,omitempty" json:"lastFeatureId,omitempty" xml:"lastFeatureId,omitempty"`
	Name             string     `form:"name" json:"name" xml:"name"`
	NumberOfFeatures *int       `form:"numberOfFeatures,omitempty" json:"numberOfFeatures,omitempty" xml:"numberOfFeatures,omitempty"`
	Radius           *float64   `form:"radius,omitempty" json:"radius,omitempty" xml:"radius,omitempty"`
	Slug             string     `form:"slug" json:"slug" xml:"slug"`
	StartTime        *time.Time `form:"startTime,omitempty" json:"startTime,omitempty" xml:"startTime,omitempty"`
}

// Validate validates the ExpeditionDetailed media type instance.
func (mt *ExpeditionDetailed) Validate() (err error) {

	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if mt.Slug == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "slug"))
	}
	if mt.Description == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "description"))
	}
	if ok := goa.ValidatePattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`, mt.Slug); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.slug`, mt.Slug, `^[[:alnum:]]+(-[[:alnum:]]+)*$`))
	}
	if utf8.RuneCountInString(mt.Slug) > 40 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.slug`, mt.Slug, utf8.RuneCountInString(mt.Slug), 40, false))
	}
	return
}

// ExpeditionCollection is the media type for an array of Expedition (default view)
//
// Identifier: application/vnd.app.expedition+json; type=collection; view=default
type ExpeditionCollection []*Expedition

// Validate validates the ExpeditionCollection media type instance.
func (mt ExpeditionCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ExpeditionCollection is the media type for an array of Expedition (detailed view)
//
// Identifier: application/vnd.app.expedition+json; type=collection; view=detailed
type ExpeditionDetailedCollection []*ExpeditionDetailed

// Validate validates the ExpeditionDetailedCollection media type instance.
func (mt ExpeditionDetailedCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// Expeditions media type (default view)
//
// Identifier: application/vnd.app.expeditions+json; view=default
type Expeditions struct {
	Expeditions ExpeditionCollection `form:"expeditions" json:"expeditions" xml:"expeditions"`
}

// Validate validates the Expeditions media type instance.
func (mt *Expeditions) Validate() (err error) {
	if mt.Expeditions == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "expeditions"))
	}
	if err2 := mt.Expeditions.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// GeoJSON media type (default view)
//
// Identifier: application/vnd.app.geojson+json; view=default
type GeoJSON struct {
	Features GeoJSONFeatureCollection `form:"features" json:"features" xml:"features"`
	Type     string                   `form:"type" json:"type" xml:"type"`
}

// Validate validates the GeoJSON media type instance.
func (mt *GeoJSON) Validate() (err error) {
	if mt.Type == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "type"))
	}
	if mt.Features == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "features"))
	}
	if err2 := mt.Features.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// GeoJSONFeature media type (default view)
//
// Identifier: application/vnd.app.geojson-feature+json; view=default
type GeoJSONFeature struct {
	Geometry   *GeoJSONGeometry       `form:"geometry" json:"geometry" xml:"geometry"`
	Properties map[string]interface{} `form:"properties" json:"properties" xml:"properties"`
	Type       string                 `form:"type" json:"type" xml:"type"`
}

// Validate validates the GeoJSONFeature media type instance.
func (mt *GeoJSONFeature) Validate() (err error) {
	if mt.Type == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "type"))
	}
	if mt.Geometry == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "geometry"))
	}
	if mt.Properties == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "properties"))
	}
	if mt.Geometry != nil {
		if err2 := mt.Geometry.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// GeoJSONFeatureCollection is the media type for an array of GeoJSONFeature (default view)
//
// Identifier: application/vnd.app.geojson-feature+json; type=collection; view=default
type GeoJSONFeatureCollection []*GeoJSONFeature

// Validate validates the GeoJSONFeatureCollection media type instance.
func (mt GeoJSONFeatureCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// GeoJSONGeometry media type (default view)
//
// Identifier: application/vnd.app.geojson-geometry+json; view=default
type GeoJSONGeometry struct {
	Coordinates []float64 `form:"coordinates" json:"coordinates" xml:"coordinates"`
	Type        string    `form:"type" json:"type" xml:"type"`
}

// Validate validates the GeoJSONGeometry media type instance.
func (mt *GeoJSONGeometry) Validate() (err error) {
	if mt.Type == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "type"))
	}
	if mt.Coordinates == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "coordinates"))
	}
	return
}

// ClusterSummary media type (default view)
//
// Identifier: application/vnd.app.geometry_cluster_summary+json; view=default
type ClusterSummary struct {
	Centroid         []float64 `form:"centroid" json:"centroid" xml:"centroid"`
	EndTime          time.Time `form:"endTime" json:"endTime" xml:"endTime"`
	ID               int       `form:"id" json:"id" xml:"id"`
	NumberOfFeatures int       `form:"numberOfFeatures" json:"numberOfFeatures" xml:"numberOfFeatures"`
	Radius           float64   `form:"radius" json:"radius" xml:"radius"`
	SourceID         int       `form:"sourceId" json:"sourceId" xml:"sourceId"`
	StartTime        time.Time `form:"startTime" json:"startTime" xml:"startTime"`
}

// Validate validates the ClusterSummary media type instance.
func (mt *ClusterSummary) Validate() (err error) {

	if mt.Centroid == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "centroid"))
	}

	return
}

// ClusterSummaryCollection is the media type for an array of ClusterSummary (default view)
//
// Identifier: application/vnd.app.geometry_cluster_summary+json; type=collection; view=default
type ClusterSummaryCollection []*ClusterSummary

// Validate validates the ClusterSummaryCollection media type instance.
func (mt ClusterSummaryCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// Location media type (default view)
//
// Identifier: application/vnd.app.location+json; view=default
type Location struct {
	Location string `form:"location" json:"location" xml:"location"`
}

// Validate validates the Location media type instance.
func (mt *Location) Validate() (err error) {
	if mt.Location == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "location"))
	}
	if err2 := goa.ValidateFormat(goa.FormatURI, mt.Location); err2 != nil {
		err = goa.MergeErrors(err, goa.InvalidFormatError(`response.location`, mt.Location, goa.FormatURI, err2))
	}
	return
}

// MapFeatures media type (default view)
//
// Identifier: application/vnd.app.map_features+json; view=default
type MapFeatures struct {
	GeoJSON    *PagedGeoJSON                    `form:"geoJSON" json:"geoJSON" xml:"geoJSON"`
	Geometries ClusterGeometrySummaryCollection `form:"geometries" json:"geometries" xml:"geometries"`
	Spatial    ClusterSummaryCollection         `form:"spatial" json:"spatial" xml:"spatial"`
	Temporal   ClusterSummaryCollection         `form:"temporal" json:"temporal" xml:"temporal"`
}

// Validate validates the MapFeatures media type instance.
func (mt *MapFeatures) Validate() (err error) {
	if mt.GeoJSON == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "geoJSON"))
	}
	if mt.Temporal == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "temporal"))
	}
	if mt.Spatial == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "spatial"))
	}
	if mt.Geometries == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "geometries"))
	}
	if mt.GeoJSON != nil {
		if err2 := mt.GeoJSON.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if err2 := mt.Geometries.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	if err2 := mt.Spatial.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	if err2 := mt.Temporal.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// TeamMember media type (default view)
//
// Identifier: application/vnd.app.member+json; view=default
type TeamMember struct {
	Role   string `form:"role" json:"role" xml:"role"`
	TeamID int    `form:"teamId" json:"teamId" xml:"teamId"`
	UserID int    `form:"userId" json:"userId" xml:"userId"`
}

// Validate validates the TeamMember media type instance.
func (mt *TeamMember) Validate() (err error) {

	if mt.Role == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "role"))
	}
	return
}

// TeamMemberCollection is the media type for an array of TeamMember (default view)
//
// Identifier: application/vnd.app.member+json; type=collection; view=default
type TeamMemberCollection []*TeamMember

// Validate validates the TeamMemberCollection media type instance.
func (mt TeamMemberCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// TeamMembers media type (default view)
//
// Identifier: application/vnd.app.members+json; view=default
type TeamMembers struct {
	Members TeamMemberCollection `form:"members" json:"members" xml:"members"`
}

// Validate validates the TeamMembers media type instance.
func (mt *TeamMembers) Validate() (err error) {
	if mt.Members == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "members"))
	}
	if err2 := mt.Members.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// PagedGeoJSON media type (default view)
//
// Identifier: application/vnd.app.paged_geojson+json; view=default
type PagedGeoJSON struct {
	Geo         *GeoJSON `form:"geo" json:"geo" xml:"geo"`
	HasMore     bool     `form:"hasMore" json:"hasMore" xml:"hasMore"`
	NextURL     string   `form:"nextUrl" json:"nextUrl" xml:"nextUrl"`
	PreviousURL *string  `form:"previousUrl,omitempty" json:"previousUrl,omitempty" xml:"previousUrl,omitempty"`
}

// Validate validates the PagedGeoJSON media type instance.
func (mt *PagedGeoJSON) Validate() (err error) {
	if mt.NextURL == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "nextUrl"))
	}
	if mt.Geo == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "geo"))
	}

	if mt.Geo != nil {
		if err2 := mt.Geo.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// Project media type (default view)
//
// Identifier: application/vnd.app.project+json; view=default
type Project struct {
	Description string `form:"description" json:"description" xml:"description"`
	ID          int    `form:"id" json:"id" xml:"id"`
	Name        string `form:"name" json:"name" xml:"name"`
	Slug        string `form:"slug" json:"slug" xml:"slug"`
}

// Validate validates the Project media type instance.
func (mt *Project) Validate() (err error) {

	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if mt.Slug == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "slug"))
	}
	if mt.Description == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "description"))
	}
	if ok := goa.ValidatePattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`, mt.Slug); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.slug`, mt.Slug, `^[[:alnum:]]+(-[[:alnum:]]+)*$`))
	}
	if utf8.RuneCountInString(mt.Slug) > 40 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.slug`, mt.Slug, utf8.RuneCountInString(mt.Slug), 40, false))
	}
	return
}

// ProjectCollection is the media type for an array of Project (default view)
//
// Identifier: application/vnd.app.project+json; type=collection; view=default
type ProjectCollection []*Project

// Validate validates the ProjectCollection media type instance.
func (mt ProjectCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// Projects media type (default view)
//
// Identifier: application/vnd.app.projects+json; view=default
type Projects struct {
	Projects ProjectCollection `form:"projects" json:"projects" xml:"projects"`
}

// Validate validates the Projects media type instance.
func (mt *Projects) Validate() (err error) {
	if mt.Projects == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "projects"))
	}
	if err2 := mt.Projects.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// QueryData media type (default view)
//
// Identifier: application/vnd.app.queried+json; view=default
type QueryData struct {
	Series SeriesDataCollection `form:"series" json:"series" xml:"series"`
}

// Validate validates the QueryData media type instance.
func (mt *QueryData) Validate() (err error) {
	if mt.Series == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "series"))
	}
	if err2 := mt.Series.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// ReadingSummary media type (default view)
//
// Identifier: application/vnd.app.reading_summary+json; view=default
type ReadingSummary struct {
	Name string `form:"name" json:"name" xml:"name"`
}

// Validate validates the ReadingSummary media type instance.
func (mt *ReadingSummary) Validate() (err error) {
	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	return
}

// ReadingSummaryCollection is the media type for an array of ReadingSummary (default view)
//
// Identifier: application/vnd.app.reading_summary+json; type=collection; view=default
type ReadingSummaryCollection []*ReadingSummary

// Validate validates the ReadingSummaryCollection media type instance.
func (mt ReadingSummaryCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// SeriesData media type (default view)
//
// Identifier: application/vnd.app.series+json; view=default
type SeriesData struct {
	Name string        `form:"name" json:"name" xml:"name"`
	Rows []interface{} `form:"rows" json:"rows" xml:"rows"`
}

// Validate validates the SeriesData media type instance.
func (mt *SeriesData) Validate() (err error) {
	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if mt.Rows == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "rows"))
	}
	return
}

// SeriesDataCollection is the media type for an array of SeriesData (default view)
//
// Identifier: application/vnd.app.series+json; type=collection; view=default
type SeriesDataCollection []*SeriesData

// Validate validates the SeriesDataCollection media type instance.
func (mt SeriesDataCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// Source media type (default view)
//
// Identifier: application/vnd.app.source+json; view=default
type Source struct {
	Active       *bool  `form:"active,omitempty" json:"active,omitempty" xml:"active,omitempty"`
	ExpeditionID int    `form:"expeditionId" json:"expeditionId" xml:"expeditionId"`
	ID           int    `form:"id" json:"id" xml:"id"`
	Name         string `form:"name" json:"name" xml:"name"`
	TeamID       *int   `form:"teamId,omitempty" json:"teamId,omitempty" xml:"teamId,omitempty"`
	UserID       *int   `form:"userId,omitempty" json:"userId,omitempty" xml:"userId,omitempty"`
}

// Validate validates the Source media type instance.
func (mt *Source) Validate() (err error) {

	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	return
}

// SourceSummary media type (default view)
//
// Identifier: application/vnd.app.source_summary+json; view=default
type SourceSummary struct {
	ID       int                      `form:"id" json:"id" xml:"id"`
	Name     string                   `form:"name" json:"name" xml:"name"`
	Readings ReadingSummaryCollection `form:"readings,omitempty" json:"readings,omitempty" xml:"readings,omitempty"`
	Spatial  ClusterSummaryCollection `form:"spatial" json:"spatial" xml:"spatial"`
	Temporal ClusterSummaryCollection `form:"temporal" json:"temporal" xml:"temporal"`
}

// Validate validates the SourceSummary media type instance.
func (mt *SourceSummary) Validate() (err error) {

	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if mt.Temporal == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "temporal"))
	}
	if mt.Spatial == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "spatial"))
	}
	if err2 := mt.Readings.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	if err2 := mt.Spatial.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	if err2 := mt.Temporal.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// SourceToken media type (default view)
//
// Identifier: application/vnd.app.source_token+json; view=default
type SourceToken struct {
	ExpeditionID int    `form:"expeditionId" json:"expeditionId" xml:"expeditionId"`
	ID           int    `form:"id" json:"id" xml:"id"`
	Token        string `form:"token" json:"token" xml:"token"`
}

// Validate validates the SourceToken media type instance.
func (mt *SourceToken) Validate() (err error) {

	if mt.Token == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "token"))
	}

	return
}

// SourceTokenCollection is the media type for an array of SourceToken (default view)
//
// Identifier: application/vnd.app.source_token+json; type=collection; view=default
type SourceTokenCollection []*SourceToken

// Validate validates the SourceTokenCollection media type instance.
func (mt SourceTokenCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// SourceTokens media type (default view)
//
// Identifier: application/vnd.app.source_tokens+json; view=default
type SourceTokens struct {
	SourceTokens SourceTokenCollection `form:"sourceTokens" json:"sourceTokens" xml:"sourceTokens"`
}

// Validate validates the SourceTokens media type instance.
func (mt *SourceTokens) Validate() (err error) {
	if mt.SourceTokens == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "sourceTokens"))
	}
	if err2 := mt.SourceTokens.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// Sources media type (default view)
//
// Identifier: application/vnd.app.sources+json; view=default
type Sources struct {
	DeviceSources         DeviceSourceCollection         `form:"deviceSources,omitempty" json:"deviceSources,omitempty" xml:"deviceSources,omitempty"`
	TwitterAccountSources TwitterAccountSourceCollection `form:"twitterAccountSources,omitempty" json:"twitterAccountSources,omitempty" xml:"twitterAccountSources,omitempty"`
}

// Validate validates the Sources media type instance.
func (mt *Sources) Validate() (err error) {
	if err2 := mt.DeviceSources.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	if err2 := mt.TwitterAccountSources.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// Team media type (default view)
//
// Identifier: application/vnd.app.team+json; view=default
type Team struct {
	Description string `form:"description" json:"description" xml:"description"`
	ID          int    `form:"id" json:"id" xml:"id"`
	Name        string `form:"name" json:"name" xml:"name"`
	Slug        string `form:"slug" json:"slug" xml:"slug"`
}

// Validate validates the Team media type instance.
func (mt *Team) Validate() (err error) {

	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if mt.Slug == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "slug"))
	}
	if mt.Description == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "description"))
	}
	if ok := goa.ValidatePattern(`\S`, mt.Name); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.name`, mt.Name, `\S`))
	}
	if utf8.RuneCountInString(mt.Name) > 256 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.name`, mt.Name, utf8.RuneCountInString(mt.Name), 256, false))
	}
	if ok := goa.ValidatePattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`, mt.Slug); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.slug`, mt.Slug, `^[[:alnum:]]+(-[[:alnum:]]+)*$`))
	}
	if utf8.RuneCountInString(mt.Slug) > 40 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.slug`, mt.Slug, utf8.RuneCountInString(mt.Slug), 40, false))
	}
	return
}

// TeamCollection is the media type for an array of Team (default view)
//
// Identifier: application/vnd.app.team+json; type=collection; view=default
type TeamCollection []*Team

// Validate validates the TeamCollection media type instance.
func (mt TeamCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// Teams media type (default view)
//
// Identifier: application/vnd.app.teams+json; view=default
type Teams struct {
	Teams TeamCollection `form:"teams" json:"teams" xml:"teams"`
}

// Validate validates the Teams media type instance.
func (mt *Teams) Validate() (err error) {
	if mt.Teams == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "teams"))
	}
	if err2 := mt.Teams.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// TwitterAccountSources media type (default view)
//
// Identifier: application/vnd.app.twitter_account_intputs+json; view=default
type TwitterAccountSources struct {
	TwitterAccountSources TwitterAccountSourceCollection `form:"twitterAccountSources" json:"twitterAccountSources" xml:"twitterAccountSources"`
}

// Validate validates the TwitterAccountSources media type instance.
func (mt *TwitterAccountSources) Validate() (err error) {
	if mt.TwitterAccountSources == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "twitterAccountSources"))
	}
	if err2 := mt.TwitterAccountSources.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// TwitterAccountSource media type (default view)
//
// Identifier: application/vnd.app.twitter_account_source+json; view=default
type TwitterAccountSource struct {
	ExpeditionID     int    `form:"expeditionId" json:"expeditionId" xml:"expeditionId"`
	ID               int    `form:"id" json:"id" xml:"id"`
	Name             string `form:"name" json:"name" xml:"name"`
	ScreenName       string `form:"screenName" json:"screenName" xml:"screenName"`
	TeamID           *int   `form:"teamId,omitempty" json:"teamId,omitempty" xml:"teamId,omitempty"`
	TwitterAccountID int    `form:"twitterAccountId" json:"twitterAccountId" xml:"twitterAccountId"`
	UserID           *int   `form:"userId,omitempty" json:"userId,omitempty" xml:"userId,omitempty"`
}

// Validate validates the TwitterAccountSource media type instance.
func (mt *TwitterAccountSource) Validate() (err error) {

	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}

	if mt.ScreenName == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "screenName"))
	}
	return
}

// TwitterAccountSourceCollection is the media type for an array of TwitterAccountSource (default view)
//
// Identifier: application/vnd.app.twitter_account_source+json; type=collection; view=default
type TwitterAccountSourceCollection []*TwitterAccountSource

// Validate validates the TwitterAccountSourceCollection media type instance.
func (mt TwitterAccountSourceCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// User media type (default view)
//
// Identifier: application/vnd.app.user+json; view=default
type User struct {
	Bio   string `form:"bio" json:"bio" xml:"bio"`
	Email string `form:"email" json:"email" xml:"email"`
	ID    int    `form:"id" json:"id" xml:"id"`
	Name  string `form:"name" json:"name" xml:"name"`
	// Username
	Username string `form:"username" json:"username" xml:"username"`
}

// Validate validates the User media type instance.
func (mt *User) Validate() (err error) {

	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if mt.Username == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "username"))
	}
	if mt.Email == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "email"))
	}
	if mt.Bio == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "bio"))
	}
	if err2 := goa.ValidateFormat(goa.FormatEmail, mt.Email); err2 != nil {
		err = goa.MergeErrors(err, goa.InvalidFormatError(`response.email`, mt.Email, goa.FormatEmail, err2))
	}
	if ok := goa.ValidatePattern(`\S`, mt.Name); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.name`, mt.Name, `\S`))
	}
	if utf8.RuneCountInString(mt.Name) > 256 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.name`, mt.Name, utf8.RuneCountInString(mt.Name), 256, false))
	}
	if ok := goa.ValidatePattern(`^[\dA-Za-z]+(?:-[\dA-Za-z]+)*$`, mt.Username); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.username`, mt.Username, `^[\dA-Za-z]+(?:-[\dA-Za-z]+)*$`))
	}
	if utf8.RuneCountInString(mt.Username) > 40 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.username`, mt.Username, utf8.RuneCountInString(mt.Username), 40, false))
	}
	return
}

// UserCollection is the media type for an array of User (default view)
//
// Identifier: application/vnd.app.user+json; type=collection; view=default
type UserCollection []*User

// Validate validates the UserCollection media type instance.
func (mt UserCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// Users media type (default view)
//
// Identifier: application/vnd.app.users+json; view=default
type Users struct {
	Users UserCollection `form:"users" json:"users" xml:"users"`
}

// Validate validates the Users media type instance.
func (mt *Users) Validate() (err error) {
	if mt.Users == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "users"))
	}
	if err2 := mt.Users.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}
