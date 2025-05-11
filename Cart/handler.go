package Cart

import (
	"SEProject/Cart/types"
	"SEProject/Middleware"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(e *echo.Echo, service *Service) {
	h := &Handler{service: service}
	group := e.Group("/cart", Middleware.RequireAuth)

	group.GET("/get-cart-items", h.GetCart)
	group.POST("/add-to-cart", h.AddToCart)
	group.DELETE("/delete/:menuId", h.RemoveFromCart)
}

// @Summary Add item to cart
// @Description Adds a product to the authenticated user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item body types.CartItemRequest true "Item to add"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cart/add-to-cart [post]
func (h *Handler) AddToCart(c echo.Context) error {
	var req types.CartItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	uid := c.Get("userID").(int)

	item := &types.CartItem{
		UserID:    uid,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := h.service.AddItem(c.Request().Context(), item); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not add item"})
	}
	return c.JSON(http.StatusCreated, echo.Map{"message": "item added"})
}

// @Summary Remove item from cart
// @Description Removes a menu item from the authenticated user's cart
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Param menuId path int true "Menu ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cart/delete/{menuId} [delete]
func (h *Handler) RemoveFromCart(c echo.Context) error {
	menuID, _ := strconv.Atoi(c.Param("menuId"))
	uid, _ := c.Get("userID").(int)

	if err := h.service.RemoveItem(c.Request().Context(), uid, menuID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not remove item"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "item removed"})
}

// @Summary Get user's cart items
// @Description Retrieves all items in the authenticated user's cart
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Success 200 {array} types.CartItem
// @Failure 500 {object} map[string]string
// @Router /cart/get-cart-items [get]
func (h *Handler) GetCart(c echo.Context) error {
	uid, _ := c.Get("userID").(int)

	items, err := h.service.GetItems(c.Request().Context(), uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not retrieve items"})
	}
	return c.JSON(http.StatusOK, items)
}
