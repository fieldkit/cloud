package api

import (
	"context"
	"fmt"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"goa.design/goa/v3/security"

	goahttp "goa.design/goa/v3/http"
	httpmdlwr "goa.design/goa/v3/http/middleware"
	"goa.design/goa/v3/middleware"

	"github.com/fieldkit/cloud/server/logging"

	testsvr "github.com/fieldkit/cloud/server/api/gen/http/test/server"
	test "github.com/fieldkit/cloud/server/api/gen/test"

	taskssvr "github.com/fieldkit/cloud/server/api/gen/http/tasks/server"
	tasks "github.com/fieldkit/cloud/server/api/gen/tasks"

	modulessvr "github.com/fieldkit/cloud/server/api/gen/http/modules/server"
	modules "github.com/fieldkit/cloud/server/api/gen/modules"
)

func CreateGoaV3Handler(ctx context.Context, options *ControllerOptions) http.Handler {
	debug := false

	testSvc := NewTestSevice(ctx, options)
	testEndpoints := test.NewEndpoints(testSvc)

	tasksSvc := NewTasksService(ctx, options)
	tasksEndpoints := tasks.NewEndpoints(tasksSvc)

	modulesSvc := NewModulesService(ctx, options)
	modulesEndpoints := modules.NewEndpoints(modulesSvc)

	// Provide the transport specific request decoder and response encoder.
	// The goa http package has built-in support for JSON, XML and gob.
	// Other encodings can be used by providing the corresponding functions,
	// see goa.design/implement/encoding.
	dec := goahttp.RequestDecoder
	enc := goahttp.ResponseEncoder
	mux := goahttp.NewMuxer()

	eh := errorHandler()

	tasksServer := taskssvr.New(tasksEndpoints, mux, dec, enc, eh, nil)
	testServer := testsvr.New(testEndpoints, mux, dec, enc, eh, nil)
	modulesServer := modulessvr.New(modulesEndpoints, mux, dec, enc, eh, nil)
	if debug {
		servers := goahttp.Servers{
			tasksServer,
			testServer,
			modulesServer,
		}
		servers.Use(httpmdlwr.Debug(mux, os.Stdout))
	}

	taskssvr.Mount(mux, tasksServer)
	testsvr.Mount(mux, testServer)
	modulessvr.Mount(mux, modulesServer)

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
		return ctx, ErrInvalidToken
	}

	if claims["scopes"] == nil {
		return ctx, ErrInvalidTokenScopes
	}
	scopes, ok := claims["scopes"].([]interface{})
	if !ok {
		return ctx, ErrInvalidTokenScopes
	}

	scopesInToken := make([]string, len(scopes))
	for _, scp := range scopes {
		scopesInToken = append(scopesInToken, scp.(string))
	}
	if err := a.Scheme.Validate(scopesInToken); err != nil {
		return ctx, ErrInvalidTokenScopes
	}

	newCtx := logging.WithUserID(ctx, fmt.Sprintf("%v", claims["sub"]))

	return newCtx, nil
}
