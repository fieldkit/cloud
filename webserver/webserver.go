package webserver

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/O-C-R/fieldkit/config"
)

func NewWebserver(c *config.Config) (*http.Server, error) {
	handler := mux.NewRouter()
	handler.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "ok")
	})

	server := &http.Server{
		Addr:    c.Addr,
		Handler: handler,
	}

	return server, nil
}
