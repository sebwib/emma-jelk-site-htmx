package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sebwib/emma-site-htmx/components/layout"
	"github.com/sebwib/emma-site-htmx/components/pages"
	"github.com/sebwib/emma-site-htmx/components/reusable"
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
		r.Get("/edit/storedtext/modal/{id}", h.storedTextModal)
		r.Put("/edit/storedtext/{id}", h.updateStoredText)
		r.Patch("/edit/art/{id}", h.patchArt)
		r.Patch("/edit/art/{id}/{field}", h.patchArtField)
		r.Post("/edit/upload", h.uploadImage)
		r.Post("/edit/art", h.createArt)
		r.Delete("/edit/art/{id}", h.deleteArt)

		r.Get("/edit/print/modal/new", h.addPrintModal)
		r.Get("/edit/print/modal/{id}", h.editPrintModal)
		r.Post("/edit/print", h.createPrint)
		r.Delete("/edit/print/{id}", h.deletePrint)
		r.Patch("/edit/print/{id}", h.patchPrint)
		r.Patch("/edit/print/{id}/{field}", h.patchPrintField)
	})
}

func (h *Handler) patchPrintField(w http.ResponseWriter, r *http.Request) {
	printID := chi.URLParam(r, "id")
	field := chi.URLParam(r, "field")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	value := r.FormValue(field)

	log.Printf("Patching print ID %s field %s to value %s", printID, field, value)
	// Convert value to appropriate type based on field
	var err error
	switch field {
	case "show_in_store":
		boolValue := value == "true" || value == "1" || value == "on"
		log.Printf("Converted value to bool: %v", boolValue)
		err = h.DB.UpdatePrintField(printID, field, boolValue)
	default:
		err = h.DB.UpdatePrintField(printID, field, value)
	}

	if err != nil {
		h.handleError(w, "Failed to update print field", http.StatusInternalServerError, err)
		return
	}

	print, err := h.DB.GetPrintById(printID)
	if err != nil {
		h.handleError(w, "Failed to load updated print", http.StatusInternalServerError, err)
		return
	}

	// Return updated print row
	h.render(w, r, pages.PrintRow(*print, 0), true)
}

func (h *Handler) patchPrint(w http.ResponseWriter, r *http.Request) {
	printID := chi.URLParam(r, "id")

	var patch db.PrintPatch
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
		if priceStr := r.FormValue("price"); priceStr != "" {
			if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
				patch.Price = &price
			}
		}
		if quantityLeftStr := r.FormValue("quantity_left"); quantityLeftStr != "" {
			if quantityLeft, err := strconv.Atoi(quantityLeftStr); err == nil {
				patch.QuantityLeft = &quantityLeft
			}
		}
		if imgURL := r.FormValue("img_url"); imgURL != "" {
			patch.ImgURL = &imgURL
		}
		if thumbURL := r.FormValue("thumb_url"); thumbURL != "" {
			patch.ThumbURL = &thumbURL
		}
		if soldStr := r.FormValue("show_in_store"); soldStr != "" {
			showInStore := soldStr == "true" || soldStr == "1" || soldStr == "on"
			patch.ShowInStore = &showInStore
		}
	}

	// Update art with only the fields that are set
	if err := h.DB.UpdatePrint(printID, patch); err != nil {
		h.handleError(w, "Failed to update print", http.StatusInternalServerError, err)
		return
	}

	if replaceParam := r.URL.Query().Get("replace"); replaceParam == "true" {
		w.Header().Set("HX-Redirect", "/edit")
		return
	}

	print, err := h.DB.GetPrintById(printID)
	if err != nil {
		h.handleError(w, "Failed to load updated print", http.StatusInternalServerError, err)
		return
	}

	// Return updated print row for drag-and-drop updates
	h.render(w, r, pages.PrintRow(*print, 0), true)
}

func (h *Handler) addPrintModal(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, pages.AddPrintModal(), true)
}

func (h *Handler) editPrintModal(w http.ResponseWriter, r *http.Request) {
	printID := chi.URLParam(r, "id")

	print, err := h.DB.GetPrintById(printID)
	if err != nil {
		h.handleError(w, "Failed to load print", http.StatusInternalServerError, err)
		return
	}

	h.render(w, r, pages.EditPrintModal(print), true)
}

func (h *Handler) createPrint(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Parse form values
	width, _ := strconv.Atoi(r.FormValue("width"))
	height, _ := strconv.Atoi(r.FormValue("height"))
	showInStore, _ := strconv.ParseBool(r.FormValue("show_in_store"))

	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	quantityLeft, _ := strconv.Atoi(r.FormValue("quantity_left"))

	print := db.Print{
		ImgURL:       r.FormValue("img_url"),
		ThumbURL:     r.FormValue("thumb_url"),
		Title:        r.FormValue("title"),
		Medium:       r.FormValue("medium"),
		Width:        width,
		Height:       height,
		Year:         r.FormValue("year"),
		Description:  r.FormValue("description"),
		Price:        price,
		QuantityLeft: quantityLeft,
		ShowInStore:  showInStore,
	}

	if err := h.DB.AddPrint(print); err != nil {
		h.handleError(w, "Failed to create print", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("HX-Redirect", "/edit")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) deletePrint(w http.ResponseWriter, r *http.Request) {
	printID := chi.URLParam(r, "id")

	if err := h.DB.DeletePrint(printID); err != nil {
		h.handleError(w, "Failed to delete print", http.StatusInternalServerError, err)
		return
	}

	h.render(w, r, reusable.Empty(), true)
}

func (h *Handler) updateStoredText(w http.ResponseWriter, r *http.Request) {
	referenceID := chi.URLParam(r, "id")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")

	storedText := db.StoredText{
		ReferenceID: referenceID,
		Content:     content,
	}

	if err := h.DB.AddStoredText(storedText); err != nil {
		h.handleError(w, "Failed to update stored text", http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) storedTextModal(w http.ResponseWriter, r *http.Request) {
	referenceID := chi.URLParam(r, "id")

	storedText, err := h.DB.GetStoredTextByReferenceID(referenceID)
	if err != nil {
		h.handleError(w, "Failed to load stored text", http.StatusInternalServerError, err)
		return
	}

	h.render(w, r, pages.EditStoredTextModal(referenceID, storedText[0].Content), true)
}

func (h *Handler) deleteArt(w http.ResponseWriter, r *http.Request) {
	artID := chi.URLParam(r, "id")

	if err := h.DB.DeleteArt(artID); err != nil {
		h.handleError(w, "Failed to delete art", http.StatusInternalServerError, err)
		return
	}

	h.render(w, r, reusable.Empty(), true)
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

	prints, err := h.DB.GetAllPrints()
	if err != nil {
		h.handleError(w, "Failed to load prints", http.StatusInternalServerError, err)
		return
	}

	references, err := h.DB.GetReferences()
	if err != nil {
		h.handleError(w, "Failed to load references", http.StatusInternalServerError, err)
		return
	}

	// oob update background
	if h.isHTMX(r) {
		h.render(w, r, layout.Background(r.URL.Path, true), true)
	}
	h.render(w, r, pages.Edit(arts, prints, references), false)
}

func (h *Handler) edit(w http.ResponseWriter, r *http.Request) {
	arts, err := h.DB.GetArts()
	if err != nil {
		h.handleError(w, "Failed to load gallery", http.StatusInternalServerError, err)
		return
	}

	prints, err := h.DB.GetAllPrints()
	if err != nil {
		h.handleError(w, "Failed to load prints", http.StatusInternalServerError, err)
		return
	}

	references, err := h.DB.GetReferences()
	if err != nil {
		h.handleError(w, "Failed to load references", http.StatusInternalServerError, err)
		return
	}

	// oob update background
	//h.render(w, r, layout.Background(r.URL.Path, true), true)
	h.render(w, r, pages.Edit(arts, prints, references), false)
}

func (h *Handler) patchArtField(w http.ResponseWriter, r *http.Request) {
	artID := chi.URLParam(r, "id")
	field := chi.URLParam(r, "field")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	value := r.FormValue(field)

	// Convert value to appropriate type based on field
	var err error
	switch field {
	case "sold", "show_in_gallery":
		boolValue := value == "true" || value == "1" || value == "on"
		err = h.DB.UpdateArtField(artID, field, boolValue)
	default:
		err = h.DB.UpdateArtField(artID, field, value)
	}

	if err != nil {
		h.handleError(w, "Failed to update art field", http.StatusInternalServerError, err)
		return
	}

	art, err := h.DB.GetArtById(artID)
	if err != nil {
		h.handleError(w, "Failed to load updated art", http.StatusInternalServerError, err)
		return
	}

	// Return updated art row
	h.render(w, r, pages.ArtRow(*art, 0), true)
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
			sold := soldStr == "true" || soldStr == "1" || soldStr == "on"
			patch.Sold = &sold
		}
		if showInGalleryStr := r.FormValue("show_in_gallery"); showInGalleryStr != "" {
			showInGallery := showInGalleryStr == "true" || showInGalleryStr == "1" || showInGalleryStr == "on"
			patch.ShowInGallery = &showInGallery
		}
	}

	// Update art with only the fields that are set
	if err := h.DB.UpdateArt(artID, patch); err != nil {
		h.handleError(w, "Failed to update art", http.StatusInternalServerError, err)
		return
	}

	if replaceParam := r.URL.Query().Get("replace"); replaceParam == "true" {
		w.Header().Set("HX-Redirect", "/edit")
		return
	}

	art, err := h.DB.GetArtById(artID)
	if err != nil {
		h.handleError(w, "Failed to load updated art", http.StatusInternalServerError, err)
		return
	}

	// Return updated art row for drag-and-drop updates
	h.render(w, r, pages.ArtRow(*art, 0), true)
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
	url, thumbURL, err := h.ImageUploader.UploadImage(file, header)
	if err != nil {
		log.Printf("Failed to upload image: %v", err)
		http.Error(w, "Failed to upload image", http.StatusInternalServerError)
		return
	}

	// Return the URL as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": url, "thumb_url": thumbURL})
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
	}

	if err := h.DB.AddArt(art); err != nil {
		h.handleError(w, "Failed to create art", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("HX-Redirect", "/edit")
	w.WriteHeader(http.StatusOK)
}
