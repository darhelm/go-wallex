package wallex

import (
	"encoding/json"
	"fmt"
	"strconv"

	t "go-wallex/types"
)

type GoWallexError struct {
	Message string
	Err     error
}

func (e *GoWallexError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *GoWallexError) Unwrap() error { return e.Err }

// RequestError represents failures that occur before the Wallex server responds.
//
// These include:
//   - JSON marshal/unmarshal failures
//   - HTTP request creation errors
//   - Network/transport errors
//   - Response body read failures
//
// These errors indicate the request was never successfully processed by Wallex.
type RequestError struct {
	GoWallexError
	Operation string
}

// APIError represents any non-2xx error response returned by the Wallex API.
//
// Wallex generally returns one of the following shapes:
//
//	{
//	  "success": false,
//	  "code": 1201,
//	  "message": "invalid API key format",
//	  "result": {}
//	}
//
//	{ "detail": "missing required parameter" }
//
//	{ "message": "something went wrong" }
//
//	{ ...undocumented fields... }
//
// This struct attempts to capture all information:
//   - Standard Wallex fields (success, code, message, result)
//   - Arbitrary fields via Fields map
//   - HTTP status code via StatusCode
//
// Result is left as json.RawMessage because Wallex may return:
//   - {}  (empty object)
//   - null
//   - array
//   - string
//   - undocumented structures
type APIError struct {
	GoWallexError

	StatusCode int

	// Message provides a human-readable description of the status.
	Message string `json:"message"`

	// Status indicates failure state, typically "false" or "true".
	Success bool `json:"success"`

	// Code is the server status code indicating request status.
	Code int16 `json:"code"`

	// Result array or object of available response
	Result json.RawMessage `json:"result"`

	// Map of all parsed key->values for inspection (similar to go-bitpin)
	Fields map[string][]string
}

// parseErrorResponse creates an APIError from a Wallex non-2xx response.
//
// It attempts the most complete extraction possible by processing documented and
// undocumented error shapes observed in Wallex responses.
//
// Extraction logic:
//
//  1. Try decoding into the standard BaseResponse-like struct (t.ErrorResponse).
//  2. Parse all fields of the raw JSON into a generic map[string]any.
//  3. Capture string, array, and non-string values into the Fields map.
//  4. Preserve the "detail" or "message" fields if no message is found.
//  5. Leave apiErr.Result as raw JSON to avoid type assumptions.
//
// Always returns an *APIError that is safe to present to the caller.
func parseErrorResponse(statusCode int, respBody []byte) *APIError {
	apiErr := &APIError{
		StatusCode: statusCode,
		Fields:     make(map[string][]string),
	}

	// #1 — Attempt to parse official Wallex error format
	var base t.ErrorResponse
	_ = json.Unmarshal(respBody, &base)

	if base.Message != "" {
		apiErr.Message = base.Message
	}
	if base.Code > 0 {
		apiErr.Code = base.Code
		apiErr.Fields["code"] = []string{strconv.Itoa(int(base.Code))}
	}
	if base.Message != "" {
		apiErr.Message = base.Message
		apiErr.Fields["message"] = []string{base.Message}
	}

	// #2 — Parse raw JSON object for extra fields, including "detail"
	raw := map[string]any{}
	_ = json.Unmarshal(respBody, &raw)

	for k, v := range raw {
		switch val := v.(type) {
		case string:
			apiErr.Fields[k] = []string{val}
			if k == "detail" {
				apiErr.Result = json.RawMessage(val)
				if apiErr.Message == "" {
					apiErr.Message = val
				}
			}
		case []any:
			// convert []any -> []string
			strs := make([]string, 0, len(val))
			for _, item := range val {
				strs = append(strs, fmt.Sprintf("%v", item))
			}
			apiErr.Fields[k] = strs
		default:
			apiErr.Fields[k] = []string{fmt.Sprintf("%v", v)}
		}
	}

	// #3 — If message is still empty, fallback
	if apiErr.Message == "" {
		apiErr.Message = fmt.Sprintf("Wallex API error (%d)", statusCode)
	}

	apiErr.GoWallexError.Message = apiErr.Message
	return apiErr
}
