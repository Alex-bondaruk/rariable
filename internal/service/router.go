package service

import "net/http"

// NewRouter wires the handlers onto their routes.
func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.Healthz)
	mux.HandleFunc("GET /ownerships/{id}", h.GetOwnership)
	mux.HandleFunc("POST /traits/rarity", h.TraitsRarity)

	return withRecovery(h.log, withLogging(h.log, mux))
}
