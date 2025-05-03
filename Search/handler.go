package Search

import (
	"SEProject/Restaurant"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	service *Restaurant.Service
}

func NewHandler(e *echo.Echo, service *Restaurant.Service) {
	h := &Handler{service: service}
	e.GET("/restaurants/search", h.SearchRestaurants)
}

func (h *Handler) SearchRestaurants(c echo.Context) error {
	ctx := c.Request().Context()

	all, err := h.service.GetAllRestaurants(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Restoranlar getirilemedi")
	}

	keyword := c.QueryParam("q")
	cuisine := c.QueryParam("cuisine")
	location := c.QueryParam("location") // ✅ eklendi

	result := all
	if keyword != "" {
		result = SearchRestaurants(result, keyword)
	}
	if cuisine != "" {
		result = FilterByCuisine(result, cuisine)
	}
	if location != "" {
		result = FilterByLocation(result, location) // ✅ eklendi
	}

	return c.JSON(http.StatusOK, result)
}
