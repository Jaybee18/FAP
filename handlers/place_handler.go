package handlers

import (
	"encoding/json"
	"fap-server/models"
	"fap-server/pkg"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

/*
Handler for the /setStandort route
Requires method PUT
Requires Content-Type application/json
Returns Content-Type application/json
*/
func SetStandortHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPut {
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Method Not Allowed"), http.StatusMethodNotAllowed)
		return
	}

	if request.Header.Get("Accept") != "application/json" {
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Not acceptable"), http.StatusNotAcceptable)
		return
	}

	if request.Header.Get("Content-Type") != "application/json" {
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Unsupported Media Type"), http.StatusUnsupportedMediaType)
		return
	}

	rawBody, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Println(err)
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Internal Server Error"), http.StatusInternalServerError)
		return
	}

	var req models.SetStandortRequest
	err = json.Unmarshal(rawBody, &req)
	if err != nil {
		fmt.Println(err)
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Internal Server Error"), http.StatusInternalServerError)
		return
	}

	if err := validator.Struct(req); err != nil {
		fmt.Println(err)
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Bad request"), http.StatusBadRequest)
		return
	}

	if !userService.ValidSession(req.LoginName, req.SessionId) {
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Session ist invalide"), http.StatusUnauthorized)
		return
	}

	err = userService.SetStandortOfUser(req.LoginName, req.Location)
	if err != nil {
		fmt.Println(err)
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Internal Server Error"), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(map[string]string{
		"ergebnis": "erfolgreich",
	})
}

/*
Handler for the /getStandort route
Requires method GET
Returns Content-Type application/json
*/
func GetStandortHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Method Not Allowed"), http.StatusMethodNotAllowed)
		return
	}

	if request.Header.Get("Accept") != "application/json" {
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Not acceptable"), http.StatusNotAcceptable)
		return
	}

	query := request.URL.Query()
	loginName := query.Get("login")
	sessionId := query.Get("session")
	searchName := query.Get("id")

	if loginName == "" || sessionId == "" || searchName == "" {
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "login, session und id sind erforderliche query parameter"), http.StatusBadRequest)
		return
	}

	if !userService.ValidSession(loginName, sessionId) {
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Session ist invalide"), http.StatusUnauthorized)
		return
	}

	location, err := userService.GetStandortOfUser(searchName)
	if err != nil {
		fmt.Println(err)
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Internal Server Error"), http.StatusInternalServerError)
		return
	}

	standort := models.Standort{
		Location: *location,
	}
	rawJson, err := json.Marshal(standort)
	if err != nil {
		fmt.Println(err)
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Internal Server Error"), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(rawJson)
}

/*
Handler for the /getStandortPerAdresse route
Requires method GET
Returns Content-Type application/json
*/
func GetStandortPerAdresseHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Method Not Allowed"), http.StatusMethodNotAllowed)
		return
	}

	query := request.URL.Query()
	country := query.Get("land")
	postalCode := query.Get("plz")
	place := query.Get("ort")
	street := query.Get("strasse")

	if country == "" || postalCode == "" || place == "" || street == "" {
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "land, plz, ort und strasse sind erforderliche query parameter"), http.StatusBadRequest)
		return
	}

	coordinates, err := pkg.GetLocationByAdressGeoapify(country, postalCode, place, street)
	if err != nil {
		fmt.Println(err)
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Internal Server Error"), http.StatusInternalServerError)
		return
	}

	coordinatesJson, err := json.Marshal(coordinates)
	if err != nil {
		fmt.Println(err)
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Internal Server Error"), http.StatusInternalServerError)
		return
	}

	response.Header().Set("content-type", "application/json")
	response.Write(coordinatesJson)
}

// GetOrtHandler just calls the geonames api and forwards the response body without
// making any changes or processing it in any way. This means that the responses are
// probably different from the rest of the api, but the way the endpoint is described
// it is only meant to be a proxy to the geonames api.
// Handler for the /getOrt route
// Requires method GET
// Returns Content-Type application/json
func GetOrtHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Method Not Allowed"), http.StatusMethodNotAllowed)
		return
	}

	query := request.URL.Query()
	postalcode := query.Get("postalcode")
	username := query.Get("username")

	baseUrl := "http://api.geonames.org/postalCodeSearchJSON"
	params := url.Values{}
	params.Add("postalcode", postalcode)
	params.Add("username", username)
	requestUrl := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	resp, err := http.Get(requestUrl)
	defer resp.Body.Close()
	// only give out internal server errors that are actually caused on our side
	if err != nil {
		fmt.Println(err)
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Internal Server Error"), http.StatusInternalServerError)
		return
	}
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Internal Server Error"), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(rawBody)
}
