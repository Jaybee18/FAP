package models

import "time"

type User struct {
	LoginName string    `json:"loginName" validate:"required"`
	Password  Password  `json:"passwort" validate:"required"`
	FirstName string    `json:"vorname" validate:"required"`
	LastName  string    `json:"nachname" validate:"required"`
	Street    *string   `json:"strasse,omitempty"`
	ZipCode   *string   `json:"plz,omitempty"`
	City      *string   `json:"ort,omitempty"`
	Country   *string   `json:"land,omitempty"`
	Phone     *string   `json:"telefon,omitempty"`
	Email     *Email    `json:"email,omitempty"`
	Location  *Location `json:"-"`
}

type Password struct {
	Password string `json:"passwort" validate:"required,min=6"`
}

type Email struct {
	Address string `json:"adresse" validate:"email"`
}

type GenericResponse struct {
	Result  string `json:"ergebnis"`
	Message string `json:"meldung"`
}

type LoginRequest struct {
	LoginName string `json:"loginName" validate:"required"`
	Password  struct {
		Password string `json:"passwort" validate:"required"`
	} `json:"passwort" validate:"required"`
}

type LoginResponse struct {
	SessionID string `json:"sessionID"`
}

type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

type GetUserRequest struct {
	LoginName string    `json:"loginName" validate:"required"`
	Session   string    `json:"sitzung" validate:"required"`
	Location  *Location `json:"standort,omitempty"`
}

type GetUserResponseUser struct {
	FirstName string `json:"vorname"`
	LastName  string `json:"nachname"`
	LoginName string `json:"loginName"`
}

type GetUsersResponse struct {
	UserList []GetUserResponseUser `json:"benutzerListe"`
}

type Standort struct {
	Location Location `json:"standort"`
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

type SetStandortRequest struct {
	LoginName string   `json:"loginName"`
	SessionId string   `json:"sitzung"`
	Location  Location `json:"standort"`
}
