package data

import (
	"github.com/O-C-R/auth/id"
)

type Project struct {
	ID   id.ID  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Slug string `db:"slug" json:"slug"`
}

func NewProject(name, slug string) (*Project, error) {
	projectID, err := id.New()
	if err != nil {
		return nil, err
	}

	return &Project{
		ID:   projectID,
		Name: name,
		Slug: slug,
	}, nil
}
