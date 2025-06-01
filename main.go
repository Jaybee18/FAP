package main

import (
	"fap-server/handlers"
	"fap-server/services"
	"fmt"
	"net/http"
	"time"
)

func main() {
	// Initialize dependencies
	userService := services.NewUserService()
	authHandler := handlers.NewAuthHandler(userService)
	userHandler := handlers.NewUserHandler(userService)
	placeHandler := handlers.NewPlaceHandler(userService)

	// Setup routes
	http.HandleFunc("/FAPServer/service/fapservice/login", authHandler.Login)
	http.HandleFunc("/FAPServer/service/fapservice/logout", authHandler.Logout)

	http.HandleFunc("/FAPServer/service/fapservice/addUser", userHandler.AddUser)
	http.HandleFunc("/FAPServer/service/fapservice/getBenutzer", userHandler.GetUser) // todo falsch!
	http.HandleFunc("/FAPServer/service/fapservice/checkLoginName", userHandler.CheckLoginName)

	http.HandleFunc("/FAPServer/service/fapservice/getStandort", placeHandler.GetStandortHandler)
	http.HandleFunc("/FAPServer/service/fapservice/setStandort", placeHandler.SetStandortHandler)
	http.HandleFunc("/FAPServer/service/fapservice/getStandortPerAdresse", placeHandler.GetStandortPerAdresseHandler)
	http.HandleFunc("/FAPServer/service/fapservice/getOrt", placeHandler.GetOrtHandler)

	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			userService.CleanupSessions()
		}
	}()

	// Start server
	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
