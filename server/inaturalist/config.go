package inaturalist

import (
	"github.com/fieldkit/cloud/server/jobs"
)

type INaturalistConfig struct {
	ApplicationId string
	Secret        string
	AccessToken   string
	RedirectUrl   string
	RootUrl       string
}

var (
	INaturalistObservationsQueue = &jobs.QueueDef{
		Name: "inaturalist_observations",
	}
)
