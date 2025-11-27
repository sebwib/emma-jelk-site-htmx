package handlers

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) RegisterAPIRoutes(r chi.Router) {
	r.Get("/api/db", h.downloadDB)
	r.Get("/api/db-shm", h.downloadDHShm)
	r.Get("/api/db-wal", h.downloadDHWal)
}

func (h *Handler) downloadDB(w http.ResponseWriter, r *http.Request) {
	h.downloadFile(w, r, "database.db")
}

func (h *Handler) downloadDHShm(w http.ResponseWriter, r *http.Request) {
	h.downloadFile(w, r, "database.db-shm")
}

func (h *Handler) downloadDHWal(w http.ResponseWriter, r *http.Request) {
	h.downloadFile(w, r, "database.db-wal")
}

func (h *Handler) downloadFile(w http.ResponseWriter, r *http.Request, fileName string) {
	apiToken := os.Getenv("API_TOKEN")

	if r.Header.Get("Authorization") != "Bearer "+apiToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	filePath := "./" + fileName

	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")

	http.ServeFile(w, r, filePath)
}
