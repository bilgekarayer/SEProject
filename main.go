// @title SEProject API
// @version 1.0
// @description RESTful API for restaurant ordering system
// @host localhost:8080
// @BasePath /

package main

import (
	"context"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"

	menu "SEProject/Menu"
	order "SEProject/Order"
	restaurant "SEProject/Restaurant"
	user "SEProject/User"
	"SEProject/config"

	_ "SEProject/docs"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	echoSwagger "github.com/swaggo/echo-swagger"
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

	// 2. Klasik InitDB bağlantısı (lib/pq)
	config.InitDB()
	//defer config.DB.Close()
	// 3. Echo Web Framework başlatılıyor
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowCredentials: true,
	}))
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
	for _, r := range e.Routes() {
		log.Println("Route:", r.Method, r.Path)
	}

	// PORT
	log.Println("Sunucu 8080 portunda çalışıyor...")
	e.Logger.Fatal(e.Start(":8080"))
}
