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

type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

type Standort struct {
	Location Location `json:"standort"`
}

type Location struct {
	Latitude  float64 `json:"breitengrad"`
	Longitude float64 `json:"laengengrad"`
}
