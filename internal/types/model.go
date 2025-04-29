package types

type User struct {
	ID       int    `bson:"_id" json:"id"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

type Restaurant struct {
	ID          int
	Name        string
	Description string
}

type Menu struct {
	ID           int
	RestaurantID int
	Name         string
	Price        float64
}

type Order struct {
	ID           int
	UserID       int
	RestaurantID int
	TotalAmount  float64
}
