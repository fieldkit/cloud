package health

import (
	"context"
	"encoding/json"
	_ "fmt"
	"net/http"
	"os"

	"github.com/fieldkit/cloud/server/common/logging"
)

type GitStatus struct {
	Hash string `json:"hash,omitempty"`
}

type StatusResponse struct {
	ServerName string    `json:"server_name,omitempty"`
	Version    string    `json:"version"`
	Tag        string    `json:"tag"`
	Name       string    `json:"name"`
	Git        GitStatus `json:"git"`
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
		Git: GitStatus{
			Hash: getEnv("GIT_HASH", ""),
		},
		Version:    getEnv("FIELDKIT_VERSION", ""),
		ServerName: getEnv("FIELDKIT_SERVER_NAME", ""),
		Tag:        getEnv("TAG", ""),
		Name:       name,
	}

	return
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
}

func StatusHandler(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		setupResponse(&w, req)

		if (*req).Method == "OPTIONS" {
			return
		}

		ctx := req.Context()
		log := logging.Logger(ctx).Named("status").Sugar()

		req.Header.Del("Authorization")

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
