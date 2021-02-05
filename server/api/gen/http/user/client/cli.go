// Code generated by goa v3.2.4, DO NOT EDIT.
//
// user HTTP client CLI support package
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"encoding/json"
	"fmt"
	"strconv"
	"unicode/utf8"

	user "github.com/fieldkit/cloud/server/api/gen/user"
	goa "goa.design/goa/v3/pkg"
)

// BuildRolesPayload builds the payload for the user roles endpoint from CLI
// flags.
func BuildRolesPayload(userRolesAuth string) (*user.RolesPayload, error) {
	var auth string
	{
		auth = userRolesAuth
	}
	v := &user.RolesPayload{}
	v.Auth = auth

	return v, nil
}

// BuildDeletePayload builds the payload for the user delete endpoint from CLI
// flags.
func BuildDeletePayload(userDeleteUserID string, userDeleteAuth string) (*user.DeletePayload, error) {
	var err error
	var userID int32
	{
		var v int64
		v, err = strconv.ParseInt(userDeleteUserID, 10, 32)
		userID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for userID, must be INT32")
		}
	}
	var auth string
	{
		auth = userDeleteAuth
	}
	v := &user.DeletePayload{}
	v.UserID = userID
	v.Auth = auth

	return v, nil
}

// BuildUploadPhotoPayload builds the payload for the user upload photo
// endpoint from CLI flags.
func BuildUploadPhotoPayload(userUploadPhotoContentType string, userUploadPhotoContentLength string, userUploadPhotoAuth string) (*user.UploadPhotoPayload, error) {
	var err error
	var contentType string
	{
		contentType = userUploadPhotoContentType
	}
	var contentLength int64
	{
		contentLength, err = strconv.ParseInt(userUploadPhotoContentLength, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value for contentLength, must be INT64")
		}
	}
	var auth string
	{
		auth = userUploadPhotoAuth
	}
	v := &user.UploadPhotoPayload{}
	v.ContentType = contentType
	v.ContentLength = contentLength
	v.Auth = auth

	return v, nil
}

// BuildDownloadPhotoPayload builds the payload for the user download photo
// endpoint from CLI flags.
func BuildDownloadPhotoPayload(userDownloadPhotoUserID string, userDownloadPhotoSize string, userDownloadPhotoIfNoneMatch string) (*user.DownloadPhotoPayload, error) {
	var err error
	var userID int32
	{
		var v int64
		v, err = strconv.ParseInt(userDownloadPhotoUserID, 10, 32)
		userID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for userID, must be INT32")
		}
	}
	var size *int32
	{
		if userDownloadPhotoSize != "" {
			var v int64
			v, err = strconv.ParseInt(userDownloadPhotoSize, 10, 32)
			val := int32(v)
			size = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for size, must be INT32")
			}
		}
	}
	var ifNoneMatch *string
	{
		if userDownloadPhotoIfNoneMatch != "" {
			ifNoneMatch = &userDownloadPhotoIfNoneMatch
		}
	}
	v := &user.DownloadPhotoPayload{}
	v.UserID = userID
	v.Size = size
	v.IfNoneMatch = ifNoneMatch

	return v, nil
}

// BuildLoginPayload builds the payload for the user login endpoint from CLI
// flags.
func BuildLoginPayload(userLoginBody string) (*user.LoginPayload, error) {
	var err error
	var body LoginRequestBody
	{
		err = json.Unmarshal([]byte(userLoginBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"email\": \"9ta\",\n      \"password\": \"fwe\"\n   }'")
		}
		err = goa.MergeErrors(err, goa.ValidateFormat("body.email", body.Email, goa.FormatEmail))

		if utf8.RuneCountInString(body.Email) > 40 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.email", body.Email, utf8.RuneCountInString(body.Email), 40, false))
		}
		if utf8.RuneCountInString(body.Password) < 10 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.password", body.Password, utf8.RuneCountInString(body.Password), 10, true))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &user.LoginFields{
		Email:    body.Email,
		Password: body.Password,
	}
	res := &user.LoginPayload{
		Login: v,
	}

	return res, nil
}

// BuildRecoveryLookupPayload builds the payload for the user recovery lookup
// endpoint from CLI flags.
func BuildRecoveryLookupPayload(userRecoveryLookupBody string) (*user.RecoveryLookupPayload, error) {
	var err error
	var body RecoveryLookupRequestBody
	{
		err = json.Unmarshal([]byte(userRecoveryLookupBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"email\": \"Earum sint ut veritatis est optio.\"\n   }'")
		}
	}
	v := &user.RecoveryLookupFields{
		Email: body.Email,
	}
	res := &user.RecoveryLookupPayload{
		Recovery: v,
	}

	return res, nil
}

// BuildRecoveryPayload builds the payload for the user recovery endpoint from
// CLI flags.
func BuildRecoveryPayload(userRecoveryBody string) (*user.RecoveryPayload, error) {
	var err error
	var body RecoveryRequestBody
	{
		err = json.Unmarshal([]byte(userRecoveryBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"password\": \"u2m\",\n      \"token\": \"A et sit libero tempora voluptas delectus.\"\n   }'")
		}
		if utf8.RuneCountInString(body.Password) < 10 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.password", body.Password, utf8.RuneCountInString(body.Password), 10, true))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &user.RecoveryFields{
		Token:    body.Token,
		Password: body.Password,
	}
	res := &user.RecoveryPayload{
		Recovery: v,
	}

	return res, nil
}

// BuildResumePayload builds the payload for the user resume endpoint from CLI
// flags.
func BuildResumePayload(userResumeBody string) (*user.ResumePayload, error) {
	var err error
	var body struct {
		Token *string `form:"token" json:"token" xml:"token"`
	}
	{
		err = json.Unmarshal([]byte(userResumeBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"token\": \"Ullam odit sequi sunt reiciendis sint.\"\n   }'")
		}
	}
	v := &user.ResumePayload{}
	if body.Token != nil {
		v.Token = *body.Token
	}

	return v, nil
}

// BuildLogoutPayload builds the payload for the user logout endpoint from CLI
// flags.
func BuildLogoutPayload(userLogoutAuth string) (*user.LogoutPayload, error) {
	var auth string
	{
		auth = userLogoutAuth
	}
	v := &user.LogoutPayload{}
	v.Auth = auth

	return v, nil
}

// BuildRefreshPayload builds the payload for the user refresh endpoint from
// CLI flags.
func BuildRefreshPayload(userRefreshBody string) (*user.RefreshPayload, error) {
	var err error
	var body RefreshRequestBody
	{
		err = json.Unmarshal([]byte(userRefreshBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"refreshToken\": \"Molestias natus dolorem a possimus rem.\"\n   }'")
		}
	}
	v := &user.RefreshPayload{
		RefreshToken: body.RefreshToken,
	}

	return v, nil
}

// BuildSendValidationPayload builds the payload for the user send validation
// endpoint from CLI flags.
func BuildSendValidationPayload(userSendValidationUserID string) (*user.SendValidationPayload, error) {
	var err error
	var userID int32
	{
		var v int64
		v, err = strconv.ParseInt(userSendValidationUserID, 10, 32)
		userID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for userID, must be INT32")
		}
	}
	v := &user.SendValidationPayload{}
	v.UserID = userID

	return v, nil
}

// BuildValidatePayload builds the payload for the user validate endpoint from
// CLI flags.
func BuildValidatePayload(userValidateToken string) (*user.ValidatePayload, error) {
	var token string
	{
		token = userValidateToken
	}
	v := &user.ValidatePayload{}
	v.Token = token

	return v, nil
}

// BuildAddPayload builds the payload for the user add endpoint from CLI flags.
func BuildAddPayload(userAddBody string) (*user.AddPayload, error) {
	var err error
	var body AddRequestBody
	{
		err = json.Unmarshal([]byte(userAddBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"etag\": \"Sed et maxime dolore.\",\n      \"meta\": \"Praesentium atque quibusdam tempora assumenda sed eos.\",\n      \"module\": \"Distinctio nesciunt aut.\",\n      \"profile\": \"Cupiditate qui repudiandae et blanditiis dolor mollitia.\",\n      \"url\": \"Quo hic a molestias culpa a pariatur.\"\n   }'")
		}
		err = goa.MergeErrors(err, goa.ValidatePattern("body.name", body.Name, "\\S"))
		if utf8.RuneCountInString(body.Name) > 256 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.name", body.Name, utf8.RuneCountInString(body.Name), 256, false))
		}
		err = goa.MergeErrors(err, goa.ValidateFormat("body.email", body.Email, goa.FormatEmail))

		if utf8.RuneCountInString(body.Email) > 40 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.email", body.Email, utf8.RuneCountInString(body.Email), 40, false))
		}
		if utf8.RuneCountInString(body.Password) < 10 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.password", body.Password, utf8.RuneCountInString(body.Password), 10, true))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &user.AddUserFields{
		Name:        body.Name,
		Email:       body.Email,
		Password:    body.Password,
		InviteToken: body.InviteToken,
	}
	res := &user.AddPayload{
		User: v,
	}

	return res, nil
}

// BuildUpdatePayload builds the payload for the user update endpoint from CLI
// flags.
func BuildUpdatePayload(userUpdateBody string, userUpdateUserID string, userUpdateAuth string) (*user.UpdatePayload, error) {
	var err error
	var body UpdateRequestBody
	{
		err = json.Unmarshal([]byte(userUpdateBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"description\": \"Voluptatem adipisci dicta recusandae dolores assumenda.\",\n      \"endTime\": \"Consequuntur totam impedit nobis inventore nesciunt.\",\n      \"goal\": \"Qui doloremque nemo similique fuga voluptatem.\",\n      \"location\": \"Non quae voluptatem.\",\n      \"name\": \"Corporis enim harum.\",\n      \"privacy\": 2011573778,\n      \"startTime\": \"Optio qui architecto perspiciatis nulla inventore ut.\",\n      \"tags\": \"Eius aut.\"\n   }'")
		}
		err = goa.MergeErrors(err, goa.ValidatePattern("body.name", body.Name, "\\S"))
		if utf8.RuneCountInString(body.Name) > 256 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.name", body.Name, utf8.RuneCountInString(body.Name), 256, false))
		}
		err = goa.MergeErrors(err, goa.ValidateFormat("body.email", body.Email, goa.FormatEmail))

		if utf8.RuneCountInString(body.Email) > 40 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.email", body.Email, utf8.RuneCountInString(body.Email), 40, false))
		}
		if err != nil {
			return nil, err
		}
	}
	var userID int32
	{
		var v int64
		v, err = strconv.ParseInt(userUpdateUserID, 10, 32)
		userID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for userID, must be INT32")
		}
	}
	var auth string
	{
		auth = userUpdateAuth
	}
	v := &user.UpdateUserFields{
		Name:  body.Name,
		Email: body.Email,
		Bio:   body.Bio,
	}
	res := &user.UpdatePayload{
		Update: v,
	}
	res.UserID = userID
	res.Auth = auth

	return res, nil
}

// BuildChangePasswordPayload builds the payload for the user change password
// endpoint from CLI flags.
func BuildChangePasswordPayload(userChangePasswordBody string, userChangePasswordUserID string, userChangePasswordAuth string) (*user.ChangePasswordPayload, error) {
	var err error
	var body ChangePasswordRequestBody
	{
		err = json.Unmarshal([]byte(userChangePasswordBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"newPassword\": \"k8q\",\n      \"oldPassword\": \"t4l\"\n   }'")
		}
		if utf8.RuneCountInString(body.OldPassword) < 10 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.oldPassword", body.OldPassword, utf8.RuneCountInString(body.OldPassword), 10, true))
		}
		if utf8.RuneCountInString(body.NewPassword) < 10 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.newPassword", body.NewPassword, utf8.RuneCountInString(body.NewPassword), 10, true))
		}
		if err != nil {
			return nil, err
		}
	}
	var userID int32
	{
		var v int64
		v, err = strconv.ParseInt(userChangePasswordUserID, 10, 32)
		userID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for userID, must be INT32")
		}
	}
	var auth string
	{
		auth = userChangePasswordAuth
	}
	v := &user.UpdateUserPasswordFields{
		OldPassword: body.OldPassword,
		NewPassword: body.NewPassword,
	}
	res := &user.ChangePasswordPayload{
		Change: v,
	}
	res.UserID = userID
	res.Auth = auth

	return res, nil
}

// BuildGetCurrentPayload builds the payload for the user get current endpoint
// from CLI flags.
func BuildGetCurrentPayload(userGetCurrentAuth string) (*user.GetCurrentPayload, error) {
	var auth string
	{
		auth = userGetCurrentAuth
	}
	v := &user.GetCurrentPayload{}
	v.Auth = auth

	return v, nil
}

// BuildListByProjectPayload builds the payload for the user list by project
// endpoint from CLI flags.
func BuildListByProjectPayload(userListByProjectProjectID string, userListByProjectAuth string) (*user.ListByProjectPayload, error) {
	var err error
	var projectID int32
	{
		var v int64
		v, err = strconv.ParseInt(userListByProjectProjectID, 10, 32)
		projectID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for projectID, must be INT32")
		}
	}
	var auth string
	{
		auth = userListByProjectAuth
	}
	v := &user.ListByProjectPayload{}
	v.ProjectID = projectID
	v.Auth = auth

	return v, nil
}

// BuildIssueTransmissionTokenPayload builds the payload for the user issue
// transmission token endpoint from CLI flags.
func BuildIssueTransmissionTokenPayload(userIssueTransmissionTokenAuth string) (*user.IssueTransmissionTokenPayload, error) {
	var auth string
	{
		auth = userIssueTransmissionTokenAuth
	}
	v := &user.IssueTransmissionTokenPayload{}
	v.Auth = auth

	return v, nil
}

// BuildAdminDeletePayload builds the payload for the user admin delete
// endpoint from CLI flags.
func BuildAdminDeletePayload(userAdminDeleteBody string, userAdminDeleteAuth string) (*user.AdminDeletePayload, error) {
	var err error
	var body AdminDeleteRequestBody
	{
		err = json.Unmarshal([]byte(userAdminDeleteBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"email\": \"Voluptatem qui.\",\n      \"password\": \"Labore qui a.\"\n   }'")
		}
	}
	var auth string
	{
		auth = userAdminDeleteAuth
	}
	v := &user.AdminDeleteFields{
		Email:    body.Email,
		Password: body.Password,
	}
	res := &user.AdminDeletePayload{
		Delete: v,
	}
	res.Auth = auth

	return res, nil
}

// BuildAdminSearchPayload builds the payload for the user admin search
// endpoint from CLI flags.
func BuildAdminSearchPayload(userAdminSearchQuery string, userAdminSearchAuth string) (*user.AdminSearchPayload, error) {
	var query string
	{
		query = userAdminSearchQuery
	}
	var auth string
	{
		auth = userAdminSearchAuth
	}
	v := &user.AdminSearchPayload{}
	v.Query = query
	v.Auth = auth

	return v, nil
}
