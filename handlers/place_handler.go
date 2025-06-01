package handlers

import (
	"encoding/json"
	"fap-server/models"
	"fap-server/pkg"
	"fap-server/services"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-playground/validator/v10"
)

type PlaceHandler struct {
	service  *services.UserService
	validate *validator.Validate
}

func NewPlaceHandler(service *services.UserService) *PlaceHandler {
	return &PlaceHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (p *PlaceHandler) SetStandortHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPut {
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Method Not Allowed"), http.StatusMethodNotAllowed)
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

	if err := p.validate.Struct(req); err != nil {
		fmt.Println(err)
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Bad request"), http.StatusBadRequest)
		return
	}

	if !p.service.ValidSession(req.LoginName, req.SessionId) {
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Session ist invalide"), http.StatusUnauthorized)
		return
	}

	err = p.service.SetStandortOfUser(req.LoginName, req.Location)
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

func (p *PlaceHandler) GetStandortHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Method Not Allowed"), http.StatusMethodNotAllowed)
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

	if !p.service.ValidSession(loginName, sessionId) {
		pkg.JsonError(response, pkg.GenericResponseJson("Fehler", "Session ist invalide"), http.StatusUnauthorized)
		return
	}

	location, err := p.service.GetStandortOfUser(searchName)
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

func (h *PlaceHandler) GetStandortPerAdresseHandler(response http.ResponseWriter, request *http.Request) {
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
func (h *PlaceHandler) GetOrtHandler(response http.ResponseWriter, request *http.Request) {
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
