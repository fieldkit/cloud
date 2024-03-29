// Code generated by goa v3.2.4, DO NOT EDIT.
//
// project HTTP client CLI support package
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"encoding/json"
	"fmt"
	"strconv"

	project "github.com/fieldkit/cloud/server/api/gen/project"
	goa "goa.design/goa/v3/pkg"
)

// BuildAddUpdatePayload builds the payload for the project add update endpoint
// from CLI flags.
func BuildAddUpdatePayload(projectAddUpdateBody string, projectAddUpdateProjectID string, projectAddUpdateAuth string) (*project.AddUpdatePayload, error) {
	var err error
	var body AddUpdateRequestBody
	{
		err = json.Unmarshal([]byte(projectAddUpdateBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"body\": \"Repellendus porro.\"\n   }'")
		}
	}
	var projectID int32
	{
		var v int64
		v, err = strconv.ParseInt(projectAddUpdateProjectID, 10, 32)
		projectID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for projectID, must be INT32")
		}
	}
	var auth string
	{
		auth = projectAddUpdateAuth
	}
	v := &project.AddUpdatePayload{
		Body: body.Body,
	}
	v.ProjectID = projectID
	v.Auth = auth

	return v, nil
}

// BuildDeleteUpdatePayload builds the payload for the project delete update
// endpoint from CLI flags.
func BuildDeleteUpdatePayload(projectDeleteUpdateProjectID string, projectDeleteUpdateUpdateID string, projectDeleteUpdateAuth string) (*project.DeleteUpdatePayload, error) {
	var err error
	var projectID int32
	{
		var v int64
		v, err = strconv.ParseInt(projectDeleteUpdateProjectID, 10, 32)
		projectID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for projectID, must be INT32")
		}
	}
	var updateID int64
	{
		updateID, err = strconv.ParseInt(projectDeleteUpdateUpdateID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value for updateID, must be INT64")
		}
	}
	var auth string
	{
		auth = projectDeleteUpdateAuth
	}
	v := &project.DeleteUpdatePayload{}
	v.ProjectID = projectID
	v.UpdateID = updateID
	v.Auth = auth

	return v, nil
}

// BuildModifyUpdatePayload builds the payload for the project modify update
// endpoint from CLI flags.
func BuildModifyUpdatePayload(projectModifyUpdateBody string, projectModifyUpdateProjectID string, projectModifyUpdateUpdateID string, projectModifyUpdateAuth string) (*project.ModifyUpdatePayload, error) {
	var err error
	var body ModifyUpdateRequestBody
	{
		err = json.Unmarshal([]byte(projectModifyUpdateBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"body\": \"Aut quia illo voluptas.\"\n   }'")
		}
	}
	var projectID int32
	{
		var v int64
		v, err = strconv.ParseInt(projectModifyUpdateProjectID, 10, 32)
		projectID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for projectID, must be INT32")
		}
	}
	var updateID int64
	{
		updateID, err = strconv.ParseInt(projectModifyUpdateUpdateID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value for updateID, must be INT64")
		}
	}
	var auth string
	{
		auth = projectModifyUpdateAuth
	}
	v := &project.ModifyUpdatePayload{
		Body: body.Body,
	}
	v.ProjectID = projectID
	v.UpdateID = updateID
	v.Auth = auth

	return v, nil
}

// BuildInvitesPayload builds the payload for the project invites endpoint from
// CLI flags.
func BuildInvitesPayload(projectInvitesAuth string) (*project.InvitesPayload, error) {
	var auth string
	{
		auth = projectInvitesAuth
	}
	v := &project.InvitesPayload{}
	v.Auth = auth

	return v, nil
}

// BuildLookupInvitePayload builds the payload for the project lookup invite
// endpoint from CLI flags.
func BuildLookupInvitePayload(projectLookupInviteToken string, projectLookupInviteAuth string) (*project.LookupInvitePayload, error) {
	var token string
	{
		token = projectLookupInviteToken
	}
	var auth string
	{
		auth = projectLookupInviteAuth
	}
	v := &project.LookupInvitePayload{}
	v.Token = token
	v.Auth = auth

	return v, nil
}

// BuildAcceptProjectInvitePayload builds the payload for the project accept
// project invite endpoint from CLI flags.
func BuildAcceptProjectInvitePayload(projectAcceptProjectInviteProjectID string, projectAcceptProjectInviteAuth string) (*project.AcceptProjectInvitePayload, error) {
	var err error
	var projectID int32
	{
		var v int64
		v, err = strconv.ParseInt(projectAcceptProjectInviteProjectID, 10, 32)
		projectID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for projectID, must be INT32")
		}
	}
	var auth string
	{
		auth = projectAcceptProjectInviteAuth
	}
	v := &project.AcceptProjectInvitePayload{}
	v.ProjectID = projectID
	v.Auth = auth

	return v, nil
}

// BuildRejectProjectInvitePayload builds the payload for the project reject
// project invite endpoint from CLI flags.
func BuildRejectProjectInvitePayload(projectRejectProjectInviteProjectID string, projectRejectProjectInviteAuth string) (*project.RejectProjectInvitePayload, error) {
	var err error
	var projectID int32
	{
		var v int64
		v, err = strconv.ParseInt(projectRejectProjectInviteProjectID, 10, 32)
		projectID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for projectID, must be INT32")
		}
	}
	var auth string
	{
		auth = projectRejectProjectInviteAuth
	}
	v := &project.RejectProjectInvitePayload{}
	v.ProjectID = projectID
	v.Auth = auth

	return v, nil
}

// BuildAcceptInvitePayload builds the payload for the project accept invite
// endpoint from CLI flags.
func BuildAcceptInvitePayload(projectAcceptInviteID string, projectAcceptInviteToken string, projectAcceptInviteAuth string) (*project.AcceptInvitePayload, error) {
	var err error
	var id int64
	{
		id, err = strconv.ParseInt(projectAcceptInviteID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value for id, must be INT64")
		}
	}
	var token *string
	{
		if projectAcceptInviteToken != "" {
			token = &projectAcceptInviteToken
		}
	}
	var auth string
	{
		auth = projectAcceptInviteAuth
	}
	v := &project.AcceptInvitePayload{}
	v.ID = id
	v.Token = token
	v.Auth = auth

	return v, nil
}

// BuildRejectInvitePayload builds the payload for the project reject invite
// endpoint from CLI flags.
func BuildRejectInvitePayload(projectRejectInviteID string, projectRejectInviteToken string, projectRejectInviteAuth string) (*project.RejectInvitePayload, error) {
	var err error
	var id int64
	{
		id, err = strconv.ParseInt(projectRejectInviteID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value for id, must be INT64")
		}
	}
	var token *string
	{
		if projectRejectInviteToken != "" {
			token = &projectRejectInviteToken
		}
	}
	var auth string
	{
		auth = projectRejectInviteAuth
	}
	v := &project.RejectInvitePayload{}
	v.ID = id
	v.Token = token
	v.Auth = auth

	return v, nil
}

// BuildAddPayload builds the payload for the project add endpoint from CLI
// flags.
func BuildAddPayload(projectAddBody string, projectAddAuth string) (*project.AddPayload, error) {
	var err error
	var body AddRequestBody
	{
		err = json.Unmarshal([]byte(projectAddBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"etag\": \"Aut aut sed aut et non accusamus.\",\n      \"logicalAddress\": 1510282003840073426,\n      \"meta\": \"Delectus quibusdam consequuntur nulla fugit voluptatem.\",\n      \"module\": \"Ab praesentium.\",\n      \"profile\": \"Quia molestiae repellendus et distinctio repellendus atque.\",\n      \"url\": \"Quaerat dolore sed.\",\n      \"version\": \"Qui sed.\"\n   }'")
		}
		if body.Bounds != nil {
			if err2 := ValidateProjectBoundsRequestBodyRequestBody(body.Bounds); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
		if err != nil {
			return nil, err
		}
	}
	var auth string
	{
		auth = projectAddAuth
	}
	v := &project.AddProjectFields{
		Name:         body.Name,
		Description:  body.Description,
		Goal:         body.Goal,
		Location:     body.Location,
		Tags:         body.Tags,
		Privacy:      body.Privacy,
		StartTime:    body.StartTime,
		EndTime:      body.EndTime,
		ShowStations: body.ShowStations,
	}
	if body.Bounds != nil {
		v.Bounds = marshalProjectBoundsRequestBodyRequestBodyToProjectProjectBounds(body.Bounds)
	}
	res := &project.AddPayload{
		Project: v,
	}
	res.Auth = auth

	return res, nil
}

// BuildUpdatePayload builds the payload for the project update endpoint from
// CLI flags.
func BuildUpdatePayload(projectUpdateBody string, projectUpdateProjectID string, projectUpdateAuth string) (*project.UpdatePayload, error) {
	var err error
	var body UpdateRequestBody
	{
		err = json.Unmarshal([]byte(projectUpdateBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"bounds\": {\n         \"max\": [\n            0.595389826652958,\n            0.6948214963772193,\n            0.837650615118303,\n            0.09669303147164357\n         ],\n         \"min\": [\n            0.11372084988344884,\n            0.5240321213868652,\n            0.3935396254118812\n         ]\n      },\n      \"description\": \"Aperiam ratione.\",\n      \"endTime\": \"Ut optio aut.\",\n      \"goal\": \"Quidem dolor architecto illo dolores quisquam.\",\n      \"location\": \"Enim natus perspiciatis quis.\",\n      \"name\": \"Dolores a ipsa.\",\n      \"privacy\": 513503072,\n      \"showStations\": true,\n      \"startTime\": \"Voluptas repudiandae cumque.\",\n      \"tags\": \"Eum a voluptatem et aspernatur unde.\"\n   }'")
		}
		if body.Bounds != nil {
			if err2 := ValidateProjectBoundsRequestBodyRequestBody(body.Bounds); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
		if err != nil {
			return nil, err
		}
	}
	var projectID int32
	{
		var v int64
		v, err = strconv.ParseInt(projectUpdateProjectID, 10, 32)
		projectID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for projectID, must be INT32")
		}
	}
	var auth string
	{
		auth = projectUpdateAuth
	}
	v := &project.AddProjectFields{
		Name:         body.Name,
		Description:  body.Description,
		Goal:         body.Goal,
		Location:     body.Location,
		Tags:         body.Tags,
		Privacy:      body.Privacy,
		StartTime:    body.StartTime,
		EndTime:      body.EndTime,
		ShowStations: body.ShowStations,
	}
	if body.Bounds != nil {
		v.Bounds = marshalProjectBoundsRequestBodyRequestBodyToProjectProjectBounds(body.Bounds)
	}
	res := &project.UpdatePayload{
		Project: v,
	}
	res.ProjectID = projectID
	res.Auth = auth

	return res, nil
}

// BuildGetPayload builds the payload for the project get endpoint from CLI
// flags.
func BuildGetPayload(projectGetProjectID string, projectGetAuth string) (*project.GetPayload, error) {
	var err error
	var projectID int32
	{
		var v int64
		v, err = strconv.ParseInt(projectGetProjectID, 10, 32)
		projectID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for projectID, must be INT32")
		}
	}
	var auth *string
	{
		if projectGetAuth != "" {
			auth = &projectGetAuth
		}
	}
	v := &project.GetPayload{}
	v.ProjectID = projectID
	v.Auth = auth

	return v, nil
}

// BuildListCommunityPayload builds the payload for the project list community
// endpoint from CLI flags.
func BuildListCommunityPayload(projectListCommunityAuth string) (*project.ListCommunityPayload, error) {
	var auth *string
	{
		if projectListCommunityAuth != "" {
			auth = &projectListCommunityAuth
		}
	}
	v := &project.ListCommunityPayload{}
	v.Auth = auth

	return v, nil
}

// BuildListMinePayload builds the payload for the project list mine endpoint
// from CLI flags.
func BuildListMinePayload(projectListMineAuth string) (*project.ListMinePayload, error) {
	var auth string
	{
		auth = projectListMineAuth
	}
	v := &project.ListMinePayload{}
	v.Auth = auth

	return v, nil
}

// BuildInvitePayload builds the payload for the project invite endpoint from
// CLI flags.
func BuildInvitePayload(projectInviteBody string, projectInviteProjectID string, projectInviteAuth string) (*project.InvitePayload, error) {
	var err error
	var body InviteRequestBody
	{
		err = json.Unmarshal([]byte(projectInviteBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"email\": \"Aut quia reiciendis laborum.\",\n      \"role\": 1665659981\n   }'")
		}
	}
	var projectID int32
	{
		var v int64
		v, err = strconv.ParseInt(projectInviteProjectID, 10, 32)
		projectID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for projectID, must be INT32")
		}
	}
	var auth string
	{
		auth = projectInviteAuth
	}
	v := &project.InviteUserFields{
		Email: body.Email,
		Role:  body.Role,
	}
	res := &project.InvitePayload{
		Invite: v,
	}
	res.ProjectID = projectID
	res.Auth = auth

	return res, nil
}

// BuildEditUserPayload builds the payload for the project edit user endpoint
// from CLI flags.
func BuildEditUserPayload(projectEditUserBody string, projectEditUserProjectID string, projectEditUserAuth string) (*project.EditUserPayload, error) {
	var err error
	var body EditUserRequestBody
	{
		err = json.Unmarshal([]byte(projectEditUserBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"email\": \"Non voluptate tempora dolor velit.\",\n      \"role\": 377366315\n   }'")
		}
	}
	var projectID int32
	{
		var v int64
		v, err = strconv.ParseInt(projectEditUserProjectID, 10, 32)
		projectID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for projectID, must be INT32")
		}
	}
	var auth string
	{
		auth = projectEditUserAuth
	}
	v := &project.EditUserFields{
		Email: body.Email,
		Role:  body.Role,
	}
	res := &project.EditUserPayload{
		Edit: v,
	}
	res.ProjectID = projectID
	res.Auth = auth

	return res, nil
}

// BuildRemoveUserPayload builds the payload for the project remove user
// endpoint from CLI flags.
func BuildRemoveUserPayload(projectRemoveUserBody string, projectRemoveUserProjectID string, projectRemoveUserAuth string) (*project.RemoveUserPayload, error) {
	var err error
	var body RemoveUserRequestBody
	{
		err = json.Unmarshal([]byte(projectRemoveUserBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"email\": \"Est nobis voluptates pariatur autem eum quasi.\"\n   }'")
		}
	}
	var projectID int32
	{
		var v int64
		v, err = strconv.ParseInt(projectRemoveUserProjectID, 10, 32)
		projectID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for projectID, must be INT32")
		}
	}
	var auth string
	{
		auth = projectRemoveUserAuth
	}
	v := &project.RemoveUserFields{
		Email: body.Email,
	}
	res := &project.RemoveUserPayload{
		Remove: v,
	}
	res.ProjectID = projectID
	res.Auth = auth

	return res, nil
}

// BuildAddStationPayload builds the payload for the project add station
// endpoint from CLI flags.
func BuildAddStationPayload(projectAddStationProjectID string, projectAddStationStationID string, projectAddStationAuth string) (*project.AddStationPayload, error) {
	var err error
	var projectID int32
	{
		var v int64
		v, err = strconv.ParseInt(projectAddStationProjectID, 10, 32)
		projectID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for projectID, must be INT32")
		}
	}
	var stationID int32
	{
		var v int64
		v, err = strconv.ParseInt(projectAddStationStationID, 10, 32)
		stationID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for stationID, must be INT32")
		}
	}
	var auth string
	{
		auth = projectAddStationAuth
	}
	v := &project.AddStationPayload{}
	v.ProjectID = projectID
	v.StationID = stationID
	v.Auth = auth

	return v, nil
}

// BuildRemoveStationPayload builds the payload for the project remove station
// endpoint from CLI flags.
func BuildRemoveStationPayload(projectRemoveStationProjectID string, projectRemoveStationStationID string, projectRemoveStationAuth string) (*project.RemoveStationPayload, error) {
	var err error
	var projectID int32
	{
		var v int64
		v, err = strconv.ParseInt(projectRemoveStationProjectID, 10, 32)
		projectID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for projectID, must be INT32")
		}
	}
	var stationID int32
	{
		var v int64
		v, err = strconv.ParseInt(projectRemoveStationStationID, 10, 32)
		stationID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for stationID, must be INT32")
		}
	}
	var auth string
	{
		auth = projectRemoveStationAuth
	}
	v := &project.RemoveStationPayload{}
	v.ProjectID = projectID
	v.StationID = stationID
	v.Auth = auth

	return v, nil
}

// BuildDeletePayload builds the payload for the project delete endpoint from
// CLI flags.
func BuildDeletePayload(projectDeleteProjectID string, projectDeleteAuth string) (*project.DeletePayload, error) {
	var err error
	var projectID int32
	{
		var v int64
		v, err = strconv.ParseInt(projectDeleteProjectID, 10, 32)
		projectID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for projectID, must be INT32")
		}
	}
	var auth string
	{
		auth = projectDeleteAuth
	}
	v := &project.DeletePayload{}
	v.ProjectID = projectID
	v.Auth = auth

	return v, nil
}

// BuildUploadPhotoPayload builds the payload for the project upload photo
// endpoint from CLI flags.
func BuildUploadPhotoPayload(projectUploadPhotoProjectID string, projectUploadPhotoContentType string, projectUploadPhotoContentLength string, projectUploadPhotoAuth string) (*project.UploadPhotoPayload, error) {
	var err error
	var projectID int32
	{
		var v int64
		v, err = strconv.ParseInt(projectUploadPhotoProjectID, 10, 32)
		projectID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for projectID, must be INT32")
		}
	}
	var contentType string
	{
		contentType = projectUploadPhotoContentType
	}
	var contentLength int64
	{
		contentLength, err = strconv.ParseInt(projectUploadPhotoContentLength, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value for contentLength, must be INT64")
		}
	}
	var auth string
	{
		auth = projectUploadPhotoAuth
	}
	v := &project.UploadPhotoPayload{}
	v.ProjectID = projectID
	v.ContentType = contentType
	v.ContentLength = contentLength
	v.Auth = auth

	return v, nil
}

// BuildDownloadPhotoPayload builds the payload for the project download photo
// endpoint from CLI flags.
func BuildDownloadPhotoPayload(projectDownloadPhotoProjectID string, projectDownloadPhotoSize string, projectDownloadPhotoIfNoneMatch string, projectDownloadPhotoAuth string) (*project.DownloadPhotoPayload, error) {
	var err error
	var projectID int32
	{
		var v int64
		v, err = strconv.ParseInt(projectDownloadPhotoProjectID, 10, 32)
		projectID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for projectID, must be INT32")
		}
	}
	var size *int32
	{
		if projectDownloadPhotoSize != "" {
			var v int64
			v, err = strconv.ParseInt(projectDownloadPhotoSize, 10, 32)
			val := int32(v)
			size = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for size, must be INT32")
			}
		}
	}
	var ifNoneMatch *string
	{
		if projectDownloadPhotoIfNoneMatch != "" {
			ifNoneMatch = &projectDownloadPhotoIfNoneMatch
		}
	}
	var auth *string
	{
		if projectDownloadPhotoAuth != "" {
			auth = &projectDownloadPhotoAuth
		}
	}
	v := &project.DownloadPhotoPayload{}
	v.ProjectID = projectID
	v.Size = size
	v.IfNoneMatch = ifNoneMatch
	v.Auth = auth

	return v, nil
}

// BuildProjectsStationPayload builds the payload for the project projects
// station endpoint from CLI flags.
func BuildProjectsStationPayload(projectProjectsStationID string, projectProjectsStationAuth string) (*project.ProjectsStationPayload, error) {
	var err error
	var id int32
	{
		var v int64
		v, err = strconv.ParseInt(projectProjectsStationID, 10, 32)
		id = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for id, must be INT32")
		}
	}
	var auth string
	{
		auth = projectProjectsStationAuth
	}
	v := &project.ProjectsStationPayload{}
	v.ID = id
	v.Auth = auth

	return v, nil
}
