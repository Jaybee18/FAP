package main

import (
	"fap-server/handlers"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/FAPServer/service/fapservice/addUser", handlers.AddUserHandler)

	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}