package types

import "encoding/json"

// BaseResponse defines the base response type for all wallex endpoints
type BaseResponse struct {
	// Message provides a human-readable description of the status.
	Message string `json:"message"`

	// Status indicates failure state, typically "false" or "true".
	Success bool `json:"success"`

	// Result array or object of available response
	Result json.RawMessage `json:"result"`
}
