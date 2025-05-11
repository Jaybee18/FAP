package handlers

import (
	"encoding/json"
	"fap-server/models"
	"fap-server/services"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service *services.UserService
	validate *validator.Validate
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
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

func (h *UserHandler) Login( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

    var loginReq models.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Invalid request format",
        })
        return
    }

    sessionID, err := h.service.Login(loginReq.LoginName, loginReq.Password.Password)
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{
            "error": err.Error(),
        })
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(models.LoginResponse{
        SessionID: sessionID,
    })
}