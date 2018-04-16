package gonaturalist

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type CommentParentType string

const (
	AssessmentSection CommentParentType = "AssessmentSection"
	ListedTaxon       CommentParentType = "ListedTaxon"
	Observation       CommentParentType = "Observation"
	ObservationField  CommentParentType = "ObservationField"
	Post              CommentParentType = "Post"
	TaxonChange       CommentParentType = "TaxonChange"
)

type AddCommentOpt struct {
	ParentType CommentParentType `json:"parent_type"`
	ParentId   int64             `json:"parent_id"`
	Body       string            `json:"body"`
}

type UpdateCommentOpt struct {
	Id   int64  `json:"-"`
	Body string `json:"body"`
}

func (c *Client) GetObservationComments(observationId int64) (comments []*Comment, err error) {
	var result FullObservation

	u := c.buildUrl("/observations/%d.json", observationId)
	_, err = c.get(u, &result)
	if err != nil {
		return nil, err
	}

	return result.Comments, nil
}

func (c *Client) AddComment(opt *AddCommentOpt) error {
	u := c.buildUrl("/comments.json")

	bodyJson, err := json.Marshal(opt)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", u, bytes.NewReader(bodyJson))
	if err != nil {
		return err
	}
	var p interface{}
	err = c.execute(req, &p, http.StatusCreated)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateCommentBody(id int64, body string) error {
	updateCommentOpt := UpdateCommentOpt{
		Id:   id,
		Body: body,
	}
	return c.UpdateComment(&updateCommentOpt)
}

func (c *Client) UpdateComment(opt *UpdateCommentOpt) error {
	u := c.buildUrl("/comments/%d.json", opt.Id)

	bodyJson, err := json.Marshal(opt)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", u, bytes.NewReader(bodyJson))
	if err != nil {
		return err
	}
	var p interface{}
	err = c.execute(req, &p, http.StatusCreated)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteComment(id int64) error {
	u := c.buildUrl("/comments/%d.json", id)

	empty := make([]byte, 0)

	req, err := http.NewRequest("DELETE", u, bytes.NewReader(empty))
	if err != nil {
		return err
	}
	err = c.execute(req, nil, http.StatusCreated)
	if err != nil {
		return err
	}

	return nil
}
