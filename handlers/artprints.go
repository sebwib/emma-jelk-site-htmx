package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sebwib/emma-site-htmx/components/layout"
	"github.com/sebwib/emma-site-htmx/components/pages"
)

func (h *Handler) RegisterArtPrintRoutes(r chi.Router) {
	r.Get("/prints", h.prints)
}

func (h *Handler) prints(w http.ResponseWriter, r *http.Request) {
	printTitle, err := h.DB.GetStoredTextByReferenceID("prints_title")
	if err != nil {
		h.handleError(w, "Failed to load prints title", http.StatusInternalServerError, err)
		return
	}
	printText, err := h.DB.GetStoredTextByReferenceID("prints_text")
	if err != nil {
		h.handleError(w, "Failed to load about me text", http.StatusInternalServerError, err)
		return
	}

	prints, err := h.DB.GetPrintsForStore()
	if err != nil {
		h.handleError(w, "Failed to load prints", http.StatusInternalServerError, err)
		return
	}

	h.render(w, r, pages.Prints(printTitle[0].Content, printText[0].Content, prints), false)

	// oob update background
	if h.isHTMX(r) {
		h.render(w, r, layout.Background(r.URL.Path, true), true)
	}
}
