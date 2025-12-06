package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sebwib/emma-site-htmx/components/pages"
	"github.com/sebwib/emma-site-htmx/middleware"
)

func (h *Handler) RegisterOrderRoutes(r chi.Router, store *middleware.SessionStore) {
	r.Group(func(r chi.Router) {
		//r.Use(middleware.RequireAuth(store))
		r.Get("/orders", h.ordersPage)
		r.Post("/orders/{orderID}/update_status", h.updateOrderStatus)
		r.Post("/orders/{orderID}/row/{rowID}/toggle_paid", h.toggleOrderRowPaid)
	})
}

func (h *Handler) ordersPage(w http.ResponseWriter, r *http.Request) {
	orders, err := h.DB.GetAllOrders()
	if err != nil {
		h.handleError(w, "Failed to get orders", 500, err)
		return
	}

	h.render(w, r, pages.Orders(orders), false)
}

func (h *Handler) updateOrderStatus(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderID")

	orderStatus := r.FormValue("order_status")
	err := h.DB.UpdateOrderStatus(orderID, orderStatus)
	if err != nil {
		h.handleError(w, "Failed to update order status", 500, err)
		return
	}

	order, err := h.DB.GetOrderByID(orderID)
	if err != nil {
		h.handleError(w, "Failed to get order", 500, err)
		return
	}

	h.render(w, r, pages.OrderSingle(order), true)
}

func (h *Handler) toggleOrderRowPaid(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderID")
	rowID := chi.URLParam(r, "rowID")
	hasPaidStr := r.FormValue("paid")

	hasPaid := false
	if hasPaidStr == "true" || hasPaidStr == "on" || hasPaidStr == "1" {
		hasPaid = true
	}
	err := h.DB.UpdateOrderPaidStatus(rowID, hasPaid)
	if err != nil {
		h.handleError(w, "Failed to toggle order row paid status", 500, err)
		return
	}

	order, err := h.DB.GetOrderByID(orderID)
	if err != nil {
		h.handleError(w, "Failed to get order", 500, err)
		return
	}

	h.render(w, r, pages.OrderSingle(order), true)
	w.WriteHeader(http.StatusOK)
}
