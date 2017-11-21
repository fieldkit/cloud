package backend

import (
	"fmt"
	"strconv"
	"strings"
)

type PagingToken struct {
	time int64
	page int32
}

func NewPagingTokenFromString(s string) *PagingToken {
	fields := strings.Split(s, ",")
	if len(fields) != 2 {
		return nil
	}
	time, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return nil
	}
	page, err := strconv.ParseInt(fields[1], 10, 32)
	if err != nil {
		return nil
	}

	return &PagingToken{
		time: time,
		page: int32(page),
	}
}

func (t *PagingToken) String() string {
	return fmt.Sprintf("%d,%d", t.time, t.page)
}
