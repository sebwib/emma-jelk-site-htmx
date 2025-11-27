package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/sebwib/emma-site-htmx/components/partial"
	"github.com/sebwib/emma-site-htmx/db"
	"github.com/sebwib/emma-site-htmx/handlers"
	authmw "github.com/sebwib/emma-site-htmx/middleware"
	"github.com/sebwib/emma-site-htmx/services"
)

var routes = []partial.Route{
	{Path: "/", Name: "ref:home_title", InSidebar: false},
	{Path: "/gallery", Name: "ref:gallery_title", InSidebar: true},
	{Path: "/buyart", Name: "ref:buy_art_title", InSidebar: true},
	{Path: "/about", Name: "ref:about_me_title", InSidebar: true},
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, proceeding without it")
	}

	// Initialize handlers with services
	r := chi.NewRouter()

	log.Println("v0.0.9")

	db, err := db.New("./database.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize image uploader
	imageUploader, err := services.NewImageUploader()
	if err != nil {
		log.Fatalf("Failed to initialize image uploader: %v", err)
	}

	// Initialize session store
	sessionStore := authmw.NewSessionStore()

	h := handlers.NewHandler(db, routes, imageUploader)
	registerMiddlewares(r)
	registerRoutes(h, r, sessionStore)

	port := ":8080"
	log.Printf("Server starting on %s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}

func registerRoutes(h *handlers.Handler, r chi.Router, sessionStore *authmw.SessionStore) {
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	h.RegisterHomeRoutes(r)
	h.RegisterModalRoutes(r)
	h.RegisterGalleryRoutes(r)
	h.RegisterBuyArtRoutes(r)
	h.RegisterAboutRoutes(r)
	h.RegisterAPIRoutes(r)
	h.RegisterAuthRoutes(r, sessionStore)
	h.RegisterEditRoutes(r, sessionStore)
}

func registerMiddlewares(r chi.Router) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
}
