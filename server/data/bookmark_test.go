package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBookmark(t *testing.T) {
	assert := assert.New(t)
	b, err := ParseBookmark(`{"v":1,"g":[[[[[159],[2],[-8640000000000000,8640000000000000],[],0,0]]]],"s":[]}`)
	assert.NoError(err)
	assert.NotNil(b)
}
