package Menu

import (
	"SEProject/Menu/types"
	"SEProject/Middleware"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Handler struct{ service *Service }

func NewHandler(e *echo.Echo, svc *Service) {
	h := &Handler{svc}

	// Listeleme
	e.GET("/restaurants/:rid/menu", h.GetMenuByRestaurant)

	// Oluştur + görsel (multipart)
	e.POST("/restaurant/menu", h.CreateMenuItem,
		Middleware.RequireRoles("admin", "restaurant_admin"))

	// Güncelle / Sil
	e.PUT("/restaurant/menu/:id", h.UpdateMenuItem,
		Middleware.RequireRoles("admin", "restaurant_admin"))
	e.DELETE("/restaurant/menu/:id", h.DeleteMenuItem,
		Middleware.RequireRoles("admin", "restaurant_admin"))
}

// GetMenuByRestaurant godoc
// @Summary      List menu items
// @Tags         Menu
// @Produce      json
// @Param        rid path int true "Restaurant ID"
// @Success      200 {array} types.Menu
// @Router       /restaurants/{rid}/menu [get]
func (h *Handler) GetMenuByRestaurant(c echo.Context) error {
	rid, err := strconv.Atoi(c.Param("rid"))
	if err != nil {
		return c.String(http.StatusBadRequest, "bad id")
	}
	list, err := h.service.GetMenuByRestaurantID(c.Request().Context(), rid)
	if err != nil {
		return c.String(http.StatusInternalServerError, "fail")
	}
	return c.JSON(http.StatusOK, list)
}

// CreateMenuItem godoc
// @Summary      Create menu item (multipart)
// @Tags         Menu
// @Security     BearerAuth
// @Accept       mpfd
// @Produce      json
// @Param        restaurant_id formData int     true  "Restaurant ID"
// @Param        name          formData string  true  "Name"
// @Param        price         formData number  true  "Price"
// @Param        image         formData file    true  "Image file"
// @Success      201 {object}  types.Menu
// @Router       /restaurant/menu [post]
func (h *Handler) CreateMenuItem(c echo.Context) error {
	rid, _ := strconv.Atoi(c.FormValue("restaurant_id"))
	price, _ := strconv.ParseFloat(c.FormValue("price"), 64)

	item := types.Menu{
		RestaurantID: rid,
		Name:         c.FormValue("name"),
		Price:        price,
	}

	// 1) DB’ye kaydet, id al
	id, err := h.service.CreateMenuItem(c.Request().Context(), &item)
	if err != nil {
		return c.String(http.StatusInternalServerError, "db fail")
	}

	// 2) Dosya işle
	fh, err := c.FormFile("image")
	if err != nil {
		return c.String(http.StatusBadRequest, "file missing")
	}
	src, _ := fh.Open()
	defer src.Close()

	ext := filepath.Ext(fh.Filename)
	dstPath := filepath.Join("uploads", "menu", strconv.Itoa(id)+ext)
	if err = os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return c.String(http.StatusInternalServerError, "mkdir fail")
	}
	dst, _ := os.Create(dstPath)
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return c.String(http.StatusInternalServerError, "copy fail")
	}

	url := "/static/menu/" + strconv.Itoa(id) + ext
	if err = h.service.UpdateMenuItemImage(c.Request().Context(), id, url); err != nil {
		return c.String(http.StatusInternalServerError, "img update fail")
	}

	item.ID, item.ImageURL = id, url
	return c.JSON(http.StatusCreated, item)
}

// UpdateMenuItem godoc
// @Summary      Update menu item (multipart)
// @Tags         Menu
// @Security     BearerAuth
// @Accept       mpfd
// @Produce      json
// @Param        id            path     int    true   "Menu ID"
// @Param        name          formData string false  "Name"
// @Param        price         formData number false  "Price"
// @Param        image         formData file   false  "New image"
// @Success      200 {object}  types.Menu
// @Failure      400 {string}  string "bad"
// @Failure      500 {string}  string "fail"
// @Router       /restaurant/menu/{id} [put]
func (h *Handler) UpdateMenuItem(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "bad id")
	}

	// ---- form alanları (opsiyonel) ----
	price, _ := strconv.ParseFloat(c.FormValue("price"), 64)
	m := types.Menu{
		Name:  c.FormValue("name"),
		Price: price,
	}

	// güncelle metin/sayı
	if err := h.service.UpdateMenuItem(c.Request().Context(), id, &m); err != nil {
		return c.String(http.StatusInternalServerError, "fail")
	}

	// ---- resim var mı? ----
	if fh, err := c.FormFile("image"); err == nil {
		src, _ := fh.Open()
		defer src.Close()

		ext := filepath.Ext(fh.Filename)
		dstPath := filepath.Join("uploads", "menu", strconv.Itoa(id)+ext)
		if err = os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return c.String(http.StatusInternalServerError, "mkdir fail")
		}
		dst, _ := os.Create(dstPath)
		defer dst.Close()
		if _, err = io.Copy(dst, src); err != nil {
			return c.String(http.StatusInternalServerError, "copy fail")
		}

		url := "/static/menu/" + strconv.Itoa(id) + ext
		if err = h.service.UpdateMenuItemImage(c.Request().Context(), id, url); err != nil {
			return c.String(http.StatusInternalServerError, "img update fail")
		}
		m.ImageURL = url
	}

	m.ID = id
	return c.JSON(http.StatusOK, m)
}

// DeleteMenuItem godoc
// @Summary Delete menu item
// @Tags    Menu
// @Security BearerAuth
// @Param   id path int true "Menu ID"
// @Success 200 {string} string "ok"
// @Router  /restaurant/menu/{id} [delete]
func (h *Handler) DeleteMenuItem(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "bad id")
	}
	if err = h.service.DeleteMenuItem(c.Request().Context(), id); err != nil {
		return c.String(http.StatusInternalServerError, "fail")
	}
	return c.String(http.StatusOK, "ok")
}
