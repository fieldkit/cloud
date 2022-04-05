package social

import (
	"context"
	"net/http"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/repositories"
)

type FacebookContext struct {
	SocialContext
}

type FacebookSchema struct {
}

func (s *FacebookSchema) SharedProject(ctx context.Context, w http.ResponseWriter, req *http.Request, payload *SharedProjectPayload) ([]*Meta, error) {
	meta := make([]*Meta, 0)

	// meta["og:url"] = req.URL.String()
	// meta["og:title"] = payload.project.Name
	// meta["og:description"] = payload.project.Description
	// meta["og:image"] = payload.photoUrl
	// meta["og:image:alt"] = payload.project.Description

	return meta, nil
}

func (s *FacebookSchema) SharedWorkspace(ctx context.Context, w http.ResponseWriter, req *http.Request, payload *SharedWorkspacePayload) ([]*Meta, error) {
	meta := make([]*Meta, 0)

	// size := 1080

	// meta["og:url"] = req.URL.String()
	// meta["og:image"] = fmt.Sprintf("%s&w=%d&h=%d", payload.photoUrl, size, size)
	// meta["og:image:alt"] = payload.description
	// meta["og:title"] = payload.title
	// meta["og:description"] = payload.description
	// meta["og:image:width"] = fmt.Sprintf("%d", size)
	// meta["og:image:height"] = meta["og:image:width"]

	return meta, nil
}

func NewFacebookContext(db *sqlxcache.DB, baseApiUrl string, rootPath string) (fb *FacebookContext) {
	return &FacebookContext{
		SocialContext{
			db:                db,
			projectRepository: repositories.NewProjectRepository(db),
			baseApiUrl:        baseApiUrl,
			rootPath:          rootPath,
			schema:            &FacebookSchema{},
		},
	}
}

/*
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
*/
