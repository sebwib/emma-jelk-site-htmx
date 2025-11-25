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
	buyArtTitle, err := h.DB.GetStoredTextByReferenceID("buy_art_title")
	if err != nil {
		h.handleError(w, "Failed to load buy art title", http.StatusInternalServerError, err)
		return
	}
	buyArtText, err := h.DB.GetStoredTextByReferenceID("buy_art_text")
	if err != nil {
		h.handleError(w, "Failed to load buy art text", http.StatusInternalServerError, err)
		return
	}

	// oob update background
	h.render(w, r, layout.Background(r.URL.Path, true), true)
	h.render(w, r, pages.BuyArt(buyArtTitle[0].Content, buyArtText[0].Content), false)
}
