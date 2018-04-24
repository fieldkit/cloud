package inaturalist

import (
	"context"
	"fmt"
	"strings"

	"github.com/Conservify/gonaturalist"
)

type INaturalistCommentor struct {
	Correlator *INaturalistCorrelator
	Clients    *Clients
}

func NewINaturalistCommentor(correlator *INaturalistCorrelator, clients *Clients) (*INaturalistCommentor, error) {
	return &INaturalistCommentor{
		Correlator: correlator,
		Clients:    clients,
	}, nil
}

func (nc *INaturalistCommentor) Handle(ctx context.Context, co *CachedObservation) error {
	log := Logger(ctx).Sugar()

	clientAndConfig := nc.Clients.GetClientAndConfig(co.SiteID)
	if clientAndConfig == nil {
		return fmt.Errorf("No Naturalist client for SiteID: %d", co.SiteID)
	}

	client := clientAndConfig.Client

	correlation, err := nc.Correlator.Correlate(ctx, co)
	if err != nil {
		return err
	}

	if correlation == nil {
		return nil
	}

	body := nc.CreateSummary(co, correlation)

	user, err := client.GetCurrentUser()
	if err != nil {
		return err
	}

	log.Infow("GetObservationComments", "observationId", co.ID)
	full, err := client.GetObservationComments(co.ID)
	if err != nil {
		return err
	}

	updated := false

	for _, c := range full {
		if c.UserId == user.Id {
			updated = true

			if c.Body == body {
				continue
			}

			log.Infow("Updating comment", "observationId", co.ID, "siteId", co.SiteID, "commentId", c.Id, "mutable", clientAndConfig.Config.Mutable)

			if clientAndConfig.Config.Mutable {
				err = client.UpdateCommentBody(c.Id, body)
				if err != nil {
					return err
				}
			}
		}
	}

	if !updated {
		log.Infow("Adding comment", "observationId", co.ID, "siteId", co.SiteID, "mutable", clientAndConfig.Config.Mutable)

		addComment := gonaturalist.AddCommentOpt{
			ParentType: gonaturalist.Observation,
			ParentId:   co.ID,
			Body:       body,
		}
		if clientAndConfig.Config.Mutable {
			err = client.AddComment(&addComment)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (nc *INaturalistCommentor) CreateSummary(co *CachedObservation, c *Correlation) string {
	sb := strings.Builder{}

	sb.WriteString("FieldKit Readings\n")
	sb.WriteString("\n")
	for key, records := range c.RecordsByTier {
		sb.WriteString(fmt.Sprintf("Tier < %f\n", key.Threshold))
		for _, record := range records {
			distance := record.Location.Distance(co.Location)
			coords := record.Location.Coordinates()
			time := record.Timestamp.Format("Mon Jan 2 15:04:05 -0700 MST 2006")

			sb.WriteString(fmt.Sprintf("\tLocation: %v, %v (distance = %f)\n", coords[1], coords[0], distance))
			sb.WriteString(fmt.Sprintf("\tTime: %v\n", time))

			fields, _ := record.GetRawFields()
			if fields != nil {
				values := make([]string, 0)
				for key, value := range fields {
					values = append(values, fmt.Sprintf("%s = %v", key, value))
				}
				sb.WriteString(fmt.Sprintf("\tData: %s\n", strings.Join(values, " ")))
			}
		}
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("https://www.fkdev.org/summaries/inaturalist/%d/%d\n", co.SiteID, co.ID))

	return sb.String()
}
