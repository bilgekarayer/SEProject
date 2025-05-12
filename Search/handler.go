package Search

import (
	"SEProject/Restaurant"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
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
		log.Println("GetAllRestaurants error:", err) // 🛠️ Konsola bas
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
	price := c.QueryParam("price")   // örn. "100"
	rating := c.QueryParam("rating") // örn. "4.5"

	if price != "" {
		maxPrice, _ := strconv.Atoi(price)
		result = FilterByPrice(result, maxPrice)
	}
	if rating != "" {
		minRating, _ := strconv.ParseFloat(rating, 64)
		result = FilterByRating(result, minRating)
	}

	return c.JSON(http.StatusOK, result)
}
