// Code generated by goa v3.2.4, DO NOT EDIT.
//
// HTTP request path constructors for the project service.
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"fmt"
)

// AddUpdateProjectPath returns the URL path to the project service add update HTTP endpoint.
func AddUpdateProjectPath(projectID int32) string {
	return fmt.Sprintf("/projects/%v/updates", projectID)
}

// DeleteUpdateProjectPath returns the URL path to the project service delete update HTTP endpoint.
func DeleteUpdateProjectPath(projectID int32, updateID int64) string {
	return fmt.Sprintf("/projects/%v/updates/%v", projectID, updateID)
}

// ModifyUpdateProjectPath returns the URL path to the project service modify update HTTP endpoint.
func ModifyUpdateProjectPath(projectID int32, updateID int64) string {
	return fmt.Sprintf("/projects/%v/updates/%v", projectID, updateID)
}

// InvitesProjectPath returns the URL path to the project service invites HTTP endpoint.
func InvitesProjectPath() string {
	return "/projects/invites/pending"
}

// LookupInviteProjectPath returns the URL path to the project service lookup invite HTTP endpoint.
func LookupInviteProjectPath(token string) string {
	return fmt.Sprintf("/projects/invites/%v", token)
}

// AcceptProjectInviteProjectPath returns the URL path to the project service accept project invite HTTP endpoint.
func AcceptProjectInviteProjectPath(projectID int32) string {
	return fmt.Sprintf("/projects/%v/invites/accept", projectID)
}

// RejectProjectInviteProjectPath returns the URL path to the project service reject project invite HTTP endpoint.
func RejectProjectInviteProjectPath(projectID int32) string {
	return fmt.Sprintf("/projects/%v/invites/reject", projectID)
}

// AcceptInviteProjectPath returns the URL path to the project service accept invite HTTP endpoint.
func AcceptInviteProjectPath(id int64) string {
	return fmt.Sprintf("/projects/invites/%v/accept", id)
}

// RejectInviteProjectPath returns the URL path to the project service reject invite HTTP endpoint.
func RejectInviteProjectPath(id int64) string {
	return fmt.Sprintf("/projects/invites/%v/reject", id)
}

// AddProjectPath returns the URL path to the project service add HTTP endpoint.
func AddProjectPath() string {
	return "/projects"
}

// UpdateProjectPath returns the URL path to the project service update HTTP endpoint.
func UpdateProjectPath(projectID int32) string {
	return fmt.Sprintf("/projects/%v", projectID)
}

// GetProjectPath returns the URL path to the project service get HTTP endpoint.
func GetProjectPath(projectID int32) string {
	return fmt.Sprintf("/projects/%v", projectID)
}

// ListCommunityProjectPath returns the URL path to the project service list community HTTP endpoint.
func ListCommunityProjectPath() string {
	return "/projects"
}

// ListMineProjectPath returns the URL path to the project service list mine HTTP endpoint.
func ListMineProjectPath() string {
	return "/user/projects"
}

// InviteProjectPath returns the URL path to the project service invite HTTP endpoint.
func InviteProjectPath(projectID int32) string {
	return fmt.Sprintf("/projects/%v/invite", projectID)
}

// EditUserProjectPath returns the URL path to the project service edit user HTTP endpoint.
func EditUserProjectPath(projectID int32) string {
	return fmt.Sprintf("/projects/%v/roles", projectID)
}

// RemoveUserProjectPath returns the URL path to the project service remove user HTTP endpoint.
func RemoveUserProjectPath(projectID int32) string {
	return fmt.Sprintf("/projects/%v/members", projectID)
}

// AddStationProjectPath returns the URL path to the project service add station HTTP endpoint.
func AddStationProjectPath(projectID int32, stationID int32) string {
	return fmt.Sprintf("/projects/%v/stations/%v", projectID, stationID)
}

// RemoveStationProjectPath returns the URL path to the project service remove station HTTP endpoint.
func RemoveStationProjectPath(projectID int32, stationID int32) string {
	return fmt.Sprintf("/projects/%v/stations/%v", projectID, stationID)
}

// DeleteProjectPath returns the URL path to the project service delete HTTP endpoint.
func DeleteProjectPath(projectID int32) string {
	return fmt.Sprintf("/projects/%v", projectID)
}

// UploadPhotoProjectPath returns the URL path to the project service upload photo HTTP endpoint.
func UploadPhotoProjectPath(projectID int32) string {
	return fmt.Sprintf("/projects/%v/media", projectID)
}

// DownloadPhotoProjectPath returns the URL path to the project service download photo HTTP endpoint.
func DownloadPhotoProjectPath(projectID int32) string {
	return fmt.Sprintf("/projects/%v/media", projectID)
}

// GetProjectsForStationProjectPath returns the URL path to the project service get projects for station HTTP endpoint.
func GetProjectsForStationProjectPath(id int32) string {
	return fmt.Sprintf("/projects/station/%v", id)
}
