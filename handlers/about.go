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
	aboutMeTitle, err := h.DB.GetStoredTextByReferenceID("about_me_title")
	if err != nil {
		h.handleError(w, "Failed to load about me title", http.StatusInternalServerError, err)
		return
	}
	aboutMeText, err := h.DB.GetStoredTextByReferenceID("about_me_text")
	if err != nil {
		h.handleError(w, "Failed to load about me text", http.StatusInternalServerError, err)
		return
	}

	h.render(w, r, pages.About(aboutMeTitle[0].Content, aboutMeText[0].Content), false)

	// oob update background
	if h.isHTMX(r) {
		h.render(w, r, layout.Background(r.URL.Path, true), true)
	}
}
