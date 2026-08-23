package rarible

import "time"

// Ownership is the response of GET /v0.1/ownerships/{ownershipId}.
type Ownership struct {
	ID         string `json:"id"`
	Blockchain string `json:"blockchain"`
	ItemID     string `json:"itemId"`
	Contract   string `json:"contract"`
	Collection string `json:"collection"`
	Owner      string `json:"owner"`

	// uint256 values, serialized by Rarible as strings.
	TokenID   string `json:"tokenId"`
	Value     string `json:"value"`
	LazyValue string `json:"lazyValue"`

	CreatedAt     time.Time `json:"createdAt"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`

	// nil when the token is not listed for sale.
	BestSellOrder *SellOrder `json:"bestSellOrder"`

	Version int `json:"version"`
}

// SellOrder is the best active sell offer; only meaningful fields are mapped.
type SellOrder struct {
	ID           string `json:"id"`
	Platform     string `json:"platform"`
	Status       string `json:"status"`
	Maker        string `json:"maker"`
	MakePrice    string `json:"makePrice"`
	MakePriceUsd string `json:"makePriceUsd"`
}
