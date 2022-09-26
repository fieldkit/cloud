package jobs

import (
	"github.com/fieldkit/cloud/server/common/logging"
)

type CompletionIDs struct {
	generator *logging.IdGenerator
	ids       []string
}

func NewCompletionIDs() *CompletionIDs {
	return &CompletionIDs{
		generator: logging.NewIdGenerator(),
		ids:       make([]string, 0),
	}
}

func (c *CompletionIDs) Generate() string {
	id := c.generator.Generate()
	c.ids = append(c.ids, id)
	return id
}

func (c *CompletionIDs) IDs() []string {
	return c.ids
}
