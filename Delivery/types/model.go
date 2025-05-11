package types

type DeliveryOrder struct {
	OrderID        int    `json:"order_id"`
	CustomerName   string `json:"customer_name"`
	RestaurantName string `json:"restaurant_name"`
	Address        string `json:"address"`
	Status         string `json:"status"`
}
