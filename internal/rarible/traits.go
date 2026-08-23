package rarible

import (
	"context"
	"errors"
	"net/http"
)

// QueryTraitsWithRarity returns the rarity of the requested traits within a collection.
func (c *Client) QueryTraitsWithRarity(ctx context.Context, req TraitsRarityRequest) (*TraitsRarityResponse, error) {
	if req.CollectionID == "" {
		return nil, errors.New("rarible: collectionId is required")
	}
	if len(req.Properties) == 0 {
		return nil, errors.New("rarible: at least one property is required")
	}

	var out TraitsRarityResponse
	if err := c.do(ctx, http.MethodPost, "/v0.1/items/traits/rarity", req, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
