package types

// Menu/types/model.go
type Menu struct {
	ID           int     `json:"id"`
	RestaurantID int     `json:"restaurant_id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	ImageURL     string  `json:"image_url"` // NEW
}
