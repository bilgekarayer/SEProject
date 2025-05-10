package types

type CartItem struct {
	UserID    int `bson:"_id" json:"user_id"`
	ProductID int `bson:"product_id" json:"product_id"`
	Quantity  int `bson:"quantity" json:"quantity"`
}

type PlaceOrderRequest struct {
	UserID       int    `json:"user_id"`
	RestaurantID int    `json:"restaurant_id"`
	Address      string `json:"address"`
}

type Order struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	RestaurantID int    `json:"restaurant_id"`
	Address      string `json:"address"`
	Status       string `json:"status"`
}
