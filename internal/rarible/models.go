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

// TraitsRarityRequest is the body of POST /v0.1/items/traits/rarity.
type TraitsRarityRequest struct {
	CollectionID string          `json:"collectionId"`
	Properties   []TraitProperty `json:"properties"`
}

// TraitProperty is a single trait key/value pair.
type TraitProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// TraitsRarityResponse holds the requested traits enriched with rarity.
type TraitsRarityResponse struct {
	Traits []TraitRarity `json:"traits"`
}

// TraitRarity is a trait with its rarity percentage, serialized as a string.
type TraitRarity struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Rarity string `json:"rarity"`
}
