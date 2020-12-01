package model

type DataLogin struct {
	ID           string `json:"id,omitempty"`
	Names        string `json:"names,omitempty"`
	LastNames    string `json:"last_names,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}