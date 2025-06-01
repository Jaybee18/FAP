package pkg

import (
	"encoding/json"
	"fap-server/models"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func GetLocationByAdressGeoapify(country string, postalCode string, place string, street string) (models.Location, error) {
	baseUrl := "https://api.geoapify.com/v1/geocode/search"
	params := url.Values{}
	params.Add("text", fmt.Sprintf("%s, %s, %s, %s", street, postalCode, place, country))
	params.Add("apiKey", "14c70f396ee04ffeab069ef7167d37ea")
	requestUrl := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	resp, err := http.Get(requestUrl)
	if err != nil {
		return models.Location{}, err
	}
	rawBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return models.Location{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return models.Location{}, fmt.Errorf("could not get location from geoapify: %s", string(rawBody))
	}

	var respJson models.GeoJSONResponse
	err = json.Unmarshal(rawBody, &respJson)
	if err != nil {
		return models.Location{}, err
	}

	// No coordinates were found for the given adress, but there is also no
	// error, so an empty struct is returned. The handler should handle this
	if len(respJson.Features) == 0 {
		return models.Location{}, nil
	}

	coordinates := respJson.Features[0].Geometry.Coordinates
	return models.Location{
		Longitude: coordinates[0],
		Latitude:  coordinates[1],
	}, nil
}
