package social

import (
	"context"
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
	return nil

}

func (s *FacebookSchema) SharedWorkspace(ctx context.Context, w http.ResponseWriter, req *http.Request, payload *SharedWorkspacePayload) error {
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
	s := r.NewRoute().HeadersRegexp("User-Agent", ".*Facebook.*").Subrouter()
	s.HandleFunc("/dashboard/projects/{id:[0-9]+}", fb.SharedProject)
	s.HandleFunc("/dashboard/projects/{id:[0-9]+}/public", fb.SharedProject)
	s.HandleFunc("/dashboard/explore/{bookmark}", fb.SharedWorkspace)
	s.HandleFunc("/dashboard/share/{bookmark}", fb.SharedWorkspace)
	s.HandleFunc("/dashboard/explore", fb.SharedWorkspace)
	s.HandleFunc("/dashboard/share", fb.SharedWorkspace)
}
