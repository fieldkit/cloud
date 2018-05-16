package api

import (
	"context"
	"time"

	"github.com/Conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/email"
)

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
	log := Logger(ctx).Sugar()

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

		details := log.With("source_id", summary.SourceID, "end_time", summary.EndTime, "age", age, "last_notification", lastNotification)

		if age < ThirtyMinutes {
			if wasNotified {
				details.Infof("Online!")
				if err := n.SendOnline(ctx, int32(summary.SourceID), age, notification); err != nil {
					return err
				}
			} else {
				details.Infof("Have readings")
			}
			continue
		}

		if age > TwoDays {
			details.Infof("Too old")
			continue
		}

		for _, interval := range []time.Duration{SixHours, TwoHours, ThirtyMinutes} {
			if age > interval {
				if lastNotification < interval {
					details.Infof("Already notified (%v)", interval)
				} else {
					details.Infof("Offline! (%v)", interval)
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
