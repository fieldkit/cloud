package api

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/O-C-R/sqlxcache"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware/security/jwt"

	"github.com/O-C-R/fieldkit/server/api/app"
	"github.com/O-C-R/fieldkit/server/data"
	"github.com/O-C-R/fieldkit/server/email"
)

func UserType(user *data.User) *app.User {
	return &app.User{
		ID:       int(user.ID),
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Bio:      user.Bio,
	}
}

func UsersType(users []*data.User) *app.Users {
	usersCollection := make([]*app.User, len(users))
	for i, user := range users {
		usersCollection[i] = UserType(user)
	}

	return &app.Users{
		Users: usersCollection,
	}
}

func NewToken(now time.Time, user *data.User, refreshToken *data.RefreshToken) *jwtgo.Token {
	token := jwtgo.New(jwtgo.SigningMethodHS512)
	token.Claims = jwtgo.MapClaims{
		"iat":           now.Unix(),
		"exp":           now.Add(time.Hour).Unix(),
		"sub":           user.ID,
		"username":      user.Username,
		"refresh_token": refreshToken.Token.String(),
		"scopes":        "api:access",
	}

	return token
}

type UserControllerOptions struct {
	Database   *sqlxcache.DB
	Emailer    email.Emailer
	JWTHMACKey []byte
	Domain     string
}

// UserController implements the user resource.
type UserController struct {
	*goa.Controller
	options UserControllerOptions
}

func NewUserController(service *goa.Service, options UserControllerOptions) (*UserController, error) {
	return &UserController{
		Controller: service.NewController("UserController"),
		options:    options,
	}, nil
}

func (c *UserController) Add(ctx *app.AddUserContext) error {
	user := &data.User{
		Name:     ctx.Payload.Name,
		Username: ctx.Payload.Username,
		Email:    ctx.Payload.Email,
		Bio:      ctx.Payload.Bio,
	}

	if err := user.SetPassword(ctx.Payload.Password); err != nil {
		return err
	}

	if err := c.options.Database.NamedGetContext(ctx, user, "INSERT INTO fieldkit.user (name, username, email, password, bio) VALUES (:name, :username, :email, :password, :bio) RETURNING *", user); err != nil {
		return err
	}

	validationToken, err := data.NewValidationToken(user.ID, 20, time.Now().Add(time.Duration(72)*time.Hour))
	if err != nil {
		return err
	}

	if _, err := c.options.Database.NamedExecContext(ctx, "INSERT INTO fieldkit.validation_token (token, user_id, expires) VALUES (:token, :user_id, :expires)", validationToken); err != nil {
		return err
	}

	if err := c.options.Emailer.SendValidationToken(user, validationToken); err != nil {
		return err
	}

	return ctx.OK(UserType(user))
}

func (c *UserController) Validate(ctx *app.ValidateUserContext) error {
	fmt.Println(ctx.Token)

	validationToken := &data.ValidationToken{}
	if err := validationToken.Token.UnmarshalText([]byte(ctx.Token)); err != nil {
		return err
	}

	err := c.options.Database.GetContext(ctx, validationToken, "SELECT * FROM fieldkit.validation_token WHERE token = $1", validationToken.Token)
	if err == sql.ErrNoRows {
		return ctx.Unauthorized()
	}

	if err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, "UPDATE fieldkit.user SET valid = true WHERE id = $1", validationToken.UserID); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, "DELETE FROM fieldkit.validation_token WHERE token = $1", validationToken.Token); err != nil {
		return err
	}

	ctx.ResponseData.Header().Set("Location", "https://"+c.options.Domain+"/admin/signin")
	return ctx.Found()
}

func (c *UserController) Login(ctx *app.LoginUserContext) error {
	now := time.Now()

	user := &data.User{}
	err := c.options.Database.GetContext(ctx, user, "SELECT * FROM fieldkit.user WHERE username = $1", ctx.Payload.Username)
	if err == sql.ErrNoRows {
		return ctx.Unauthorized(goa.ErrUnauthorized("invalid username or password"))
	}

	if err != nil {
		return err
	}

	// if !user.Valid {
	// 	return ctx.Unauthorized()
	// }

	err = user.CheckPassword(ctx.Payload.Password)
	if err == data.IncorrectPasswordError {
		return ctx.Unauthorized(goa.ErrUnauthorized("invalid username or password"))
	}

	if err != nil {
		return err
	}

	refreshToken, err := data.NewRefreshToken(user.ID, 20, now.Add(time.Duration(72)*time.Hour))
	if err != nil {
		return err
	}

	if _, err := c.options.Database.NamedExecContext(ctx, "INSERT INTO fieldkit.refresh_token (token, user_id, expires) VALUES (:token, :user_id, :expires)", refreshToken); err != nil {
		return ctx.Unauthorized(fmt.Errorf("invalid username or password"))
	}

	token := NewToken(now, user, refreshToken)
	signedToken, err := token.SignedString(c.options.JWTHMACKey)
	if err != nil {
		return fmt.Errorf("failed to sign token: %s", err) // internal error
	}

	// Send response
	ctx.ResponseData.Header().Set("Authorization", "Bearer "+signedToken)
	return ctx.NoContent()
}

func (c *UserController) Logout(ctx *app.LogoutUserContext) error {
	token := jwt.ContextJWT(ctx)
	if token == nil {
		return fmt.Errorf("JWT token is missing from context") // internal error
	}

	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return fmt.Errorf("JWT claims error") // internal error
	}

	refreshToken := data.Token{}
	if err := refreshToken.UnmarshalText([]byte(claims["refresh_token"].(string))); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, "DELETE FROM fieldkit.refresh_token WHERE token = $1", refreshToken); err != nil {
		return err
	}

	return ctx.NoContent()
}

func (c *UserController) Refresh(ctx *app.RefreshUserContext) error {
	now := time.Now()

	token := data.Token{}
	if err := token.UnmarshalText([]byte(ctx.Payload.RefreshToken)); err != nil {
		return err
	}

	refreshToken := &data.RefreshToken{}
	err := c.options.Database.GetContext(ctx, refreshToken, "SELECT * FROM fieldkit.refresh_token WHERE token = $1", token)
	if err == sql.ErrNoRows {
		return ctx.Unauthorized()
	}

	if err != nil {
		return err
	}

	if now.After(refreshToken.Expires) {
		return ctx.Unauthorized()
	}

	user := &data.User{}
	err = c.options.Database.GetContext(ctx, user, "SELECT * FROM fieldkit.user WHERE id = $1", refreshToken.UserID)
	if err == sql.ErrNoRows {
		return ctx.Unauthorized()
	}

	if err != nil {
		return err
	}

	jwtToken := NewToken(now, user, refreshToken)
	signedToken, err := jwtToken.SignedString(c.options.JWTHMACKey)
	if err != nil {
		return fmt.Errorf("failed to sign token: %s", err) // internal error
	}

	// Send response
	ctx.ResponseData.Header().Set("Authorization", "Bearer "+signedToken)
	return ctx.NoContent()
}

func (c *UserController) GetCurrent(ctx *app.GetCurrentUserContext) error {
	token := jwt.ContextJWT(ctx)
	if token == nil {
		return fmt.Errorf("JWT token is missing from context") // internal error
	}

	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return fmt.Errorf("JWT claims error") // internal error
	}

	user := &data.User{}
	if err := c.options.Database.GetContext(ctx, user, "SELECT * FROM fieldkit.user WHERE id = $1", claims["sub"]); err != nil {
		return err
	}

	return ctx.OK(UserType(user))
}

func (c *UserController) Get(ctx *app.GetUserContext) error {
	user := &data.User{}
	if err := c.options.Database.GetContext(ctx, user, "SELECT * FROM fieldkit.user WHERE username = $1 LIMIT 1", ctx.Username); err != nil {
		return err
	}

	return ctx.OK(UserType(user))
}

func (c *UserController) GetID(ctx *app.GetIDUserContext) error {
	user := &data.User{}
	if err := c.options.Database.GetContext(ctx, user, "SELECT * FROM fieldkit.user WHERE id = $1", ctx.UserID); err != nil {
		return err
	}

	return ctx.OK(UserType(user))
}

func (c *UserController) List(ctx *app.ListUserContext) error {
	users := []*data.User{}
	if err := c.options.Database.SelectContext(ctx, &users, "SELECT * FROM fieldkit.user"); err != nil {
		return err
	}

	return ctx.OK(UsersType(users))
}
