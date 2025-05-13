package User

import (
	"SEProject/Middleware"
	"SEProject/User/auth"
	"SEProject/User/types"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	service *Service
}

func NewHandler(e *echo.Echo, service *Service) {
	h := &Handler{service: service}

	api := e.Group("/user")
	e.POST("/register", h.Register)
	e.POST("/login", h.Login)
	api.GET("", h.GetUserByUsername, Middleware.RequireAuth)           // /user?username=abc
	api.POST("", h.CreateUser)                                         // POST /user
	api.PUT("/:id", h.UpdateUser, Middleware.RequireRoles("admin"))    // PUT /user/5
	api.DELETE("/:id", h.DeleteUser, Middleware.RequireRoles("admin")) // DELETE /user/5
	api.GET("/all", h.GetAllUsers, Middleware.RequireAuth, Middleware.RequireRoles("admin"))
	api.GET("/me", h.GetCurrentUser, Middleware.RequireAuth)

}

// internal/handler.go
// Register godoc
// @Summary Register a new user
// @Description Creates a new user with hashed password
// @Tags User
// @Accept json
// @Produce json
// @Param user body types.RegisterRequest true "User credentials"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /register [post]
func (h *Handler) Register(c echo.Context) error {
	var req types.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Geçersiz istek"})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Şifre hashlenemedi"})
	}

	u := &types.User{
		Username:  req.Username,
		Password:  string(hashed),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		RoleID:    1, // default user
	}

	if err := h.service.Register(c.Request().Context(), u); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Kayıt başarısız"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "Kayıt başarılı"})
}

// @Summary Login user
// @Description Authenticates user and returns JWT
// @Tags User
// @Accept json
// @Produce json
// @Param user body types.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (h *Handler) Login(c echo.Context) error {
	var req types.LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Geçersiz giriş verisi"})
	}

	user, err := h.service.GetUserByUsername(c.Request().Context(), req.Username, "")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Kullanıcı bulunamadı",
			"debug": err.Error(),
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Şifre hatalı"})
	}

	// ✅ Burada 5 argüman veriyoruz
	token, err := auth.GenerateJWT(
		user.ID,
		user.Username,
		user.RoleName,
		user.FirstName,
		user.LastName,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Token oluşturulamadı"})
	}
	cookie := &http.Cookie{
		Name:     "Authorization",
		Value:    token,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: false,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	// JSON response
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Giriş başarılı",
	})
}

// @Summary Get user by username
// @Description Get a user by providing their username
// @Tags User
// @Produce json
// @Param username query string true "Username to search"
// @Success 200 {object} types.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /user [get]
func (h *Handler) GetUserByUsername(c echo.Context) error {
	username := c.QueryParam("username")
	if username == "" {
		return c.String(http.StatusBadRequest, "Username parametresi gerekli")
	}

	user, err := h.service.GetUserByUsername(c.Request().Context(), username, "")
	if err != nil {
		return c.String(http.StatusNotFound, "Kullanıcı bulunamadı")
	}
	return c.JSON(http.StatusOK, user)
}

// @Summary Create a new user
// @Description Creates a new user in the system
// @Tags User
// @Accept json
// @Produce json
// @Param user body types.User true "User object"
// @Success 201 {string} string "User created"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Failed to create user"
// @Router /user [post]
func (h *Handler) CreateUser(c echo.Context) error {
	var user types.User
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, "Invalid input")
	}

	if err := h.service.CreateUser(c.Request().Context(), &user); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create user: "+err.Error())
	}

	return c.String(http.StatusCreated, "User created")
}

// @Summary Update user
// @Description Updates mutable user fields (username, password, first_name, last_name, role_id)
// @Tags User
// @Accept json
// @Produce json
// @Param id   path int                    true  "User ID"
// @Param body body types.UpdateUserRequest true "Updatable fields"
// @Success 200 {string} string "User updated"
// @Failure 400 {string} string "Invalid id or input"
// @Failure 500 {string} string "Failed to update user"
// @Router /user/{id} [put]
func (h *Handler) UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid id")
	}

	var req types.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid input")
	}

	// password boşsa şifreyi koru
	if req.Password == "" {
		if err := h.service.UpdateUserPartial(c.Request().Context(), id, &req); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to update user")
		}
	} else { // yeni şifre hash’le
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Password hash error")
		}
		req.Password = string(hashed)
		if err := h.service.UpdateUserPartial(c.Request().Context(), id, &req); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to update user")
		}
	}

	return c.String(http.StatusOK, "User updated")
}

// @Summary Delete user
// @Description Deletes a user by ID
// @Tags User
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {string} string "User deleted"
// @Failure 400 {string} string "Invalid id"
// @Failure 500 {string} string "Failed to delete user"
// @Router /user/{id} [delete]
func (h *Handler) DeleteUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid id")
	}

	if err := h.service.DeleteUser(c.Request().Context(), id); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete user")
	}

	return c.String(http.StatusOK, "User deleted")
}

// @Summary Get all users
// @Description Returns all users
// @Tags User
// @Produce json
// @Success 200 {array} types.User
// @Failure 500 {object} map[string]string
// @Router /user/all [get]
func (h *Handler) GetAllUsers(c echo.Context) error {
	users, err := h.service.GetAllUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Kullanıcılar alınamadı"})
	}
	return c.JSON(http.StatusOK, users)
}

// internal/handler.go
// GetCurrentUser godoc
// @Summary Get current user
// @Description Returns the currently authenticated user's info
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /user/me [get]
func (h *Handler) GetCurrentUser(c echo.Context) error {
	// Token'ı al
	userToken := c.Get("user").(*jwt.Token)
	if userToken == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Geçerli kullanıcı bilgisi bulunamadı"})
	}

	// Claims kısmını al
	claims := userToken.Claims.(*Middleware.Claims)

	// Kullanıcı bilgilerini yanıt olarak döndür
	response := map[string]interface{}{
		"id":        claims.UserID,
		"username":  claims.Username,
		"role":      claims.Role,
		"firstName": claims.FirstName, // Eklenen first name
		"lastName":  claims.LastName,  // Eklenen last name
	}

	return c.JSON(http.StatusOK, response)
}
