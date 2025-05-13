package Delivery

import (
	"SEProject/Middleware"
	"SEProject/Order/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(g *echo.Group, service *Service) {
	h := &Handler{service: service}
	g.GET("/orders", h.GetOrders)
}

func uid(c echo.Context) int {
	return c.Get("user").(*jwt.Token).Claims.(*Middleware.Claims).UserID
}

// GetOrders godoc
// @Summary      Kuryenin tüm siparişlerini listeler
// @Tags         Delivery
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   types.OrderResponse
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /delivery/orders [get]
func (h *Handler) GetOrders(c echo.Context) error {
	orders, err := h.service.GetOrdersForDeliveryPerson(c.Request().Context(), uid(c))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Siparişler alınamadı"})
	}
	if orders == nil {
		orders = []types.OrderResponse{}
	}
	return c.JSON(http.StatusOK, orders)
}
