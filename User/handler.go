// internal/handler.go

package User

import (
	"SEProject/User/types"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(e *echo.Echo, service *Service) {
	h := &Handler{service: service}

	api := e.Group("/user")
	e.POST("/register", h.Register)
	e.POST("/login", h.Login)
	api.GET("", h.GetUserByUsername) // /user?username=abc
	api.POST("", h.CreateUser)       // POST /user
	api.PUT("/:id", h.UpdateUser)    // PUT /user/5
	api.DELETE("/:id", h.DeleteUser) // DELETE /user/5

}

func (h *Handler) Register(c echo.Context) error {
	var req types.User
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Geçersiz veri")
	}
	if req.Username == "" || req.Password == "" {
		return c.String(http.StatusBadRequest, "Kullanıcı adı ve şifre zorunludur")
	}
	if err := h.service.Register(c.Request().Context(), &req); err != nil {
		return c.String(http.StatusInternalServerError, "Kayıt başarısız: "+err.Error())
	}
	return c.String(http.StatusCreated, "Kullanıcı kaydı başarıyla oluşturuldu")
}

func (h *Handler) Login(c echo.Context) error {
	var req types.User
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Geçersiz giriş verisi")
	}
	ok, err := h.service.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil || !ok {
		return c.String(http.StatusUnauthorized, "Giriş başarısız")
	}
	return c.String(http.StatusOK, "Giriş başarılı")
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
