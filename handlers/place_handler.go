package handlers

import (
	"encoding/json"
	"fap-server/services"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/go-playground/validator/v10"
)

type GeoJSONResponse struct {
	Features []Feature `json:"features"`
}

type Feature struct {
	Geometry Geometry `json:"geometry"`
}

type Geometry struct {
	Coordinates []float64 `json:"coordinates"`
}

type EncodedAddress struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

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

func (h *PlaceHandler) GetStandortPerAdresseHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(response, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	query := request.URL.Query()

	country := query.Get("land")
	postalCode := query.Get("plz")
	place := query.Get("ort")
	street := query.Get("strasse")

	if country == "" || postalCode == "" || place == "" || street == "" {
		http.Error(response, "land, plz, ort und strasse sind erforderliche query parameter", http.StatusBadRequest)
		return
	}

	geoapifyBaseUrl := "https://api.geoapify.com/v1/geocode/search"
	params := url.Values{}
	params.Add("text", fmt.Sprintf("%s, %s, %s, %s", street, postalCode, place, country))
	params.Add("apiKey", "14c70f396ee04ffeab069ef7167d37ea")
	geoapifyUrl := fmt.Sprintf("%s?%s", geoapifyBaseUrl, params.Encode())

	geoapifyResponse, geoapifyError := http.Get(geoapifyUrl)
	defer geoapifyResponse.Body.Close()

	if geoapifyResponse.StatusCode != http.StatusOK || geoapifyError != nil {
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	bodyAsBytes, readingError := io.ReadAll(geoapifyResponse.Body)

	if readingError != nil {
		log.Fatalf("Error retrieving coordinates %v", readingError)
	}

	var geoapifyJson GeoJSONResponse
	parsingError := json.Unmarshal(bodyAsBytes, &geoapifyJson)

	if parsingError != nil {
		log.Fatalf("Error retrieving information from body %v", parsingError)
	}

	features := geoapifyJson.Features

	if len(features) < 1 {
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	coordinates := features[0].Geometry.Coordinates

	encodedAdress := &EncodedAddress{
		Longitude: coordinates[0],
		Latitude:  coordinates[1],
	}

	coordinatesJson, marshalError := json.Marshal(encodedAdress)

	if marshalError != nil {
		http.Error(response, "Internal server error", http.StatusInternalServerError)
		return
	}

	response.Header().Set("content-type", "application/json")
	response.Write(coordinatesJson)
}

func (h *PlaceHandler) GetOrtHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(response, "Method Not Allowed", http.StatusMethodNotAllowed)
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
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(rawBody)
}
