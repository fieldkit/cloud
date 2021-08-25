// Code generated by goa v3.2.4, DO NOT EDIT.
//
// HTTP request path constructors for the user service.
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package server

import (
	"fmt"
)

// RolesUserPath returns the URL path to the user service roles HTTP endpoint.
func RolesUserPath() string {
	return "/roles"
}

// DeleteUserPath returns the URL path to the user service delete HTTP endpoint.
func DeleteUserPath(userID int32) string {
	return fmt.Sprintf("/admin/users/%v", userID)
}

// UploadPhotoUserPath returns the URL path to the user service upload photo HTTP endpoint.
func UploadPhotoUserPath() string {
	return "/user/media"
}

// DownloadPhotoUserPath returns the URL path to the user service download photo HTTP endpoint.
func DownloadPhotoUserPath(userID int32) string {
	return fmt.Sprintf("/user/%v/media", userID)
}

// LoginUserPath returns the URL path to the user service login HTTP endpoint.
func LoginUserPath() string {
	return "/login"
}

// RecoveryLookupUserPath returns the URL path to the user service recovery lookup HTTP endpoint.
func RecoveryLookupUserPath() string {
	return "/user/recovery/lookup"
}

// RecoveryUserPath returns the URL path to the user service recovery HTTP endpoint.
func RecoveryUserPath() string {
	return "/user/recovery"
}

// ResumeUserPath returns the URL path to the user service resume HTTP endpoint.
func ResumeUserPath() string {
	return "/user/resume"
}

// LogoutUserPath returns the URL path to the user service logout HTTP endpoint.
func LogoutUserPath() string {
	return "/logout"
}

// RefreshUserPath returns the URL path to the user service refresh HTTP endpoint.
func RefreshUserPath() string {
	return "/refresh"
}

// SendValidationUserPath returns the URL path to the user service send validation HTTP endpoint.
func SendValidationUserPath(userID int32) string {
	return fmt.Sprintf("/users/%v/validate-email", userID)
}

// ValidateUserPath returns the URL path to the user service validate HTTP endpoint.
func ValidateUserPath() string {
	return "/validate"
}

// AddUserPath returns the URL path to the user service add HTTP endpoint.
func AddUserPath() string {
	return "/users"
}

// UpdateUserPath returns the URL path to the user service update HTTP endpoint.
func UpdateUserPath(userID int32) string {
	return fmt.Sprintf("/users/%v", userID)
}

// ChangePasswordUserPath returns the URL path to the user service change password HTTP endpoint.
func ChangePasswordUserPath(userID int32) string {
	return fmt.Sprintf("/users/%v/password", userID)
}

// AcceptTncUserPath returns the URL path to the user service accept tnc HTTP endpoint.
func AcceptTncUserPath(userID int32) string {
	return fmt.Sprintf("/users/%v/accept-tnc", userID)
}

// GetCurrentUserPath returns the URL path to the user service get current HTTP endpoint.
func GetCurrentUserPath() string {
	return "/user"
}

// ListByProjectUserPath returns the URL path to the user service list by project HTTP endpoint.
func ListByProjectUserPath(projectID int32) string {
	return fmt.Sprintf("/users/project/%v", projectID)
}

// IssueTransmissionTokenUserPath returns the URL path to the user service issue transmission token HTTP endpoint.
func IssueTransmissionTokenUserPath() string {
	return "/user/transmission-token"
}

// ProjectRolesUserPath returns the URL path to the user service project roles HTTP endpoint.
func ProjectRolesUserPath() string {
	return "/projects/roles"
}

// AdminTermsAndConditionsUserPath returns the URL path to the user service admin terms and conditions HTTP endpoint.
func AdminTermsAndConditionsUserPath() string {
	return "/admin/user/tnc"
}

// AdminDeleteUserPath returns the URL path to the user service admin delete HTTP endpoint.
func AdminDeleteUserPath() string {
	return "/admin/user"
}

// AdminSearchUserPath returns the URL path to the user service admin search HTTP endpoint.
func AdminSearchUserPath() string {
	return "/admin/users/search"
}
