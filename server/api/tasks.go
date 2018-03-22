package api

import (
	"log"
	"time"

	"github.com/conservify/sqlxcache"
	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/email"
)

type TasksControllerOptions struct {
	Database *sqlxcache.DB
	Backend  *backend.Backend
	Emailer  email.Emailer
}

type TasksController struct {
	*goa.Controller
	options TasksControllerOptions
}

func NewTasksController(service *goa.Service, options TasksControllerOptions) *TasksController {
	return &TasksController{
		Controller: service.NewController("TasksController"),
		options:    options,
	}
}

var (
	ThirtyMinutes = 30 * time.Minute
	TwoHours      = 2 * time.Hour
	SixHours      = 6 * time.Hour
	OneDay        = 24 * time.Hour
	TwoDays       = 2 * 24 * time.Hour
)

func (c *TasksController) Check(ctx *app.CheckTasksContext) error {
	summaries := []*backend.FeatureSummary{}
	if err := c.options.Database.SelectContext(ctx, &summaries, `SELECT c.source_id, c.end_time FROM fieldkit.sources_summaries c`); err != nil {
		return err
	}

	notifications := []*backend.NotificationStatus{}
	if err := c.options.Database.SelectContext(ctx, &notifications, `SELECT * FROM fieldkit.notification n`); err != nil {
		return err
	}

	notificationsBySource := make(map[int]*backend.NotificationStatus)
	for _, s := range summaries {
		notificationsBySource[s.SourceID] = &backend.NotificationStatus{
			SourceID:  s.SourceID,
			UpdatedAt: time.Time{},
		}
	}
	for _, n := range notifications {
		notificationsBySource[n.SourceID] = n
	}

	now := time.Now()
	for _, summary := range summaries {
		notification := notificationsBySource[summary.SourceID]
		age := now.Sub(summary.EndTime)
		lastNotification := now.Sub(notification.UpdatedAt)

		if age < ThirtyMinutes {
			log.Printf("Have readings: sourceId=%v endTime=%v age=%v", summary.SourceID, summary.EndTime, age)
			_, err := c.options.Database.NamedExecContext(ctx, `DELETE FROM fieldkit.notification WHERE source_id = :source_id`, notification)
			if err != nil {
				return err
			}
			continue
		}

		if age > TwoDays {
			log.Printf("Two old to notify: sourceId=%v endTime=%v age=%v", summary.SourceID, summary.EndTime, age)
			continue
		}

		for _, interval := range []time.Duration{SixHours, TwoHours, ThirtyMinutes} {
			if age > interval {
				if lastNotification < interval {
					log.Printf("Already notified (%v): sourceId=%v endTime=%v age=%v lastNotification=%v", interval, summary.SourceID, summary.EndTime, age, lastNotification)
				} else {
					source, err := c.options.Backend.GetSourceByID(ctx, int32(summary.SourceID))
					if err != nil {
						return err
					}

					log.Printf("Notifying! (%v): sourceId=%v endTime=%v age=%v lastNotification=%v source=%v", interval, summary.SourceID, summary.EndTime, age, lastNotification, source)

					err = c.options.Emailer.SendSourceSilenceWarning(source, age)
					if err != nil {
						return err
					}

					_, err = c.options.Database.NamedExecContext(ctx,
						`INSERT INTO fieldkit.notification (source_id, updated_at) VALUES (:source_id, NOW()) ON CONFLICT (source_id) DO UPDATE SET updated_at = :updated_at`,
						notification)
					if err != nil {
						return err
					}
				}
				break
			}
		}
	}

	return ctx.OK([]byte("Ok"))
}
