package models

type Order struct {
	ID         string      `json:"id,omitempty"`
	UserID     string      `json:"user_id"`
	Items      []OrderItem `json:"items"`
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status"`     // e.g., "pending", "confirmed", "delivered"
	OrderedAt  int64       `json:"ordered_at"` // Unix timestamp
	UpdatedAt  string      `json:"updated_at"`
}

type OrderItem struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}
