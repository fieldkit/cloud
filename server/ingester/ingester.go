package ingester

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/goadesign/goa"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/common"
	"github.com/fieldkit/cloud/server/common/goahelpers"
	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/logging"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/messages"
)

var (
	ids = logging.NewIdGenerator()
)

type IngesterOptions struct {
	Database   *sqlxcache.DB
	Files      files.FileArchive
	Publisher  jobs.MessagePublisher
	Metrics    *logging.Metrics
	JwtHMACKey []byte
}

func Ingester(ctx context.Context, o *IngesterOptions) http.Handler {
	errorHandler := goahelpers.ErrorHandler(true)

	handler := useMiddleware(errorHandler, func(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
		startedAt := time.Now()

		userID, err := getUserID(ctx, o, req)
		if err != nil {
			serviceError := NewUnauthorizedServiceError()
			body, err := json.Marshal(serviceError)
			if err != nil {
				return err
			}
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(body)
			return nil
		}

		if req.Method == http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("{}"))
			return nil
		}

		log := Logger(ctx).Sugar().With("user_id", userID)

		o.Metrics.UserID(userID)

		headers, err := newIncomingHeaders(req)
		if err != nil {
			if err := writeInvalidHeaders(ctx, w, req); err != nil {
				return err
			}
			return nil
		}

		o.Metrics.IngestionDevice(headers.FkDeviceID)

		log.Infow("receiving", "device_id", headers.FkDeviceID, "blocks", headers.FkBlocks, "user_id", userID)

		metaMap := make(map[string]string)
		metaMap[common.FkDeviceIdHeaderName] = hex.EncodeToString(headers.FkDeviceID)
		metaMap[common.FkGenerationHeaderName] = hex.EncodeToString(headers.FkGenerationID)
		metaMap[common.FkBlocksHeaderName] = fmt.Sprintf("%v", headers.FkBlocks)
		metaMap[common.FkFlagsIdHeaderName] = fmt.Sprintf("%v", headers.FkFlags)

		saved, err := o.Files.Archive(ctx, headers.ContentType, metaMap, req.Body)
		if err != nil {
			return err
		}

		if saved.BytesRead != int(headers.ContentLength) {
			log.Warnw("size mismatch", "expected", headers.ContentLength, "actual", saved.BytesRead)
		}

		ingestion := &data.Ingestion{
			Time:         time.Now(),
			URL:          saved.URL,
			UploadID:     saved.Key,
			UserID:       userID,
			DeviceID:     headers.FkDeviceID,
			GenerationID: headers.FkGenerationID,
			Type:         headers.FkType,
			Size:         int64(saved.BytesRead),
			Blocks:       data.Int64Range(headers.FkBlocks),
			Flags:        pq.Int64Array([]int64{}),
		}

		ir, err := repositories.NewIngestionRepository(o.Database)
		if err != nil {
			return err
		}

		queuedID, err := ir.AddAndQueue(ctx, ingestion)
		if err != nil {
			return err
		}

		b := ingestion.Blocks[1] - ingestion.Blocks[0]

		o.Metrics.Ingested(int(b), saved.BytesRead)

		log.Infow("saved", "device_id", headers.FkDeviceID, "file_id", saved.Key, "time", time.Since(startedAt).String(), "size", saved.BytesRead,
			"type", ingestion.Type, "ingestion_id", ingestion.ID, "generation_id", ingestion.GenerationID, "user_id", userID,
			"device_name", headers.FkDeviceName, "blocks", headers.FkBlocks)

		if err := o.Publisher.Publish(ctx, &messages.IngestionReceived{
			QueuedID: queuedID,
			UserID:   userID,
		}); err != nil {
			log.Warnw("publishing", "err", err)
		}

		return writeSuccess(ctx, w, req, ingestion)
	})

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		if err := handler(ctx, w, req); err != nil {
			log := Logger(ctx).Sugar()
			log.Errorw("error", "error", err)
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
	FkDeviceID      []byte
	FkGenerationID  []byte
	FkDeviceName    string
	FkBlocks        []int64
	FkFlags         []int64
}

func newIncomingHeaders(req *http.Request) (*IncomingHeaders, error) {
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

	deviceIDRaw := req.Header.Get(common.FkDeviceIdHeaderName)
	if len(deviceIDRaw) == 0 {
		return nil, fmt.Errorf("invalid %s (no header)", common.FkDeviceIdHeaderName)
	}

	generationIDRaw := req.Header.Get(common.FkGenerationHeaderName)
	if len(generationIDRaw) == 0 {
		return nil, fmt.Errorf("invalid %s (no header)", common.FkGenerationHeaderName)
	}

	typeRaw := req.Header.Get(common.FkTypeHeaderName)
	if len(typeRaw) == 0 {
		return nil, fmt.Errorf("invalid %s (no header)", common.FkTypeHeaderName)
	}

	deviceID, err := data.DecodeBinaryString(deviceIDRaw)
	if err != nil {
		return nil, fmt.Errorf("invalid %s (%v)", common.FkDeviceIdHeaderName, err)
	}

	generationID, err := data.DecodeBinaryString(generationIDRaw)
	if err != nil {
		return nil, fmt.Errorf("invalid %s (%v)", common.FkGenerationHeaderName, err)
	}

	blocks, err := data.ParseBlocks(req.Header.Get(common.FkBlocksHeaderName))
	if err != nil {
		return nil, fmt.Errorf("invalid %s (%v)", common.FkBlocksHeaderName, err)
	}

	name := req.Header.Get(common.FkDeviceNameHeaderName)

	headers := &IncomingHeaders{
		ContentType:     contentType,
		ContentLength:   int32(contentLength),
		MediaType:       mediaType,
		MediaTypeParams: mediaTypeParams,
		XForwardedFor:   req.Header.Get(common.XForwardedForHeaderName),
		FkType:          typeRaw,
		FkDeviceID:      deviceID,
		FkGenerationID:  generationID,
		FkDeviceName:    name,
		FkBlocks:        blocks,
	}

	return headers, nil
}

type IngestionSuccessful struct {
	ID       int64  `json:"id"`
	UploadID string `json:"upload_id"`
}

func writeSuccess(ctx context.Context, w http.ResponseWriter, req *http.Request, ingestion *data.Ingestion) error {
	payload, err := json.Marshal(IngestionSuccessful{
		ID:       ingestion.ID,
		UploadID: ingestion.UploadID,
	})

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(payload)

	return nil
}

func writeInvalidHeaders(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	log := Logger(ctx).Sugar()

	log.Errorw("headers", "headers", req.Header)

	payload, err := json.Marshal(
		struct {
			Message string      `json:"message"`
			Headers interface{} `json:"headers"`
		}{
			Message: "invalid headers",
			Headers: req.Header,
		},
	)

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write(payload)

	return nil
}

func useMiddleware(middleware goa.Middleware, next goa.Handler) goa.Handler {
	return func(ctx context.Context, res http.ResponseWriter, req *http.Request) error {
		return middleware(next)(ctx, res, req)
	}
}

func getUserID(ctx context.Context, o *IngesterOptions, req *http.Request) (int32, error) {
	log := Logger(ctx).Sugar()

	header := req.Header.Get("Authorization")
	if header == "" {
		return 0, fmt.Errorf("no authorization header")
	}

	token := strings.ReplaceAll(header, "Bearer ", "")
	claims := make(jwtgo.MapClaims)
	_, err := jwtgo.ParseWithClaims(token, claims, func(t *jwtgo.Token) (interface{}, error) {
		return o.JwtHMACKey, nil
	})
	if err != nil {
		return 0, err
	}

	id := int32(claims["sub"].(float64))

	log.Infow("scopes", "scopes", claims["scopes"], "user_id", id)

	return id, nil
}

type ServiceError struct {
	ID     string `json:"id"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
	Status int32  `json:"status"`
}

func NewUnauthorizedServiceError() *ServiceError {
	return &ServiceError{
		ID:     newErrorID(),
		Code:   "jwt_security_error",
		Detail: "",
		Status: 401,
	}
}

// https://github.com/goadesign/goa/blob/master/error.go#L312
func newErrorID() string {
	b := make([]byte, 6)
	io.ReadFull(rand.Reader, b)
	return base64.StdEncoding.EncodeToString(b)
}
