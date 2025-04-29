package main

import (
	"SEProject/config"
	"SEProject/internal"
	_ "github.com/lib/pq"
)

func main() {
	config.InitDB()
	defer config.DB.Close()
	userRepo := internal.NewRepository(config.DB)
	userService := internal.NewService(userRepo)
	internal.NewHandler(userService)
}
