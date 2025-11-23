package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sebwib/emma-site-htmx/components/layout"
	"github.com/sebwib/emma-site-htmx/components/pages"
)

func (h *Handler) RegisterAboutRoutes(r chi.Router) {
	r.Get("/about", h.about)
}

func (h *Handler) about(w http.ResponseWriter, r *http.Request) {
	// oob update background
	h.render(w, r, layout.Background(r.URL.Path, true), true)
	h.render(w, r, pages.About(), false)
}
