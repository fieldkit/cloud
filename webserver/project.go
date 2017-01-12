package webserver

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/O-C-R/fieldkit/config"
	"github.com/O-C-R/fieldkit/data"
)

func ProjectsHandler(c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		projects, err := c.Backend.Projects()
		if err != nil {
			Error(w, err, 500)
			return
		}

		WriteJSON(w, projects)
	})
}

func ProjectHandler(c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		project, err := c.Backend.ProjectBySlug(vars["project"])
		if err != nil {
			Error(w, err, 500)
			return
		}

		WriteJSON(w, project)
	})
}

func ProjectAddHandler(c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		userID, valid := ContextUserID(req.Context())
		if !valid {
			Error(w, InvalidUserError, 500)
			return
		}

		errs := data.NewErrors()

		// Validate the project name.
		name, slug := data.Name(req.FormValue("name"))
		if name == "" || slug == "" {
			errs.Error("name", "invalid Name")
		} else {
			slugInUse, err := c.Backend.ProjectSlugInUse(slug)
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

		project, err := data.NewProject(name, slug)
		if err != nil {
			Error(w, err, 500)
			return
		}

		if err := c.Backend.AddProjectWithOwner(project, userID); err != nil {
			WriteJSONStatusCode(w, err, 500)
			return
		}

		WriteJSON(w, project)
	})
}

func AuthProjectHandler(c *config.Config, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		userID, valid := ContextUserID(req.Context())
		if !valid {
			Error(w, InvalidUserError, 500)
			return
		}

		vars := mux.Vars(req)
		if _, err := c.Backend.ProjectUserRoleByProjectSlugAndUserID(vars["project"], userID); err != nil {
			Error(w, err, 500)
			return
		}

		handler.ServeHTTP(w, req)
	})
}
