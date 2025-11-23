package types

import (
	"time"
)

// Direction represents the buy/sell directional distribution of recent trades
// for a symbol. The Wallex API returns this structure as part of the market
// statistics block under GET /v1/markets.
//
// Example JSON:
//
//	"direction": { "SELL": 51, "BUY": 49 }
//
// Both fields represent percentages (0–100). The JSON keys are uppercase.
type Direction struct {
	Sell int `json:"SELL"`
	Buy  int `json:"BUY"`
}

// Stats contains 24-hour and 7-day market statistics for a trading pair,
// including last trade data, price extremes, quote volume, and trade counts.
// This object appears inside the SymbolInfo struct in GET /v1/markets.
//
// Numeric price and volume values are returned as **number strings** in the
// Wallex API, while percentage changes are returned as numeric values.
//
// Relevant endpoint:
//
//	GET /v1/markets (result.symbols[*].stats)
type Stats struct {
	BidPrice       string    `json:"bidPrice"`
	AskPrice       string    `json:"askPrice"`
	DayCh          float64   `json:"24h_ch"`
	WeekCh         float64   `json:"7d_ch"`
	DayVolume      string    `json:"24h_volume"`
	WeekVolume     string    `json:"7d_volume"`
	QuoteVolumeDay string    `json:"24h_quoteVolume"`
	HighPriceDay   string    `json:"24h_highPrice"`
	LowPriceDay    string    `json:"24h_lowPrice"`
	LastPrice      string    `json:"lastPrice"`
	LastQty        string    `json:"lastQty"`
	LastTradeSide  string    `json:"lastTradeSide"`
	BidVolume      string    `json:"bidVolume"`
	AskVolume      string    `json:"askVolume"`
	BidCount       int8      `json:"bidCount"`
	AskCount       int8      `json:"askCount"`
	Direction      Direction `json:"direction"`
}

// SymbolInfo represents the complete metadata of a trading symbol (market pair)
// on Wallex. This structure includes asset identifiers, precision rules, minimum
// trading constraints, tick sizes, and statistical market data.
// Returned inside result.symbols from:
//
//	GET /v1/markets
//
// Notes:
//   - stepSize and tickSize specify decimal precision (number of digits).
//   - minQty, maxQty, and minNotional are numeric and may include fractions.
//   - stats contains real-time 24h/7d market information.
type SymbolInfo struct {
	Symbol             string    `json:"symbol"`
	BaseAsset          string    `json:"baseAsset"`
	BaseAssetPrecision int8      `json:"baseAssetPrecision"`
	QuoteAsset         string    `json:"quoteAsset"`
	QuotePrecision     int8      `json:"quotePrecision"`
	FaName             string    `json:"faName"`
	FaBaseAsset        string    `json:"faBaseAsset"`
	FaQuoteAsset       string    `json:"faQuoteAsset"`
	StepSize           int64     `json:"stepSize"`
	TickSize           int64     `json:"tickSize"`
	MinQty             int64     `json:"minQty"`
	MinNotional        int64     `json:"minNotional"`
	Stats              Stats     `json:"stats"`
	CreatedAt          time.Time `json:"createdAt"`
}

// Symbols is a container type wrapping a map of symbol identifiers to their
// corresponding metadata.
// It appears under result.symbols in GET /v1/markets.
type Symbols struct {
	Symbols map[string]SymbolInfo `json:"symbols"`
}

// MarketInformation wraps the result from /markets
type MarketInformation struct {
	BaseResponse
	Result map[string]Symbols `json:"result"`
}

// Order represents a single order book level (price level).
// The price and sum fields are number-strings, while quantity is numeric.
//
// Appears under:
//
//	GET /v1/depth
//	GET /v2/depth/all
type Order struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
	Sum      string  `json:"sum"`
}

// OrderBook holds the full depth for a symbol: arrays of bid and ask levels.
// Returned as result for GET /v1/depth and as map values for GET /v2/depth/all
type OrderBook struct {
	Ask []Order `json:"ask"`
	Bid []Order `json:"bid"`
}

// Depth wraps an orderbook response for a single market.
//
// Response shape:
//
//	{ "success": true, "result": { "ask": [...], "bid": [...] } }
type Depth struct {
	BaseResponse
	Result OrderBook `json:"result"`
}

// AllDepths wraps orderbooks for all markets simultaneously.
//
// Response shape:
//
//	{ "success": true, "result": { "BTCUSDT": { ... }, "ETHUSDT": { ... }, ... } }
type AllDepths struct {
	BaseResponse
	Result map[string]Depth `json:"result"`
}

// Trade represents a single executed trade on Wallex, returned in the recent
// trades endpoint. Prices and quantities are number-strings; timestamp is ISO8601.
type Trade struct {
	Symbol     string    `json:"symbol"`     // Symbol the trade belongs to
	Quantity   string    `json:"quantity"`   // Traded quantity (number-string)
	Price      string    `json:"price"`      // Trade price (number-string)
	Sum        string    `json:"sum"`        // price * quantity (number-string)
	IsBuyOrder bool      `json:"isBuyOrder"` // true = taker was buying
	Timestamp  time.Time `json:"timestamp"`  // Execution timestamp
}

// LatestTrades wraps the trades array under result.latestTrades.
type LatestTrades struct {
	LatestTrades []Trade `json:"latestTrades"`
}

// Trades wraps the full response for GET /v1/trades.
type Trades struct {
	BaseResponse
	Result LatestTrades `json:"result"`
}
