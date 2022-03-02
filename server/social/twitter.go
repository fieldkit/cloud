package social

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"

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

	parsed, err := data.ParseBookmark(bookmark)
	if err != nil {
		log.Errorw("error-internal", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sensorRows := []*data.Sensor{}
	if err := tw.db.SelectContext(ctx, &sensorRows, `SELECT * FROM fieldkit.aggregated_sensor ORDER BY key`); err != nil {
		log.Errorw("error-internal", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mmr := repositories.NewModuleMetaRepository(tw.db)

	sensorIdToKey := make(map[int64]string)
	for _, row := range sensorRows {
		sensorIdToKey[row.ID] = row.Key
	}

	sr := repositories.NewStationRepository(tw.db)

	sensorKeys := make([]string, 0)
	sensorLabels := make([]string, 0)
	stationNames := make([]string, 0)
	ranges := make([]string, 0)

	meta := make(map[string]string)

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

		meta["twitter:title"] = fmt.Sprintf("%s: %s", stationNames[0], sensorLabels[0])

		if !v.ExtremeTime {
			meta["twitter:description"] = ranges[0]
		}

		break // NOTE Single out first Viz.
	}

	log.Infow("viz", "sensors", sensorKeys)

	// NOTE TODO We're casually assuming https everywhere.
	photoUrl := fmt.Sprintf("%s/charting/rendered?bookmark=%v", tw.baseApiUrl, url.QueryEscape(bookmark))
	meta["twitter:card"] = "summary_large_image"
	meta["twitter:site"] = "@FieldKitOrg"
	meta["twitter:image"] = photoUrl
	meta["twitter:image:alt"] = meta["twitter:description"]

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
	s.HandleFunc("/dashboard/share/{bookmark}", tw.SharedWorkspace)
}
