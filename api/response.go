package api

// Response represents an operation.
type Response struct {
	Type ResponseType `json:"type" yaml:"type"`

	// Valid only for Sync responses
	Status     string `json:"status" yaml:"status"`
	StatusCode int    `json:"status_code" yaml:"status_code"`

	// Valid only for Error responses
	Code  int    `json:"error_code" yaml:"error_code"`
	Error string `json:"error" yaml:"error"`

	// Valid for Sync and Error responses
	Message interface{} `json:"message" yaml:"message"`
}

// ResponseType represents a valid LXD response type.
type ResponseType string

// Response types.
const (
	SyncResponse  ResponseType = "sync"
	ErrorResponse ResponseType = "error"
)
