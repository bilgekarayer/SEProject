package types

type Menu struct {
	ID           int     `json:"id"`
	RestaurantID int     `json:"restaurant_id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
}
