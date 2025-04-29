package main

import (
	"SEProject/config"
	"SEProject/internal"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	config.InitDB()
	//defer config.DB.Close()
	userRepo := internal.NewRepository(config.DB)
	userService := internal.NewService(userRepo)
	userHandler := internal.NewHandler(userService)

	// ✅ Router yapılandırması burada
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			userHandler.CreateUser(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetUserByUsername(w, r)
		case http.MethodPut:
			userHandler.UpdateUser(w, r)
		case http.MethodDelete:
			userHandler.DeleteUser(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server 8080 portunda çalışıyor...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
