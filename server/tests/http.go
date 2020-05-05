package tests

import (
	"net/http"
	"net/http/httptest"
)

func ExecuteRequest(req *http.Request, router http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}
