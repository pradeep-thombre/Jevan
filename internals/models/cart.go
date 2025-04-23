package models

// CartItem represents an item in the cart
type CartItem struct {
	ItemID   string  `json:"item_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Name     string  `json:"name"`
}

// Cart represents the structure of a user's cart
type Cart struct {
	ID     string     `json:"id"`
	UserID string     `json:"user_id"`
	Items  []CartItem `json:"items"`
}
