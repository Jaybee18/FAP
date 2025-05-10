package models

type User struct {
	LoginName string   `json:"loginName" validate:"required"`
	Password  Password `json:"passwort" validate:"required"`
	FirstName string   `json:"vorname" validate:"required"`
	LastName  string   `json:"nachname" validate:"required"`
	Street    *string  `json:"strasse,omitempty"`
	ZipCode   *string  `json:"plz,omitempty"`
	City      *string  `json:"ort,omitempty"`
	Country   *string  `json:"land,omitempty"`
	Phone     *string  `json:"telefon,omitempty"`
	Email     *Email   `json:"email,omitempty"`
}

type Password struct {
	Password string `json:"passwort" validate:"required,min=6"`
}

type Email struct {
	Address string `json:"adresse" validate:"email"`
}

type Response struct {
	Result  bool   `json:"ergebnis"`
	Message string `json:"meldung"`
}