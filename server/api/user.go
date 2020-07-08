package api

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware/security/jwt"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

func UserType(user *data.User) *app.User {
	userType := &app.User{
		ID:    int(user.ID),
		Name:  user.Name,
		Email: user.Email,
		Bio:   user.Bio,
	}

	if user.MediaURL != nil {
		userType.Photo = &app.UserPhoto{
			URL: makePhotoURL(fmt.Sprintf("/user/%d/media", user.ID), user.MediaURL),
		}
	}

	return userType
}

func ProjectUserType(user *data.ProjectUserAndUser) *app.ProjectUser {
	return &app.ProjectUser{
		User:       UserType(&user.User),
		Role:       user.RoleName(),
		Membership: data.MembershipAccepted,
	}
}

func ProjectUsersType(users []*data.ProjectUserAndUser, invites []*data.ProjectInvite) *app.ProjectUsers {
	usersCollection := make([]*app.ProjectUser, 0, len(users)+len(invites))
	for _, user := range users {
		usersCollection = append(usersCollection, ProjectUserType(user))
	}
	for _, invite := range invites {
		membership := data.MembershipPending
		if invite.RejectedTime != nil {
			membership = data.MembershipRejected
		}
		usersCollection = append(usersCollection, &app.ProjectUser{
			User: &app.User{
				Name:  invite.InvitedEmail,
				Email: invite.InvitedEmail,
			},
			Role:       invite.LookupRole().Name,
			Membership: membership,
		})
	}

	return &app.ProjectUsers{
		Users: usersCollection,
	}
}

type UserController struct {
	*goa.Controller
	options *ControllerOptions
}

func NewUserController(service *goa.Service, options *ControllerOptions) *UserController {
	return &UserController{
		Controller: service.NewController("UserController"),
		options:    options,
	}
}

func (c *UserController) userExists(ctx context.Context, email string) (bool, error) {
	existing := []*data.User{}
	err := c.options.Database.SelectContext(ctx, &existing, `SELECT * FROM fieldkit.user WHERE LOWER(email) = LOWER($1)`, email)
	if err != nil {
		return false, err
	}
	if len(existing) > 0 {
		return true, nil
	}
	return false, nil
}

func (c *UserController) Add(ctx *app.AddUserContext) error {
	log := Logger(ctx).Sugar()

	if yes, err := c.userExists(ctx, ctx.Payload.Email); err != nil {
		return err
	} else if yes {
		return ctx.BadRequest()
	}

	user := &data.User{
		Name:     data.Name(ctx.Payload.Name),
		Email:    ctx.Payload.Email,
		Username: ctx.Payload.Email,
		Bio:      "",
	}

	if err := user.SetPassword(ctx.Payload.Password); err != nil {
		return err
	}

	if err := c.options.Database.NamedGetContext(ctx, user, `
		INSERT INTO fieldkit.user (name, username, email, password, bio) VALUES (:name, :email, :email, :password, :bio) RETURNING *
		`, user); err != nil {
		return err
	}

	validationToken, err := data.NewValidationToken(user.ID, 20, time.Now().Add(time.Duration(72)*time.Hour))
	if err != nil {
		return err
	}

	if _, err := c.options.Database.NamedExecContext(ctx, `
		INSERT INTO fieldkit.validation_token (token, user_id, expires) VALUES (:token, :user_id, :expires)
		`, validationToken); err != nil {
		return err
	}

	if err := c.options.Emailer.SendValidationToken(user, validationToken); err != nil {
		return err
	}

	log.Infow("validation", "token", validationToken.Token)

	c.options.Metrics.UserAdded()

	c.options.Metrics.EmailVerificationSent()

	return ctx.OK(UserType(user))
}

func (c *UserController) Update(ctx *app.UpdateUserContext) error {
	user := &data.User{
		ID:    int32(ctx.UserID),
		Name:  data.Name(ctx.Payload.Name),
		Email: ctx.Payload.Email,
		Bio:   ctx.Payload.Bio,
	}

	if err := c.options.Database.NamedGetContext(ctx, user, `
		UPDATE fieldkit.user SET name = :name, username = :email, email = :email, bio = :bio WHERE id = :id RETURNING *
		`, user); err != nil {
		return err
	}

	return ctx.OK(UserType(user))
}

func (c *UserController) ChangePassword(ctx *app.ChangePasswordUserContext) error {
	token := jwt.ContextJWT(ctx)
	if token == nil {
		return fmt.Errorf("JWT token is missing from context") // internal error
	}

	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return fmt.Errorf("JWT claims error") // internal error
	}

	log := Logger(ctx).Sugar()

	log.Infow("password", "authorized_user_id", claims["sub"], "user_id", ctx.UserID)

	user := &data.User{}
	err := c.options.Database.GetContext(ctx, user, `SELECT * FROM fieldkit.user WHERE id = $1`, claims["sub"])
	if err == sql.ErrNoRows {
		return data.IncorrectPasswordError
	}
	if err != nil {
		return err
	}

	if user.ID != int32(ctx.UserID) {
		return fmt.Errorf("invalid user id")
	}

	err = user.CheckPassword(ctx.Payload.OldPassword)
	if err == data.IncorrectPasswordError {
		return data.IncorrectPasswordError
	}
	if err != nil {
		return err
	}

	if err := user.SetPassword(ctx.Payload.NewPassword); err != nil {
		return err
	}

	if err := c.options.Database.NamedGetContext(ctx, user, `
		UPDATE fieldkit.user SET password = :password WHERE id = :id RETURNING *
		`, user); err != nil {
		return err
	}

	return ctx.OK(UserType(user))
}

func (c *UserController) SendValidation(ctx *app.SendValidationUserContext) error {
	log := Logger(ctx).Sugar()

	user := &data.User{}
	if err := c.options.Database.GetContext(ctx, user, `
		SELECT u.* FROM fieldkit.user AS u WHERE u.id = $1
		`, ctx.UserID); err != nil {
		return err
	}

	// TODO Rate limit the number of these we can send?
	if !user.Valid {
		validationToken, err := data.NewValidationToken(user.ID, 20, time.Now().Add(time.Duration(72)*time.Hour))
		if err != nil {
			return err
		}

		if _, err := c.options.Database.ExecContext(ctx, `
			DELETE FROM fieldkit.validation_token WHERE user_id = $1
			`, user.ID); err != nil {
			return err
		}

		if _, err := c.options.Database.NamedExecContext(ctx, `
			INSERT INTO fieldkit.validation_token (token, user_id, expires) VALUES (:token, :user_id, :expires)
			`, validationToken); err != nil {
			return err
		}

		if err := c.options.Emailer.SendValidationToken(user, validationToken); err != nil {
			return err
		}

		log.Infow("validation", "token", validationToken.Token)

		c.options.Metrics.EmailVerificationSent()
	}

	return ctx.NoContent()
}

func (c *UserController) Validate(ctx *app.ValidateUserContext) error {
	log := Logger(ctx).Sugar()

	validationToken := &data.ValidationToken{}
	if err := validationToken.Token.UnmarshalText([]byte(ctx.Token)); err != nil {
		return err
	}

	err := c.options.Database.GetContext(ctx, validationToken, `
		SELECT * FROM fieldkit.validation_token WHERE token = $1
		`, validationToken.Token)
	if err == sql.ErrNoRows {
		log.Infow("invalid", "token", ctx.Token)
		ctx.ResponseData.Header().Set("Location", "https://"+c.options.PortalDomain+"?bad_token=true")
		return ctx.Found()
	}
	if err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, `UPDATE fieldkit.user SET valid = true WHERE id = $1`, validationToken.UserID); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, `DELETE FROM fieldkit.validation_token WHERE token = $1`, validationToken.Token); err != nil {
		return err
	}

	log.Infow("verified", "token", ctx.Token)

	c.options.Metrics.UserValidated()

	ctx.ResponseData.Header().Set("Location", "https://"+c.options.PortalDomain+"/")

	return ctx.Found()
}

func (c *UserController) authenticateOrSpoof(ctx context.Context, email, password string) (*data.User, error) {
	user := &data.User{}
	err := c.options.Database.GetContext(ctx, user, `SELECT u.* FROM fieldkit.user AS u WHERE LOWER(u.email) = LOWER($1)`, email)
	if err == sql.ErrNoRows {
		return nil, data.IncorrectPasswordError
	}
	if err != nil {
		return nil, err
	}

	parts := strings.Split(password, " ")
	if len(parts) == 2 {
		log := Logger(ctx).Sugar()

		// NOTE Not logging password here, may be a real one.
		log.Infow("spoofing", "email", email)

		adminUser := &data.User{}
		err := c.options.Database.GetContext(ctx, adminUser, `SELECT u.* FROM fieldkit.user AS u WHERE LOWER(u.email) = LOWER($1) AND u.admin`, parts[0])
		if err == nil {
			// We can safely log the user doing the spoofing here.
			err = adminUser.CheckPassword(parts[1])
			if err == nil {
				log.Infow("spoofed", "email", email, "spoofer", adminUser.Email)
				return user, nil
			} else {
				log.Infow("denied", "email", email, "spoofer", adminUser.Email, "spoofing", true)
			}
		}
	}

	err = user.CheckPassword(password)
	if err == data.IncorrectPasswordError {
		return nil, data.IncorrectPasswordError
	}
	if err != nil {
		return nil, err
	}

	if !user.Valid {
		return nil, data.IncorrectPasswordError
	}

	return user, nil
}

func (c *UserController) Login(ctx *app.LoginUserContext) error {
	now := time.Now()

	c.options.Metrics.AuthTry()

	user, err := c.authenticateOrSpoof(ctx, ctx.Payload.Email, ctx.Payload.Password)
	if err == data.IncorrectPasswordError {
		return ctx.Unauthorized(goa.ErrUnauthorized("invalid email or password"))
	}
	if err != nil {
		return err
	}
	if user == nil {
		return ctx.Unauthorized(goa.ErrUnauthorized("invalid email or password"))
	}

	refreshToken, err := data.NewRefreshToken(user.ID, 20, now.Add(time.Duration(72)*time.Hour))
	if err != nil {
		return err
	}

	if _, err := c.options.Database.NamedExecContext(ctx, `
		INSERT INTO fieldkit.refresh_token (token, user_id, expires) VALUES (:token, :user_id, :expires)
		`, refreshToken); err != nil {
		return ctx.Unauthorized(fmt.Errorf("invalid email or password"))
	}

	token := user.NewToken(now, refreshToken)
	signedToken, err := token.SignedString(c.options.JWTHMACKey)
	if err != nil {
		return fmt.Errorf("failed to sign token: %s", err) // internal error
	}

	c.options.Metrics.AuthSuccess()

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

	if _, err := c.options.Database.ExecContext(ctx, `DELETE FROM fieldkit.refresh_token WHERE token = $1`, refreshToken); err != nil {
		return err
	}

	return ctx.NoContent()
}

func (c *UserController) RecoveryLookup(ctx *app.RecoveryLookupUserContext) error {
	log := Logger(ctx).Sugar()

	user := &data.User{}
	err := c.options.Database.GetContext(ctx, user, `SELECT * FROM fieldkit.user WHERE email = $1`, ctx.Payload.Email)
	if err == sql.ErrNoRows {
		log.Infow("recovery, no user")
		return ctx.OK([]byte("{}"))
	}
	if err != nil {
		log.Errorw("recovery", "error", err)
		return ctx.OK([]byte("{}"))
	}

	now := time.Now().UTC()

	recoveryToken, err := data.NewRecoveryToken(user.ID, 20, now.Add(time.Duration(1)*time.Hour))
	if err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, `DELETE FROM fieldkit.recovery_token WHERE user_id = $1`, user.ID); err != nil {
		return err
	}

	if _, err := c.options.Database.NamedExecContext(ctx, `
		INSERT INTO fieldkit.recovery_token (token, user_id, expires) VALUES (:token, :user_id, :expires)
		`, recoveryToken); err != nil {
		return err
	}

	if err := c.options.Emailer.SendRecoveryToken(user, recoveryToken); err != nil {
		return err
	}

	log.Infow("recovery", "token", recoveryToken.Token)

	c.options.Metrics.EmailRecoverPasswordSent()

	return ctx.OK([]byte("{}"))
}

func (c *UserController) Recovery(ctx *app.RecoveryUserContext) error {
	log := Logger(ctx).Sugar()

	token := data.Token{}
	if err := token.UnmarshalText([]byte(ctx.Payload.Token)); err != nil {
		return err
	}

	log.Infow("recovery", "token_raw", ctx.Payload.Token, "token", token)

	recoveryToken := &data.RecoveryToken{}
	err := c.options.Database.GetContext(ctx, recoveryToken, `SELECT * FROM fieldkit.recovery_token WHERE token = $1`, token)
	if err == sql.ErrNoRows {
		log.Infow("recovery, token bad")
		return ctx.Unauthorized()
	}
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	if now.After(recoveryToken.Expires) {
		log.Infow("recovery, token expired", "token_expires", recoveryToken.Expires, "now", now)
		return ctx.Unauthorized()
	}

	user := &data.User{}
	err = c.options.Database.GetContext(ctx, user, `
		SELECT * FROM fieldkit.user WHERE id = $1
		`, recoveryToken.UserID)
	if err == sql.ErrNoRows {
		return ctx.OK([]byte("{}"))
	}
	if err != nil {
		return err
	}

	if err := user.SetPassword(ctx.Payload.Password); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, `
		DELETE FROM fieldkit.recovery_token WHERE user_id = $1
		`, user.ID); err != nil {
		return err
	}

	if err := c.options.Database.NamedGetContext(ctx, user, `
		UPDATE fieldkit.user SET password = :password, valid = true WHERE id = :id RETURNING *
		`, user); err != nil {
		return err
	}

	return ctx.OK([]byte("{}"))
}

func (c *UserController) Refresh(ctx *app.RefreshUserContext) error {
	c.options.Metrics.AuthRefreshTry()

	token := data.Token{}
	if err := token.UnmarshalText([]byte(ctx.Payload.RefreshToken)); err != nil {
		return err
	}

	refreshToken := &data.RefreshToken{}
	err := c.options.Database.GetContext(ctx, refreshToken, `SELECT * FROM fieldkit.refresh_token WHERE token = $1`, token)
	if err == sql.ErrNoRows {
		return ctx.Unauthorized()
	}
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	if now.After(refreshToken.Expires) {
		return ctx.Unauthorized()
	}

	user := &data.User{}
	err = c.options.Database.GetContext(ctx, user, `SELECT * FROM fieldkit.user WHERE id = $1`, refreshToken.UserID)
	if err == sql.ErrNoRows {
		return ctx.Unauthorized()
	}

	if err != nil {
		return err
	}

	jwtToken := user.NewToken(now, refreshToken)
	signedToken, err := jwtToken.SignedString(c.options.JWTHMACKey)
	if err != nil {
		return fmt.Errorf("failed to sign token: %s", err) // internal error
	}

	c.options.Metrics.AuthRefreshSuccess()

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

	log := Logger(ctx).Sugar()

	user := &data.User{}
	if err := c.options.Database.GetContext(ctx, user, `SELECT u.* FROM fieldkit.user AS u WHERE u.id = $1`, claims["sub"]); err != nil {
		log.Infow("user", "user_id", claims["sub"])
		return err
	}

	log.Infow("user", "user_id", claims["sub"], "email", user.Email)

	return ctx.OK(UserType(user))
}

func (c *UserController) GetID(ctx *app.GetIDUserContext) error {
	user := &data.User{}
	if err := c.options.Database.GetContext(ctx, user, `SELECT u.* FROM fieldkit.user AS u WHERE u.id = $1`, ctx.UserID); err != nil {
		return err
	}

	return ctx.OK(UserType(user))
}

func (c *UserController) ListByProject(ctx *app.ListByProjectUserContext) error {
	users := []*data.ProjectUserAndUser{}
	if err := c.options.Database.SelectContext(ctx, &users, `
		SELECT pu.*, u.* FROM fieldkit.user AS u JOIN fieldkit.project_user AS pu ON pu.user_id = u.id WHERE pu.project_id = $1 ORDER BY pu.role_id DESC, u.id
		`, ctx.ProjectID); err != nil {
		return err
	}

	invites := []*data.ProjectInvite{}
	if err := c.options.Database.SelectContext(ctx, &invites, `
		SELECT * FROM fieldkit.project_invite WHERE project_id = $1 AND accepted_time IS NULL ORDER BY invited_time
		`, ctx.ProjectID); err != nil {
		return err
	}

	return ctx.OK(ProjectUsersType(users, invites))
}

func (c *UserController) SaveCurrentUserImage(ctx *app.SaveCurrentUserImageUserContext) error {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return err
	}

	mr := repositories.NewMediaRepository(c.options.MediaFiles)
	saved, err := mr.Save(ctx, ctx.RequestData)
	if err != nil {
		return err
	}

	user := &data.User{}
	if err := c.options.Database.GetContext(ctx, user, `UPDATE fieldkit.user SET media_url = $1, media_content_type = $2 WHERE id = $3 RETURNING *`, saved.URL, saved.MimeType, p.UserID()); err != nil {
		return err
	}

	return ctx.OK(UserType(user))
}

func (c *UserController) GetCurrentUserImage(ctx *app.GetCurrentUserImageUserContext) error {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return err
	}

	user := &data.User{}
	if err := c.options.Database.GetContext(ctx, user, `SELECT media_url FROM fieldkit.user WHERE id = $1`, p.UserID); err != nil {
		return err
	}

	if user.MediaURL != nil {
		mr := repositories.NewMediaRepository(c.options.MediaFiles)

		lm, err := mr.LoadByURL(ctx, *user.MediaURL)
		if err != nil {
			return err
		}

		if lm != nil {
			sendLoadedMedia(ctx.ResponseData, lm)
		}

		return nil
	}

	return ctx.OK(nil)
}

func (c *UserController) GetUserImage(ctx *app.GetUserImageUserContext) error {
	user := &data.User{}
	if err := c.options.Database.GetContext(ctx, user, `SELECT media_url FROM fieldkit.user WHERE id = $1`, ctx.UserID); err != nil {
		return err
	}

	if user.MediaURL != nil {
		mr := repositories.NewMediaRepository(c.options.MediaFiles)

		lm, err := mr.LoadByURL(ctx, *user.MediaURL)
		if err != nil {
			return err
		}

		if lm != nil {
			sendLoadedMedia(ctx.ResponseData, lm)
		}

		return nil
	}

	return ctx.OK(nil)
}

func NewTransmissionToken(now time.Time, user *data.User) *jwtgo.Token {
	scopes := []string{"api:transmission"}

	token := jwtgo.New(jwtgo.SigningMethodHS512)
	token.Claims = jwtgo.MapClaims{
		"iat":    now.Unix(),
		"exp":    now.Add(time.Hour * 24 * 365).Unix(),
		"sub":    user.ID,
		"email":  user.Email,
		"scopes": scopes,
	}

	return token
}

func (c *UserController) TransmissionToken(ctx *app.TransmissionTokenUserContext) error {
	token := jwt.ContextJWT(ctx)
	if token == nil {
		return fmt.Errorf("JWT token is missing from context") // internal error
	}

	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return fmt.Errorf("JWT claims error") // internal error
	}

	user := &data.User{}
	if err := c.options.Database.GetContext(ctx, user, `SELECT u.* FROM fieldkit.user AS u WHERE u.id = $1`, claims["sub"]); err != nil {
		return err
	}

	now := time.Now().UTC()
	transmissionToken := NewTransmissionToken(now, user)
	signedToken, err := transmissionToken.SignedString(c.options.JWTHMACKey)
	if err != nil {
		return fmt.Errorf("failed to sign token: %s", err) // internal error
	}

	return ctx.OK(&app.TransmissionToken{
		Token: signedToken,
	})
}

func (c *UserController) ProjectRoles(ctx *app.ProjectRolesUserContext) error {
	roles := make([]*app.ProjectRole, 0)
	for _, role := range data.Roles {
		roles = append(roles, &app.ProjectRole{
			ID:   int(role.ID),
			Name: role.Name,
		})
	}
	return ctx.OK(roles)
}
