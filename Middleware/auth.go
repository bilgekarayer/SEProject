// Register godoc
// @Summary Register a new user
// @Description Creates a new user with hashed password
// @Tags User
// @Accept json
// @Produce json
// @Param user body struct{Username string `json:"username"`; Password string `json:"password"`} true "User credentials"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /register [post]

package Middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	UserID    int    `json:"uid"`
	Username  string `json:"uname"`
	Role      string `json:"role"`
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
	jwt.RegisteredClaims
}

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 1. Cookie'den JWT'yi al
		cookie, err := c.Cookie("Authorization")
		if err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		tokenStr := cookie.Value

		// 2. Token'ı parse et ve doğrula
		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			// HMAC kontrolü
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Beklenmeyen imzalama yöntemi: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.NoContent(http.StatusUnauthorized)
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			return c.NoContent(http.StatusUnauthorized)
		}

		// 3. Expiration kontrolü
		if claims.ExpiresAt.Time.Before(time.Now()) {
			return c.NoContent(http.StatusUnauthorized)
		}

		// 4. Kullanıcıyı sub claim'den bul (örnek olarak sub = user ID)
		userID := strconv.Itoa(claims.UID)
		if userID == "" {
			return c.NoContent(http.StatusUnauthorized)
		}
		c.Set("userID", userID)

		// 5. İsteğe kullanıcıyı bağla
		c.Set("userID", userID)


func RequireRole(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("Authorization")
			if err != nil {
				return c.NoContent(http.StatusUnauthorized)
			}

			tokenStr := cookie.Value
			token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
			if err != nil || !token.Valid {
				return c.NoContent(http.StatusUnauthorized)
			}

			claims, ok := token.Claims.(*Claims)
			if !ok {
				return c.NoContent(http.StatusUnauthorized)
			}

			fmt.Println("🛡 Gelen rol:", claims.Role, "| Beklenen rol:", requiredRole)

			if strings.ToLower(claims.Role) != strings.ToLower(requiredRole) {
				fmt.Println("⛔ ROL ENGELLENDİ")
				return c.NoContent(http.StatusForbidden)
			}

			c.Set("userID", claims.UID)
			c.Set("userRole", claims.Role)

			return next(c)
		}
	}
}
