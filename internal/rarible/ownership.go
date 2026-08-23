package rarible

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// GetOwnershipByID returns the ownership identified by BLOCKCHAIN:contract:tokenId:owner.
func (c *Client) GetOwnershipByID(ctx context.Context, id string) (*Ownership, error) {
	if id == "" {
		return nil, errors.New("rarible: ownership id is required")
	}

	var out Ownership
	if err := c.do(ctx, http.MethodGet, "/v0.1/ownerships/"+url.PathEscape(id), nil, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
