package handlers

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	"github.com/sebwib/emma-site-htmx/components/layout"
	"github.com/sebwib/emma-site-htmx/components/pages"
	"github.com/sebwib/emma-site-htmx/components/partial"
	"github.com/sebwib/emma-site-htmx/components/reusable"
	"github.com/sebwib/emma-site-htmx/db"
	"github.com/sebwib/emma-site-htmx/services"
)

type Handler struct {
	DB            *db.DB
	Routes        []partial.Route
	ImageUploader *services.ImageUploader
}

func (h *Handler) getRoutesWithReferences(routes []partial.Route) []partial.Route {
	references, err := h.DB.GetReferences()
	if err != nil {
		log.Println("get references", err)
		return routes
	}

	_routes := make([]partial.Route, len(routes))
	copy(_routes, routes)
	for i, route := range _routes {
		for _, ref := range references {
			if route.Name == "ref:"+ref.ReferenceID {
				_routes[i].Name = ref.Content
			}
		}
	}
	return _routes
}

func NewHandler(database *db.DB, routes []partial.Route, imageUploader *services.ImageUploader) *Handler {
	return &Handler{
		DB:            database,
		Routes:        routes,
		ImageUploader: imageUploader,
	}
}

func (h *Handler) RegisterModalRoutes(r chi.Router) {
	r.Get("/modal/close", h.closeModal)
}

func (h *Handler) RegisterHomeRoutes(r chi.Router) {
	r.Get("/", h.home)
}

func (h *Handler) isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func (h *Handler) render(w http.ResponseWriter, r *http.Request, content templ.Component, raw bool) {
	if raw {
		templ.Handler(content).ServeHTTP(w, r)
		return
	}

	var component templ.Component

	if h.isHTMX(r) {
		component = layout.Content(h.getRoutesWithReferences(h.Routes), h.DB, r.URL.Path, content)
	} else {
		// Regular request - return full page with sidebar showing current path
		component = layout.Base(h.getRoutesWithReferences(h.Routes), h.DB, r.URL.Path, content)
	}

	templ.Handler(component).ServeHTTP(w, r)
}

// handleError logs the error and sends an HTTP error response
func (h *Handler) handleError(w http.ResponseWriter, message string, statusCode int, err error) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
	http.Error(w, message, statusCode)
}

func (h *Handler) handleServerError(w http.ResponseWriter, message string, err error) {
	h.handleError(w, message, http.StatusInternalServerError, err)
}

func (h *Handler) handleBadRequest(w http.ResponseWriter, message string) {
	h.handleError(w, message, http.StatusBadRequest, nil)
}

func (h *Handler) handleNotFound(w http.ResponseWriter, message string) {
	h.handleError(w, message, http.StatusNotFound, nil)
}

func (h *Handler) handleMethodNotAllowed(w http.ResponseWriter) {
	h.handleError(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, pages.Index(h.getRoutesWithReferences(h.Routes)), false)
	h.render(w, r, layout.Background(r.URL.Path, true), true)
}

func (h *Handler) closeModal(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, reusable.Empty(), true)
}
