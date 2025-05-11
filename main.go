package main

import (
	"fap-server/handlers"
	"fap-server/services"
	"fmt"
	"net/http"
)

func main() {
		// Initialize dependencies
		userService := services.NewUserService()
		userHandler := handlers.NewUserHandler(userService)
	
		// Setup routes
		http.HandleFunc("/FAPServer/service/fapservice/login", userHandler.Login)
		http.HandleFunc("/FAPServer/service/fapservice/addUser", userHandler.AddUser)
	
		// Start server
		fmt.Println("Server starting on :8080...")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
}