package models

type User struct {
	// Required fields
	LoginName string `json:"loginName"`
	Password  struct {
		Password string `json:"passwort"`
	} `json:"passwort"`
	FirstName string `json:"vorname"`
	LastName  string `json:"nachname"`
	
	// Optional fields
	Street    *string `json:"strasse,omitempty"`
	ZipCode   *string `json:"plz,omitempty"`
	City      *string `json:"ort,omitempty"`
	Country   *string `json:"land,omitempty"`
	Phone     *string `json:"telefon,omitempty"`
	Email     *struct {
		Address string `json:"adresse,omitempty"`
	} `json:"email,omitempty"`
}

type Response struct {
	Result  bool   `json:"ergebnis"`
	Message string `json:"meldung"`
}