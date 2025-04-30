package main

import (
	menu "SEProject/Menu"
	order "SEProject/Order"
	restaurant "SEProject/Restaurant"
	user "SEProject/User"
	"SEProject/config"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	// DB başlat
	config.InitDB()
	//defer config.DB.Close()

	// Echo başlat
	e := echo.New()

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

	// PORT
	log.Println("Sunucu 8080 portunda çalışıyor...")
	e.Logger.Fatal(e.Start(":8080"))
}
