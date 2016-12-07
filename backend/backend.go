package backend

import (
	"crypto/tls"
	"errors"
	"net"
	"time"

	"github.com/O-C-R/auth/id"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/O-C-R/fieldkit/data"
)

var (
	DuplicateKeyError = errors.New("duplicate key")
	NotFoundError     = errors.New("not found")
	dialer            = &net.Dialer{
		// Default timeout that mgo uses.
		Timeout:   10 * time.Second,
		KeepAlive: time.Second,
	}
	indexes = map[string][]mgo.Index{
		"user": []mgo.Index{
			mgo.Index{
				Key:    []string{"email"},
				Unique: true,
			},
			mgo.Index{
				Key:    []string{"username"},
				Unique: true,
			},
			mgo.Index{
				Key:    []string{"validation_token"},
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

func newBackend(url string, tls bool) (*Backend, error) {
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

	return backend, nil
}

func NewBackend(url string, tls bool) (*Backend, error) {
	backend, err := newBackend(url, tls)
	if err != nil {
		return nil, err
	}

	if err := backend.init(); err != nil {
		return nil, err
	}

	return backend, nil
}

func NewTestBackend() (*Backend, error) {
	backend, err := newBackend("mongodb://localhost/test", false)
	if err != nil {
		return nil, err
	}

	if err := backend.dropDatabase(); err != nil {
		return nil, err
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

	session.ResetIndexCache()
	for collection, indexes := range indexes {
		for _, index := range indexes {
			if err := session.DB("").C(collection).EnsureIndex(index); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Backend) dropDatabase() error {
	session := b.newSession()
	defer session.Close()

	return session.DB("").DropDatabase()
}

func (b *Backend) Ping() error {
	session := b.newSession()
	defer session.Close()

	return session.Ping()
}

func (b *Backend) AddUser(user *data.User) error {
	session := b.newSession()
	defer session.Close()

	err := session.DB("").C("user").Insert(user)
	if mgo.IsDup(err) {
		return DuplicateKeyError
	}

	return err
}

func (b *Backend) UpdateUser(user *data.User) error {
	session := b.newSession()
	defer session.Close()

	err := session.DB("").C("user").Update(bson.M{"_id": user.ID}, user)
	if mgo.IsDup(err) {
		return DuplicateKeyError
	}

	if err == mgo.ErrNotFound {
		return NotFoundError
	}

	return err
}

func (b *Backend) UserByID(userID id.ID) (*data.User, error) {
	session := b.newSession()
	defer session.Close()

	user := &data.User{}
	err := session.DB("").C("user").Find(bson.M{"_id": userID}).One(user)
	if err == mgo.ErrNotFound {
		return nil, NotFoundError
	}

	if err != nil {
		return nil, err
	}

	return user, err
}

func (b *Backend) UserByEmail(email string) (*data.User, error) {
	session := b.newSession()
	defer session.Close()

	user := &data.User{}
	err := session.DB("").C("user").Find(bson.M{"email": email}).One(user)
	if err == mgo.ErrNotFound {
		return nil, NotFoundError
	}

	if err != nil {
		return nil, err
	}

	return user, err
}

func (b *Backend) UserByUsername(username string) (*data.User, error) {
	session := b.newSession()
	defer session.Close()

	user := &data.User{}
	err := session.DB("").C("user").Find(bson.M{"username": username}).One(user)
	if err == mgo.ErrNotFound {
		return nil, NotFoundError
	}

	if err != nil {
		return nil, err
	}

	return user, err
}

func (b *Backend) UserByValidationToken(token id.ID) (*data.User, error) {
	session := b.newSession()
	defer session.Close()

	user := &data.User{}
	err := session.DB("").C("user").Find(bson.M{"validation_token": token}).One(user)
	if err == mgo.ErrNotFound {
		return nil, NotFoundError
	}

	if err != nil {
		return nil, err
	}

	return user, err
}

func (b *Backend) UserEmailInUse(email string) (bool, error) {
	session := b.newSession()
	defer session.Close()

	n, err := session.DB("").C("user").Find(bson.M{"email": email}).Count()
	if err != nil {
		return false, err
	}

	return n > 0, nil
}

func (b *Backend) UserUsernameInUse(username string) (bool, error) {
	session := b.newSession()
	defer session.Close()

	n, err := session.DB("").C("user").Find(bson.M{"username": username}).Count()
	if err != nil {
		return false, err
	}

	return n > 0, nil
}
