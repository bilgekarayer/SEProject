package Order

import (
	"SEProject/Middleware"
	"SEProject/Order/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Handler struct {
	s *Service
}

func NewHandler(e *echo.Echo, s *Service) {
	h := &Handler{s}

	auth := e.Group("", Middleware.RequireAuth)

	auth.POST("/cart/add", h.AddToCart)
	auth.GET("/cart", h.GetCart)
	auth.POST("/order/place", h.PlaceOrder)
	auth.GET("/user/orders", h.MyOrders)
	auth.GET("/orders", h.AllOrders, Middleware.RequireRoles("admin", "restaurant_admin"))

	auth.GET("/restaurant/orders", h.RestaurantOrders)
	auth.PUT("/restaurant/orders/:id/prepare", h.MarkPrepared, Middleware.RequireRoles("admin", "restaurant_admin"))
	auth.PUT("/restaurant/orders/:id/delivered", h.MarkDelivered, Middleware.RequireRoles("admin", "delivery_person"))
}

func uid(c echo.Context) int {
	t := c.Get("user").(*jwt.Token)
	return t.Claims.(*Middleware.Claims).UserID
}

// AddToCart godoc
// @Summary      Add item to cart
// @Tags         Cart
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        item  body      types.CartItem  true  "Cart item payload"
// @Success      201   {string}  string          "ok"
// @Failure      400   {string}  string          "bad"
// @Failure      500   {string}  string          "fail"
// @Router       /cart/add [post]
func (h *Handler) AddToCart(c echo.Context) error {
	var it types.CartItem
	if err := c.Bind(&it); err != nil {
		return c.String(http.StatusBadRequest, "bad")
	}
	it.UserID = uid(c)
	if err := h.s.AddToCart(c.Request().Context(), &it); err != nil {
		return c.String(http.StatusInternalServerError, "fail")
	}
	return c.String(http.StatusCreated, "ok")
}

// GetCart godoc
// @Summary      Get current user's cart
// @Tags         Cart
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   types.CartItem
// @Failure      500  {string}  string  "fail"
// @Router       /cart [get]
func (h *Handler) GetCart(c echo.Context) error {
	items, err := h.s.GetCart(c.Request().Context(), uid(c))
	if err != nil {
		return c.String(http.StatusInternalServerError, "fail")
	}
	return c.JSON(http.StatusOK, items)
}

// AllOrders godoc
// @Summary      List all orders (admin)
// @Tags         Order
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   types.OrderResponse
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /orders [get]
func (h *Handler) AllOrders(c echo.Context) error {
	orders, err := h.s.GetAllOrders(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Siparişler alınamadı"})
	}
	if orders == nil {
		orders = []types.OrderResponse{}
	}
	return c.JSON(http.StatusOK, orders)
}

// PlaceOrder godoc
// @Summary      Place a new order
// @Tags         Order
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        order  body      types.PlaceOrderRequest  true  "Order payload"
// @Success      201    {string}  string  "ok"
// @Failure      400    {string}  string  "bad"
// @Failure      500    {string}  string  "fail"
// @Router       /order/place [post]
func (h *Handler) PlaceOrder(c echo.Context) error {
	var req types.PlaceOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad")
	}
	if err := h.s.PlaceOrder(c.Request().Context(), uid(c), &req); err != nil {
		return c.String(http.StatusInternalServerError, "fail")
	}
	return c.String(http.StatusCreated, "ok")
}

// MyOrders godoc
// @Summary      List current user's orders
// @Tags         Order
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   types.OrderResponse
// @Failure      500  {string}  string  "fail"
// @Router       /user/orders [get]
func (h *Handler) MyOrders(c echo.Context) error {
	res, err := h.s.GetOrdersByUser(c.Request().Context(), uid(c))
	if err != nil {
		return c.String(http.StatusInternalServerError, "fail")
	}
	return c.JSON(http.StatusOK, res)
}

// RestaurantOrders godoc
// @Summary      List orders for a restaurant
// @Tags         Order
// @Security     BearerAuth
// @Produce      json
// @Param        restaurant_id  query     int  true  "Restaurant ID"
// @Success      200            {array}   types.OrderResponse
// @Failure      400            {string}  string  "bad"
// @Failure      500            {string}  string  "fail"
// @Router       /restaurant/orders [get]
func (h *Handler) RestaurantOrders(c echo.Context) error {
	ridStr := c.QueryParam("restaurant_id")
	rid, err := strconv.Atoi(ridStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad")
	}
	res, err := h.s.GetOrdersByRestaurant(c.Request().Context(), rid)
	if err != nil {
		return c.String(http.StatusInternalServerError, "fail")
	}
	return c.JSON(http.StatusOK, res)
}

// MarkPrepared godoc
// @Summary      Mark order as prepared
// @Tags         Order
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Order ID"
// @Success      200  {string}  string  "ok"
// @Failure      400  {string}  string  "bad"
// @Failure      500  {string}  string  "fail"
// @Router       /restaurant/orders/{id}/prepare [put]
func (h *Handler) MarkPrepared(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "bad")
	}
	if err = h.s.UpdateOrderStatus(c.Request().Context(), id, "prepared"); err != nil {
		return c.String(http.StatusInternalServerError, "fail")
	}
	return c.String(http.StatusOK, "ok")
}

// MarkSent godoc
// @Summary      Mark order as delivered
// @Tags         Order
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Order ID"
// @Success      200  {string}  string  "ok"
// @Failure      400  {string}  string  "bad"
// @Failure      500  {string}  string  "fail"
// @Router       /restaurant/orders/{id}/delivered [put]
func (h *Handler) MarkDelivered(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "bad")
	}
	if err = h.s.UpdateOrderStatus(c.Request().Context(), id, "delivered"); err != nil {
		return c.String(http.StatusInternalServerError, "fail")
	}
	return c.String(http.StatusOK, "ok")
}
