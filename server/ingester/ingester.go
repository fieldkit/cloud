package ingester

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"time"

	"github.com/lib/pq"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware/security/jwt"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/common"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/jobs"
	"github.com/fieldkit/cloud/server/logging"
	"github.com/fieldkit/cloud/server/messages"
)

var (
	ids = logging.NewIdGenerator()
)

type IngesterOptions struct {
	Database                 *sqlxcache.DB
	AwsSession               *session.Session
	AuthenticationMiddleware goa.Middleware
	Archiver                 files.StreamArchiver
	Publisher                jobs.MessagePublisher
	Metrics                  *logging.Metrics
}

func Ingester(ctx context.Context, o *IngesterOptions) http.Handler {
	errors := ErrorHandler()

	handler := errorHandling(errors, authentication(o.AuthenticationMiddleware, func(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
		startedAt := time.Now()

		log := Logger(ctx).Sugar()

		token := jwt.ContextJWT(ctx)
		if token == nil {
			return fmt.Errorf("JWT token is missing from context")
		}

		claims, ok := token.Claims.(jwtgo.MapClaims)
		if !ok {
			return fmt.Errorf("JWT claims error")
		}

		userID := int32(claims["sub"].(float64))

		o.Metrics.UserID(userID)

		headers, err := NewIncomingHeaders(req)
		if err != nil {
			log.Infow("headers error", "headers", req.Header)
			return err
		}

		o.Metrics.IngestionDevice(headers.FkDeviceId)

		log.Infow("receiving", "device_id", headers.FkDeviceId, "blocks", headers.FkBlocks)

		fileMeta := &files.FileMeta{
			ContentType: headers.ContentType,
			DeviceID:    headers.FkDeviceId,
			Generation:  headers.FkGeneration,
			Blocks:      headers.FkBlocks,
			Flags:       headers.FkFlags,
		}
		if saved, err := o.Archiver.Archive(ctx, fileMeta, req.Body); err != nil {
			return err
		} else {
			if saved != nil {
				if saved.BytesRead != int(headers.ContentLength) {
					log.Warnw("size mismatch", "expected", headers.ContentLength, "actual", saved.BytesRead)
				}

				ingestion := &data.Ingestion{
					URL:        saved.URL,
					UploadID:   saved.ID,
					UserID:     userID,
					DeviceID:   headers.FkDeviceId,
					Generation: headers.FkGeneration,
					Type:       headers.FkType,
					Size:       int64(saved.BytesRead),
					Blocks:     data.Int64Range(headers.FkBlocks),
					Flags:      pq.Int64Array([]int64{}),
				}

				if err := o.Database.NamedGetContext(ctx, ingestion, `
				    INSERT INTO fieldkit.ingestion
					(time, upload_id, user_id, device_id, generation, type, size, url, blocks, flags) VALUES
					(NOW(), :upload_id, :user_id, :device_id, :generation, :type, :size, :url, :blocks, :flags)
				    RETURNING *`, ingestion); err != nil {
					return err
				}

				o.Publisher.Publish(ctx, &messages.IngestionReceived{
					Time: ingestion.Time,
					ID:   ingestion.ID,
					URL:  saved.URL,
				})

				b := ingestion.Blocks[1] - ingestion.Blocks[0]
				o.Metrics.Ingested(int(b), saved.BytesRead)

				log.Infow("saved", "device_id", headers.FkDeviceId, "stream_id", saved.ID, "time", time.Since(startedAt).String(), "size", saved.BytesRead, "type", ingestion.Type, "ingestion_id", ingestion.ID, "generation", ingestion.Generation)
			} else {
				log.Infow("unsaved", "device_id", headers.FkDeviceId, "stream_id", saved.ID, "time", time.Since(startedAt).String(), "size", saved.BytesRead)
			}

			// TODO Give information.
			w.WriteHeader(http.StatusOK)
		}

		_ = claims
		_ = headers
		_ = log

		return nil
	}))

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		err := handler(ctx, w, req)
		if err != nil {
			log := Logger(ctx).Sugar()
			log.Errorw("completed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

type IncomingHeaders struct {
	ContentType     string
	ContentLength   int32
	MediaType       string
	MediaTypeParams map[string]string
	XForwardedFor   string
	FkType          string
	FkDeviceId      []byte
	FkGeneration    []byte
	FkBlocks        []int64
	FkFlags         []int64
}

func NewIncomingHeaders(req *http.Request) (*IncomingHeaders, error) {
	contentType := req.Header.Get(common.ContentTypeHeaderName)
	mediaType, mediaTypeParams, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, fmt.Errorf("invalid %s (%s)", common.ContentTypeHeaderName, contentType)
	}

	contentLengthString := req.Header.Get(common.ContentLengthHeaderName)
	contentLength, err := strconv.Atoi(contentLengthString)
	if err != nil {
		return nil, fmt.Errorf("invalid %s (%s)", common.ContentLengthHeaderName, contentLengthString)
	}

	if contentLength <= 0 {
		return nil, fmt.Errorf("invalid %s (%v)", common.ContentLengthHeaderName, contentLength)
	}

	deviceIdRaw := req.Header.Get(common.FkDeviceIdHeaderName)
	if len(deviceIdRaw) == 0 {
		return nil, fmt.Errorf("invalid %s (no header)", common.FkDeviceIdHeaderName)
	}

	generationRaw := req.Header.Get(common.FkGenerationHeaderName)
	if len(generationRaw) == 0 {
		return nil, fmt.Errorf("invalid %s (no header)", common.FkGenerationHeaderName)
	}

	typeRaw := req.Header.Get(common.FkTypeHeaderName)
	if len(typeRaw) == 0 {
		return nil, fmt.Errorf("invalid %s (no header)", common.FkTypeHeaderName)
	}

	deviceId, err := data.DecodeBinaryString(deviceIdRaw)
	if err != nil {
		return nil, fmt.Errorf("invalid %s (%v)", common.FkDeviceIdHeaderName, err)
	}

	generation, err := data.DecodeBinaryString(generationRaw)
	if err != nil {
		return nil, fmt.Errorf("invalid %s (%v)", common.FkGenerationHeaderName, err)
	}

	blocks, err := data.ParseBlocks(req.Header.Get(common.FkBlocksIdHeaderName))
	if err != nil {
		return nil, fmt.Errorf("invalid %s (%v)", common.FkBlocksIdHeaderName, err)
	}

	headers := &IncomingHeaders{
		ContentType:     contentType,
		ContentLength:   int32(contentLength),
		MediaType:       mediaType,
		MediaTypeParams: mediaTypeParams,
		XForwardedFor:   req.Header.Get(common.XForwardedForHeaderName),
		FkType:          typeRaw,
		FkDeviceId:      deviceId,
		FkGeneration:    generation,
		FkBlocks:        blocks,
	}

	return headers, nil
}

func errorHandling(middleware goa.Middleware, next goa.Handler) goa.Handler {
	return func(ctx context.Context, res http.ResponseWriter, req *http.Request) error {
		return middleware(next)(ctx, res, req)
	}
}

func authentication(middleware goa.Middleware, next goa.Handler) goa.Handler {
	return func(ctx context.Context, res http.ResponseWriter, req *http.Request) error {
		ctx = goa.WithRequiredScopes(ctx, []string{"api:access"})
		return middleware(next)(ctx, res, req)
	}
}
