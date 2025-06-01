package handlers

import (
	"encoding/json"
	"fap-server/models"
	"fap-server/pkg"
	"fap-server/services"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	service  *services.UserService
	validate *validator.Validate
}

func NewAuthHandler(service *services.UserService) *AuthHandler {
	return &AuthHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	var loginReq models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		fmt.Println(err)
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Malformed json in request body"), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(loginReq); err != nil {
		fmt.Println(err)
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Invalid body"), http.StatusBadRequest)
		return
	}

	sessionID := h.service.Login(loginReq.LoginName, loginReq.Password.Password)
	if sessionID == "" {
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Fehler beim einloggen"), http.StatusBadRequest)
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
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	var logoutReq models.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&logoutReq); err != nil {
		fmt.Println(err)
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Malformed json in request body"), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(logoutReq); err != nil {
		fmt.Println(err)
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Invalid body"), http.StatusBadRequest)
		return
	}

	if !h.service.ValidSession(logoutReq.LoginName, logoutReq.Session) {
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Invalid session id or username"), http.StatusBadRequest)
		return
	}

	success := h.service.Logout(logoutReq.Session, logoutReq.LoginName)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.LogoutResponse{
		Result: success,
	})
}
