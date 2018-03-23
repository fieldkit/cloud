package formatting

import (
	"strings"
	"time"

	"github.com/fieldkit/cloud/server/backend/ingestion"
)

const (
	ParticleProviderName          = "PARTICLE"
	ParticleFormCoreId            = "coreid"
	ParticleFormData              = "data"
	ParticleFormPublishedAt       = "published_at"
	ParticleFormPublishedAtLayout = "2006-01-02T15:04:05Z"
)

type ParticleMessageProvider struct {
	MessageProviderBase
}

func (i *ParticleMessageProvider) CanFormatMessage(raw *RawMessage) bool {
	return raw.Form.Get(ParticleFormCoreId) != ""
}

func (i *ParticleMessageProvider) FormatMessage(raw *RawMessage) (fm *ingestion.FormattedMessage, err error) {
	coreId := strings.TrimSpace(raw.Form.Get(ParticleFormCoreId))
	trimmed := strings.TrimSpace(raw.Form.Get(ParticleFormData))
	fields := strings.Split(trimmed, ",")

	publishedAt, err := time.Parse(ParticleFormPublishedAtLayout, strings.TrimSpace(raw.Form.Get(ParticleFormPublishedAt)))
	if err != nil {
		return nil, err
	}

	fm = &ingestion.FormattedMessage{
		MessageId:   ingestion.MessageId(raw.RequestId),
		SchemaId:    ingestion.NewSchemaId(ingestion.NewProviderDeviceId(ParticleProviderName, coreId), ""),
		Time:        &publishedAt,
		ArrayValues: fields,
	}

	return
}
