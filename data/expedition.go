package data

import (
	"github.com/O-C-R/auth/id"
)

type Expedition struct {
	ID        id.ID  `db:"id" json:"id"`
	ProjectID id.ID  `db:"project_id" json:"id"`
	Name      string `db:"name" json:"name"`
	Slug      string `db:"slug" json:"slug"`
}

func NewExpedition(projectID id.ID, name, slug string) (*Expedition, error) {
	expeditionID, err := id.New()
	if err != nil {
		return nil, err
	}

	return &Expedition{
		ID:        expeditionID,
		ProjectID: projectID,
		Name:      name,
		Slug:      slug,
	}, nil
}
