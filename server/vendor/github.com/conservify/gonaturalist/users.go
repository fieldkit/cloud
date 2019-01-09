package gonaturalist

import (
	"time"
)

type PrivateUser struct {
	Name      string
	Email     string
	Id        int64
	Login     string
	Uri       string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ObservationsCount int32 `json:"observations_count"`

	LifeListId        int32  `json:"life_list_id"`
	LifeListTaxaCount int32  `json:"life_list_taxa_count"`
	TimeZone          string `json:"time_zone"`

	IconUrl         string `json:"icon_url"`
	IconContentType string `json:"icon_content_type"`
	IconFileName    string `json:"icon_file_name"`
	IconFileSize    int32  `json:"icon_file_size"`
}

func (c *Client) GetCurrentUser() (*PrivateUser, error) {
	var result PrivateUser

	_, err := c.get(c.buildUrl("/users/edit.json"), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) AddUser() error {
	return nil
}
