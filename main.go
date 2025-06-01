package main

import (
	"fap-server/handlers"
	"fmt"
	"net/http"
)

func main() {
	// Setup routes
	http.HandleFunc("/FAPServer/service/fapservice/login", handlers.LoginHandler)
	http.HandleFunc("/FAPServer/service/fapservice/logout", handlers.LogoutHandler)

	http.HandleFunc("/FAPServer/service/fapservice/addUser", handlers.AddUserHandler)
	http.HandleFunc("/FAPServer/service/fapservice/getBenutzer", handlers.GetUsersHandler)
	http.HandleFunc("/FAPServer/service/fapservice/checkLoginName", handlers.CheckLoginNameHandler)

	http.HandleFunc("/FAPServer/service/fapservice/getStandort", handlers.GetStandortHandler)
	http.HandleFunc("/FAPServer/service/fapservice/setStandort", handlers.SetStandortHandler)
	http.HandleFunc("/FAPServer/service/fapservice/getStandortPerAdresse", handlers.GetStandortPerAdresseHandler)
	http.HandleFunc("/FAPServer/service/fapservice/getOrt", handlers.GetOrtHandler)

	// Start server
	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
