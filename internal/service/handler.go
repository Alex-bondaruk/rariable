package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Alex-bondaruk/rariable/internal/rarible"

	"go.uber.org/zap"
)

const maxRequestBody = 1 << 20 // 1MB

// raribleClient is the client subset the handlers depend on.
type raribleClient interface {
	GetOwnershipByID(ctx context.Context, id string) (*rarible.Ownership, error)
	QueryTraitsWithRarity(ctx context.Context, req rarible.TraitsRarityRequest) (*rarible.TraitsRarityResponse, error)
}

// Handler exposes the Rarible client over HTTP.
type Handler struct {
	client raribleClient
	log    *zap.Logger
}

// NewHandler returns a Handler backed by client.
func NewHandler(client raribleClient, log *zap.Logger) *Handler {
	return &Handler{client: client, log: log}
}

// GetOwnership handles GET /ownerships/{id}.
func (h *Handler) GetOwnership(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	o, err := h.client.GetOwnershipByID(r.Context(), id)
	if err != nil {
		h.writeUpstreamError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, o)
}

// TraitsRarity handles POST /traits/rarity.
func (h *Handler) TraitsRarity(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBody)

	var req rarible.TraitsRarityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.CollectionID == "" {
		h.writeError(w, http.StatusBadRequest, "collectionId is required")
		return
	}
	if len(req.Properties) == 0 {
		h.writeError(w, http.StatusBadRequest, "at least one property is required")
		return
	}

	res, err := h.client.QueryTraitsWithRarity(r.Context(), req)
	if err != nil {
		h.writeUpstreamError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, res)
}

// Healthz handles GET /healthz.
func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.log.Error("encode response", zap.Error(err))
	}
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}

// writeUpstreamError maps a Rarible error onto our HTTP response.
func (h *Handler) writeUpstreamError(w http.ResponseWriter, err error) {
	var apiErr *rarible.APIError
	if !errors.As(err, &apiErr) {
		h.log.Error("rarible request failed", zap.Error(err))
		h.writeError(w, http.StatusBadGateway, "upstream request failed")
		return
	}

	switch apiErr.StatusCode {
	case http.StatusBadRequest:
		h.log.Debug("rarible rejected the request", zap.Error(apiErr))
		h.writeError(w, http.StatusBadRequest, apiErr.Message)
	case http.StatusNotFound:
		h.log.Debug("rarible rejected the request", zap.Error(apiErr))
		h.writeError(w, http.StatusNotFound, apiErr.Message)
	case http.StatusTooManyRequests:
		h.log.Warn("rarible rate limited us", zap.Error(apiErr))
		h.writeError(w, http.StatusServiceUnavailable, "rate limited, try again later")
	default:
		// 401/403 (our key), 5xx, or anything unmapped: not the caller's fault.
		h.log.Error("rarible request failed", zap.Error(apiErr))
		h.writeError(w, http.StatusBadGateway, "upstream request failed")
	}
}
