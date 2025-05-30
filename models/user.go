package models

import "time"

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

type AddUserResponse struct {
	Result  bool   `json:"ergebnis"`
	Message string `json:"meldung"`
}

type LoginRequest struct {
	LoginName string `json:"loginName" validate:"required"`
	Password struct {
		Password string `json:"passwort" validate:"required"`
	} `json:"passwort" validate:"required"`
}

type LoginResponse struct {
	SessionID string `json:"sessionID"`
}

type Session struct {
	ID string
	UserID string
	ExpiresAt time.Time
}

type GetUserRequest struct {
	LoginName string `json:"loginName" validate:"required"`
	Session string `json:"sitzung" validate:"required"`
	Location *Location `json:"standort,omitempty"` 
}

type GetUserResponse struct {
	UserList []User `json:"benutzerListe"`
}

type Location struct {
    Latitude  float64 `json:"breitengrad"`
    Longitude float64 `json:"laengengrad"`
}

type LogoutRequest struct {
    LoginName string `json:"loginName" validate:"required"`
    Session   string `json:"sitzung" validate:"required"`
}

type LogoutResponse struct {
    Result bool `json:"ergebnis"`
}