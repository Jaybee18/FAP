package handlers

import (
	"encoding/json"
	"fap-server/models"
	"fap-server/pkg"
	"fmt"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Accept") != "application/json" {
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Not acceptable"), http.StatusNotAcceptable)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Unsupported Media Type"), http.StatusUnsupportedMediaType)
		return
	}

	var loginReq models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		fmt.Println(err)
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Malformed json in request body"), http.StatusBadRequest)
		return
	}

	if err := validator.Struct(loginReq); err != nil {
		fmt.Println(err)
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Invalid body"), http.StatusBadRequest)
		return
	}

	sessionID := userService.Login(loginReq.LoginName, loginReq.Password.Password)
	if sessionID == "" {
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Fehler beim einloggen"), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.LoginResponse{
		SessionID: sessionID,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Accept") != "application/json" {
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Not acceptable"), http.StatusNotAcceptable)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Unsupported Media Type"), http.StatusUnsupportedMediaType)
		return
	}

	var logoutReq models.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&logoutReq); err != nil {
		fmt.Println(err)
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Malformed json in request body"), http.StatusBadRequest)
		return
	}

	if err := validator.Struct(logoutReq); err != nil {
		fmt.Println(err)
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Invalid body"), http.StatusBadRequest)
		return
	}

	if !userService.ValidSession(logoutReq.LoginName, logoutReq.Session) {
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Invalid session id or username"), http.StatusBadRequest)
		return
	}

	success := userService.Logout(logoutReq.Session, logoutReq.LoginName)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.LogoutResponse{
		Result: success,
	})
}
