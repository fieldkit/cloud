package social

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

type MetaSchema interface {
	SharedProject(ctx context.Context, w http.ResponseWriter, req *http.Request, payload *SharedProjectPayload) error
	SharedWorkspace(ctx context.Context, w http.ResponseWriter, req *http.Request, payload *SharedWorkspacePayload) error
}

type SocialContext struct {
	db                *sqlxcache.DB
	projectRepository *repositories.ProjectRepository
	baseApiUrl        string
	schema            MetaSchema
}

type SharedProjectPayload struct {
	project  *data.Project
	photoUrl string
}

type SharedWorkspacePayload struct {
	bookmark    *data.Bookmark
	title       string
	description string
	photoUrl    string
}

func matchUserAgent(partial string) mux.MatcherFunc {
	return func(req *http.Request, match *mux.RouteMatch) bool {
		if userAgent, ok := req.Header[http.CanonicalHeaderKey("user-agent")]; ok {
			if len(userAgent) == 0 {
				return false
			}
			return strings.Contains(strings.ToLower(userAgent[0]), strings.ToLower(partial))
		}
		return false
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

func serveMeta(w http.ResponseWriter, req *http.Request, meta map[string]string) error {
	t, err := template.New("social-meta").Parse(metaOnlyTemplate)
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

func (sc *SocialContext) SharedProject(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := Logger(ctx).Sugar()
	vars := mux.Vars(req)

	projectId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Errorw("error-internal", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	project, err := sc.projectRepository.QueryByID(ctx, int32(projectId))
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
	photoUrl := fmt.Sprintf("%s/projects/%d/media", sc.baseApiUrl, project.ID)

	log.Infow("social-project-card", "project_id", project.ID, "url", req.URL)

	sharedPayload := &SharedProjectPayload{
		project:  project,
		photoUrl: photoUrl,
	}

	err = sc.schema.SharedProject(ctx, w, req, sharedPayload)
	if err != nil {
		log.Errorw("error", "project_id", projectId)
		w.WriteHeader(http.StatusForbidden)
		return
	}
}

func (sc *SocialContext) SharedWorkspace(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := Logger(ctx).Sugar()
	vars := mux.Vars(req)

	bookmark := ""
	if vars["bookmark"] != "" {
		bookmark = vars["bookmark"] // Deprecated
	}

	qs := req.URL.Query()

	if qs.Get("bookmark") != "" {
		bookmark = qs.Get("bookmark")
	}

	log.Infow("social-workspace-card", "bookmark", bookmark)

	parsed, err := data.ParseBookmark(bookmark)
	if err != nil {
		log.Errorw("error-internal", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sensorRows := []*data.Sensor{}
	if err := sc.db.SelectContext(ctx, &sensorRows, `SELECT * FROM fieldkit.aggregated_sensor ORDER BY key`); err != nil {
		log.Errorw("error-internal", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mmr := repositories.NewModuleMetaRepository(sc.db)

	sensorIdToKey := make(map[int64]string)
	for _, row := range sensorRows {
		sensorIdToKey[row.ID] = row.Key
	}

	sr := repositories.NewStationRepository(sc.db)

	sensorKeys := make([]string, 0)
	sensorLabels := make([]string, 0)
	stationNames := make([]string, 0)
	ranges := make([]string, 0)

	title := ""
	description := ""

	for _, v := range parsed.Vizes() {
		log.Infow("viz:parsed", "v", v)

		for _, s := range v.Sensors {
			station, err := sr.QueryStationByID(ctx, s.StationID)
			if err != nil {
				log.Errorw("error-internal", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			sensorKey := sensorIdToKey[s.SensorID]

			sensorMeta, err := mmr.FindByFullKey(ctx, sensorKey)
			if err != nil {
				log.Errorw("error-internal", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			label := sensorMeta.Sensor.Strings["en-us"]["label"]
			stationNames = append(stationNames, station.Name)
			sensorKeys = append(sensorKeys, sensorKey)
			sensorLabels = append(sensorLabels, label)
		}

		start := v.Start.Format(time.RFC1123)
		end := v.End.Format(time.RFC1123)
		ranges = append(ranges, fmt.Sprintf("%s to %s", start, end))

		title = fmt.Sprintf("%s: %s", stationNames[0], sensorLabels[0])

		if !v.ExtremeTime {
			description = ranges[0]
		}

		break // NOTE Single out first Viz.
	}

	// We're assuming https here.
	now := time.Now()
	photoUrl := fmt.Sprintf("%s/charting/rendered?bookmark=%v&ts=%v", sc.baseApiUrl, url.QueryEscape(bookmark), now.Unix())

	sharedPayload := &SharedWorkspacePayload{
		photoUrl:    photoUrl,
		title:       title,
		description: description,
		bookmark:    parsed,
	}

	err = sc.schema.SharedWorkspace(ctx, w, req, sharedPayload)
	if err != nil {
		log.Errorw("error")
		w.WriteHeader(http.StatusForbidden)
		return
	}
}
