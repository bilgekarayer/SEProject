package User

import (
	"SEProject/Middleware"
	"SEProject/User/auth"
	"SEProject/User/types"
	"net/http"
	"strconv"
	"strings"
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
	api.POST("/register", h.Register)
	api.POST("/login", h.Login)
	api.GET("", h.GetUserByUsername, Middleware.RequireAuth) // /user?username=abc
	api.POST("", h.CreateUser)                               // POST /user
	api.PUT("/:id", h.UpdateUser)                            // PUT /user/5
	api.DELETE("/:id", h.DeleteUser)                         // DELETE /user/5

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
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  string(hashed),
	}
	if err := h.service.Register(c.Request().Context(), u); err != nil {
		// UNIQUE constraint hatası kontrolü
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "already exists") {
			return c.JSON(http.StatusConflict, echo.Map{
				"error": "Bu kullanıcı adı zaten kullanılıyor",
			})
		}

		// Diğer tüm hatalar için genel hata
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Kayıt başarısız",
		})
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
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Kullanıcı bulunamadı"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Şifre hatalı"})
	}

	token, err := auth.GenerateJWT(user.ID, user.Username)
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

func (h *Handler) UpdateUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid id")
	}

	var user types.User
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, "Invalid input")
	}

	if err := h.service.UpdateUser(c.Request().Context(), id, &user); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to update user")
	}

	return c.String(http.StatusOK, "User updated")
}

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
