package social

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/repositories"
)

type TwitterContext struct {
	SocialContext
}

type TwitterSchema struct {
}

func (s *TwitterSchema) SharedProject(ctx context.Context, w http.ResponseWriter, req *http.Request, payload *SharedProjectPayload) error {
	meta := make(map[string]string)

	meta["twitter:card"] = "summary_large_image"
	meta["twitter:site"] = "@FieldKitOrg"
	meta["twitter:title"] = payload.project.Name
	meta["twitter:description"] = payload.project.Description
	meta["twitter:image:alt"] = payload.project.Description
	meta["twitter:image"] = payload.photoUrl

	if err := serveMeta(w, req, meta); err != nil {
		return err
	}

	return nil
}

func (s *TwitterSchema) SharedWorkspace(ctx context.Context, w http.ResponseWriter, req *http.Request, payload *SharedWorkspacePayload) error {
	meta := make(map[string]string)

	meta["twitter:card"] = "summary_large_image"
	meta["twitter:site"] = "@FieldKitOrg"
	meta["twitter:title"] = payload.title
	meta["twitter:description"] = payload.description
	meta["twitter:image:alt"] = payload.description
	meta["twitter:image"] = payload.photoUrl

	if err := serveMeta(w, req, meta); err != nil {
		return err
	}

	return nil
}

func NewTwitterContext(db *sqlxcache.DB, baseApiUrl string) (tw *TwitterContext) {
	return &TwitterContext{
		SocialContext{
			db:                db,
			projectRepository: repositories.NewProjectRepository(db),
			baseApiUrl:        baseApiUrl,
			schema:            &TwitterSchema{},
		},
	}
}

func (tw *TwitterContext) Register(r *mux.Router) {
	s := r.NewRoute().HeadersRegexp("User-Agent", ".*Twitterbot.*").Subrouter()
	s.HandleFunc("/dashboard/projects/{id:[0-9]+}", tw.SharedProject)
	s.HandleFunc("/dashboard/projects/{id:[0-9]+}/public", tw.SharedProject)
	s.HandleFunc("/dashboard/explore/{bookmark}", tw.SharedWorkspace)
	s.HandleFunc("/dashboard/share/{bookmark}", tw.SharedWorkspace)
	s.HandleFunc("/dashboard/explore", tw.SharedWorkspace)
	s.HandleFunc("/dashboard/share", tw.SharedWorkspace)
}
