package Delivery

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	service *Service
}

func NewHandler(g *echo.Group, service *Service) {
	h := &Handler{service: service}
	g.GET("/orders", h.GetOrders)
	g.PUT("/orders/:id", h.UpdateOrderStatus) // ✅ Sipariş güncelleme route'u
}

func (h *Handler) GetOrders(c echo.Context) error {
	uid := c.Get("userID")

	var deliveryPersonID int
	switch v := uid.(type) {
	case int:
		deliveryPersonID = v
	case float64:
		deliveryPersonID = int(v)
	default:
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Kullanıcı ID alınamadı"})
	}

	orders, err := h.service.GetOrdersForDeliveryPerson(c.Request().Context(), deliveryPersonID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Siparişler alınamadı"})
	}

	return c.JSON(http.StatusOK, orders)
}

// ✅ Sipariş durumu güncelleme fonksiyonu

func (h *Handler) UpdateOrderStatus(c echo.Context) error {
	fmt.Println("🔥 UpdateOrderStatus'a girildi")
	idParam := strings.TrimSpace(c.Param("id")) // 👈 \n karakterlerini temizle
	fmt.Println("ID PARAM:", idParam)

	orderID, err := strconv.Atoi(idParam)
	if err != nil {
		fmt.Println("HATA:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Geçersiz sipariş ID"})
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Geçersiz veri"})
	}

	if err := h.service.UpdateOrderStatus(c.Request().Context(), orderID, req.Status); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Sipariş durumu güncellenemedi"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Sipariş durumu güncellendi"})
}
