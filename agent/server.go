package agent

import (
	"fmt"
	"net/http"

	"grace/response"

	"github.com/gorilla/mux"
)

// APIEndpoint represents a URL in our API.
type APIEndpoint struct {
	Name   string // Name for this endpoint.
	Path   string // Path pattern for this endpoint.
	Get    APIEndpointAction
	Put    APIEndpointAction
	Post   APIEndpointAction
	Delete APIEndpointAction
	Patch  APIEndpointAction
}

// APIEndpointAction represents an action on an API endpoint.
type APIEndpointAction struct {
	Handler func(req *http.Request) response.Response
}

func CreateCmd(restAPI *mux.Router, version string, c APIEndpoint) {
	var uri string
	if c.Path == "" {
		uri = fmt.Sprintf("/%s", version)
	} else {
		uri = fmt.Sprintf("/%s/%s", version, c.Path)
	}

	restAPI.HandleFunc(uri, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Actually process the request
		var resp response.Response

		handleRequest := func(action APIEndpointAction) response.Response {
			if action.Handler == nil {
				return response.NotImplemented(nil)
			}

			return action.Handler(req)
		}

		switch req.Method {
		case "GET":
			resp = handleRequest(c.Get)
		case "PUT":
			resp = handleRequest(c.Put)
		case "POST":
			resp = handleRequest(c.Post)
		case "DELETE":
			resp = handleRequest(c.Delete)
		case "PATCH":
			resp = handleRequest(c.Patch)
		default:
			resp = response.NotFound(fmt.Errorf("Method %q not found", req.Method))
		}

		// Handle errors
		err := resp.Render(w)
		if err != nil {
			_ = response.InternalError(err).Render(w)
		}
	}).Methods("GET", "PUT", "POST", "DELETE", "PATCH")
}
