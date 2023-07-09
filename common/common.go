package common

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type Route struct {
	Path string
	Handler func(http.ResponseWriter, *http.Request, *testing.T)
}

func Server(t *testing.T, routes []Route) *httptest.Server {
	cases := make(map[string]func(http.ResponseWriter, *http.Request, *testing.T))

	for _, route := range routes {
		cases[route.Path] = route.Handler
	}

	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		cases[req.URL.Path](res, req, t)
	}))
}


