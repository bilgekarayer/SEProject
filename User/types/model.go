package types

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	RoleID    int    `json:"role_id"`
	RoleName  string `json:"role_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CreatedAt string `json:"created_at"`
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
	FirstName string `json:"firstName"` // ⬅︎ değişti
	LastName  string `json:"lastName"`  // ⬅︎ değişti
}

type LoginRequest struct {
	Username string `json:"username" example:"john_doe"`
	Password string `json:"password" example:"123456"`
}
