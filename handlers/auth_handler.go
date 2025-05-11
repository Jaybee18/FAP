package handlers

import (
	"encoding/json"
	"fap-server/models"
	"fap-server/services"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
    service *services.UserService
	validate *validator.Validate
}

func NewAuthHandler(service *services.UserService) *AuthHandler {
    return &AuthHandler{
		service: service,
		validate: validator.New(),
	}
}

func (h *AuthHandler) Login( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.AddUserResponse{
			Result:  false,
			Message: "Method not allowed",
		})
		return
	}

    var loginReq models.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Invalid request format",
        })
        return
    }

	if err := h.validate.Struct(loginReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": err.Error(),
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

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.AddUserResponse{
			Result:  false,
			Message: "Method not allowed",
		})
		return
	}

    var logoutReq models.LogoutRequest
    if err := json.NewDecoder(r.Body).Decode(&logoutReq); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(models.LogoutResponse{
            Result: false,
        })
        return
    }

	if err := h.validate.Struct(logoutReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": err.Error(),
        })
        return
	}

    success := h.service.Logout(logoutReq.Session, logoutReq.LoginName)
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(models.LogoutResponse{
        Result: success,
    })
}