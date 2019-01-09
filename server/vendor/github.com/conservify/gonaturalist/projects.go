package gonaturalist

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type GetProjectsOpt struct {
	Page *int
}

type SimpleProject struct {
	Id          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Terms       string    `json:"terms"`
	Description string    `json:"description"`
	Title       string    `json:"title"`
	UserId      int64     `json:"user_id"`
	IconUrl     string    `json:"icon_url"`
}

type FullProject struct {
	SimpleProject
	Terms                    string `json:"terms"`
	ProjectObservationsCount int    `json:"project_observations_count"`
	ProjectType              string `json:"project_type"`
}

type ProjectsPage struct {
	Paging   *PageHeaders
	Projects []SimpleProject
}

func (c *Client) GetProjects(opt *GetProjectsOpt) (*ProjectsPage, error) {
	var result []SimpleProject

	u := c.buildUrl("/projects.json")
	if opt != nil {
		v := url.Values{}
		if opt.Page != nil {
			v.Set("page", strconv.Itoa(*opt.Page))
		}
		if params := v.Encode(); params != "" {
			u += "?" + params
		}
	}
	p, err := c.get(u, &result)
	if err != nil {
		return nil, err
	}

	return &ProjectsPage{
		Projects: result,
		Paging:   p,
	}, nil
}

func (c *Client) GetProject(id interface{}) (*FullProject, error) {
	var result FullProject

	u := c.buildUrl("/projects/%s.json", id)
	_, err := c.get(u, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetProjectsByLogin(login string) (*ProjectsPage, error) {
	var result []SimpleProject

	u := c.buildUrl("/projects/user/%s.json", login)
	p, err := c.get(u, &result)
	if err != nil {
		return nil, err
	}

	return &ProjectsPage{
		Projects: result,
		Paging:   p,
	}, nil
}

func (c *Client) GetProjectMembers(id interface{}) error {
	return fmt.Errorf("Unimplemented")
}

func (c *Client) JoinProject(id interface{}) error {
	return fmt.Errorf("Unimplemented")
}

func (c *Client) LeaveProject(id interface{}) error {
	return fmt.Errorf("Unimplemented")
}

func (c *Client) AddObservationToProject(projectId interface{}, observationId interface{}) error {
	return fmt.Errorf("Unimplemented")
}
