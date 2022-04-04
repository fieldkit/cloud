package social

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/repositories"
)

type FacebookContext struct {
	SocialContext
}

type FacebookSchema struct {
}

func (s *FacebookSchema) SharedProject(ctx context.Context, w http.ResponseWriter, req *http.Request, payload *SharedProjectPayload) error {
	meta := make(map[string]string)

	meta["og:url"] = req.URL.String()
	meta["og:title"] = payload.project.Name
	meta["og:description"] = payload.project.Description
	meta["og:image"] = payload.photoUrl
	meta["og:image:alt"] = payload.project.Description

	if err := serveMeta(w, req, meta); err != nil {
		return err
	}

	return nil
}

func (s *FacebookSchema) SharedWorkspace(ctx context.Context, w http.ResponseWriter, req *http.Request, payload *SharedWorkspacePayload) error {
	meta := make(map[string]string)

	size := 1080

	meta["og:url"] = req.URL.String()
	meta["og:image"] = fmt.Sprintf("%s&w=%d&h=%d", payload.photoUrl, size, size)
	meta["og:image:alt"] = payload.description
	meta["og:title"] = payload.title
	meta["og:description"] = payload.description
	meta["og:image:width"] = fmt.Sprintf("%d", size)
	meta["og:image:height"] = meta["og:image:width"]

	if err := serveMeta(w, req, meta); err != nil {
		return err
	}

	return nil
}

func NewFacebookContext(db *sqlxcache.DB, baseApiUrl string) (fb *FacebookContext) {
	return &FacebookContext{
		SocialContext{
			db:                db,
			projectRepository: repositories.NewProjectRepository(db),
			baseApiUrl:        baseApiUrl,
			schema:            &FacebookSchema{},
		},
	}
}

func (fb *FacebookContext) Register(r *mux.Router) {
	s := r.NewRoute().MatcherFunc(matchUserAgent("facebook")).Subrouter()
	s.HandleFunc("/dashboard/projects/{id:[0-9]+}", fb.SharedProject)
	s.HandleFunc("/dashboard/projects/{id:[0-9]+}/public", fb.SharedProject)
	s.HandleFunc("/dashboard/explore/{bookmark}", fb.SharedWorkspace)
	s.HandleFunc("/dashboard/share/{bookmark}", fb.SharedWorkspace)
	s.HandleFunc("/dashboard/explore", fb.SharedWorkspace)
	s.HandleFunc("/dashboard/share", fb.SharedWorkspace)
	s.HandleFunc("/viz", fb.SharedWorkspace)
}
