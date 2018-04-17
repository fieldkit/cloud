package api

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/goadesign/goa"

	"github.com/Conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/email"
	"github.com/fieldkit/cloud/server/naturalist"
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
	TwoDays       = 2 * 24 * time.Hour
)

type Notifier struct {
	Database *sqlxcache.DB
	Backend  *backend.Backend
	Emailer  email.Emailer
}

func NewNotifier(backend *backend.Backend, database *sqlxcache.DB, emailer email.Emailer) *Notifier {
	return &Notifier{
		Database: database,
		Backend:  backend,
		Emailer:  emailer,
	}
}

func (n *Notifier) SendOffline(ctx context.Context, sourceId int32, age time.Duration, notification *backend.NotificationStatus) error {
	source, err := n.Backend.GetSourceByID(ctx, int32(sourceId))
	if err != nil {
		return err
	}

	if err = n.Emailer.SendSourceOfflineWarning(source, age); err != nil {
		return err
	}

	if err := n.MarkNotified(ctx, notification); err != nil {
		return err
	}

	return nil
}

func (n *Notifier) SendOnline(ctx context.Context, sourceId int32, age time.Duration, notification *backend.NotificationStatus) error {
	source, err := n.Backend.GetSourceByID(ctx, int32(sourceId))
	if err != nil {
		return err
	}

	if err := n.Emailer.SendSourceOnlineWarning(source, age); err != nil {
		return err
	}

	if err := n.ClearNotifications(ctx, notification); err != nil {
		return err
	}

	return nil
}

// TODO: Should narrow this query down eventually.
func (n *Notifier) Check(ctx context.Context) error {
	summaries := []*backend.FeatureSummary{}
	if err := n.Database.SelectContext(ctx, &summaries, `SELECT c.source_id, c.end_time FROM fieldkit.sources_summaries c ORDER BY c.end_time`); err != nil {
		return err
	}

	notifications := []*backend.NotificationStatus{}
	if err := n.Database.SelectContext(ctx, &notifications, `SELECT * FROM fieldkit.notification n`); err != nil {
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
		wasNotified := !notification.UpdatedAt.IsZero()
		age := now.Sub(summary.EndTime)
		lastNotification := now.Sub(notification.UpdatedAt)

		prefix := fmt.Sprintf("SourceId: %d EndTime: %v: Age: %v LastNotification: %v", summary.SourceID, summary.EndTime, age, lastNotification)

		if age < ThirtyMinutes {
			if wasNotified {
				log.Printf("%s: Online!", prefix)
				if err := n.SendOnline(ctx, int32(summary.SourceID), age, notification); err != nil {
					return err
				}
			} else {
				log.Printf("%s: Have readings", prefix)
			}
			continue
		}

		if age > TwoDays {
			log.Printf("%s: Too old", prefix)
			continue
		}

		for _, interval := range []time.Duration{SixHours, TwoHours, ThirtyMinutes} {
			if age > interval {
				if lastNotification < interval {
					log.Printf("%s: Already notified (%v)", prefix, interval)
				} else {
					log.Printf("%s: Offline! (%v)", prefix, interval)
					if err := n.SendOffline(ctx, int32(summary.SourceID), age, notification); err != nil {
						return err
					}
				}
				break
			}
		}
	}

	return nil
}

func (n *Notifier) ClearNotifications(ctx context.Context, notification *backend.NotificationStatus) error {
	_, err := n.Database.NamedExecContext(ctx, `DELETE FROM fieldkit.notification WHERE source_id = :source_id`, notification)
	if err != nil {
		return err
	}
	return nil
}

func (n *Notifier) MarkNotified(ctx context.Context, notification *backend.NotificationStatus) error {
	_, err := n.Database.NamedExecContext(ctx, `INSERT INTO fieldkit.notification (source_id, updated_at) VALUES (:source_id, NOW()) ON CONFLICT (source_id) DO UPDATE SET updated_at = excluded.updated_at`, notification)
	if err != nil {
		return err
	}
	return nil
}

func (c *TasksController) Check(ctx *app.CheckTasksContext) error {
	notifier := NewNotifier(c.options.Backend, c.options.Database, c.options.Emailer)
	if err := notifier.Check(ctx); err != nil {
		return err
	}
	return ctx.OK([]byte("Ok"))
}

func (c *TasksController) Five(ctx *app.FiveTasksContext) error {
	go naturalist.RefreshNaturalistObservations(context.Background(), c.options.Database, c.options.Backend)

	return ctx.OK([]byte("Ok"))
}
