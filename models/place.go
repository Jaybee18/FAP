package models

type PostalCode struct {
	CountryCode string  `json:"countryCode"`
	PostalCode  string  `json:"postalCode"`
	Place       string  `json:"placeName"`
	Latitude    float64 `json:"lat"`
	Longitude   float64 `json:"lng"`
}

type GeonamesResponse struct {
	PostalCodes []PostalCode `json:"postalCodes"`
}
