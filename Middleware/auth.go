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
	UID   int    `json:"uid"`
	UName string `json:"uname"`
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

		// 6. Devam et
		return next(c)
	}
}
