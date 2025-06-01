package pkg

import (
	"encoding/json"
	"fap-server/models"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func GetLocationByAdress(postalcode string, country string) (models.Location, error) {
	var res models.Location

	baseUrl := "http://api.geonames.org/postalCodeSearchJSON"
	params := url.Values{}
	params.Add("postalcode", postalcode)
	params.Add("username", "mikelong")
	requestUrl := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	resp, err := http.Get(requestUrl)
	defer resp.Body.Close()
	if err != nil {
		return res, err
	}
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	var geonamesResp models.GeonamesResponse
	err = json.Unmarshal(rawBody, &geonamesResp)
	if err != nil {
		return res, err
	}

	if len(geonamesResp.PostalCodes) == 0 {
		return res, nil
	}

	foundLocation := geonamesResp.PostalCodes[0]
	res.Latitude = foundLocation.Latitude
	res.Longitude = foundLocation.Longitude
	return res, nil
}
