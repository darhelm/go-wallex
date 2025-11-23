package types

import "time"

type CancelOrder struct {
	Symbol          string    `json:"symbol"`
	Type            string    `json:"type"`
	Side            string    `json:"side"`
	ClientOrderID   string    `json:"clientOrderId"`
	Price           string    `json:"price"`
	OrigQty         string    `json:"origQty"`
	OrigSum         string    `json:"origSum"`
	ExecutedSum     string    `json:"executedSum"`
	ExecutedQty     string    `json:"executedQty"`
	ExecutedPrice   string    `json:"executedPrice"`
	Sum             string    `json:"sum"`
	Fee             string    `json:"fee"`
	ExecutedPercent string    `json:"executedPercent"`
	Status          string    `json:"status"`
	Active          bool      `json:"active"`
	Fills           []any     `json:"fills"`
	TransactTime    int64     `json:"transactTime"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CancelOrderResponse is always {"status": "ok"} if nothing other than status code 200 is returned
type CancelOrderResponse struct {
	BaseResponse
	Result CancelOrder `json:"result"`
}
