package types

import "time"

// BaseOrder represents a single user order as returned by the Wallex account
// order endpoints. This model appears in create-order responses, open-orders
// queries, and order-status queries.
//
// This structure mirrors Wallex's data model for both active and historical
// user orders. All numeric values such as price, quantities, sums are returned
// as **number-strings**.
//
// Endpoint examples:
//
//	POST /v1/account/orders/add
//	GET  /v1/account/orders/open
//
// A typical Wallex order example:
//
//	{
//	  "symbol": "BTCUSDT",
//	  "type": "LIMIT",
//	  "side": "BUY",
//	  "price": "20950.00000000",
//	  "origQty": "0.0100",
//	  "executedQty": "0.0040",
//	  "executedPercent": 40.0,
//	  "status": "PARTIALLY_FILLED",
//	  "created_at": "2022-06-17T11:53:02Z"
//	}
type BaseOrder struct {
	Symbol          string    `json:"symbol"`
	Type            string    `json:"type"`
	Side            string    `json:"side"`
	Price           string    `json:"price"`
	OrigQty         string    `json:"origQty"`
	OrigSum         string    `json:"origSum"`
	ExecutedPrice   string    `json:"executedPrice"`
	ExecutedQty     string    `json:"executedQty"`
	ExecutedSum     string    `json:"executedSum"`
	ExecutedPercent float64   `json:"executedPercent"`
	Status          string    `json:"status"`
	Active          bool      `json:"active"`
	ClientOrderId   string    `json:"clientOrderId"`
	CreatedAt       time.Time `json:"created_at"`
}

// BaseOrderResponse wraps a single order object returned by Wallex.
//
// This is the standard response type for:
//
//	POST /v1/account/orders/add
//	GET  /v1/account/orders/<id>
//
// Response shape:
//
//	{ "success": true, "result": { ...order fields... } }
type BaseOrderResponse struct {
	BaseResponse
	Result BaseOrder `json:"result"`
}

// CreateOrderParams defines the request payload used to place a new order via:
//
//	POST /v1/account/orders/add
//
// Required fields:
//   - symbol
//   - type
//   - side
//   - price (for LIMIT)
//   - quantity
//
// clientOrderId is optional and can be assigned to track orders across systems.
type CreateOrderParams struct {
	Symbol        string `json:"symbol"`
	Type          string `json:"type"`
	Side          string `json:"side"`
	Price         string `json:"price"`
	Quantity      string `json:"quantity"`
	ClientOrderId string `json:"clientOrderId,omitempty"`
}

// OpenOrdersResponse contains all currently active (open) user orders.
// Returned by:
//
//	GET /v1/account/orders/open
//
// Response shape:
//
//	{
//	  "success": true,
//	  "result": {
//	    "orders": [ ...list of BaseOrder... ]
//	  }
//	}
type OpenOrdersResponse struct {
	BaseResponse
	Result struct {
		Orders []BaseOrder `json:"orders"`
	} `json:"result"`
}

// UserTradesParams defines the query parameters for retrieving private trade
// history from the Wallex account trade endpoint.
//
// Endpoint:
//
//	GET /v1/account/trades
//
// Both fields are optional. If provided, filtering is applied server-side.
type UserTradesParams struct {
	Symbol string `json:"symbol"`
	Side   string `json:"side"`
}

// UserTrade represents a trade execution belonging to the authenticated user.
// Returned by:
//
//	GET /v1/account/trades
//
// Unlike the public trades endpoint, this includes fee information.
type UserTrade struct {
	Symbol         string    `json:"symbol"`         // Market symbol
	Quantity       string    `json:"quantity"`       // Executed amount (string number)
	Price          string    `json:"price"`          // Executed price
	Sum            string    `json:"sum"`            // price * quantity
	Fee            string    `json:"fee"`            // Exact fee paid
	FeeCoefficient string    `json:"feeCoefficient"` // Fee rate (e.g. "0.001")
	FeeAsset       string    `json:"feeAsset"`       // Asset fee was deducted in
	IsBuyer        bool      `json:"isBuyer"`        // true if user was the buyer
	Timestamp      time.Time `json:"timestamp"`      // Execution timestamp
}

// UserTradesResponse wraps the array of account trade executions returned by:
//
//	GET /v1/account/trades
//
// Response shape:
//
//	{
//	  "success": true,
//	  "result": {
//	    "accountLatestTrades": [ ...list of UserTrade... ]
//	  }
//	}
type UserTradesResponse struct {
	BaseResponse
	Result struct {
		AccountLatestTrades []UserTrade `json:"accountLatestTrades"`
	} `json:"result"`
}
