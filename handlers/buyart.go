package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sebwib/emma-site-htmx/components/layout"
	"github.com/sebwib/emma-site-htmx/components/pages"
)

func (h *Handler) RegisterBuyArtRoutes(r chi.Router) {
	r.Get("/buyart", h.buyArt)
}

func (h *Handler) buyArt(w http.ResponseWriter, r *http.Request) {

	// oob update background
	h.render(w, r, layout.Background(r.URL.Path, true), true)
	h.render(w, r, pages.BuyArt(), false)
}
