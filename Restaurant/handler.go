package Restaurant

import (
	"SEProject/Restaurant/types"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(e *echo.Echo, service *Service) {
	h := &Handler{service: service}
	e.GET("/restaurants", h.GetAllRestaurants)
	e.POST("/admin/restaurant", h.CreateRestaurant)
	e.PUT("/admin/restaurant/:id", h.UpdateRestaurant)
	e.DELETE("/admin/restaurant/:id", h.DeleteRestaurant)
}

func (h *Handler) GetAllRestaurants(c echo.Context) error {
	restaurants, err := h.service.GetAllRestaurants(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, "Restoranlar getirilemedi")
	}
	return c.JSON(http.StatusOK, restaurants)
}

func (h *Handler) CreateRestaurant(c echo.Context) error {
	var r types.Restaurant
	if err := c.Bind(&r); err != nil {
		return c.String(http.StatusBadRequest, "Geçersiz veri")
	}
	if err := h.service.CreateRestaurant(c.Request().Context(), &r); err != nil {
		return c.String(http.StatusInternalServerError, "Restoran oluşturulamadı")
	}
	return c.String(http.StatusCreated, "Restoran başarıyla eklendi")
}

func (h *Handler) UpdateRestaurant(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Geçersiz ID")
	}
	var r types.Restaurant
	if err := c.Bind(&r); err != nil {
		return c.String(http.StatusBadRequest, "Geçersiz veri")
	}
	if err := h.service.UpdateRestaurant(c.Request().Context(), id, &r); err != nil {
		return c.String(http.StatusInternalServerError, "Güncelleme başarısız")
	}
	return c.String(http.StatusOK, "Restoran güncellendi")
}

func (h *Handler) DeleteRestaurant(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Geçersiz ID")
	}
	if err := h.service.DeleteRestaurant(c.Request().Context(), id); err != nil {
		return c.String(http.StatusInternalServerError, "Silme başarısız")
	}
	return c.String(http.StatusOK, "Restoran silindi")
}
