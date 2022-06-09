package backend

import (
	_ "context"

	_ "github.com/lib/pq"
	_ "github.com/paulmach/go.geo"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	_ "github.com/fieldkit/cloud/server/data"
)

type Backend struct {
	db  *sqlxcache.DB
	url string
}

func OpenDatabase(url string) (*sqlxcache.DB, error) {
	return sqlxcache.Open("postgres", url)
}

func New(url string) (*Backend, error) {
	db, err := sqlxcache.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &Backend{
		url: url,
		db:  db,
	}, nil
}

func (b *Backend) URL() string {
	return b.url
}
