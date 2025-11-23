package types

// CancelOrderResponse is always {"status": "ok"} if nothing other than status code 200 is returned
type CancelOrderResponse struct {
	Status string `json:"status"`
}
