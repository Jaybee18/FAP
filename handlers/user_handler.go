package handlers

import (
	"encoding/json"
	"fap-server/models"
	"fmt"
	"net/http"
)

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.Response{
			Result:  false,
			Message: "Method not allowed",
		})
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Result:  false,
			Message: "Invalid request body",
		})
		return
	}

	if user.LoginName == "" || user.Password.Password == "" || 
	   user.FirstName == "" || user.LastName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Result:  false,
			Message: "Missing required fields",
		})
		return
	}

	// todo call the user service to add the user
	// For now, just print the user to the console
	fmt.Printf("Received user: %+v\n", user)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{
		Result:  true,
		Message: "",
	})
}