package types

type Restaurant struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Location    string  `json:"location"`
	Cuisine     string  `json:"cuisine"`
	AvgPrice    int     `json:"avg_price"`
	Rating      float64 `json:"rating"`
	ImageURL    string  `json:"image_url"`
}
