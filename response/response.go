package response

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grace/api"
	"io"
	"net/http"
)

type Response interface {
	Render(w http.ResponseWriter) error
	String() string
}

// REST API response.
type syncResponse struct {
	success     bool
	content     interface{}
	code        int
	contentType string
}

// WriteJSON encodes the body as JSON and sends it back to the client
// Accepts optional debugLogger that activates debug logging if non-nil.
func WriteJSON(w http.ResponseWriter, body any) error {
	output := w

	enc := json.NewEncoder(output)
	enc.SetEscapeHTML(false)
	err := enc.Encode(body)

	return err
}

func (r *syncResponse) Render(w http.ResponseWriter) error {
	var err error

	// Prepare the JSON response
	status := api.Success
	if !r.success {
		status = api.Failure
	}

	w.WriteHeader(r.code)

	if r.contentType == "json" {
		w.Header().Set("Content-Type", "application/json")
	}

	if err != nil {
		return err
	}

	// Handle JSON responses.
	resp := api.Response{
		Type:       api.SyncResponse,
		Status:     status.String(),
		StatusCode: int(status),
		Message:    r.content,
	}

	return WriteJSON(w, resp)
}

// ResponseNoContent return a new Response with no content.
func ResponseNoContent(code int) Response {
	return &syncResponse{success: true, code: code}
}

// ResponseJson return a new Response with json.
func ResponseJsonCustomStatusCode(code int, content interface{}) Response {
	return &syncResponse{success: true, code: code, content: content, contentType: "json"}
}

// ResponseJson return a new Response with json.
func ResponseJson(content interface{}) Response {
	return &syncResponse{success: true, code: http.StatusOK, content: content, contentType: "json"}
}

func (r *syncResponse) String() string {
	if r.code == http.StatusOK {
		return "success"
	}

	return "failure"
}

type errorResponse struct {
	code int    // Code to return in both the HTTP header and Code field of the response body.
	msg  string // Message to return in the Error field of the response body.
}

func (r *errorResponse) String() string {
	return r.msg
}

func (r *errorResponse) Render(w http.ResponseWriter) error {
	var output io.Writer

	buf := &bytes.Buffer{}
	output = buf

	resp := api.Response{
		Type:  api.ErrorResponse,
		Error: r.msg,
		Code:  r.code, // Set the error code in the Code field of the response body.
	}

	err := json.NewEncoder(output).Encode(resp)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(r.code) // Set the error code in the HTTP header response.

	_, err = fmt.Fprintln(w, buf.String())

	return err
}

// ErrorResponse returns an error response with the given code and msg.
func ErrorResponse(code int, msg string) Response {
	return &errorResponse{code, msg}
}

// BadRequest returns a bad request response (400) with the given error.
func BadRequest(err error) Response {
	return &errorResponse{http.StatusBadRequest, err.Error()}
}

// NotFound returns a not found response (404) with the given error.
func NotFound(err error) Response {
	message := "not found"
	if err != nil {
		message = err.Error()
	}

	return &errorResponse{http.StatusNotFound, message}
}

// InternalError returns an internal error response (500) with the given error.
func InternalError(err error) Response {
	return &errorResponse{http.StatusInternalServerError, err.Error()}
}

// NotImplemented returns a not implemented response (501) with the given error.
func NotImplemented(err error) Response {
	message := "not implemented"
	if err != nil {
		message = err.Error()
	}

	return &errorResponse{http.StatusNotImplemented, message}
}
