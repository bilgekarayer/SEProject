package Restaurant

import (
	"SEProject/Middleware"
	"SEProject/Restaurant/types"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(e *echo.Echo, service *Service) {
	h := &Handler{service: service}
	e.GET("/restaurants", h.GetAllRestaurants)
	e.POST("/admin/restaurant", h.CreateWithImage, Middleware.RequireRoles("admin", "restaurant_admin"))
	e.PUT("/admin/restaurant/:id", h.UpdateRestaurant, Middleware.RequireRoles("admin", "restaurant_admin"))
	e.PUT("/admin/restaurant/:id/image", h.UploadImage, Middleware.RequireRoles("admin", "restaurant_admin"))
	e.DELETE("/admin/restaurant/:id", h.DeleteRestaurant, Middleware.RequireRoles("admin", "restaurant_admin"))
}

// GetAllRestaurants godoc
// @Summary      List all restaurants
// @Tags         Restaurant
// @Produce      json
// @Success      200  {array}   types.Restaurant
// @Failure      500  {object}  map[string]string
// @Router       /restaurants [get]
func (h *Handler) GetAllRestaurants(c echo.Context) error {
	list, err := h.service.GetAllRestaurants(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, "error")
	}
	return c.JSON(http.StatusOK, list)
}

// CreateWithImage godoc
// @Summary      Create restaurant (multipart)
// @Tags         Restaurant
// @Security     BearerAuth
// @Accept       mpfd
// @Produce      json
// @Param        name         formData  string  true   "Name"
// @Param        description  formData  string  true   "Description"
// @Param        location     formData  string  true   "Location"
// @Param        cuisine      formData  string  true   "Cuisine"
// @Param        avg_price    formData  int     true   "Avg price"
// @Param        rating       formData  number  true   "Rating"
// @Param        image        formData  file    true   "Image file"
// @Success      201 {object} types.Restaurant
// @Failure      400 {string} string  "bad"
// @Failure      500 {string} string  "fail"
// @Router       /admin/restaurant [post]
func (h *Handler) CreateWithImage(c echo.Context) error {
	// 1) form alanları
	avgPrice, _ := strconv.Atoi(c.FormValue("avg_price"))
	rating, _ := strconv.ParseFloat(c.FormValue("rating"), 64)

	r := types.Restaurant{
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		Location:    c.FormValue("location"),
		Cuisine:     c.FormValue("cuisine"),
		AvgPrice:    avgPrice,
		Rating:      rating,
	}

	// 2) DB’ye kaydet ve id al
	id, err := h.service.CreateRestaurant(c.Request().Context(), &r)
	if err != nil {
		return c.String(http.StatusInternalServerError, "fail")
	}

	// 3) dosya
	file, err := c.FormFile("image")
	if err != nil {
		return c.String(http.StatusBadRequest, "image missing")
	}
	src, _ := file.Open()
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	dstPath := filepath.Join("uploads", "restaurants",
		strconv.Itoa(id)+ext)
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return c.String(http.StatusInternalServerError, "mkdir fail")
	}
	dst, _ := os.Create(dstPath)
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return c.String(http.StatusInternalServerError, "copy fail")
	}

	url := "/static/restaurants/" + strconv.Itoa(id) + ext
	if err = h.service.UpdateRestaurantImage(c.Request().Context(), id, url); err != nil {
		return c.String(http.StatusInternalServerError, "db img fail")
	}
	r.ID, r.ImageURL = id, url
	return c.JSON(http.StatusCreated, r)
}

// UpdateRestaurant godoc
// @Summary      Update restaurant (multipart)
// @Tags         Restaurant
// @Security     BearerAuth
// @Accept       mpfd
// @Produce      json
// @Param  id           path     int     true   "Restaurant ID"
// @Param  name         formData string  false  "Name"
// @Param  description  formData string  false  "Description"
// @Param  location     formData string  false  "Location"
// @Param  cuisine      formData string  false  "Cuisine"
// @Param  avg_price    formData int     false  "Avg price"
// @Param  rating       formData number  false  "Rating"
// @Param  image        formData file    false  "New image"
// @Success 200 {object} types.Restaurant
// @Failure 400 {string} string "bad"
// @Failure 500 {string} string "fail"
// @Router /admin/restaurant/{id} [put]
func (h *Handler) UpdateRestaurant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "bad id")
	}

	// --- form alanları (opsiyonel) ---
	avgPrice, _ := strconv.Atoi(c.FormValue("avg_price"))
	rating, _ := strconv.ParseFloat(c.FormValue("rating"), 64)

	r := types.Restaurant{
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		Location:    c.FormValue("location"),
		Cuisine:     c.FormValue("cuisine"),
	}

	if avgPrice != 0 {
		r.AvgPrice = avgPrice
	}
	if rating != 0 {
		r.Rating = rating
	}

	// DB güncelle (metin/sayı alanları)
	if err := h.service.UpdateRestaurant(c.Request().Context(), id, &r); err != nil {
		return c.String(http.StatusInternalServerError, "fail")
	}

	// --- resim var mı? ---
	if fh, err := c.FormFile("image"); err == nil {
		src, _ := fh.Open()
		defer src.Close()

		ext := filepath.Ext(fh.Filename)
		dstPath := fmt.Sprintf("uploads/restaurants/%d%s", id, ext)
		if err = os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return c.String(http.StatusInternalServerError, "mkdir fail")
		}
		dst, _ := os.Create(dstPath)
		defer dst.Close()
		if _, err = io.Copy(dst, src); err != nil {
			return c.String(http.StatusInternalServerError, "copy fail")
		}
		url := "/static/restaurants/" + strconv.Itoa(id) + ext
		if err = h.service.UpdateRestaurantImage(c.Request().Context(), id, url); err != nil {
			return c.String(http.StatusInternalServerError, "img update fail")
		}
		r.ImageURL = url
	}

	r.ID = id
	return c.JSON(http.StatusOK, r)
}

// DeleteRestaurant godoc
// @Summary      Delete restaurant
// @Tags         Restaurant
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Restaurant ID"
// @Success      200  {string}  string  "ok"
// @Failure      400  {string}  string  "bad"
// @Failure      500  {string}  string  "fail"
// @Router       /admin/restaurant/{id} [delete]
func (h *Handler) DeleteRestaurant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "bad")
	}
	if err := h.service.DeleteRestaurant(c.Request().Context(), id); err != nil {
		return c.String(http.StatusInternalServerError, "fail")
	}
	return c.String(http.StatusOK, "ok")
}

// UploadImage godoc
// @Summary      Upload restaurant image
// @Tags         Restaurant
// @Security     BearerAuth
// @Accept       mpfd
// @Produce      json
// @Param        id    path      int     true  "Restaurant ID"
// @Param        image formData  file    true  "Image file"
// @Success      200   {object}  map[string]string
// @Failure      400   {string}  string  "bad"
// @Failure      500   {string}  string  "fail"
// @Router       /admin/restaurant/{id}/image [put]
func (h *Handler) UploadImage(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "bad id")
	}
	file, err := c.FormFile("image")
	if err != nil {
		return c.String(http.StatusBadRequest, "file missing")
	}
	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusInternalServerError, "open fail")
	}
	defer src.Close()
	ext := filepath.Ext(file.Filename)
	dstPath := fmt.Sprintf("uploads/restaurants/%d%s", id, ext)
	if err = os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return c.String(http.StatusInternalServerError, "mkdir fail")
	}
	dst, err := os.Create(dstPath)
	if err != nil {
		return c.String(http.StatusInternalServerError, "create fail")
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return c.String(http.StatusInternalServerError, "copy fail")
	}
	url := "/static/restaurants/" + strconv.Itoa(id) + ext
	if err = h.service.UpdateRestaurantImage(c.Request().Context(), id, url); err != nil {
		return c.String(http.StatusInternalServerError, "db fail")
	}
	return c.JSON(http.StatusOK, map[string]string{"image_url": url})
}
