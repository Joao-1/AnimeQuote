package common

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	routes := []Route{
		{
			Path: "/route/1",
			Handler: func(res http.ResponseWriter, req *http.Request, t *testing.T) {
				res.WriteHeader(http.StatusOK)
			},
		},
		{
			Path: "/route/2",
			Handler: func(res http.ResponseWriter, req *http.Request, t *testing.T) {
				res.WriteHeader(http.StatusInternalServerError)
			},
		},
		{
			Path: "/route/3",
			Handler: func(res http.ResponseWriter, req *http.Request, t *testing.T) {
				res.WriteHeader(http.StatusNotFound)
			},
		},
	}

	server := Server(t, routes)

	expectedStatus := []int{200, 500, 404}
	receivedStatus := []int{}

	for _, route := range routes {
		res, err := http.Get(server.URL + route.Path)
		assert.Nil(t, err)

		receivedStatus = append(receivedStatus, res.StatusCode)
	}
	
	assert.Equal(t, expectedStatus, receivedStatus)
}