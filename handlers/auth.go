package handlers

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/sebwib/emma-site-htmx/components/pages"
	"github.com/sebwib/emma-site-htmx/middleware"
)

func (h *Handler) RegisterAuthRoutes(r chi.Router, store *middleware.SessionStore) {
	r.Get("/login", h.showLogin)
	r.Post("/login", h.login(store))
	r.Post("/logout", h.logout(store))
}

func (h *Handler) showLogin(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, pages.Login(), false)
}

func (h *Handler) login(store *middleware.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		// Get credentials from environment variables
		validUsername := os.Getenv("ADMIN_USERNAME")
		validPassword := os.Getenv("ADMIN_PASSWORD")

		// Set defaults if not in environment
		if validUsername == "" {
			validUsername = "admin"
		}
		if validPassword == "" {
			validPassword = "password"
		}

		if username == validUsername && password == validPassword {
			// Create session
			token, err := store.CreateSession(username)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			// Set cookie
			http.SetCookie(w, &http.Cookie{
				Name:     middleware.SessionCookieName,
				Value:    token,
				Path:     "/",
				HttpOnly: true,
				Secure:   false, // Set to true in production with HTTPS
				SameSite: http.SameSiteLaxMode,
				MaxAge:   86400, // 24 hours
			})

			// Redirect to edit page
			w.Header().Set("HX-Redirect", "/edit")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Invalid credentials - re-render login form with error
		h.render(w, r, pages.LoginError(), true)
	}
}

func (h *Handler) logout(store *middleware.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(middleware.SessionCookieName)
		if err == nil {
			store.DeleteSession(cookie.Value)
		}

		// Clear cookie
		http.SetCookie(w, &http.Cookie{
			Name:     middleware.SessionCookieName,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
