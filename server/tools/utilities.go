package utilities

import (
	"context"
	"net/http"
)

type FkClient struct {
	scheme string
	host   string
	http   *http.Client
}

func CreateAndAuthenticate(ctx context.Context, host, scheme, email, password string) (*FkClient, error) {
	fkc := &FkClient{
		scheme: scheme,
		host:   host,
		http:   http.DefaultClient,
	}

	/*
		res, err := fkc.Login(ctx, email, password)
		if err != nil {
			return nil, err
		}

		key := res.Header.Get("Authorization")
	*/

	return fkc, nil
}
