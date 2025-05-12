package types

type CartItem struct {
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type OrderItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type PlaceOrderRequest struct {
	RestaurantID int         `json:"restaurant_id"`
	Address      string      `json:"address"`
	Items        []OrderItem `json:"items"`
}

type ItemResponse struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type OrderResponse struct {
	ID         int            `json:"id"`
	User       string         `json:"user"`
	Restaurant string         `json:"restaurant"`
	Address    string         `json:"address"`
	Status     string         `json:"status"`
	Items      []ItemResponse `json:"items"`
	Total      string         `json:"total"`
}
