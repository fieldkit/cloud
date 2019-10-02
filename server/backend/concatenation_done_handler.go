package backend

import (
	"context"

	"github.com/fieldkit/cloud/server/messages"
)

type ConcatenationDoneHandler struct {
	Backend *Backend
}

func (h *ConcatenationDoneHandler) Handle(ctx context.Context, m *messages.ConcatenationDone) error {
	log := Logger(ctx).Sugar()

	log.Infow("Concatenation done", "device_id", m.DeviceID)

	deviceSource, err := h.Backend.GetDeviceSourceByKey(ctx, m.DeviceID)
	if err != nil {
		log.Errorw("Error finding DeviceByKey: %v", err)
		return err
	}

	generator := NewPregenerator(h.Backend)
	_, err = generator.Pregenerate(ctx, int64(deviceSource.ID))
	if err != nil {
		return err
	}

	log.Infow("Concatenation done handled", "device_id", m.DeviceID)

	return nil
}
