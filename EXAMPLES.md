# Wallex SDK Examples

This document provides full practical usage examples for all available SDK methods.

---

# Authentication

```go
client, err := wallex.NewClient(wallex.ClientOptions{
    ApiKey: "YOUR_WALLEX_API_KEY",
})
```

---

# Market Information

## Get Markets Info

```go
markets, err := client.GetMarketsInfo()
fmt.Println(markets.Result["symbols"]["BTCUSDT"])
```

## Get Order Book

```go
ob, err := client.GetOrderBook("BTCUSDT")
fmt.Println("Best Ask:", ob.Result.Ask[0])
fmt.Println("Best Bid:", ob.Result.Bid[0])
```

## Get Recent Trades

```go
trades, err := client.GetRecentTrades("BTCUSDT")
fmt.Println(trades.Result.LatestTrades[0])
```

---

# Wallet Operations

## Get Wallet Balances

```go
balances, err := client.GetWallets()
fmt.Println("USDT Balance:", balances.Wallets["USDT"].Balance)
```

---

# Trading

## Create Order

```go
order, err := client.CreateOrder(types.CreateOrderParams{
    Symbol:   "BTCUSDT",
    Type:     "LIMIT",
    Side:     "BUY",
    Price:    "9500",
    Quantity: "0.002",
})
fmt.Println(order.Result.Status)
```

## Cancel Order

```go
cancel, err := client.CancelOrder("my-client-order-id")
fmt.Println(cancel.Result.Status)
```

## Get Open Orders

```go
openOrders, err := client.GetOpenOrders("BTCUSDT")
fmt.Println(openOrders.Result.Orders)
```

## Get Order Status

```go
status, err := client.GetOrderStatus("my-client-order-id")
fmt.Println(status.Result)
```

# User Trades

```go
trades, err := client.GetUserTrades(types.UserTradesParams{
    Symbol: "BTCUSDT",
    Side:   "BUY",
})
fmt.Println(trades.Result.AccountLatestTrades)
```

---

# Error Handling

```go
_, err := client.GetOrderBook("INVALID")
if err != nil {
    if apiErr, ok := err.(*wallex.APIError); ok {
        fmt.Println("Status Code:", apiErr.StatusCode)
       	fmt.Println("Message:", apiErr.Message)
       	fmt.Println("Fields:", apiErr.Fields)
    }
}
```
