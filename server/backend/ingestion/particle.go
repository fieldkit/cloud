package ingestion

import (
	"strings"
	"time"
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

func (i *ParticleMessageProvider) CanProcessMessage(raw *RawMessage) bool {
	return raw.Form.Get(ParticleFormCoreId) != ""
}

func (i *ParticleMessageProvider) ProcessMessage(raw *RawMessage) (pm *ProcessedMessage, err error) {
	coreId := strings.TrimSpace(raw.Form.Get(ParticleFormCoreId))
	trimmed := strings.TrimSpace(raw.Form.Get(ParticleFormData))
	fields := strings.Split(trimmed, ",")

	publishedAt, err := time.Parse(ParticleFormPublishedAtLayout, strings.TrimSpace(raw.Form.Get(ParticleFormPublishedAt)))
	if err != nil {
		return nil, err
	}

	pm = &ProcessedMessage{
		MessageId:   MessageId(raw.RequestId),
		SchemaId:    NewSchemaId(ParticleProviderName, coreId, ""),
		Time:        &publishedAt,
		ArrayValues: fields,
	}

	return
}
