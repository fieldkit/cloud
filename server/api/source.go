package api

import (
	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"
)

func SourceType(source *data.Source) *app.Source {
	sourceType := &app.Source{
		ID:           int(source.ID),
		ExpeditionID: int(source.ExpeditionID),
		Name:         source.Name,
		Active:       &source.Active,
	}

	if source.TeamID != nil {
		teamID := int(*source.TeamID)
		sourceType.TeamID = &teamID
	}

	if source.UserID != nil {
		userID := int(*source.UserID)
		sourceType.UserID = &userID
	}

	return sourceType
}

func ClusterSummaryType(s *backend.GeometryClusterSummary) *app.ClusterSummary {
	return &app.ClusterSummary{
		ID:               s.ClusterID,
		SourceID:         s.SourceID,
		NumberOfFeatures: s.NumberOfFeatures,
		StartTime:        s.StartTime,
		EndTime:          s.EndTime,
		Envelope:         s.Envelope.Coordinates(),
		Centroid:         s.Centroid.Coordinates(),
		Radius:           s.Radius,
	}
}

func ClusterSummariesType(s []*backend.GeometryClusterSummary) []*app.ClusterSummary {
	summaries := make([]*app.ClusterSummary, len(s))
	for i, summary := range s {
		summaries[i] = ClusterSummaryType(summary)
	}
	return summaries
}

func ReadingSummariesType(s []*backend.ReadingSummary) []*app.ReadingSummary {
	summaries := make([]*app.ReadingSummary, len(s))
	for i, summary := range s {
		summaries[i] = &app.ReadingSummary{
			Name: summary.Name,
		}
	}
	return summaries
}

func SourceSummaryType(source *data.DeviceSource, spatial, temporal []*backend.GeometryClusterSummary, readings []*backend.ReadingSummary) *app.SourceSummary {
	return &app.SourceSummary{
		ID:       int(source.ID),
		Name:     source.Name,
		Readings: ReadingSummariesType(readings),
		Temporal: ClusterSummariesType(temporal),
		Spatial:  ClusterSummariesType(spatial),
	}
}

func ClusterGeometryType(source *backend.GeometryClusterSummary) *app.ClusterGeometrySummary {
	return &app.ClusterGeometrySummary{
		ID:       source.ClusterID,
		SourceID: source.SourceID,
		Geometry: source.Geometry.Coordinates(),
	}
}

func ClusterGeometriesType(source []*backend.GeometryClusterSummary) []*app.ClusterGeometrySummary {
	summaries := make([]*app.ClusterGeometrySummary, len(source))
	for i, summary := range source {
		summaries[i] = ClusterGeometryType(summary)

	}
	return summaries
}

type SourceControllerOptions struct {
	Backend *backend.Backend
}

type SourceController struct {
	*goa.Controller
	options SourceControllerOptions
}

func NewSourceController(service *goa.Service, options SourceControllerOptions) *SourceController {
	return &SourceController{
		Controller: service.NewController("SourceController"),
		options:    options,
	}
}

func (c *SourceController) Update(ctx *app.UpdateSourceContext) error {
	source, err := c.options.Backend.Source(ctx, int32(ctx.SourceID))
	if err != nil {
		return err
	}

	source.Name = ctx.Payload.Name

	if ctx.Payload.TeamID != nil {
		teamID := int32(*ctx.Payload.TeamID)
		source.TeamID = &teamID
	} else {
		source.TeamID = nil
	}

	if ctx.Payload.UserID != nil {
		userID := int32(*ctx.Payload.UserID)
		source.UserID = &userID
	} else {
		source.UserID = nil
	}

	if err := c.options.Backend.UpdateSource(ctx, source); err != nil {
		return err
	}

	return ctx.OK(SourceType(source))
}

func (c *SourceController) List(ctx *app.ListSourceContext) error {
	twitterAccountSources, err := c.options.Backend.ListTwitterAccountSources(ctx, ctx.Project, ctx.Expedition)
	if err != nil {
		return err
	}

	deviceSources, err := c.options.Backend.ListDeviceSources(ctx, ctx.Project, ctx.Expedition)
	if err != nil {
		return err
	}

	return ctx.OK(&app.Sources{
		TwitterAccountSources: TwitterAccountSourcesType(twitterAccountSources).TwitterAccountSources,
		DeviceSources:         DeviceSourcesType(deviceSources).DeviceSources,
	})
}

func (c *SourceController) ListID(ctx *app.ListIDSourceContext) error {
	deviceSource, err := c.options.Backend.GetDeviceSourceByID(ctx, int32(ctx.SourceID))
	if err != nil {
		return err
	}
	summary, err := c.options.Backend.FeatureSummaryBySourceID(ctx, ctx.SourceID)
	if err != nil {
		return err
	}
	return ctx.OKPublic(DeviceSourcePublicType(deviceSource, summary))
}

func (c *SourceController) SummaryByID(ctx *app.SummaryByIDSourceContext) error {
	deviceSource, err := c.options.Backend.GetDeviceSourceByID(ctx, int32(ctx.SourceID))
	if err != nil {
		return err
	}

	spatial, err := c.options.Backend.SpatialClustersBySourceID(ctx, ctx.SourceID)
	if err != nil {
		return err
	}

	temporal, err := c.options.Backend.TemporalClustersBySourceID(ctx, ctx.SourceID)
	if err != nil {
		return err
	}

	readings, err := c.options.Backend.ReadingsBySourceID(ctx, ctx.SourceID)
	if err != nil {
		return err
	}

	return ctx.OK(SourceSummaryType(deviceSource, spatial, temporal, readings))
}

func (c *SourceController) TemporalClusterGeometryByID(ctx *app.TemporalClusterGeometryByIDSourceContext) error {
	temporal, err := c.options.Backend.TemporalClustersBySourceID(ctx, ctx.SourceID)
	if err != nil {
		return err
	}

	for _, cluster := range temporal {
		if cluster.ClusterID == ctx.ClusterID {
			return ctx.OK(ClusterGeometryType(cluster))
		}
	}

	return ctx.BadRequest()
}

func (c *SourceController) ListExpeditionID(ctx *app.ListExpeditionIDSourceContext) error {
	twitterAccountSources, err := c.options.Backend.ListTwitterAccountSourcesByExpeditionID(ctx, int32(ctx.ExpeditionID))
	if err != nil {
		return err
	}

	deviceSources, err := c.options.Backend.ListDeviceSourcesByExpeditionID(ctx, int32(ctx.ExpeditionID))
	if err != nil {
		return err
	}

	return ctx.OK(&app.Sources{
		TwitterAccountSources: TwitterAccountSourcesType(twitterAccountSources).TwitterAccountSources,
		DeviceSources:         DeviceSourcesType(deviceSources).DeviceSources,
	})
}
