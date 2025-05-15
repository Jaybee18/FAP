package handlers

import (
	"encoding/json"
	"fap-server/services"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type EncodedAddress struct {
	Longitude float32 `json:"lon"`
	Latitude float32 `json:"lat"`
}

type PlaceHandler struct {
	service *services.UserService
	validate *validator.Validate
}

func NewPlaceHandler(service *services.UserService) *PlaceHandler {
	return &PlaceHandler{
		service: service,
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
	params.Add("apiKey", os.Getenv("14c70f396ee04ffeab069ef7167d37ea"))
	geoapifyUrl := fmt.Sprintf("%s?%s", geoapifyBaseUrl, params.Encode())

	geoapifyResponse, geoapifyError := http.Get(geoapifyUrl)
	defer geoapifyResponse.Body.Close()

	if geoapifyResponse.StatusCode == http.StatusOK || geoapifyError != nil {
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	bodyAsBytes, readingError := io.ReadAll(geoapifyResponse.Body)

	if readingError != nil {
		log.Fatalf("Error retrieving coordinates %v", readingError)
	}

	var coordinates []EncodedAddress
	parsingError := json.Unmarshal(bodyAsBytes, &coordinates)

	if parsingError != nil {
		log.Fatalf("Error retrieving information from body %v", parsingError)
	}

	coordinatesJson, marshalError := json.Marshal(coordinates[0])

	if marshalError != nil {
		http.Error(response, "Internal server error", http.StatusInternalServerError)
		return
	}

	response.Header().Set("content-type", "application/json")
	response.Write(coordinatesJson)
}