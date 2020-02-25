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
	ServerName string `json:"server_name,omitempty"`
	Tag        string `json:"tag"`
	Name       string `json:"name"`
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func MakeStatusResponse(ctx context.Context) (sr *StatusResponse, err error) {
	name, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	sr = &StatusResponse{
		ServerName: getEnv("FIELDKIT_SERVER_NAME", ""),
		Tag:        getEnv("TAG", ""),
		Name:       name,
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

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})
}
