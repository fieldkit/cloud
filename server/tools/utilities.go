package utilities

import (
	"context"
	"net/http"

	goaclient "github.com/goadesign/goa/client"

	fk "github.com/fieldkit/cloud/server/api/client"
)

func MapInt64(v, minV, maxV int, minRange, maxRange int64) int64 {
	return int64(MapFloat64(v, minV, maxV, float64(minRange), float64(maxRange)))
}

func MapFloat64(v, minV, maxV int, minRange, maxRange float64) float64 {
	if minV == maxV {
		return minRange
	}
	return float64(v-minV)/float64(maxV-minV)*(maxRange-minRange) + minRange
}

func NewClient(ctx context.Context, host, scheme string) (*fk.Client, error) {
	httpClient := newHTTPClient()
	c := fk.New(goaclient.HTTPClientDoer(httpClient))
	c.Host = host
	c.Scheme = scheme

	return c, nil
}

func CreateAndAuthenticate(ctx context.Context, host, scheme, email, password string) (*fk.Client, error) {
	c, err := NewClient(ctx, host, scheme)
	if err != nil {
		return nil, err
	}

	loginPayload := fk.LoginPayload{}
	loginPayload.Email = email
	loginPayload.Password = password
	res, err := c.LoginUser(ctx, fk.LoginUserPath(), &loginPayload)
	if err != nil {
		return nil, err
	}

	key := res.Header.Get("Authorization")
	jwtSigner := newJWTSigner(key, "%s")
	c.SetJWTSigner(jwtSigner)

	return c, nil
}

func newHTTPClient() *http.Client {
	return http.DefaultClient
}

func newJWTSigner(key, format string) goaclient.Signer {
	return &goaclient.APIKeySigner{
		SignQuery: false,
		KeyName:   "Authorization",
		KeyValue:  key,
		Format:    format,
	}

}
