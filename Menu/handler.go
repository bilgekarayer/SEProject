package Menu

import (
	"SEProject/Menu/types"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(e *echo.Echo, service *Service) {
	h := &Handler{service: service}
	e.GET("/restaurants/:id/menu", h.GetMenuByRestaurant)
	e.POST("/restaurant/menu", h.CreateMenuItem)
	e.PUT("/restaurant/menu/:id", h.UpdateMenuItem)
	e.DELETE("/restaurant/menu/:id", h.DeleteMenuItem)
}

// GetMenuByRestaurant godoc
// @Summary Get menu by restaurant ID
// @Description Returns menu items for a given restaurant
// @Tags Menu
// @Produce json
// @Param id path int true "Restaurant ID"
// @Success 200 {array} types.Menu
// @Failure 400 {string} string "Geçersiz restoran ID"
// @Failure 500 {string} string "Menü getirilemedi"
// @Router /restaurants/{id}/menu [get]
func (h *Handler) GetMenuByRestaurant(c echo.Context) error {
	restID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Geçersiz restoran ID")
	}
	menu, err := h.service.GetMenuByRestaurantID(c.Request().Context(), restID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Menü getirilemedi")
	}
	return c.JSON(http.StatusOK, menu)
}

// CreateMenuItem godoc
// @Summary Create a new menu item
// @Description Adds a new item to a restaurant's menu
// @Tags Menu
// @Accept json
// @Produce json
// @Param item body types.Menu true "Menu item"
// @Success 201 {string} string "Menü ürünü eklendi"
// @Failure 400 {string} string "Geçersiz menü verisi"
// @Failure 500 {string} string "Menü ürünü eklenemedi"
// @Router /restaurant/menu [post]
func (h *Handler) CreateMenuItem(c echo.Context) error {
	var item types.Menu
	if err := c.Bind(&item); err != nil {
		return c.String(http.StatusBadRequest, "Geçersiz menü verisi")
	}
	if err := h.service.CreateMenuItem(c.Request().Context(), &item); err != nil {
		return c.String(http.StatusInternalServerError, "Menü ürünü eklenemedi")
	}
	return c.String(http.StatusCreated, "Menü ürünü eklendi")
}

// UpdateMenuItem godoc
// @Summary Update a menu item
// @Description Updates the details of an existing menu item
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path int true "Menu Item ID"
// @Param item body types.Menu true "Menu item data"
// @Success 200 {string} string "Menü ürünü güncellendi"
// @Failure 400 {string} string "Geçersiz veri"
// @Failure 500 {string} string "Menü ürünü güncellenemedi"
// @Router /restaurant/menu/{id} [put]
func (h *Handler) UpdateMenuItem(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Geçersiz ID")
	}
	var item types.Menu
	if err := c.Bind(&item); err != nil {
		return c.String(http.StatusBadRequest, "Geçersiz veri")
	}
	if err := h.service.UpdateMenuItem(c.Request().Context(), id, &item); err != nil {
		return c.String(http.StatusInternalServerError, "Menü ürünü güncellenemedi")
	}
	return c.String(http.StatusOK, "Menü ürünü güncellendi")
}

// DeleteMenuItem godoc
// @Summary Delete a menu item
// @Description Deletes an item from the menu
// @Tags Menu
// @Produce json
// @Param id path int true "Menu Item ID"
// @Success 200 {string} string "Menü ürünü silindi"
// @Failure 400 {string} string "Geçersiz ID"
// @Failure 500 {string} string "Menü ürünü silinemedi"
// @Router /restaurant/menu/{id} [delete]
func (h *Handler) DeleteMenuItem(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Geçersiz ID")
	}
	if err := h.service.DeleteMenuItem(c.Request().Context(), id); err != nil {
		return c.String(http.StatusInternalServerError, "Menü ürünü silinemedi")
	}
	return c.String(http.StatusOK, "Menü ürünü silindi")
}
