package social

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

type TwitterContext struct {
	db                *sqlxcache.DB
	projectRepository *repositories.ProjectRepository
	baseApiUrl        string
}

func NewTwitterContext(db *sqlxcache.DB, baseApiUrl string) (tw *TwitterContext) {
	return &TwitterContext{
		db:                db,
		projectRepository: repositories.NewProjectRepository(db),
		baseApiUrl:        baseApiUrl,
	}
}

const metaOnlyTemplate = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		{{- range $key, $value := .Metas }}
		<meta name="{{ $key }}" content="{{ $value }}" />
		{{- end -}}
	</head>
	<body>
	</body>
</html>`

func (tw *TwitterContext) serveMeta(w http.ResponseWriter, req *http.Request, meta map[string]string) error {
	t, err := template.New("twitter-meta").Parse(metaOnlyTemplate)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)

	data := struct {
		Metas map[string]string
	}{
		Metas: meta,
	}

	err = t.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

func (tw *TwitterContext) SharedProject(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := Logger(ctx).Sugar()
	vars := mux.Vars(req)

	projectId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Errorw("error-internal", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	project, err := tw.projectRepository.QueryByID(ctx, int32(projectId))
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Errorw("error-query", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if project.Privacy != data.Public {
		log.Errorw("error-permission", "project_id", projectId)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// NOTE TODO We're casually assuming https everywhere.
	photoUrl := fmt.Sprintf("%s/projects/%d/media", tw.baseApiUrl, project.ID)

	log.Infow("twitter-project-card", "project_id", project.ID, "url", req.URL)

	meta := make(map[string]string)
	meta["twitter:card"] = "summary_large_image"
	meta["twitter:site"] = "@FieldKitOrg"
	meta["twitter:title"] = project.Name
	meta["twitter:description"] = project.Description
	meta["twitter:image"] = photoUrl
	meta["twitter:image:alt"] = project.Description

	if err := tw.serveMeta(w, req, meta); err != nil {
		log.Errorw("error-internal", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (tw *TwitterContext) SharedWorkspace(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := Logger(ctx).Sugar()
	vars := mux.Vars(req)

	bookmark := vars["bookmark"]

	log.Infow("twitter-workspace-card", "bookmark", bookmark)

	// NOTE TODO We're casually assuming https everywhere.
	photoUrl := fmt.Sprintf("%s/charting/rendered?bookmark=%v", tw.baseApiUrl, url.QueryEscape(bookmark))

	meta := make(map[string]string)
	meta["twitter:card"] = "summary_large_image"
	meta["twitter:site"] = "@FieldKitOrg"
	meta["twitter:title"] = "Name"
	meta["twitter:description"] = "Description"
	meta["twitter:image"] = photoUrl
	meta["twitter:image:alt"] = "Description"

	if err := tw.serveMeta(w, req, meta); err != nil {
		log.Errorw("error-internal", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (tw *TwitterContext) Register(r *mux.Router) {
	s := r.NewRoute().HeadersRegexp("User-Agent", ".*Twitterbot.*").Subrouter()
	s.HandleFunc("/dashboard/projects/{id:[0-9]+}", tw.SharedProject)
	s.HandleFunc("/dashboard/projects/{id:[0-9]+}/public", tw.SharedProject)
	s.HandleFunc("/dashboard/explore/{bookmark}", tw.SharedWorkspace)
}
