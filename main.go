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
)

var routes = []partial.Route{
	{Path: "/", Name: "Hem", InSidebar: false},
	{Path: "/gallery", Name: "Galleri", InSidebar: true},
	{Path: "/buyart", Name: "Köp konst", InSidebar: true},
	{Path: "/about", Name: "Om mig", InSidebar: true},
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, proceeding without it")
	}

	// Initialize handlers with services
	r := chi.NewRouter()

	db, err := db.New("./db.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize session store
	sessionStore := authmw.NewSessionStore()

	h := handlers.NewHandler(db, routes)
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
	h.RegisterAuthRoutes(r, sessionStore)
	h.RegisterEditRoutes(r, sessionStore)
}

func registerMiddlewares(r chi.Router) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
}
