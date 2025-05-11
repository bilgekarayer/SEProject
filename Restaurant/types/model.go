package types

type Restaurant struct {
	ID          int
	Name        string
	Description string
	Location    string
	Cuisine     string
	AvgPrice    int     `json:"avg_price"`
	Rating      float64 `json:"rating"`
}
