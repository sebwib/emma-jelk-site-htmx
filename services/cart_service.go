package services

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
)

type CartItem struct {
	PrintID  string `json:"print_id"`
	Typ      string `json:"typ"`
	Quantity int    `json:"quantity"`
}

type CartService struct {
}

func (c *CartService) GetCart(r *http.Request) ([]CartItem, error) {
	cookie, err := r.Cookie("cart")

	if err != nil {
		log.Println("No cart cookie found:", err)
		return []CartItem{}, nil // No cart yet
	}

	// Decode base64
	decodedValue, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		log.Println("Failed to decode cart cookie:", err)
		return []CartItem{}, nil // Corrupt cart, start fresh
	}

	var cart []CartItem
	err = json.Unmarshal(decodedValue, &cart)
	if err != nil {
		log.Println("Failed to unmarshal cart cookie:", err)
		return []CartItem{}, nil // Corrupt cart, start fresh
	}
	log.Println("Cart items len", len(cart))
	return cart, nil
}

func (c *CartService) SaveCart(w http.ResponseWriter, cart []CartItem) error {
	data, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	// Encode to base64 to avoid issues with special characters in cookies
	encodedValue := base64.StdEncoding.EncodeToString(data)

	cookie := &http.Cookie{
		Name:     "cart",
		Value:    encodedValue,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	return nil
}
