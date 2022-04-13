package social

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/repositories"
)

type CombinedContext struct {
	SocialContext
}

type CombinedSchema struct {
}

func (s *CombinedSchema) SharedProject(ctx context.Context, w http.ResponseWriter, req *http.Request, payload *SharedProjectPayload) ([]*Meta, error) {
	meta := make([]*Meta, 0)

	meta = append(meta, NewMetaName("twitter:card", "summary_large_image"))
	meta = append(meta, NewMetaName("twitter:site", "@FieldKitOrg"))
	meta = append(meta, NewMetaName("twitter:title", payload.project.Name))
	meta = append(meta, NewMetaName("twitter:description", payload.project.Description))
	meta = append(meta, NewMetaName("twitter:image:alt", payload.project.Description))
	meta = append(meta, NewMetaName("twitter:image", payload.photoUrl))
	meta = append(meta, NewMetaName("twitter:url", payload.url))

	meta = append(meta, NewMetaProperty("og:url", payload.url))
	meta = append(meta, NewMetaProperty("og:type", "website"))
	meta = append(meta, NewMetaProperty("og:title", payload.project.Name))
	meta = append(meta, NewMetaProperty("og:description", payload.project.Description))
	meta = append(meta, NewMetaProperty("og:image", payload.photoUrl))
	meta = append(meta, NewMetaProperty("og:image:secure_url", payload.photoUrl))
	meta = append(meta, NewMetaProperty("og:image:alt", payload.project.Description))
	meta = append(meta, NewMetaProperty("og:image:width", fmt.Sprintf("%d", payload.width)))
	meta = append(meta, NewMetaProperty("og:image:height", fmt.Sprintf("%d", payload.height)))
	meta = append(meta, NewMetaProperty("og:image:type", "image/png"))

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

func (s *CombinedSchema) SharedWorkspace(ctx context.Context, rw http.ResponseWriter, req *http.Request, payload *SharedWorkspacePayload) ([]*Meta, error) {
	w := parseDimensionParam(req, "w", 800)
	h := parseDimensionParam(req, "w", 400)
	photoUrl := fmt.Sprintf("%s&w=%d&h=%d", payload.photoUrl, w, h)

	meta := make([]*Meta, 0)

	meta = append(meta, NewMetaName("twitter:card", "summary_large_image"))
	meta = append(meta, NewMetaName("twitter:site", "@FieldKitOrg"))
	meta = append(meta, NewMetaName("twitter:title", payload.title))
	meta = append(meta, NewMetaName("twitter:description", payload.description))
	meta = append(meta, NewMetaName("twitter:image", photoUrl))
	meta = append(meta, NewMetaName("twitter:url", payload.url))

	meta = append(meta, NewMetaProperty("og:url", payload.url))
	meta = append(meta, NewMetaProperty("og:type", "website"))
	meta = append(meta, NewMetaProperty("og:title", payload.title))
	meta = append(meta, NewMetaProperty("og:description", payload.description))
	meta = append(meta, NewMetaProperty("og:image", photoUrl))
	meta = append(meta, NewMetaProperty("og:image:secure_url", photoUrl))
	meta = append(meta, NewMetaProperty("og:image:alt", payload.description))
	meta = append(meta, NewMetaProperty("og:image:width", fmt.Sprintf("%d", w)))
	meta = append(meta, NewMetaProperty("og:image:height", fmt.Sprintf("%d", h)))
	meta = append(meta, NewMetaProperty("og:image:type", "image/png"))

	return meta, nil
}

func NewContext(db *sqlxcache.DB, rootPath string) (cc *CombinedContext) {
	return &CombinedContext{
		SocialContext{
			db:                db,
			projectRepository: repositories.NewProjectRepository(db),
			rootPath:          rootPath,
			schema:            &CombinedSchema{},
		},
	}
}
