// Code generated by goa v3.1.2, DO NOT EDIT.
//
// modules HTTP client types
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	modules "github.com/fieldkit/cloud/server/api/gen/modules"
)

// NewMetaResultOK builds a "modules" service "meta" endpoint result from a
// HTTP "OK" response.
func NewMetaResultOK(body interface{}) *modules.MetaResult {
	v := body
	res := &modules.MetaResult{
		Object: v,
	}

	return res
}
