package Product

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	Service *ProductService
}

func NewProductHandler(service *ProductService) *ProductHandler {
	return &ProductHandler{Service: service}
}

// POST /products
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var product Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Geçersiz veri"})
	}
	if err := h.Service.CreateProduct(&product); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, product)
}

// PUT /products/:id
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Geçersiz ID"})
	}

	var product Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Geçersiz veri"})
	}
	product.ID = id

	if err := h.Service.UpdateProduct(&product); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, product)
}

// DELETE /products/:id
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Geçersiz ID"})
	}

	if err := h.Service.DeleteProduct(id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Ürün silindi"})
}
