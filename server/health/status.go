package health

import (
	"context"
	"encoding/json"
	_ "fmt"
	"net/http"
	"os"

	"github.com/fieldkit/cloud/server/logging"
)

type StatusResponse struct {
	Name string `json:"name"`
}

func MakeStatusResponse(ctx context.Context) (sr *StatusResponse, err error) {
	name, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	sr = &StatusResponse{
		Name: name,
	}

	return
}

func StatusHandler(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		log := logging.Logger(ctx).Named("status").Sugar()

		log.Infow("status", "headers", req.Header)

		sr, err := MakeStatusResponse(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(sr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(b)
	})
}
