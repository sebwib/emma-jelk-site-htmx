package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sebwib/emma-site-htmx/components/layout"
	"github.com/sebwib/emma-site-htmx/components/pages"
)

func (h *Handler) RegisterGalleryRoutes(r chi.Router) {
	r.Get("/gallery", h.gallery)
	r.Get("/gallery/{id}", h.gallerySingle)
}

func (h *Handler) gallerySingle(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	art, err := h.DB.GetArtById(idParam)
	if err != nil {
		h.handleError(w, "Art not found", http.StatusNotFound, err)
		return
	}

	h.render(w, r, pages.GallerySingle(*art), true)
}

func (h *Handler) gallery(w http.ResponseWriter, r *http.Request) {
	pageSize := 8
	page := 1

	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		if parsedPage, err := strconv.Atoi(pageParam); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	arts, err := h.DB.GetArtPaged(pageSize, page*pageSize-pageSize)
	if err != nil {
		h.handleError(w, "Failed to load gallery", http.StatusInternalServerError, err)
		return
	}

	// oob update background
	h.render(w, r, layout.Background(r.URL.Path, true), true)

	for _, art := range arts {
		log.Println(art.Title)
	}
	h.render(w, r, pages.Gallery(arts, page), page > 1)
}

/*
func (h *Handler) once(w http.ResponseWriter, r *http.Request) {
	_db := h.DB

	_db.AddArt(db.Art{ImgURL: "totum.jpg", ThumbURL: "totum_thumb.jpg", Title: "Totum", Medium: "", Width: 92, Height: 119, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "mare.jpg", ThumbURL: "mare_thumb.jpg", Title: "Mare", Medium: "", Width: 66, Height: 88, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "rigor_mortis.jpg", ThumbURL: "rigor_mortis_thumb.jpg", Title: "Rigor mortis", Medium: "", Width: 80, Height: 80, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "sulfur_region.jpg", ThumbURL: "sulfur_region_thumb.jpg", Title: "Sulfur region", Medium: "", Width: 70, Height: 50, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "breaking_point.jpg", ThumbURL: "breaking_point_thumb.jpg", Title: "Breaking point", Medium: "", Width: 70, Height: 70, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "yellow_rubber_boots.jpg", ThumbURL: "yellow_rubber_boots_thumb.jpg", Title: "Yellow rubber boots", Medium: "", Width: 61, Height: 50, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "everyone_sucks.jpg", ThumbURL: "everyone_sucks_thumb.jpg", Title: "Everyone sucks", Medium: "", Width: 70, Height: 100, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "new_dawn.jpg", ThumbURL: "new_dawn_thumb.jpg", Title: "New dawn", Medium: "", Width: 100, Height: 70, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "sicadas_sang.jpg", ThumbURL: "sicadas_sang_thumb.jpg", Title: "And the cicadas sang", Medium: "", Width: 50, Height: 70, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "phylgia.jpg", ThumbURL: "phylgia_thumb.jpg", Title: "Fylgia", Medium: "", Width: 50, Height: 70, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "in_another_light.jpg", ThumbURL: "in_another_light.jpg", Title: "In another light", Medium: "", Width: 100, Height: 70, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "flays.jpg", ThumbURL: "flays_thumb.jpg", Title: "And she flays", Medium: "", Width: 93, Height: 68, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "cold_fire.jpg", ThumbURL: "cold_fire_thumb.jpg", Title: "Cold fire", Medium: "", Width: 99, Height: 69, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "fallout.jpg", ThumbURL: "fallout_thumb.jpg", Title: "Fallout", Medium: "", Width: 74, Height: 57, Year: "", Description: "", Sold: true, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "home.jpg", ThumbURL: "home_thumb.jpg", Title: "H.O.M.E", Medium: "", Width: 50, Height: 61, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "giants.jpg", ThumbURL: "giants_thumb.jpg", Title: "Giants", Medium: "", Width: 116, Height: 89, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "sojourner.jpg", ThumbURL: "sojourner_thumb.jpg", Title: "The sojourner", Medium: "", Width: 65, Height: 54, Year: "", Description: "", Sold: true, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "eyeonroad.jpg", ThumbURL: "eyeonroad_thumb.jpg", Title: "I'll keep an eye on the road", Medium: "", Width: 61, Height: 50, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "mindthegap.jpg", ThumbURL: "mindthegap_thumb.jpg", Title: "Please mind the gap", Medium: "", Width: 120, Height: 85, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "nexus.jpg", ThumbURL: "nexus_thumb.jpg", Title: "Nexus", Medium: "", Width: 70, Height: 100, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "pontus.jpg", ThumbURL: "pontus_thumb.jpg", Title: "Pontus", Medium: "", Width: 61, Height: 50, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "something_hold_on.jpg", ThumbURL: "something_hold_on_thumb.jpg", Title: "Something to hold on to", Medium: "", Width: 100, Height: 80, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "terra_mater.jpg", ThumbURL: "terra_mater_thumb.jpg", Title: "Terra Mater", Medium: "", Width: 93, Height: 68, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "truth_neverleft.jpg", ThumbURL: "truth_neverleft_thumb.jpg", Title: "The truth is, I never really left", Medium: "", Width: 62, Height: 43, Year: "", Description: "", Sold: true, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "nautilus.jpg", ThumbURL: "nautilus_thumb.jpg", Title: "Nautilus", Medium: "", Width: 89, Height: 116, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "against_the_dying.jpg", ThumbURL: "against_the_dying_thumb.jpg", Title: "Against the dying of the light", Medium: "", Width: 104, Height: 121, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "ogden.jpg", ThumbURL: "ogden_thumb.jpg", Title: "Whenever we entered the field, Mr Ogden would come out screaming at us", Medium: "", Width: 70, Height: 49, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "a_sky_of_age.jpg", ThumbURL: "a_sky_of_age_thumb.jpg", Title: "A sky of age", Medium: "", Width: 120, Height: 85, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "mejeriet.jpg", ThumbURL: "mejeriet_thumb.jpg", Title: "Mejeriet", Medium: "", Width: 50, Height: 70, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "divided.jpg", ThumbURL: "divided_thumb.jpg", Title: "Divided", Medium: "", Width: 80, Height: 80, Year: "", Description: "", Sold: true, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "landing.jpg", ThumbURL: "landing_thumb.jpg", Title: "The landing area is a total mess", Medium: "", Width: 117, Height: 69, Year: "", Description: "", Sold: true, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "interlinked.jpg", ThumbURL: "interlinked_thumb.jpg", Title: "Interlinked", Medium: "", Width: 60, Height: 50, Year: "", Description: "", Sold: true, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "onlynow.jpg", ThumbURL: "onlynow_thumb.jpg", Title: "There is only now", Medium: "", Width: 70, Height: 50, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "hope.jpg", ThumbURL: "hope_thumb.jpg", Title: "Hope dies last", Medium: "", Width: 90, Height: 106, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "villebrad.jpg", ThumbURL: "villebrad_thumb.jpg", Title: "Villebråd 2021", Medium: "", Width: 0, Height: 0, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "pandemicspring.jpg", ThumbURL: "pandemicspring_thumb.jpg", Title: "Pandemic spring", Medium: "", Width: 100, Height: 100, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "mother_n_son.jpg", ThumbURL: "mother_n_son_thumb.jpg", Title: "The mother and the son", Medium: "", Width: 90, Height: 110, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "hall_kaften_tommy.jpg", ThumbURL: "hall_kaften_tommy_thumb.jpg", Title: "Håll käften, Tommy", Medium: "", Width: 100, Height: 70, Year: "", Description: "", Sold: true, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "floating.jpg", ThumbURL: "floating_thumb.jpg", Title: "Floating on the surface of a planet", Medium: "", Width: 98, Height: 70, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "fluids.jpg", ThumbURL: "fluids_thumb.jpg", Title: "Amniotic Fluids Spa & Resort", Medium: "", Width: 95, Height: 95, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "pourguy.jpg", ThumbURL: "pourguy_thumb.jpg", Title: "Pour guy", Medium: "", Width: 90, Height: 100, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "born.jpg", ThumbURL: "born_thumb.jpg", Title: "Rebirth", Medium: "", Width: 70, Height: 100, Year: "", Description: "", Sold: true, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "100000.jpg", ThumbURL: "100000.jpg", Title: "100 000", Medium: "", Width: 70, Height: 70, Year: "", Description: "", Sold: true, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "mother.jpg", ThumbURL: "mother.jpg", Title: "Mother", Medium: "", Width: 70, Height: 70, Year: "", Description: "", Sold: true, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "grimas2.jpg", ThumbURL: "grimas2_thumb.jpg", Title: "Grimas 2", Medium: "", Width: 70, Height: 70, Year: "", Description: "", Sold: true, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "hungrig.jpg", ThumbURL: "hungrig_thumb.jpg", Title: "Hungrig", Medium: "", Width: 0, Height: 0, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "yard2.jpg", ThumbURL: "yard2_thumb.jpg", Title: "", Medium: "", Width: 0, Height: 0, Year: "", Description: "", Sold: true, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "teeth.jpg", ThumbURL: "teeth_thumb.jpg", Title: "Donut hole", Medium: "", Width: 70, Height: 70, Year: "", Description: "", Sold: true, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "things.jpg", ThumbURL: "things_thumb.jpg", Title: "The things we leave behind", Medium: "", Width: 120, Height: 95, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
	_db.AddArt(db.Art{ImgURL: "trump.jpg", ThumbURL: "trump_thumb.jpg", Title: "The anatomy of a plastic white house", Medium: "", Width: 1, Height: 1, Year: "", Description: "", Sold: false, CreatedAt: time.Now().Format(time.RFC3339)})
}
*/
