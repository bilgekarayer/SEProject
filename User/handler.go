package User

import (
	"SEProject/Middleware"
	"SEProject/User/auth"
	"SEProject/User/types"
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
	api.POST("/register", h.Register)
	api.POST("/login", h.Login)
	api.GET("", h.GetUserByUsername, Middleware.RequireAuth)
	api.POST("", h.CreateUser)
	api.PUT("/:id", h.UpdateUser)
	api.DELETE("/:id", h.DeleteUser)
	api.GET("/all", h.GetAllUsers, Middleware.RequireAuth, Middleware.RequireRole("Admin"))
}

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

	token, err := auth.GenerateJWT(user.ID, user.Username, user.RoleName)
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

func (h *Handler) GetAllUsers(c echo.Context) error {
	users, err := h.service.GetAllUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Kullanıcılar alınamadı"})
	}
	return c.JSON(http.StatusOK, users)
}
