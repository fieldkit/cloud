package social

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

type Meta struct {
	Attribute string
	Key       string
	Value     string
}

func NewMetaName(key, value string) *Meta {
	return &Meta{Attribute: "name", Key: key, Value: value}
}

func NewMetaProperty(key, value string) *Meta {
	return &Meta{Attribute: "property", Key: key, Value: value}
}

type MetaSchema interface {
	SharedProject(ctx context.Context, w http.ResponseWriter, req *http.Request, payload *SharedProjectPayload) ([]*Meta, error)
	SharedWorkspace(ctx context.Context, w http.ResponseWriter, req *http.Request, payload *SharedWorkspacePayload) ([]*Meta, error)
}

type SocialContext struct {
	db                *sqlxcache.DB
	projectRepository *repositories.ProjectRepository
	schema            MetaSchema
	rootPath          string
}

type SharedProjectPayload struct {
	project  *data.Project
	url      string
	photoUrl string
	width    int
	height   int
}

type SharedWorkspacePayload struct {
	bookmark    *data.Bookmark
	title       string
	description string
	url         string
	photoUrl    string
}

const metaOnlyTemplate = `{{- range $i, $meta := .Metas }}
<meta {{ $meta.Attribute }}="{{ $meta.Key }}" content="{{ $meta.Value }}" />
{{- end }}`

func (sc *SocialContext) getRequestHost(req *http.Request) string {
	protocol := "https" // We're never accessed over http.
	forwardedHostKey := http.CanonicalHeaderKey("x-forwarded-host")
	forwardedHeader := req.Header.Values(forwardedHostKey)
	if len(forwardedHeader) == 0 {
		return fmt.Sprintf("%s://%s", protocol, req.Host)
	}
	return fmt.Sprintf("%s://%s", protocol, forwardedHeader[0])
}

func (sc *SocialContext) SharedProject(ctx context.Context, w http.ResponseWriter, req *http.Request) ([]*Meta, error) {
	log := Logger(ctx).Sugar()
	vars := mux.Vars(req)

	projectId, err := strconv.Atoi(vars["id"])
	if err != nil {
		return nil, err
	}

	project, err := sc.projectRepository.QueryByID(ctx, int32(projectId))
	if err != nil {
		return nil, err
	}

	if project.Privacy != data.Public {
		return nil, err
	}

	size := 800
	requestHost := sc.getRequestHost(req)
	photoUrl := fmt.Sprintf("%s/projects/%d/media?size=%d", requestHost, project.ID, size)
	linkUrl := fmt.Sprintf("%s%s", requestHost, req.URL.String())

	log.Infow("social-project-card", "project_id", project.ID, "url", req.URL, "request_host", requestHost)

	sharedPayload := &SharedProjectPayload{
		project:  project,
		url:      linkUrl,
		photoUrl: photoUrl,
		width:    size,
		height:   size,
	}

	meta, err := sc.schema.SharedProject(ctx, w, req, sharedPayload)
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func (sc *SocialContext) SharedWorkspace(ctx context.Context, w http.ResponseWriter, req *http.Request) ([]*Meta, error) {
	log := Logger(ctx).Sugar()

	bookmark := ""
	// Deprecated
	vars := mux.Vars(req)
	if vars["bookmark"] != "" {
		bookmark = vars["bookmark"]
	}
	qs := req.URL.Query()
	if qs.Get("bookmark") != "" {
		bookmark = qs.Get("bookmark")
	}
	if qs.Get("v") != "" {
		token := qs.Get("v")

		repository := repositories.NewBookmarkRepository(sc.db)
		resolved, err := repository.Resolve(ctx, token)
		if err != nil {
			return nil, err
		}
		if resolved == nil {
			return nil, err
		}

		bookmark = resolved.Bookmark
	}

	if bookmark == "" {
		return nil, fmt.Errorf("empty bookmark")
	}

	requestHost := sc.getRequestHost(req)

	log.Infow("social-workspace-card", "bookmark", bookmark, "url", req.URL, "request_host", requestHost)

	parsed, err := data.ParseBookmark(bookmark)
	if err != nil {
		return nil, err
	}

	sensorRows := []*data.Sensor{}
	if err := sc.db.SelectContext(ctx, &sensorRows, `SELECT * FROM fieldkit.aggregated_sensor ORDER BY key`); err != nil {
		return nil, err
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
	description := "FieldKit Chart"

	for _, v := range parsed.Vizes() {
		log.Infow("viz:parsed", "v", v)

		for _, s := range v.Sensors {
			station, err := sr.QueryStationByID(ctx, s.StationID)
			if err != nil {
				return nil, err
			}

			sensorKey := sensorIdToKey[s.SensorID]

			sensorMeta, err := mmr.FindByFullKey(ctx, sensorKey)
			if err != nil {
				return nil, err
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
	photoUrl := fmt.Sprintf("%s/charting/rendered?bookmark=%v&ts=%v", requestHost, url.QueryEscape(bookmark), now.Unix())
	linkUrl := fmt.Sprintf("%s%s", requestHost, req.URL.String())

	sharedPayload := &SharedWorkspacePayload{
		url:         linkUrl,
		photoUrl:    photoUrl,
		title:       title,
		description: description,
		bookmark:    parsed,
	}

	meta, err := sc.schema.SharedWorkspace(ctx, w, req, sharedPayload)
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func (sc SocialContext) serveMeta(ctx context.Context, w http.ResponseWriter, req *http.Request, meta []*Meta) error {
	serving, err := os.ReadFile(filepath.Join(sc.rootPath, "index.html"))
	if err != nil {
		return err
	}

	if len(meta) > 0 {
		template, err := template.New("social-meta").Parse(metaOnlyTemplate)
		if err == nil {
			data := struct {
				Metas []*Meta
			}{
				Metas: meta,
			}

			var rendered bytes.Buffer
			err = template.Execute(&rendered, data)
			if err == nil {
				serving = []byte(strings.Replace(string(serving), "<title>", rendered.String()+"<title>", 1))
			}
		}
	} else {
		log := Logger(ctx).Sugar()
		log.Infow("social:meta:empty")
	}

	// TODO Include `error` in this either as data.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(serving))

	return err
}

func (sc SocialContext) logErrorsAndServeOriginalIndex(f func(ctx context.Context, w http.ResponseWriter, req *http.Request) ([]*Meta, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := context.Background()
		meta, err := f(ctx, w, req)
		if err != nil {
			log := Logger(ctx).Sugar()
			log.Warnw("social:meta:error", "error", err)
		}

		err = sc.serveMeta(ctx, w, req, meta)
		if err != nil {
			log := Logger(ctx).Sugar()
			log.Warnw("social:serve:error", "error", err)
		}
	}
}

func (sc *SocialContext) Register(r *mux.Router) {
	r.HandleFunc("/dashboard/projects/{id:[0-9]+}", sc.logErrorsAndServeOriginalIndex(sc.SharedProject))
	r.HandleFunc("/dashboard/projects/{id:[0-9]+}/public", sc.logErrorsAndServeOriginalIndex(sc.SharedProject))
	r.HandleFunc("/dashboard/explore/{bookmark}", sc.logErrorsAndServeOriginalIndex(sc.SharedWorkspace))
	r.HandleFunc("/dashboard/share/{bookmark}", sc.logErrorsAndServeOriginalIndex(sc.SharedWorkspace))
	r.HandleFunc("/dashboard/explore", sc.logErrorsAndServeOriginalIndex(sc.SharedWorkspace))
	r.HandleFunc("/dashboard/share", sc.logErrorsAndServeOriginalIndex(sc.SharedWorkspace))
	r.HandleFunc("/viz/share", sc.logErrorsAndServeOriginalIndex(sc.SharedWorkspace))
	r.HandleFunc("/viz/export", sc.logErrorsAndServeOriginalIndex(sc.SharedWorkspace))
	r.HandleFunc("/viz", sc.logErrorsAndServeOriginalIndex(sc.SharedWorkspace))
}
