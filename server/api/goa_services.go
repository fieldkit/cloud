package api

import (
	"context"
	"fmt"
	"net/http"

	goa "goa.design/goa/v3/pkg"

	jwt "github.com/dgrijalva/jwt-go"
	"goa.design/goa/v3/security"

	goahttp "goa.design/goa/v3/http"
	"goa.design/goa/v3/middleware"

	"github.com/fieldkit/cloud/server/logging"

	testSvr "github.com/fieldkit/cloud/server/api/gen/http/test/server"
	test "github.com/fieldkit/cloud/server/api/gen/test"

	tasksSvr "github.com/fieldkit/cloud/server/api/gen/http/tasks/server"
	tasks "github.com/fieldkit/cloud/server/api/gen/tasks"

	modulesSvr "github.com/fieldkit/cloud/server/api/gen/http/modules/server"
	modules "github.com/fieldkit/cloud/server/api/gen/modules"

	following "github.com/fieldkit/cloud/server/api/gen/following"
	followingSvr "github.com/fieldkit/cloud/server/api/gen/http/following/server"
)

func LogErrors() func(goa.Endpoint) goa.Endpoint {
	return func(e goa.Endpoint) goa.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			response, err = e(ctx, request)
			if err != nil {
				id := newErrorID()
				log := logging.Logger(ctx).Sugar()
				log.Errorw("error", "error", err, "error_id", id)
			}
			return response, err
		}
	}
}

func CreateGoaV3Handler(ctx context.Context, options *ControllerOptions) http.Handler {
	testSvc := NewTestSevice(ctx, options)
	testEndpoints := test.NewEndpoints(testSvc)

	tasksSvc := NewTasksService(ctx, options)
	tasksEndpoints := tasks.NewEndpoints(tasksSvc)

	modulesSvc := NewModulesService(ctx, options)
	modulesEndpoints := modules.NewEndpoints(modulesSvc)

	followingSvc := NewFollowingService(ctx, options)
	followingEndpoints := following.NewEndpoints(followingSvc)

	logErrors := LogErrors()

	modulesEndpoints.Use(logErrors)
	tasksEndpoints.Use(logErrors)
	testEndpoints.Use(logErrors)
	followingEndpoints.Use(logErrors)

	// Provide the transport specific request decoder and response encoder.
	// The goa http package has built-in support for JSON, XML and gob.
	// Other encodings can be used by providing the corresponding functions,
	// see goa.design/implement/encoding.
	dec := goahttp.RequestDecoder
	enc := goahttp.ResponseEncoder
	mux := goahttp.NewMuxer()

	eh := errorHandler()

	tasksServer := tasksSvr.New(tasksEndpoints, mux, dec, enc, eh, nil)
	testServer := testSvr.New(testEndpoints, mux, dec, enc, eh, nil)
	modulesServer := modulesSvr.New(modulesEndpoints, mux, dec, enc, eh, nil)
	followingServer := followingSvr.New(followingEndpoints, mux, dec, enc, eh, nil)

	tasksSvr.Mount(mux, tasksServer)
	testSvr.Mount(mux, testServer)
	modulesSvr.Mount(mux, modulesServer)
	followingSvr.Mount(mux, followingServer)

	log := Logger(ctx).Sugar()

	for _, m := range testServer.Mounts {
		log.Infow("mounted", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}
	for _, m := range tasksServer.Mounts {
		log.Infow("mounted", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}
	for _, m := range modulesServer.Mounts {
		log.Infow("mounted", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}
	for _, m := range followingServer.Mounts {
		log.Infow("mounted", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}

	return mux
}

func errorHandler() func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		log := Logger(ctx).Sugar()
		id := ctx.Value(middleware.RequestIDKey).(string)
		w.Write([]byte("[" + id + "] encoding: " + err.Error()))
		log.Errorw("fatal", "id", id, "message", err.Error())
	}
}

type AuthAttempt struct {
	Token         string
	Scheme        *security.JWTScheme
	Key           []byte
	InvalidToken  error
	InvalidScopes error
}

func Authenticate(ctx context.Context, a AuthAttempt) (context.Context, error) {
	claims := make(jwt.MapClaims)
	_, err := jwt.ParseWithClaims(a.Token, claims, func(t *jwt.Token) (interface{}, error) {
		return a.Key, nil
	})
	if err != nil {
		return ctx, a.InvalidToken
	}

	if claims["scopes"] == nil {
		return ctx, a.InvalidScopes
	}
	scopes, ok := claims["scopes"].([]interface{})
	if !ok {
		return ctx, a.InvalidScopes
	}

	scopesInToken := make([]string, len(scopes))
	for _, scp := range scopes {
		scopesInToken = append(scopesInToken, scp.(string))
	}
	if err := a.Scheme.Validate(scopesInToken); err != nil {
		return ctx, a.InvalidScopes
	}

	withClaims := addClaimsToContext(ctx, claims)
	withLogging := logging.WithUserID(withClaims, fmt.Sprintf("%v", claims["sub"]))

	return withLogging, nil
}
