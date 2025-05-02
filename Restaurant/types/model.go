package types

type Restaurant struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Cuisine  string `json:"cuisine"` // Burası mutlaka string olmalı ve json etiketi yazılmalı!
}
