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

func (b *Backend) ProjectSlugInUse(slug string) (bool, error) {
	n, err := b.database.Collection("admin.project").Find(db.Cond{"slug": slug}).Count()
	if err != nil {
		return false, Err(err)
	}

	return n > 0, nil
}

func (b *Backend) Projects() ([]*data.Project, error) {
	projects := []*data.Project{}
	if err := b.database.Collection("admin.project").Find().All(&projects); err != nil {
		return nil, Err(err)
	}

	return projects, nil
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

func (b *Backend) ExpeditionSlugInUse(projectID id.ID, slug string) (bool, error) {
	n, err := b.database.Collection("admin.expedition").Find(db.Cond{"project_id": projectID, "slug": slug}).Count()
	if err != nil {
		return false, Err(err)
	}

	return n > 0, nil
}

func (b *Backend) ExpeditionByID(expeditionID id.ID) (*data.Expedition, error) {
	expedition := &data.Expedition{}
	if err := b.database.Collection("admin.expedition").Find(expeditionID).One(expedition); err != nil {
		return nil, Err(err)
	}

	return expedition, nil
}

func (b *Backend) ExpeditionsByProjectSlug(projectSlug string) ([]*data.Expedition, error) {
	expeditions := []*data.Expedition{}
	if err := b.database.Iterator(`
		SELECT e.* FROM admin.expedition AS e
			JOIN admin.project AS p ON p.id = e.project_id
				WHERE p.slug = $1
		`, projectSlug).All(&expeditions); err != nil {
		return nil, Err(err)
	}

	return expeditions, nil
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

func (b *Backend) ProjectUserRoleByProjectSlugAndUserID(projectSlug string, userID id.ID) (string, error) {
	role := ""
	if err := b.database.Iterator(`
		SELECT u.role FROM admin.project_user AS u
			JOIN admin.project AS p ON p.id = u.project_id
				WHERE p.slug = $1 AND u.user_id = $2
		`, projectSlug, userID).ScanOne(&role); err != nil {
		return "", Err(err)
	}

	return role, nil
}

func (b *Backend) AddInput(input *data.Input) error {
	_, err := b.database.Collection("admin.input").Insert(input)
	return Err(err)
}

func (b *Backend) InputsByProjectSlugAndExpeditionSlug(projectSlug, expeditionSlug string) ([]*data.Input, error) {
	inputs := []*data.Input{}
	if err := b.database.Iterator(`
		SELECT i.* FROM admin.input AS i
			JOIN admin.expedition AS e ON e.id = i.expedition_id
			JOIN admin.project AS p ON p.id = e.project_id
				WHERE p.slug = $1 AND e.slug = $2
		`, projectSlug, expeditionSlug).All(&inputs); err != nil {
		return nil, Err(err)
	}

	return inputs, nil
}

func (b *Backend) InputByID(inputID id.ID) (*data.Input, error) {
	input := &data.Input{}
	if err := b.database.Collection("admin.input").Find(inputID).One(input); err != nil {
		return nil, Err(err)
	}

	return input, nil
}

func (b *Backend) AddRequest(request *data.Request) error {
	_, err := b.database.Collection("data.request").Insert(request)
	return Err(err)
}

func (b *Backend) AddInvite(invite *data.Invite) error {
	_, err := b.database.Collection("admin.invite").Insert(invite)
	return Err(err)
}

func (b *Backend) InviteByID(inviteID id.ID) (*data.Invite, error) {
	invite := &data.Invite{}
	if err := b.database.Collection("admin.invite").Find(inviteID).One(invite); err != nil {
		return nil, Err(err)
	}

	return invite, nil
}

func (b *Backend) DeleteInviteByID(inviteID id.ID) error {
	return Err(b.database.Collection("admin.invite").Find(inviteID).Delete())
}

func (b *Backend) AddAuthToken(authToken *data.AuthToken) error {
	_, err := b.database.Collection("admin.expedition_auth_token").Insert(authToken)
	return Err(err)
}

func (b *Backend) AuthTokenByID(authTokenID id.ID) (*data.AuthToken, error) {
	authToken := &data.AuthToken{}
	if err := b.database.Collection("admin.expedition_auth_token").Find(authTokenID).One(authToken); err != nil {
		return nil, Err(err)
	}

	return authToken, nil
}

func (b *Backend) DeleteAuthTokenByID(authTokenID id.ID) error {
	return Err(b.database.Collection("admin.expedition_auth_token").Find(authTokenID).Delete())
}

func (b *Backend) InputSlugInUse(expeditionID id.ID, slug string) (bool, error) {
	n, err := b.database.Collection("admin.input").Find(db.Cond{"expedition_id": expeditionID, "slug": slug}).Count()
	if err != nil {
		return false, Err(err)
	}

	return n > 0, nil
}

func (b *Backend) AuthTokensByProjectSlugAndExpeditionSlug(projectSlug, expeditionSlug string) ([]*data.AuthToken, error) {
	authTokens := []*data.AuthToken{}
	if err := b.database.Iterator(`
		SELECT a.* FROM admin.expedition_auth_token AS a
			JOIN admin.expedition AS e ON e.id = a.expedition_id
			JOIN admin.project AS p ON p.id = e.project_id
				WHERE p.slug = $1 AND e.slug = $2
		`, projectSlug, expeditionSlug).All(&authTokens); err != nil {
		return nil, Err(err)
	}

	return authTokens, nil
}

func (b *Backend) AuthTokenByInputIDAndID(inputID, authTokenID id.ID) (*data.AuthToken, error) {
	authToken := &data.AuthToken{}
	if err := b.database.Iterator(`
		SELECT a.* FROM admin.expedition_auth_token AS a
			JOIN admin.expedition AS e ON e.id = a.expedition_id
			JOIN admin.input AS i ON i.expedition_id = e.id
				WHERE i.id = $1 AND a.id = $2
		`, inputID, authTokenID).One(authToken); err != nil {
		return nil, Err(err)
	}

	return authToken, nil
}

func (b *Backend) AddMessage(message *data.Message) error {
	_, err := b.database.Collection("data.message").Insert(message)
	return Err(err)
}

func (b *Backend) AddDocument(document *data.Document) error {
	_, err := b.database.Collection("data.document").Insert(document)
	return Err(err)
}

func (b *Backend) DocumentsByProjectSlugAndExpeditionSlug(projectSlug, expeditionSlug string) ([]*data.Document, error) {
	documents := []*data.Document{}
	if err := b.database.Iterator(`
		SELECT d.* FROM data.document AS d
			JOIN admin.input AS i ON i.id = d.input_id
			JOIN admin.expedition AS e ON e.id = i.expedition_id
			JOIN admin.project AS p ON p.id = e.project_id
				WHERE p.slug = $1 AND e.slug = $2
		`, projectSlug, expeditionSlug).All(&documents); err != nil {
		return nil, Err(err)
	}

	return documents, nil
}
