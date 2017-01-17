package webserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/O-C-R/singlepage"
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

func NewWebserver(c *config.Config) (*http.Server, error) {
	router := mux.NewRouter()
	router.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("cache-control", "no-store")
		fmt.Fprint(w, "ok")
	})

	api := mux.NewRouter()

	api.Handle("/api/user/sign-up", UserSignUpHandler(c))
	api.Handle("/api/user/validate", UserValidateHandler(c))
	api.Handle("/api/user/sign-in", UserSignInHandler(c))
	api.Handle("/api/user/current", AuthHandler(c, UserCurrentHandler(c)))

	api.Handle("/api/projects", AuthHandler(c, ProjectsHandler(c)))
	api.Handle("/api/projects/add", AuthHandler(c, ProjectAddHandler(c)))
	api.Handle("/api/project/{project}", ProjectHandler(c))

	api.Handle("/api/project/{project}/expeditions", ExpeditionsHandler(c))
	api.Handle("/api/project/{project}/expeditions/add", AuthHandler(c, AuthProjectHandler(c, ExpeditionAddHandler(c))))
	api.Handle("/api/project/{project}/expedition/{expedition}", ExpeditionHandler(c))
	api.Handle("/api/project/{project}/expedition/{expedition}/tokens", AuthHandler(c, AuthProjectHandler(c, ExpeditionAuthTokensHandler(c))))
	api.Handle("/api/project/{project}/expedition/{expedition}/tokens/add", AuthHandler(c, AuthProjectHandler(c, ExpeditionAuthTokenAddHandler(c))))

	api.Handle("/api/project/{project}/expedition/{expedition}/inputs", AuthHandler(c, AuthProjectHandler(c, InputsHandler(c))))
	api.Handle("/api/project/{project}/expedition/{expedition}/inputs/add", AuthHandler(c, AuthProjectHandler(c, InputAddHandler(c))))
	api.Handle("/api/project/{project}/expedition/{expedition}/input/{id}", AuthHandler(c, AuthProjectHandler(c, InputHandler(c))))

	api.Handle("/api/project/{project}/expedition/{expedition}/documents", DocumentsHandler(c))
	api.Handle("/api/project/{project}/expedition/{expedition}/documents/data", DocumentDataHandler(c))

	api.Handle("/api/input/{id}/{format:(?:fieldkit|csv|json)}/{source:(?:direct|rockblock)}", InputRequestHandler(c))

	router.PathPrefix("/api").Handler(handlers.CompressHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("cache-control", "no-store")
		api.ServeHTTP(w, req)
	})))

	if c.AdminPath != "" {
		application, err := regexp.Compile(`^/.*$`)
		if err != nil {
			return nil, err
		}

		longtermCache, err := regexp.Compile(`\.(?:(?:eot|png|ttf|woff|woff2)|(?:css|js|svg)(?:.gz))$`)
		if err != nil {
			return nil, err
		}

		router.Host("fieldkit.org").Handler(singlepage.NewSinglePageApplication(singlepage.SinglePageApplicationOptions{
			Root:          http.Dir(c.AdminPath),
			Application:   application,
			LongtermCache: longtermCache,
		}))
	}

	if c.FrontendPath != "" {
		application, err := regexp.Compile(`^/.*$`)
		if err != nil {
			return nil, err
		}

		longtermCache, err := regexp.Compile(`\.(?:(?:eot|png|ttf|woff|woff2)|(?:css|js|svg)(?:.gz))$`)
		if err != nil {
			return nil, err
		}

		router.Host("{project}.fieldkit.org").Handler(singlepage.NewSinglePageApplication(singlepage.SinglePageApplicationOptions{
			Root:          http.Dir(c.FrontendPath),
			Application:   application,
			LongtermCache: longtermCache,
		}))
	}

	handler := http.Handler(router)
	handler = ConfigHandler(handler, c)
	handler = handlers.CORS(
		handlers.AllowCredentials(),
		handlers.AllowedOrigins([]string{"http://localhost:8081", "https://fieldkit.org", "https://*.fieldkit.org"}),
	)(handler)

	server := &http.Server{
		Addr:    c.Addr,
		Handler: handler,
	}

	return server, nil
}
