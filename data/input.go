package data

import (
	"github.com/O-C-R/auth/id"
)

type Input struct {
	ID           id.ID  `db:"id" json:"id"`
	ExpeditionID id.ID  `db:"expedition_id" json:"expedition_id"`
	Name         string `db:"name" json:"name"`
	Slug         string `db:"slug" json:"slug"`
}

func NewInput(expeditionID id.ID, name, slug string) (*Input, error) {
	inputID, err := id.New()
	if err != nil {
		return nil, err
	}

	return &Input{
		ID:           inputID,
		ExpeditionID: expeditionID,
		Name:         name,
		Slug:         slug,
	}, nil
}

type Request struct {
	ID       id.ID  `db:"id" json:"id"`
	InputID  id.ID  `db:"input_id" json:"input_id"`
	Format   string `db:"format" json:"format"`
	Checksum id.ID  `db:"checksum" json:"checksum"`
	Data     []byte `db:"data" json:"-"`
}

func NewRequest(inputID id.ID) (*Request, error) {
	requestID, err := id.New()
	if err != nil {
		return nil, err
	}

	return &Request{
		ID:      requestID,
		InputID: inputID,
	}, nil
}
