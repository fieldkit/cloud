package social

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/repositories"
)

type TwitterContext struct {
	SocialContext
}

type TwitterSchema struct {
}

func (s *TwitterSchema) SharedProject(ctx context.Context, w http.ResponseWriter, req *http.Request, payload *SharedProjectPayload) (map[string]string, error) {
	meta := make(map[string]string)

	meta["twitter:card"] = "summary_large_image"
	meta["twitter:site"] = "@FieldKitOrg"
	meta["twitter:title"] = payload.project.Name
	meta["twitter:description"] = payload.project.Description
	meta["twitter:image:alt"] = payload.project.Description
	meta["twitter:image"] = payload.photoUrl
	meta["twitter:url"] = req.URL.String()

	meta["og:url"] = meta["twitter:url"]
	meta["og:title"] = meta["twitter:title"]
	meta["og:description"] = meta["twitter:description"]
	meta["og:image"] = meta["twitter:image"]
	meta["og:image:alt"] = meta["twitter:image:alt"]

	return meta, nil
}

func parseDimensionParam(req *http.Request, name string, d int) int {
	if str := req.URL.Query().Get(name); str != "" {
		v, err := strconv.Atoi(str)
		if err == nil {
			return v
		}
	}
	return d
}

func (s *TwitterSchema) SharedWorkspace(ctx context.Context, rw http.ResponseWriter, req *http.Request, payload *SharedWorkspacePayload) (map[string]string, error) {
	meta := make(map[string]string)

	w := parseDimensionParam(req, "w", 800)
	h := parseDimensionParam(req, "w", 400)

	meta["twitter:card"] = "summary_large_image"
	meta["twitter:site"] = "@FieldKitOrg"
	meta["twitter:title"] = payload.title
	meta["twitter:domain"] = ""
	meta["twitter:description"] = payload.description
	meta["twitter:image:alt"] = payload.description
	meta["twitter:image"] = fmt.Sprintf("%s&w=%d&h=%d", payload.photoUrl, w, h)
	meta["twitter:url"] = req.URL.String()

	meta["og:url"] = meta["twitter:url"]
	meta["og:title"] = meta["twitter:title"]
	meta["og:description"] = meta["twitter:description"]
	meta["og:image"] = meta["twitter:image"]
	meta["og:image:alt"] = meta["twitter:image:alt"]
	meta["og:image:width"] = fmt.Sprintf("%d", w)
	meta["og:image:height"] = fmt.Sprintf("%d", h)

	return meta, nil
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
	s := r.NewRoute().MatcherFunc(matchUserAgent("twitterbot")).Subrouter()
	s.HandleFunc("/dashboard/projects/{id:[0-9]+}", tw.SharedProject)
	s.HandleFunc("/dashboard/projects/{id:[0-9]+}/public", tw.SharedProject)
	s.HandleFunc("/dashboard/explore/{bookmark}", tw.SharedWorkspace)
	s.HandleFunc("/dashboard/share/{bookmark}", tw.SharedWorkspace)
	s.HandleFunc("/dashboard/explore", tw.SharedWorkspace)
	s.HandleFunc("/dashboard/share", tw.SharedWorkspace)
	s.HandleFunc("/viz", tw.SharedWorkspace)
}
