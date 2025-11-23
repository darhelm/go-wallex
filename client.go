package wallex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	t "github.com/darhelm/go-wallex/types"
	u "github.com/darhelm/go-wallex/utils"
)

// Constants defining the API base URL and version.
const (
	// BaseUrl is the root URL for the Wallex Market API.
	BaseUrl = "https://api.wallex.ir"
)

// ClientOptions represents the configuration options for creating a new API client.
// These options allow customization of the HTTP client, authentication tokens,
// API credentials, and automatic authentication/refresh behaviors.
type ClientOptions struct {
	// HttpClient is the custom HTTP client to be used for API requests.
	// If nil, the default HTTP client is used.
	HttpClient *http.Client

	// Timeout specifies the request timeout duration for the HTTP client.
	Timeout time.Duration

	// BaseUrl is the base URL of the API. Defaults to the constant BaseUrl
	// if not provided.
	BaseUrl string

	Version string

	// ApiKey is the token used for authenticated API requests.
	ApiKey string
}

// Client represents the API client for interacting with the Wallex Market API.
// It manages authentication, base URL, and API requests.
type Client struct {
	// HttpClient is the HTTP client used for API requests.
	// Defaults to the Go standard library's http.DefaultClient.
	HttpClient *http.Client

	// BaseUrl is the base URL of the API used by this client.
	// Defaults to the constant BaseUrl.
	BaseUrl string

	Version string

	// ApiKey is the API key for authentication.
	ApiKey string
}

// NewClient creates a new Wallex API client.
//
// This client is a lightweight wrapper around the Wallex REST API.
// Wallex uses only one authentication mechanism: X-API-Key.
//
// Parameters:
//   - opts.HttpClient: Optional custom HTTP client (default: http.DefaultClient).
//   - opts.Timeout: Request timeout used if a custom client is not provided.
//   - opts.BaseUrl: Override API base URL (default: https://api.wallex.ir).
//   - opts.Version: Optional API version prefix.
//   - opts.ApiKey: API key for authenticated endpoints.
//
// Behavior:
//   - Does NOT perform login (Wallex has no login endpoint).
//   - Does NOT refresh tokens (Wallex API keys are static).
//
// Returns:
//   - *Client ready to make Wallex API requests.
func NewClient(opts ClientOptions) (*Client, error) {
	client := &Client{
		BaseUrl: BaseUrl,
		ApiKey:  opts.ApiKey,
	}

	if opts.BaseUrl != "" {
		client.BaseUrl = opts.BaseUrl
	}

	if opts.HttpClient != nil {
		client.HttpClient = opts.HttpClient
	} else {
		client.HttpClient = &http.Client{
			Timeout: opts.Timeout,
		}
	}

	return client, nil
}

// assertAuth ensures the client contains a non-empty API key.
//
// Used internally by authenticated requests.
// Returns an error if ApiKey is empty.
func assertAuth(client *Client) error {
	if client.ApiKey == "" {
		return &GoWallexError{
			Message: "API Key is empty",
			Err:     nil,
		}
	}
	return nil
}

// createApiURI constructs the full Wallex API URL by combining:
//
//	BaseUrl + "/" + version + endpoint
//
// Wallex endpoints may or may not use a version prefix (e.g. v1, v2).
//
// Example:
//
//	createApiURI("/depth?symbol=BTCUSDT", "v1")
//	→ "https://api.wallex.ir/v1/depth?symbol=BTCUSDT"
func (c *Client) createApiURI(endpoint string, version string) string {
	return fmt.Sprintf("%s/%s%s", c.BaseUrl, version, endpoint)
}

// Request performs an HTTP request to the Wallex API.
//
// Capabilities:
//   - GET: URL-encoded query parameters generated from `body`.
//   - POST: JSON-encoded request body.
//   - Adds X-API-Key header when auth=true.
//   - Parses Wallex-style success/error envelopes.
//   - Unmarshals successful JSON responses into `result`.
//
// Wallex Error Handling:
//   - Non-2xx responses are passed to parseErrorResponse(), which extracts:
//   - success=false
//   - code (if present)
//   - message
//   - detail or undocumented fields
//
// Returns:
//   - nil on success
//   - *RequestError for network/JSON failures
//   - *APIError for Wallex server-side errors
func (c *Client) Request(method string, url string, auth bool, body interface{}, result interface{}) error {
	var reqBody []byte
	var err error

	if method == "GET" {
		if body != nil {
			urlParams, err := u.StructToURLParams(body)
			if err != nil {
				return &RequestError{
					GoWallexError: GoWallexError{
						Message: "failed to convert struct to URL params",
						Err:     err,
					},
					Operation: "preparing request parameters",
				}
			}
			url += "?" + urlParams
		}
	}

	if method == "POST" {
		if body != nil {
			reqBody, err = json.Marshal(body)
			if err != nil {
				return &RequestError{
					GoWallexError: GoWallexError{
						Message: "failed to marshal request body",
						Err:     err,
					},
					Operation: "preparing request body",
				}
			}
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return &RequestError{
			GoWallexError: GoWallexError{
				Message: "failed to create request",
				Err:     err,
			},
			Operation: "creating request",
		}
	}

	req.Header.Set("Content-Type", "application/json")

	if auth {
		if err := assertAuth(c); err != nil {
			return &GoWallexError{
				Message: "authentication validation failed",
				Err:     err,
			}
		}

		req.Header.Set("X-API-Key", c.ApiKey)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return &RequestError{
			GoWallexError: GoWallexError{
				Message: "failed to send request",
				Err:     err,
			},
			Operation: "sending request",
		}
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RequestError{
			GoWallexError: GoWallexError{
				Message: "failed to read response body",
				Err:     err,
			},
			Operation: "reading response",
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseErrorResponse(resp.StatusCode, respBody)
	}

	if result != nil {
		if err = json.Unmarshal(respBody, result); err != nil {
			return &RequestError{
				GoWallexError: GoWallexError{
					Message: "failed to unmarshal response",
					Err:     err,
				},
				Operation: "parsing response",
			}
		}
	}

	return nil
}

// ApiRequest is a helper that builds the Wallex API URL using createApiURI(),
// then executes the request through Request().
//
// It simply forwards:
//   - method ("GET", "POST", "DELETE")
//   - endpoint (e.g. "/account/orders")
//   - version (e.g. "v1")
//   - auth flag
//   - body and output result pointer
//
// Most Wallex endpoints live under version "v1" unless documented otherwise.
func (c *Client) ApiRequest(method, endpoint string, version string, auth bool, body interface{}, result interface{}) error {
	url := c.createApiURI(endpoint, version)
	return c.Request(method, url, auth, body, result)
}

// GetMarketsInfo retrieves metadata for all trading symbols on Wallex.
//
// Endpoint:
//
//	GET /v1/markets
//
// This returns market specifications including:
//   - base/quote assets
//   - precisions
//   - minQty / maxQty / minNotional
//   - tickSize / stepSize
//   - 24h & 7d statistics
//
// Authentication: NOT required.
// Rate Limit: 100 requests/sec (global Wallex limit).
func (c *Client) GetMarketsInfo() (*t.MarketInformation, error) {
	var marketInfo *t.MarketInformation
	err := c.ApiRequest("GET", "/markets", "v1", false, nil, &marketInfo)
	if err != nil {
		return nil, err
	}
	return marketInfo, nil
}

// GetOrderBook retrieves the current order book for a specific market.
//
// Endpoint:
//
//	GET /v1/depth?symbol={SYMBOL}
//
// Returns aggregated bid/ask levels.
// Authentication: NOT required.
// Rate Limit: 100 requests/sec.
//
// Example:
//
//	depth, _ := client.GetOrderBook("BTCUSDT")
func (c *Client) GetOrderBook(symbol string) (*t.Depth, error) {
	var depth *t.Depth
	err := c.ApiRequest("GET", fmt.Sprintf("/depth?symbol=%s", symbol), "v1", false, nil, &depth)
	if err != nil {
		return nil, err
	}
	return depth, nil
}

// GetAllOrderBooks retrieves the order books for ALL markets in a single call.
//
// Endpoint:
//
//	GET /v2/depth/all
//
// Response:
//
//	result: map[symbol]OrderBook
//
// Authentication: NOT required.
// Rate Limit: 100 requests/sec (heavy endpoint).
func (c *Client) GetAllOrderBooks() (*t.AllDepths, error) {
	var depths *t.AllDepths
	err := c.ApiRequest("GET", "/depth/all", "v2", false, nil, &depths)
	if err != nil {
		return nil, err
	}
	return depths, nil
}

// GetRecentTrades retrieves recent executed trades for a given symbol.
//
// Endpoint:
//
//	GET /v1/trades?symbol={SYMBOL}
//
// Returned fields include:
//   - price
//   - quantity
//   - sum
//   - isBuyOrder
//   - timestamp
//
// Authentication: NOT required.
// Rate Limit: 100 requests/sec.
func (c *Client) GetRecentTrades(symbol string) (*t.Trades, error) {
	var trades *t.Trades
	err := c.ApiRequest("GET", fmt.Sprintf("/trades?symbol=%s", symbol), "v1", false, nil, &trades)
	if err != nil {
		return nil, err
	}
	return trades, nil
}

// GetWallets retrieves the authenticated user's wallet balances.
//
// Endpoint:
//
//	GET /v1/account/balances
//
// Returns balances for all assets associated with the API key.
//
// Authentication: REQUIRED (X-API-Key).
// Rate Limit: 100 requests/sec.
func (c *Client) GetWallets() (*t.Wallets, error) {
	var wallets *t.Wallets
	err := c.ApiRequest("GET", "/account/balances", "v1", true, nil, &wallets)
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

// CreateOrder submits a new trading order on Wallex.
//
// Endpoint:
//
//	POST /v1/account/orders
//
// Required fields (CreateOrderParams):
//   - symbol
//   - type      ("LIMIT" or "MARKET")
//   - side      ("BUY" or "SELL")
//   - price     (only for LIMIT)
//   - quantity
//
// Returns:
//   - BaseOrder with server-evaluated order status.
//
// Authentication: REQUIRED.
// Rate Limit: 100 req/sec.
func (c *Client) CreateOrder(params t.CreateOrderParams) (*t.BaseOrderResponse, error) {
	var orderStatus *t.BaseOrderResponse
	err := c.ApiRequest("POST", "/account/orders", "v1", true, params, &orderStatus)
	if err != nil {
		return nil, err
	}
	return orderStatus, nil
}

// CancelOrder cancels an active user order.
//
// Endpoint:
//
//	DELETE /v1/account/orders?clientOrderId={ID}
//
// On success:
//   - result contains the updated order (status="CANCELED").
//
// Authentication: REQUIRED.
// Rate Limit: 100 req/sec.
//
// If clientOrderId is invalid or order already closed,
// Wallex returns success=false with an API error.
func (c *Client) CancelOrder(clientOrderId string) (*t.CancelOrderResponse, error) {
	var cancelOrderStatus *t.CancelOrderResponse
	err := c.ApiRequest("DELETE", fmt.Sprintf("/account/orders?clientOrderId=%s", clientOrderId), "v1", true, nil, &cancelOrderStatus)
	if err != nil {
		return nil, err
	}
	return cancelOrderStatus, nil
}

// GetOpenOrders retrieves all active (not filled/cancelled) orders.
//
// Endpoint:
//
//	GET /v1/account/openOrders
//	GET /v1/account/openOrders?symbol={SYMBOL}
//
// Returns:
//   - A list of BaseOrder objects.
//
// Authentication: REQUIRED.
// Rate Limit: 100 req/sec.
func (c *Client) GetOpenOrders(symbol string) (*t.OpenOrdersResponse, error) {
	var orders *t.OpenOrdersResponse

	var endPoint = "/account/openOrders"
	if symbol != "" {
		endPoint = fmt.Sprintf("%s?symbol=%s", endPoint, symbol)
	}

	err := c.ApiRequest("GET", endPoint, "v1", true, nil, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetOrderStatus retrieves full details for a specific user order.
//
// Endpoint:
//
//	GET /v1/account/orders/{clientOrderId}
//
// Returns full BaseOrder including:
//   - status
//   - executedQty / executedSum
//   - executedPercent
//   - price, quantity
//   - timestamps
//
// Authentication: REQUIRED.
// Rate Limit: 100 req/sec.
//
// Errors:
//   - Missing or invalid clientOrderId
//   - Order does not belong to this API key
func (c *Client) GetOrderStatus(clientOrderId string) (*t.BaseOrderResponse, error) {
	var orders *t.BaseOrderResponse
	if clientOrderId == "" {
		return nil, &GoWallexError{
			Message: "client order id is required for getting order status",
			Err:     nil,
		}
	}

	err := c.ApiRequest("GET", fmt.Sprintf("/account/orders/%s", clientOrderId), "v1", true, nil, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetUserTrades retrieves private trade history for the authenticated user.
//
// Endpoint:
//
//	GET /v1/account/trades
//
// Optional filters:
//   - symbol
//   - side ("BUY" / "SELL")
//
// Each trade includes:
//   - price, quantity, sum
//   - fee, feeCoefficient, feeAsset
//   - isBuyer
//   - timestamp
//
// Authentication: REQUIRED.
// Rate Limit: 100 req/sec.
func (c *Client) GetUserTrades(params t.UserTradesParams) (*t.UserTradesResponse, error) {
	var trades *t.UserTradesResponse
	err := c.ApiRequest("GET", "/account/trades", "v1", true, params, &trades)
	if err != nil {
		return nil, err
	}
	return trades, nil
}
