# Go Wallex

[![Go Reference](https://pkg.go.dev/badge/github.com/darhelm/go-wallex.svg)](https://pkg.go.dev/github.com/darhelm/go-wallex)
[![Go Report Card](https://goreportcard.com/badge/github.com/darhelm/go-wallex)](https://goreportcard.com/report/github.com/darhelm/go-wallex)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/darhelm/go-wallex)](https://golang.org/dl/)

A comprehensive, type-safe, and fully documented Go SDK for interacting with the **Wallex** cryptocurrency exchange API.  
This SDK provides a clean and intuitive interface for accessing market data, managing balances, and performing spot trading.

## Disclaimer

This SDK is **unofficial** and not affiliated with Wallex.  
Use at your own risk — the author(s) assume no liability for financial losses or incorrect API usage.

## Features

- Full implementation of public and private Wallex endpoints
- Strongly typed request/response models
- Simple, static API-key authentication (`X-API-Key`)
- Real-time market data: order books, latest trades, market metadata
- Wallet (balance) retrieval
- Full order lifecycle: create, cancel, fetch, open orders, trade history
- Structured error handling (`APIError`, `RequestError`)
- Clean and maintainable Go codebase

## Installation

```bash
go get github.com/darhelm/go-wallex
```

## Quick Start

```go
package main

import (
    "fmt"
    wallex "github.com/darhelm/go-wallex"
)

func main() {
    client, err := wallex.NewClient(wallex.ClientOptions{
        ApiKey: "YOUR_WALLEX_API_KEY",
    })
    if err != nil {
        panic(err)
    }

    markets, err := client.GetMarketsInfo()
    if err != nil {
        panic(err)
    }

    fmt.Println("Symbols:", len(markets.Result["symbols"]))
}
```

## Documentation

- Go SDK docs: https://pkg.go.dev/github.com/darhelm/go-wallex
- Wallex API: https://api-docs.wallex.ir/
- Full examples: `EXAMPLES.md`

---

# Examples

## Authentication (API Key)
be mindful that there is currently no way to auto refresh your tokens,
you have to regenerate new tokens based on allowed set timestamp and recreate the client.
```go
client, err := wallex.NewClient(wallex.ClientOptions{
    ApiKey: "YOUR_WALLEX_API_KEY",
})
```

## Get Market Information

```go
markets, err := client.GetMarketsInfo()
fmt.Println(markets.Result["symbols"]["BTCUSDT"])
```

## Get Order Book

```go
orderBook, err := client.GetOrderBook("BTCUSDT")
fmt.Println(orderBook.Result.Ask[0], orderBook.Result.Bid[0])
```

## Get Recent Trades

```go
recentTrades, err := client.GetRecentTrades("BTCUSDT")
fmt.Println(recentTrades.Result.LatestTrades[0])
```

## Get Wallet Balances

```go
balances, err := client.GetWallets()
fmt.Println(balances.Wallets["USDT"].Balance)
```

## Create Order

```go
createOrder, err := client.CreateOrder(types.CreateOrderParams{
    Symbol:   "BTCUSDT",
    Type:     "LIMIT",
    Side:     "BUY",
    Price:    "10000",
    Quantity: "0.001",
})
fmt.Println(createOrder.Result.Status)
```

## Cancel Order

```go
cancelRes, err := client.CancelOrder("my-order-id")
fmt.Println(cancelRes.Result.Status)
```

## Open Orders

```go
openOrders, _ := client.GetOpenOrders("BTCUSDT")
fmt.Println(openOrders.Result.Orders)
```

## Order Status

```go
orderStatus, _ := client.GetOrderStatus("my-order-id")
fmt.Println(orderStatus.Result)
```

## User Trades

```go
userTrades, _ := client.GetUserTrades(types.UserTradesParams{
    Symbol: "BTCUSDT",
})
fmt.Println(userTrades.Result.AccountLatestTrades)
```

## Error Handling

```go
if err != nil {
    if apiErr, ok := err.(*wallex.APIError); ok {
        fmt.Println("HTTP Status:", apiErr.StatusCode)
        fmt.Println("Code:", apiErr.Code)
        fmt.Println("Message:", apiErr.Message)
        fmt.Println("Fields:", apiErr.Fields)
    }
}
```

## Contributing

1. Fork the repository
2. Create a branch, e.g. `feat/new-feature`
3. Commit changes
4. Open a Pull Request

Before pushing:

```bash
go vet ./...
golangci-lint run
```

## License

MIT License.