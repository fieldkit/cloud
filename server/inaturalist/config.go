package inaturalist

import (
	"github.com/Conservify/gonaturalist"

	"github.com/fieldkit/cloud/server/jobs"
)

type INaturalistConfig struct {
	ApplicationId string
	Secret        string
	AccessToken   string
	RedirectUrl   string
	RootUrl       string
	SiteID        int64
	Verbose       bool
	Mutable       bool
}

var (
	INaturalistObservationsQueue = &jobs.QueueDef{
		Name: "inaturalist_observations",
	}
)

type ClientAndConfig struct {
	Client *gonaturalist.Client
	Config *INaturalistConfig
}

type Clients struct {
	bySiteID map[int64]*ClientAndConfig
}

func NewClients() (c *Clients, err error) {
	bySiteID := make(map[int64]*ClientAndConfig)

	for _, config := range AllNaturalistConfigs {
		var authenticator = gonaturalist.NewAuthenticatorAtCustomRoot(config.ApplicationId, config.Secret, config.RedirectUrl, config.RootUrl)
		siteClient := authenticator.NewClientWithAccessToken(config.AccessToken, &gonaturalist.NoopCallbacks{})
		bySiteID[config.SiteID] = &ClientAndConfig{
			Client: siteClient,
			Config: &config,
		}
	}

	c = &Clients{
		bySiteID: bySiteID,
	}

	return
}

func (c *Clients) GetClientAndConfig(siteId int64) *ClientAndConfig {
	return c.bySiteID[siteId]
}
