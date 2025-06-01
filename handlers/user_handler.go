package handlers

import (
	"encoding/json"
	"fap-server/models"
	"fap-server/services"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service  *services.UserService
	validate *validator.Validate
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.AddUserResponse{
			Result:  false,
			Message: "Method not allowed",
		})
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.AddUserResponse{
			Result:  false,
			Message: "Invalid request body",
		})
		return
	}

	if err := h.validate.Struct(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.AddUserResponse{
			Result:  false,
			Message: err.Error(),
		})
		return
	}

	response, err := h.service.AddUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.AddUserResponse{
			Result:  false,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.AddUserResponse{
			Result:  false,
			Message: "Method not allowed",
		})
		return
	}

	// Get session from query param
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "session ID required",
		})
		return
	}

	// Get login from query param or body
	loginName := r.URL.Query().Get("login")

	if !h.service.ValidSession(loginName, sessionID) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	users := h.service.GetAllUsers()
	var resp models.GetUsersResponse
	for _, user := range users {
		resp.UserList = append(resp.UserList, models.GetUserResponseUser{
			LoginName: user.LoginName,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		})
	}
	rawJson, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(rawJson)
}

func (h *UserHandler) CheckLoginName(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	if request.Method != http.MethodGet {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := request.URL.Query().Get("id")

	if id == "" {
		http.Error(response, "id search param is required", http.StatusBadRequest)
		return
	}

	response.WriteHeader(http.StatusOK)

	if h.service.NameTaken(id) {
		json.NewEncoder(response).Encode(map[string]string{
			"ergebnis": "false",
		})
		return
	}

	json.NewEncoder(response).Encode(map[string]string{
		"ergebnis": "true",
	})
}
