package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sebwib/emma-site-htmx/components/id"
	"github.com/sebwib/emma-site-htmx/components/layout"
	"github.com/sebwib/emma-site-htmx/components/pages"
	"github.com/sebwib/emma-site-htmx/components/partial"
	"github.com/sebwib/emma-site-htmx/db"
	"github.com/sebwib/emma-site-htmx/services"
)

func (h *Handler) RegisterCartRoutes(r chi.Router) {
	r.Post("/cart/add", h.addToCartHandler)
	r.Get("/cart", h.cartPage)
	r.Post("/cart/remove", h.removeFromCartHandler)
	r.Put("/cart/{id}/quantity", h.quantityChangeHandler)
	r.Post("/cart/checkout", h.checkoutHandler)
	r.Get("/cart/thanks", h.thanksPage)
}

func (h *Handler) updateCartSymbol(w http.ResponseWriter, r *http.Request, cart []services.CartItem) {
	h.render(w, r, partial.CartSymbol(cart, true, id.CartSymbolModeDesktop), true)
	h.render(w, r, partial.CartSymbol(cart, true, id.CartSymbolModeMobile), true)
	h.render(w, r, partial.CartSymbol(cart, true, id.CartSymbolModeCount), true)
}

func (h *Handler) thanksPage(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, pages.ThanksForOrdering(true), false)
}

func (h *Handler) checkoutHandler(w http.ResponseWriter, r *http.Request) {
	cart, err := h.CartService.GetCart(r)
	if err != nil {
		h.handleError(w, "Failed to get cart for checkout", http.StatusInternalServerError, err)
		return
	}

	buyerEmail := r.FormValue("email")

	orderRows := []db.OrderRow{}
	createdAt := time.Now().Format(time.RFC3339)
	orderID := uuid.NewString()

	for _, item := range cart {
		print, err := h.DB.GetPrintById(item.PrintID)
		if err != nil {
			h.handleError(w, "Failed to get print for order item", http.StatusInternalServerError, err)
			return
		}

		orderRow := db.OrderRow{
			UUID:      uuid.NewString(),
			CreatedAt: createdAt,
			OrderID:   orderID,
			Email:     buyerEmail,
			PrintID:   item.PrintID,
			Title:     print.Title,
			Typ:       item.Typ,
			Quantity:  item.Quantity,
			Price:     print.Price,
			Status:    db.OrderStatusPlaced,
			HasPaid:   false,
		}

		err = h.DB.AddOrder(orderRow)
		if err != nil {
			h.handleError(w, "Failed to store order", http.StatusInternalServerError, err)
		}

		orderRows = append(orderRows, orderRow)
	}

	order := db.Order{
		BuyerEmail: buyerEmail,
		Rows:       orderRows,
	}

	err = services.SendOrder(buyerEmail, order)
	success := err == nil
	if err != nil {
		// store order failed, but don't crash the user experience
		fmt.Printf("Failed to send order email: %v\n", err)
	}

	h.CartService.SaveCart(w, []services.CartItem{})
	h.updateCartSymbol(w, r, []services.CartItem{})

	h.render(w, r, pages.ThanksForOrdering(success), true)
}

func (h *Handler) quantityChangeHandler(w http.ResponseWriter, r *http.Request) {
	printID := chi.URLParam(r, "id")
	typ := r.FormValue("type")
	quantityStr := r.FormValue("quantity")

	cart, _ := h.CartService.GetCart(r)
	newCart := []services.CartItem{}
	for _, item := range cart {
		if item.PrintID == printID && item.Typ == typ {
			// Update quantity
			// Convert quantityStr to int
			var quantity int
			_, err := fmt.Sscanf(quantityStr, "%d", &quantity)
			if err != nil || quantity < 1 {
				quantity = 1 // Default to 1 if invalid
			}
			item.Quantity = quantity
		}
		newCart = append(newCart, item)
	}

	h.CartService.SaveCart(w, newCart)

	h.updateCartSymbol(w, r, newCart)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) removeFromCartHandler(w http.ResponseWriter, r *http.Request) {
	printID := r.FormValue("print_id")
	typ := r.FormValue("type")

	cart, _ := h.CartService.GetCart(r)
	newCart := []services.CartItem{}
	for _, item := range cart {
		if item.PrintID == printID && item.Typ == typ {
			continue // Skip this item to remove it
		}
		newCart = append(newCart, item)
	}

	h.CartService.SaveCart(w, newCart)

	if len(newCart) == 0 {
		// If cart is empty, oob render empty cart page
		h.render(w, r, pages.Cart([]pages.CartItemView{}, true), true)
	}

	h.updateCartSymbol(w, r, newCart)

	h.cartPage(w, r)
}

func (h *Handler) cartPage(w http.ResponseWriter, r *http.Request) {
	cartItems, err := h.CartService.GetCart(r)
	if err != nil {
		h.handleError(w, "Failed to load cart items", http.StatusInternalServerError, err)
		return
	}

	cartViews := make([]pages.CartItemView, len(cartItems))
	for i, item := range cartItems {
		print, err := h.DB.GetPrintById(item.PrintID)
		if err != nil {
			h.handleError(w, "Failed to load print for cart item", http.StatusInternalServerError, err)
			return
		}
		cartViews[i] = pages.CartItemView{
			ThumbURL: print.ThumbURL,
			Quantity: item.Quantity,
			Title:    print.Title,
			ID:       item.PrintID,
			Typ:      item.Typ,
		}
	}

	h.render(w, r, pages.Cart(cartViews, false), false)

	// oob update background
	if h.isHTMX(r) {
		h.render(w, r, layout.Background(r.URL.Path, true), true)
	}
}

func (h *Handler) addToCartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	printID := r.FormValue("print_id")

	cart, _ := h.CartService.GetCart(r)
	found := false
	for i, item := range cart {
		if item.PrintID == printID {
			cart[i].Quantity++
			found = true
			break
		}
	}
	if !found {
		cart = append(cart, services.CartItem{PrintID: printID, Quantity: 1, Typ: "print"})
	}

	h.CartService.SaveCart(w, cart)

	h.updateCartSymbol(w, r, cart)

	h.render(w, r, pages.BoughtButton(), true)
	w.WriteHeader(http.StatusOK)
}
