package webserver

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/O-C-R/fieldkit/backend"
	"github.com/O-C-R/fieldkit/config"
	"github.com/O-C-R/fieldkit/data"
)

func ExpeditionsHandler(c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		expeditions, err := c.Backend.ExpeditionsByProjectSlug(vars["project"])
		if err != nil {
			Error(w, err, 500)
			return
		}

		WriteJSON(w, expeditions)
	})
}

func ExpeditionHandler(c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		expedition, err := c.Backend.ExpeditionByProjectSlugAndSlug(vars["project"], vars["expedition"])
		if err != nil {
			Error(w, err, 500)
			return
		}

		WriteJSON(w, expedition)
	})
}

func ExpeditionAddHandler(c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		project, err := c.Backend.ProjectBySlug(vars["project"])
		if err == backend.NotFoundError {
			Error(w, err, 404)
			return
		}

		if err != nil {
			Error(w, err, 500)
			return
		}

		errs := data.NewErrors()

		// Validate the project name.
		name, slug := data.Name(req.FormValue("name"))
		if name == "" || slug == "" {
			errs.Error("name", "invalid Name")
		} else {
			slugInUse, err := c.Backend.ExpeditionSlugInUse(project.ID, slug)
			if err != nil {
				Error(w, err, 500)
				return
			}

			if slugInUse {
				errs.Error("name", "Name in use")
			}
		}

		if len(errs) > 0 {
			WriteJSONStatusCode(w, errs, 400)
			return
		}

		expedition, err := data.NewExpedition(project.ID, name, slug)
		if err != nil {
			Error(w, err, 500)
			return
		}

		if err := c.Backend.AddExpedition(expedition); err != nil {
			WriteJSONStatusCode(w, err, 500)
			return
		}

		WriteJSON(w, expedition)
	})
}

func ExpeditionAuthTokenAddHandler(c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		expedition, err := c.Backend.ExpeditionByProjectSlugAndSlug(vars["project"], vars["expedition"])
		if err == backend.NotFoundError {
			Error(w, err, 404)
			return
		}

		if err != nil {
			Error(w, err, 500)
			return
		}

		authToken, err := data.NewAuthToken(expedition.ID)
		if err != nil {
			Error(w, err, 500)
			return
		}

		if err := c.Backend.AddAuthToken(authToken); err != nil {
			WriteJSONStatusCode(w, err, 500)
			return
		}

		WriteJSON(w, authToken)
	})
}

func ExpeditionAuthTokensHandler(c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		expeditions, err := c.Backend.AuthTokensByProjectSlugAndExpeditionSlug(vars["project"], vars["expedition"])
		if err != nil {
			Error(w, err, 500)
			return
		}

		WriteJSON(w, expeditions)
	})
}
