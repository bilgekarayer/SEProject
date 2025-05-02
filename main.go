package main

import (
	"context"
	"log"
	"os"

	menu "SEProject/Menu"
	order "SEProject/Order"
	product "SEProject/Product"
	restaurant "SEProject/Restaurant"
	search "SEProject/Search"
	user "SEProject/User"
	"SEProject/config"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Supabase / DATABASE_URL testi (isteğe bağlı)
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

	// USER
	userRepo := user.NewRepository(config.DB)
	userService := user.NewService(userRepo)
	user.NewHandler(e, userService)

	// RESTAURANT
	restaurantRepo := restaurant.NewRepository(config.DB)
	restaurantService := restaurant.NewService(restaurantRepo)
	restaurant.NewHandler(e, restaurantService)

	// SEARCH (Restoranları filtrelemek/arayıp getirmek için)
	search.NewHandler(e, restaurantService)

	// MENU
	menuRepo := menu.NewRepository(config.DB)
	menuService := menu.NewService(menuRepo)
	menu.NewHandler(e, menuService)

	// ORDER
	orderRepo := order.NewRepository(config.DB)
	orderService := order.NewService(orderRepo)
	order.NewHandler(e, orderService)

	// PRODUCT
	productRepo := product.NewProductRepository(config.DB)
	productService := product.NewProductService(productRepo)
	productHandler := product.NewProductHandler(productService)

	// PRODUCT endpoint tanımları
	e.POST("/products", productHandler.CreateProduct)
	e.PUT("/products/:id", productHandler.UpdateProduct)
	e.DELETE("/products/:id", productHandler.DeleteProduct)

	// SUNUCU BAŞLATILDI
	log.Println("Sunucu 8080 portunda çalışıyor...")
	e.Logger.Fatal(e.Start(":8080"))
}
