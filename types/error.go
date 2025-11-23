package types

// ErrorResponse represents the generic error object returned by Wallex.
// Different endpoints may include:
//   - code: machine-readable error identifier
//   - message: human-readable explanation
//   - success: always present on failure
//   - result: result array of parameters about the issue,
//     i,e attempted market for order book which had a type like BTCUSDD
type ErrorResponse struct {
	BaseResponse
	Code int16 `json:"code"`
}
