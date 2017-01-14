package webserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/O-C-R/fieldkit/config"
)

var (
	InvalidConfigError = errors.New("invalid config")
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

type configKey struct{}

func ContextConfig(ctx context.Context) (*config.Config, error) {
	c, ok := ctx.Value(configKey{}).(*config.Config)
	if !ok {
		return nil, InvalidConfigError
	}

	return c, nil
}

func ConfigHandler(handler http.Handler, c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req = req.WithContext(context.WithValue(req.Context(), configKey{}, c))
		handler.ServeHTTP(w, req)
	})
}

// func FormHandler(handler http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 		vars := mux.Vars(req)
// 		if req.Header.Get("content-type") == "application/json" {
// 			if err := json.NewDecoder(req.Body).Decode(vars); err != nil {
// 				Error(w, err, 400)
// 			}
// 		}

// 	})
// }

func NewWebserver(c *config.Config) (*http.Server, error) {
	router := mux.NewRouter()
	router.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "ok")
	})

	api := router.PathPrefix("/api").Subrouter()

	api.Handle("/user/sign-up", UserSignUpHandler(c))
	api.Handle("/user/validate", UserValidateHandler(c))
	api.Handle("/user/sign-in", UserSignInHandler(c))
	api.Handle("/user/current", AuthHandler(c, UserCurrentHandler(c)))

	api.Handle("/projects", AuthHandler(c, ProjectsHandler(c)))
	api.Handle("/projects/add", AuthHandler(c, ProjectAddHandler(c)))
	api.Handle("/project/{project}", AuthHandler(c, ProjectHandler(c)))

	api.Handle("/project/{project}/expeditions", AuthHandler(c, AuthProjectHandler(c, ExpeditionsHandler(c))))
	api.Handle("/project/{project}/expeditions/add", AuthHandler(c, AuthProjectHandler(c, ExpeditionAddHandler(c))))
	api.Handle("/project/{project}/expedition/{expedition}", AuthHandler(c, AuthProjectHandler(c, ExpeditionHandler(c))))

	api.Handle("/project/{project}/expedition/{expedition}/inputs", AuthHandler(c, AuthProjectHandler(c, InputsHandler(c))))
	api.Handle("/project/{project}/expedition/{expedition}/inputs/add", AuthHandler(c, AuthProjectHandler(c, InputAddHandler(c))))
	api.Handle("/project/{project}/expedition/{expedition}/input/{id}", AuthHandler(c, AuthProjectHandler(c, InputHandler(c))))

	api.Handle("/input/{id}/{format:(?:fieldkit|csv|json)}/{source:(?:direct)}", InputRequestHandler(c))

	handler := http.Handler(router)
	handler = ConfigHandler(handler, c)
	handler = handlers.CORS(
		handlers.AllowCredentials(),
		handlers.AllowedOrigins([]string{"http://localhost:8000", "https://fieldkit.org"}),
	)(handler)

	server := &http.Server{
		Addr:    c.Addr,
		Handler: handler,
	}

	return server, nil
}
