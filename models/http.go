package models

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

type SetStandortRequest struct {
	LoginName string   `json:"loginName"`
	SessionId string   `json:"sitzung"`
	Location  Location `json:"standort"`
}

type GeoJSONResponse struct {
	Features []Feature `json:"features"`
}

type Feature struct {
	Geometry Geometry `json:"geometry"`
}

type Geometry struct {
	Coordinates []float64 `json:"coordinates"`
}

type LogoutRequest struct {
	LoginName string `json:"loginName" validate:"required"`
	Session   string `json:"sitzung" validate:"required"`
}

type LogoutResponse struct {
	Result bool `json:"ergebnis"`
}
