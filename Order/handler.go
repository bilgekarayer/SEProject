package Order

import (
	"SEProject/Order/types"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(e *echo.Echo, service *Service) {
	h := &Handler{service: service}

	e.POST("/cart/add", h.AddToCart)
	e.GET("/cart", h.GetCart)
	e.POST("/order/place", h.PlaceOrder)
	e.GET("/restaurant/orders", h.GetOrdersByRestaurantID)
	e.PUT("/restaurant/orders/:id/prepare", h.MarkOrderPrepared)
	e.PUT("/restaurant/orders/:id/send", h.MarkOrderSent)
}

func (h *Handler) AddToCart(c echo.Context) error {
	var item types.CartItem
	if err := c.Bind(&item); err != nil {
		return c.String(http.StatusBadRequest, "Invalid input")
	}
	if err := h.service.AddToCart(c.Request().Context(), &item); err != nil {
		log.Println("AddToCart ERROR:", err)
		return c.String(http.StatusInternalServerError, "Failed to add to cart")
	}
	return c.String(http.StatusCreated, "Added to cart")
}

func (h *Handler) GetCart(c echo.Context) error {
	uidStr := c.QueryParam("user_id")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid user_id")
	}
	cart, err := h.service.GetCart(c.Request().Context(), uid)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch cart")
	}
	return c.JSON(http.StatusOK, cart)
}

func (h *Handler) PlaceOrder(c echo.Context) error {
	var req types.PlaceOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid input")
	}
	if err := h.service.PlaceOrder(c.Request().Context(), &req); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to place order")
	}
	return c.String(http.StatusCreated, "Order placed")
}

func (h *Handler) GetOrdersByRestaurantID(c echo.Context) error {
	ridStr := c.QueryParam("restaurant_id")
	rid, err := strconv.Atoi(ridStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid restaurant_id")
	}
	orders, err := h.service.GetOrdersByRestaurantID(c.Request().Context(), rid)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch orders")
	}
	return c.JSON(http.StatusOK, orders)
}

func (h *Handler) MarkOrderPrepared(c echo.Context) error {
	oidStr := c.Param("id")
	oid, err := strconv.Atoi(oidStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid order ID")
	}
	err = h.service.UpdateOrderStatus(c.Request().Context(), oid, "prepared")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to mark prepared")
	}
	return c.String(http.StatusOK, "Order marked as prepared")
}

func (h *Handler) MarkOrderSent(c echo.Context) error {
	oidStr := c.Param("id")
	oid, err := strconv.Atoi(oidStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid order ID")
	}
	err = h.service.UpdateOrderStatus(c.Request().Context(), oid, "sent")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to mark sent")
	}
	return c.String(http.StatusOK, "Order marked as sent")
}
