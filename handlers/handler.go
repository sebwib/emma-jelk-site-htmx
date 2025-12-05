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
	CartService   *services.CartService
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

func NewHandler(database *db.DB, routes []partial.Route, imageUploader *services.ImageUploader, cartService *services.CartService) *Handler {
	return &Handler{
		DB:            database,
		Routes:        routes,
		ImageUploader: imageUploader,
		CartService:   cartService,
	}
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

	cartItems, err := h.CartService.GetCart(r)
	if err != nil {
		log.Println("Failed to get cart items:", err)
		cartItems = []services.CartItem{}
	}

	if h.isHTMX(r) {
		component = layout.Content(h.getRoutesWithReferences(h.Routes), cartItems, h.DB, r.URL.Path, content)
	} else {
		// Regular request - return full page with sidebar showing current path
		component = layout.Base(h.getRoutesWithReferences(h.Routes), cartItems, h.DB, r.URL.Path, content)
	}

	if r.URL.Path == "/" {
		w.Header().Set("HX-Trigger", "pageThemeDark")
	} else {
		w.Header().Set("HX-Trigger", "pageThemeLight")
	}

	templ.Handler(component).ServeHTTP(w, r)
}

func (h *Handler) handleError(w http.ResponseWriter, message string, statusCode int, err error) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
	http.Error(w, message, statusCode)
}

func (h *Handler) RegisterModalRoutes(r chi.Router) {
	r.Get("/modal/close", h.closeModal)
}

func (h *Handler) RegisterHomeRoutes(r chi.Router) {
	r.Get("/", h.home)
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, pages.Index(h.getRoutesWithReferences(h.Routes)), false)

	// oob update background
	if h.isHTMX(r) {
		h.render(w, r, layout.Background(r.URL.Path, true), true)
	}
}

func (h *Handler) closeModal(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, reusable.Empty(), true)
}
