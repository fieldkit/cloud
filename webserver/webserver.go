package webserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/O-C-R/fieldkit/config"
)

func Error(w http.ResponseWriter, err error, code int) {
	if code >= 400 {
		log.Println(err)
	}

	if code >= 500 {
		http.Error(w, "an error occurred", code)
		return
	}

	http.Error(w, err.Error(), code)
}

func WriteJSON(w http.ResponseWriter, value interface{}) {
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Println(err)
	}
}

func WriteJSONStatusCode(w http.ResponseWriter, value interface{}, code int) {
	w.WriteHeader(code)
	WriteJSON(w, value)
}

func NewWebserver(c *config.Config) (*http.Server, error) {
	router := mux.NewRouter()
	router.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "ok")
	})

	router.Handle("/api/user/sign-up", UserSignUpHandler(c))
	router.Handle("/api/user/validate", UserValidateHandler(c))
	router.Handle("/api/user/sign-in", UserSignInHandler(c))

	router.Handle("/api/user/current", AuthHandler(c, UserCurrentHandler(c)))

	handler := http.Handler(router)
	handler = handlers.CORS(
		handlers.AllowCredentials(),
		handlers.AllowedOrigins([]string{"http://localhost:8000"}),
	)(handler)

	server := &http.Server{
		Addr:    c.Addr,
		Handler: handler,
	}

	return server, nil
}
