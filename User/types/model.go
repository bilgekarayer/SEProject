package types

type User struct {
	ID        int    `bson:"_id" json:"id"`
	Username  string `bson:"username" json:"username"`
	Password  string `bson:"password" json:"password"`
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
}

type Order struct {
	ID           int
	UserID       int
	RestaurantID int
	TotalAmount  float64
}

type RegisterRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type LoginRequest struct {
	Username string `json:"username" example:"john_doe"`
	Password string `json:"password" example:"123456"`
}
