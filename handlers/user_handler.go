package handlers

import (
	"encoding/json"
	"fap-server/models"
	"fap-server/pkg"
	"fap-server/services"
	"fmt"
	"net/http"
	"time"

	val "github.com/go-playground/validator/v10"
)

// Singleton validator instance for handlers package
var validator = val.New()

// Singleton user service instance for handlers package
var userService = services.NewUserService()

func init() {
	// Set up an async cleanup loop at the initialisation of the package
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			userService.CleanupSessions()
		}
	}()
}

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Malformed json in request body"), http.StatusBadRequest)
		return
	}

	if err := validator.Struct(user); err != nil {
		fmt.Println(err)
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Invalid body"), http.StatusBadRequest)
		return
	}

	userExists := userService.UserExists(user.LoginName)
	if userExists {
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", fmt.Sprintf("Benutzer mit dem Namen %q existiert bereits", user.LoginName)), http.StatusConflict)
		return
	}

	// Return value can be ignored since the user cannot already exist
	// because that was checked above
	_ = userService.AddUser(user)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"ergebnis": "erfolgreich",
	})
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	sessionId := r.URL.Query().Get("session")
	loginName := r.URL.Query().Get("login")

	if sessionId == "" || loginName == "" {
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "session und login sind erforderliche parameter"), http.StatusBadRequest)
		return
	}

	if !userService.ValidSession(loginName, sessionId) {
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Unauthorized"), http.StatusUnauthorized)
		return
	}

	users := userService.GetAllUsers()
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
		pkg.JsonError(w, pkg.GenericResponseJson("Fehler", "Internal Server Error"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(rawJson)
}

func CheckLoginNameHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	if request.Method != http.MethodGet {
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	id := request.URL.Query().Get("id")

	if id == "" {
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "id search param is required"), http.StatusBadRequest)
		return
	}

	response.WriteHeader(http.StatusOK)

	if userService.UserExists(id) {
		json.NewEncoder(response).Encode(map[string]string{
			"ergebnis": "false",
		})
		return
	}

	json.NewEncoder(response).Encode(map[string]string{
		"ergebnis": "true",
	})
}
