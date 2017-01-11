package backend

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/O-C-R/fieldkit/data"
)

const (
	testBackendURL = "postgres://localhost/test?sslmode=disable"
)

var (
	invalidDriverError      = errors.New("invalid driver")
	schemaDirUndefinedError = errors.New("SCHEMA_DIR undefined")
)

func NewTestBackend() (*Backend, error) {
	backend, err := NewBackend(testBackendURL)
	if err != nil {
		return nil, err
	}

	database, ok := backend.database.Driver().(*sql.DB)
	if !ok {
		return nil, invalidDriverError
	}

	schemaDir := os.Getenv("SCHEMA_DIR")
	if schemaDir == "" {
		return nil, schemaDirUndefinedError
	}

	filenames, err := ioutil.ReadDir(schemaDir)
	if err != nil {
		return nil, err
	}

	for _, filename := range filenames {
		file, err := os.Open(filepath.Join(schemaDir, filename.Name()))
		if err != nil {
			return nil, err
		}

		query, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		if _, err := database.Exec(string(query)); err != nil {
			return nil, err
		}

		if err := file.Close(); err != nil {
			return nil, err
		}
	}

	return backend, nil
}

func TestBackend(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	if err := backend.Ping(); err != nil {
		t.Fatal(err)
	}
}

func NewTestUser(backend *Backend) (*data.User, error) {
	user, err := data.NewUser("test@ocr.nyc", "test", "password")
	if err != nil {
		return nil, err
	}

	user.Valid = true
	if err := backend.AddUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func NewTestValidationToken(backend *Backend) (*data.User, *data.ValidationToken, error) {
	user, err := NewTestUser(backend)
	if err != nil {
		return nil, nil, err
	}

	validationToken, err := data.NewValidationToken(user.ID)
	if err != nil {
		return nil, nil, err
	}

	if err := backend.AddValidationToken(validationToken); err != nil {
		return nil, nil, err
	}

	return user, validationToken, nil
}

func NewTestProject(backend *Backend) (*data.User, *data.Project, error) {
	user, err := NewTestUser(backend)
	if err != nil {
		return nil, nil, err
	}

	name, slug := data.Name("Test Project")
	project, err := data.NewProject(name, slug)
	if err != nil {
		return nil, nil, err
	}

	if err := backend.AddProjectWithOwner(project, user.ID); err != nil {
		return nil, nil, err
	}

	return user, project, err
}

func NewTestExpedition(backend *Backend) (*data.User, *data.Project, *data.Expedition, error) {
	user, project, err := NewTestProject(backend)
	if err != nil {
		return nil, nil, nil, err
	}

	name, slug := data.Name("Test Expedition")
	expedition, err := data.NewExpedition(project.ID, name, slug)
	if err != nil {
		return nil, nil, nil, err
	}

	if err := backend.AddExpedition(expedition); err != nil {
		return nil, nil, nil, err
	}

	return user, project, expedition, nil
}

func TestBackendAddUser(t *testing.T) {
	backend, err := NewBackend(testBackendURL)
	if err != nil {
		t.Fatal(err)
	}

	user, err := NewTestUser(backend)
	if err != nil {
		t.Fatal(err)
	}

	returnedUser, err := backend.UserByID(user.ID)
	if err != nil {
		t.Error(err)
	}

	if returnedUser.ID != user.ID {
		t.Error("unexpexted user ID")
	}

	if err := backend.DeleteUserByID(user.ID); err != nil {
		t.Error(err)
	}
}

func TestBackendSetUserValidByID(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	user, err := NewTestUser(backend)
	if err != nil {
		t.Fatal(err)
	}

	if err := backend.SetUserValidByID(user.ID, false); err != nil {
		t.Error(err)
	}
}

func TestBackendUserByID(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	user, err := NewTestUser(backend)
	if err != nil {
		t.Fatal(err)
	}

	returnedUser, err := backend.UserByID(user.ID)
	if err != nil {
		t.Fatal(err)
	}

	if returnedUser.ID != user.ID {
		t.Error("incorrect returned user ID")
	}
}

func TestBackendUsersByID(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	user0, err := data.NewUser("test0@ocr.nyc", "test0", "password")
	if err != nil {
		t.Fatal(err)
	}

	if err := backend.AddUser(user0); err != nil {
		t.Fatal(err)
	}

	user1, err := data.NewUser("test1@ocr.nyc", "test1", "password")
	if err != nil {
		t.Fatal(err)
	}

	if err := backend.AddUser(user1); err != nil {
		t.Fatal(err)
	}

	user2, err := data.NewUser("test2@ocr.nyc", "test2", "password")
	if err != nil {
		t.Fatal(err)
	}

	if err := backend.AddUser(user2); err != nil {
		t.Fatal(err)
	}

	returnedUsers, err := backend.UsersByID(user0.ID, user1.ID, user2.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(returnedUsers) != 3 {
		t.Errorf("unexpected number of user IDs, %d", len(returnedUsers))
		t.Log(returnedUsers[0])
	}

	for _, user := range returnedUsers {
		if !(user.ID == user0.ID || user.ID == user1.ID || user.ID == user2.ID) {
			t.Error("unexpected user ID")
		}
	}
}

func TestBackendUserByUsername(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	user, err := NewTestUser(backend)
	if err != nil {
		t.Fatal(err)
	}

	returnedUser, err := backend.UserByUsername(user.Username)
	if err != nil {
		t.Fatal(err)
	}

	if returnedUser.ID != user.ID {
		t.Error("incorrect returned user ID")
	}
}

func TestBackendUserUsernameInUse(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	user, err := NewTestUser(backend)
	if err != nil {
		t.Fatal(err)
	}

	usernameInUse, err := backend.UserUsernameInUse(user.Username)
	if err != nil {
		t.Fatal(err)
	}

	if !usernameInUse {
		t.Error("expected username in use")
	}
}

func TestBackendUserByEmail(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	user, err := NewTestUser(backend)
	if err != nil {
		t.Fatal(err)
	}

	returnedUser, err := backend.UserByEmail(user.Email)
	if err != nil {
		t.Fatal(err)
	}

	if returnedUser.ID != user.ID {
		t.Error("incorrect returned user ID")
	}
}

func TestBackendUserEmailEmailInUse(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	user, err := NewTestUser(backend)
	if err != nil {
		t.Fatal(err)
	}

	emailInUse, err := backend.UserEmailInUse(user.Email)
	if err != nil {
		t.Fatal(err)
	}

	if !emailInUse {
		t.Error("expected email in use")
	}
}

func TestBackendAddValidationToken(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	if _, _, err := NewTestValidationToken(backend); err != nil {
		t.Fatal(err)
	}
}

func TestBackendValidationTokenByID(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	_, validationToken, err := NewTestValidationToken(backend)
	if err != nil {
		t.Fatal(err)
	}

	returnedValidationToken, err := backend.ValidationTokenByID(validationToken.ID)
	if err != nil {
		t.Error(err)
	}

	if returnedValidationToken.ID != validationToken.ID {
		t.Error("unexpexted user validation token ID")
	}
}

func TestDeleteValidationTokenByID(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	_, validationToken, err := NewTestValidationToken(backend)
	if err != nil {
		t.Fatal(err)
	}

	if err := backend.DeleteValidationTokenByID(validationToken.ID); err != nil {
		t.Fatal(err)
	}

	_, err = backend.ValidationTokenByID(validationToken.ID)
	if err == nil {
		t.Fatal("validation token not deleted")
	}

	if err != NotFoundError {
		t.Fatal(err)
	}
}

func TestBackendAddProject(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	if _, _, err := NewTestProject(backend); err != nil {
		t.Fatal(err)
	}
}

func TestBackendProjectByID(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	_, project, err := NewTestProject(backend)
	if err != nil {
		t.Fatal(err)
	}

	returnedProject, err := backend.ProjectByID(project.ID)
	if err != nil {
		t.Error(err)
	}

	if returnedProject.ID != project.ID {
		t.Error("unexpexted project ID")
	}
}

func TestBackendDeleteProjectByID(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	_, project, expedition, err := NewTestExpedition(backend)
	if err != nil {
		t.Fatal(err)
	}

	if err := backend.DeleteProjectByID(project.ID); err != nil {
		t.Fatal(err)
	}

	_, err = backend.ProjectByID(project.ID)
	if err == nil {
		t.Fatal("project not deleted")
	}

	if err != NotFoundError {
		t.Fatal(err)
	}

	_, err = backend.ExpeditionByID(expedition.ID)
	if err == nil {
		t.Fatal("project expedition not deleted")
	}

	if err != NotFoundError {
		t.Fatal(err)
	}
}

func TestBackendProjectBySlug(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	_, project, err := NewTestProject(backend)
	if err != nil {
		t.Fatal(err)
	}

	returnedProject, err := backend.ProjectBySlug(project.Slug)
	if err != nil {
		t.Error(err)
	}

	if returnedProject.ID != project.ID {
		t.Error("unexpected project ID")
	}
}

func TestBackendAddExpedition(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	if _, _, _, err := NewTestExpedition(backend); err != nil {
		t.Fatal(err)
	}
}

func TestBackendExpeditionByID(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	_, _, expedition, err := NewTestExpedition(backend)
	if err != nil {
		t.Fatal(err)
	}

	returnedExpedition, err := backend.ExpeditionByID(expedition.ID)
	if err != nil {
		t.Error(err)
	}

	if returnedExpedition.ID != expedition.ID {
		t.Error("unexpected expedition ID")
	}
}

func TestBackendExpeditionByProjectSlugAndSlug(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	_, project, expedition, err := NewTestExpedition(backend)
	if err != nil {
		t.Fatal(err)
	}

	returnedExpedition, err := backend.ExpeditionByProjectSlugAndSlug(project.Slug, expedition.Slug)
	if err != nil {
		t.Error(err)
	}

	if returnedExpedition.ID != expedition.ID {
		t.Error("unexpected expedition ID")
	}
}

func TestBackendDeleteExpeditionByID(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	_, _, expedition, err := NewTestExpedition(backend)
	if err != nil {
		t.Fatal(err)
	}

	if err := backend.DeleteExpeditionByID(expedition.ID); err != nil {
		t.Fatal(err)
	}

	_, err = backend.ExpeditionByID(expedition.ID)
	if err == nil {
		t.Fatal("project expedition not deleted")
	}

	if err != NotFoundError {
		t.Fatal(err)
	}
}
