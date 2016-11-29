package backend

import (
	"crypto/tls"
	"net"
	"time"

	"gopkg.in/mgo.v2"
)

var (
	dialer = &net.Dialer{
		// Default timeout that mgo uses.
		Timeout:   10 * time.Second,
		KeepAlive: time.Second,
	}
	indexes = map[string][]mgo.Index{
		"user": []mgo.Index{
			mgo.Index{
				Key:    []string{"id"},
				Unique: true,
			},
			mgo.Index{
				Key:    []string{"email"},
				Unique: true,
			},
		},
	}
)

func dialServerTLS(addr *mgo.ServerAddr) (net.Conn, error) {
	return tls.DialWithDialer(dialer, "tcp", addr.String(), nil)
}

type Backend struct {
	session *mgo.Session
}

func NewBackend(url string, tls bool) (*Backend, error) {
	dialInfo, err := mgo.ParseURL(url)
	if err != nil {
		return nil, err
	}

	if tls {
		dialInfo.DialServer = dialServerTLS
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}

	backend := &Backend{
		session: session,
	}

	if err := backend.init(); err != nil {
		return nil, err
	}

	return backend, nil
}

func (b *Backend) newSession() *mgo.Session {
	return b.session.Copy()
}

func (b *Backend) init() error {
	session := b.newSession()
	defer session.Close()

	for collection, indexes := range indexes {
		for _, index := range indexes {
			if err := session.DB("").C(collection).EnsureIndex(index); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Backend) Ping() error {
	session := b.newSession()
	defer session.Close()

	return session.Ping()
}
