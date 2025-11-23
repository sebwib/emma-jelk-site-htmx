package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sebwib/emma-site-htmx/components/layout"
	"github.com/sebwib/emma-site-htmx/components/pages"
	"github.com/sebwib/emma-site-htmx/db"
	"github.com/sebwib/emma-site-htmx/middleware"
)

func (h *Handler) RegisterEditRoutes(r chi.Router, store *middleware.SessionStore) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(store))
		r.Get("/edit", h.edit)
		r.Get("/edit/resetorder", h.resetArtOrder)
		r.Get("/edit/art/modal/new", h.addModal)
		r.Get("/edit/art/modal/{id}", h.editModal)
		r.Patch("/edit/art/{id}", h.patchArt)
		r.Post("/edit/upload", h.uploadImage)
		r.Post("/edit/art", h.createArt)
	})
}

func (h *Handler) addModal(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, pages.AddArtModal(), true)
}

func (h *Handler) editModal(w http.ResponseWriter, r *http.Request) {
	artID := chi.URLParam(r, "id")

	art, err := h.DB.GetArtById(artID)
	if err != nil {
		h.handleError(w, "Failed to load art", http.StatusInternalServerError, err)
		return
	}

	h.render(w, r, pages.EditArtModal(art), true)
}

func (h *Handler) resetArtOrder(w http.ResponseWriter, r *http.Request) {
	arts, err := h.DB.GetArts()
	if err != nil {
		h.handleError(w, "Failed to load gallery", http.StatusInternalServerError, err)
		return
	}

	// Reset ordering
	for i, art := range arts {
		art.Ordering = float64(i + 1)
		log.Println("Resetting art ID", art.Id, "to ordering", strconv.FormatFloat(art.Ordering, 'f', 6, 64))
		if err := h.DB.UpdateArt(art.Id, db.ArtPatch{Ordering: &art.Ordering}); err != nil {
			h.handleError(w, "Failed to reset art ordering", http.StatusInternalServerError, err)
			return
		}
	}

	// oob update background
	h.render(w, r, layout.Background(r.URL.Path, true), true)
	h.render(w, r, pages.Edit(arts), false)
}

func (h *Handler) edit(w http.ResponseWriter, r *http.Request) {
	arts, err := h.DB.GetArts()
	if err != nil {
		h.handleError(w, "Failed to load gallery", http.StatusInternalServerError, err)
		return
	}

	// oob update background
	h.render(w, r, layout.Background(r.URL.Path, true), true)
	h.render(w, r, pages.Edit(arts), false)
}

func (h *Handler) patchArt(w http.ResponseWriter, r *http.Request) {
	artID := chi.URLParam(r, "id")

	var patch db.ArtPatch
	contentType := r.Header.Get("Content-Type")

	// Handle both JSON (from drag-and-drop) and form data (from modal)
	if contentType == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
	} else {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Convert form values to ArtPatch (only set fields that are present)
		if title := r.FormValue("title"); title != "" {
			patch.Title = &title
		}
		if medium := r.FormValue("medium"); medium != "" {
			patch.Medium = &medium
		}
		if widthStr := r.FormValue("width"); widthStr != "" {
			if width, err := strconv.Atoi(widthStr); err == nil {
				patch.Width = &width
			}
		}
		if heightStr := r.FormValue("height"); heightStr != "" {
			if height, err := strconv.Atoi(heightStr); err == nil {
				patch.Height = &height
			}
		}
		if year := r.FormValue("year"); year != "" {
			patch.Year = &year
		}
		if description := r.FormValue("description"); description != "" {
			patch.Description = &description
		}
		if imgURL := r.FormValue("img_url"); imgURL != "" {
			patch.ImgURL = &imgURL
		}
		if thumbURL := r.FormValue("thumb_url"); thumbURL != "" {
			patch.ThumbURL = &thumbURL
		}
		if soldStr := r.FormValue("sold"); soldStr != "" {
			sold := soldStr == "true" || soldStr == "1"
			patch.Sold = &sold
		}
	}

	// Update art with only the fields that are set
	if err := h.DB.UpdateArt(artID, patch); err != nil {
		h.handleError(w, "Failed to update art", http.StatusInternalServerError, err)
		return
	}

	if replaceParam := r.URL.Query().Get("replace"); replaceParam == "true" {
		// Use HX-Redirect to reload the page
		w.Header().Set("HX-Redirect", "/edit")
		w.WriteHeader(http.StatusOK)
		return
	}

	// For form submissions (HTMX), close the modal
	if contentType != "application/json" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// For JSON (drag-and-drop), return success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Art updated successfully"})
}

func (h *Handler) uploadImage(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Upload image
	url, err := h.ImageUploader.UploadImage(file, header)
	if err != nil {
		log.Printf("Failed to upload image: %v", err)
		http.Error(w, "Failed to upload image", http.StatusInternalServerError)
		return
	}

	// Return the URL as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": url})
}

func (h *Handler) createArt(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Parse form values
	width, _ := strconv.Atoi(r.FormValue("width"))
	height, _ := strconv.Atoi(r.FormValue("height"))
	sold, _ := strconv.ParseBool(r.FormValue("sold"))

	log.Println("form title", r.FormValue("title"))

	// Create new art entry
	art := db.Art{
		ImgURL:      r.FormValue("img_url"),
		ThumbURL:    r.FormValue("thumb_url"),
		Title:       r.FormValue("title"),
		Medium:      r.FormValue("medium"),
		Width:       width,
		Height:      height,
		Year:        r.FormValue("year"),
		Description: r.FormValue("description"),
		Sold:        sold,
		CreatedAt:   r.FormValue("created_at"),
	}

	log.Println("title", art.Title)
	if err := h.DB.AddArt(art); err != nil {
		h.handleError(w, "Failed to create art", http.StatusInternalServerError, err)
		return
	}

	// Redirect back to edit page
	//w.Header().Set("HX-Redirect", "/edit")
	w.WriteHeader(http.StatusOK)
}
