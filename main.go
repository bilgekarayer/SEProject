package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	delivery "SEProject/Delivery"
	menu "SEProject/Menu"
	customMiddleware "SEProject/Middleware"
	order "SEProject/Order"
	restaurant "SEProject/Restaurant"
	user "SEProject/User"
	"SEProject/config"

	_ "SEProject/docs"
	_ "github.com/lib/pq"
)

func main() {
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		conn, err := pgx.Connect(context.Background(), databaseURL)
		if err != nil {
			log.Fatalf("Failed to connect to DATABASE_URL: %v", err)
		}
		defer conn.Close(context.Background())

		var version string
		if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
			log.Fatalf("Query failed: %v", err)
		}
		log.Println("Connected to Supabase (pgx):", version)
	} else {
		log.Println("Warning: DATABASE_URL not set, skipping pgx connection test")
	}

	config.InitDB()
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("🌐 Yeni istek:", c.Request().Method, c.Request().URL.Path)
			return next(c)
		}
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// USER
	userRepo := user.NewRepository(config.DB)
	userService := user.NewService(userRepo)
	user.NewHandler(e, userService)

	// RESTAURANT
	restaurantRepo := restaurant.NewRepository(config.DB)
	restaurantService := restaurant.NewService(restaurantRepo)
	restaurant.NewHandler(e, restaurantService)

	// MENU
	menuRepo := menu.NewRepository(config.DB)
	menuService := menu.NewService(menuRepo)
	menu.NewHandler(e, menuService)

	// ORDER
	orderRepo := order.NewRepository(config.DB)
	orderService := order.NewService(orderRepo)
	order.NewHandler(e, orderService)

	// DELIVERY (YENİ EKLENDİ)
	deliveryRepo := delivery.NewRepository(config.DB)
	deliveryService := delivery.NewService(deliveryRepo)

	// /delivery/... yolları sadece auth'lu ve delivery_person rolü olanlara açık
	deliveryGroup := e.Group("/delivery",
		customMiddleware.RequireAuth,
		customMiddleware.RequireRole("delivery_person"),
	)
	delivery.NewHandler(deliveryGroup, deliveryService)

	log.Println("Sunucu 8080 portunda çalışıyor...")
	e.Logger.Fatal(e.Start(":8080"))
}
