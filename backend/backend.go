package backend

import (
	"errors"

	"upper.io/db.v2"
	"upper.io/db.v2/lib/sqlbuilder"
	"upper.io/db.v2/postgresql"

	"github.com/O-C-R/auth/id"
	"github.com/O-C-R/fieldkit/data"
)

var (
	NotFoundError = errors.New("not found")
)

func Err(err error) error {
	if err == db.ErrNoMoreRows {
		return NotFoundError
	}

	return err
}

type projectUser struct {
	ProjectID id.ID  `db:"project_id"`
	UserID    id.ID  `db:"user_id"`
	Role      string `db:"role"`
}

type Backend struct {
	database sqlbuilder.Database
}

func NewBackend(url string) (*Backend, error) {
	postgresqlURL, err := postgresql.ParseURL(url)
	if err != nil {
		return nil, err
	}

	database, err := postgresql.Open(postgresqlURL)
	if err != nil {
		return nil, err
	}

	return &Backend{
		database: database,
	}, nil
}

func (b *Backend) Ping() error {
	return b.database.Ping()
}

func (b *Backend) AddUser(user *data.User) error {
	_, err := b.database.Collection("admin.user").Insert(user)
	return Err(err)
}

func (b *Backend) SetUserValidByID(userID id.ID, valid bool) error {
	return Err(b.database.Collection("admin.user").Find(userID).Update(map[string]interface{}{"valid": valid}))
}

func (b *Backend) UserByID(userID id.ID) (*data.User, error) {
	user := &data.User{}
	if err := b.database.Collection("admin.user").Find(userID).One(user); err != nil {
		return nil, Err(err)
	}

	return user, nil
}

func (b *Backend) UsersByID(userIDs ...id.ID) ([]*data.User, error) {
	users := []*data.User{}
	if err := b.database.Collection("admin.user").Find(db.Cond{"id": userIDs}).All(&users); err != nil {
		return nil, Err(err)
	}

	return users, nil
}

func (b *Backend) UserByUsername(username string) (*data.User, error) {
	user := &data.User{}
	if err := b.database.Collection("admin.user").Find(db.Cond{"username": username}).One(user); err != nil {
		return nil, Err(err)
	}

	return user, nil
}

func (b *Backend) UserUsernameInUse(username string) (bool, error) {
	n, err := b.database.Collection("admin.user").Find(db.Cond{"username": username}).Count()
	if err != nil {
		return false, Err(err)
	}

	return n > 0, nil
}

func (b *Backend) UserByEmail(email string) (*data.User, error) {
	user := &data.User{}
	if err := b.database.Collection("admin.user").Find(db.Cond{"email": email}).One(user); err != nil {
		return nil, Err(err)
	}

	return user, nil
}

func (b *Backend) UserEmailInUse(email string) (bool, error) {
	n, err := b.database.Collection("admin.user").Find(db.Cond{"email": email}).Count()
	if err != nil {
		return false, Err(err)
	}

	return n > 0, nil
}

func (b *Backend) UserByValidationTokenID(validationTokenID id.ID) (*data.User, error) {
	user := &data.User{}
	if err := b.database.Iterator(`
		SELECT u.* FROM admin.user AS u
			JOIN admin.validation_token AS v ON v.user_id = u.id
				WHERE v.id = $1
		`, validationTokenID).One(user); err != nil {
		return nil, Err(err)
	}

	return user, nil
}

func (b *Backend) DeleteUserByID(userID id.ID) error {
	return b.database.Collection("admin.user").Find(userID).Delete()
}

func (b *Backend) AddValidationToken(validationToken *data.ValidationToken) error {
	_, err := b.database.Collection("admin.user_validation_token").Insert(validationToken)
	return Err(err)
}

func (b *Backend) ValidationTokenByID(validationTokenID id.ID) (*data.ValidationToken, error) {
	validationToken := &data.ValidationToken{}
	if err := b.database.Collection("admin.user_validation_token").Find(validationTokenID).One(validationToken); err != nil {
		return nil, Err(err)
	}

	return validationToken, nil
}

func (b *Backend) DeleteValidationTokenByID(validationTokenID id.ID) error {
	return Err(b.database.Collection("admin.user_validation_token").Find(validationTokenID).Delete())
}

func (b *Backend) AddProject(project *data.Project) error {
	_, err := b.database.Collection("admin.project").Insert(project)
	return err
}

func (b *Backend) AddProjectWithOwner(project *data.Project, ownerID id.ID) error {
	return Err(b.database.Tx(func(sess sqlbuilder.Tx) error {
		if _, err := sess.Collection("admin.project").Insert(project); err != nil {
			return err
		}

		if _, err := sess.Collection("admin.project_user").Insert(&projectUser{
			ProjectID: project.ID,
			UserID:    ownerID,
			Role:      "owner",
		}); err != nil {
			return err
		}

		return nil
	}))
}

func (b *Backend) ProjectByID(projectID id.ID) (*data.Project, error) {
	project := &data.Project{}
	if err := b.database.Collection("admin.project").Find(projectID).One(project); err != nil {
		return nil, Err(err)
	}

	return project, nil
}

func (b *Backend) ProjectBySlug(slug string) (*data.Project, error) {
	project := &data.Project{}
	if err := b.database.Collection("admin.project").Find(db.Cond{"slug": slug}).One(project); err != nil {
		return nil, Err(err)
	}

	return project, nil
}

func (b *Backend) DeleteProjectByID(projectID id.ID) error {
	return Err(b.database.Tx(func(sess sqlbuilder.Tx) error {
		if err := sess.Collection("admin.project_user").Find(db.Cond{"project_id": projectID}).Delete(); err != nil {
			return err
		}

		if err := sess.Collection("admin.expedition").Find(db.Cond{"project_id": projectID}).Delete(); err != nil {
			return err
		}

		if err := sess.Collection("admin.project").Find(projectID).Delete(); err != nil {
			return err
		}

		return nil
	}))
}

func (b *Backend) AddExpedition(expedition *data.Expedition) error {
	_, err := b.database.Collection("admin.expedition").Insert(expedition)
	return Err(err)
}

func (b *Backend) ExpeditionByID(expeditionID id.ID) (*data.Expedition, error) {
	expedition := &data.Expedition{}
	if err := b.database.Collection("admin.expedition").Find(expeditionID).One(expedition); err != nil {
		return nil, Err(err)
	}

	return expedition, nil
}

func (b *Backend) ExpeditionByProjectSlugAndSlug(projectSlug, slug string) (*data.Expedition, error) {
	expedition := &data.Expedition{}
	if err := b.database.Iterator(`
		SELECT e.* FROM admin.expedition AS e
			JOIN admin.project AS p ON p.id = e.project_id
				WHERE p.slug = $1 AND e.slug = $2
		`, projectSlug, slug).One(expedition); err != nil {
		return nil, Err(err)
	}

	return expedition, nil
}

func (b *Backend) DeleteExpeditionByID(expeditionID id.ID) error {
	return Err(b.database.Collection("admin.expedition").Find(expeditionID).Delete())
}
