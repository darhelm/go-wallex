package types

// Balance represents a user’s wallet entry for a specific currency,
// including available and blocked balances.
type Balance struct {
	Asset  string `json:"asset"`
	FaName string `json:"faName"`
	Fiat   bool   `json:"fiat"`
	Value  string `json:"value"`
	Locked string `json:"locked"`
}

// Balances defines a map struct of asset: balance
type Balances struct {
	Balances map[string]Balance `json:"balances"`
}

// Wallets represents a collection of wallet entries,
// keyed by currency symbol and grouped under a Balances field under Result.
type Wallets struct {
	BaseResponse
	Results Balances `json:"results"`
}
