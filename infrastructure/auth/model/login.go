package model

// Data of DataLogin
// swagger:model
type DataLogin struct {
	ID           string `json:"id,omitempty"`
	Names        string `json:"names,omitempty"`
	LastNames    string `json:"last_names,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// Information from user
// swagger:parameters userLoginRequest
type SwaggerUserLoginRequest struct {
	// in: body
	Body struct {
		// Required: true
		Email    string `json:"email,omitempty"`
		// Required: true
		Password string `json:"password,omitempty"`
	}
}

// DataLogin It is the response of the login information.
// swagger:response SwaggerDataLogin
type SwaggerDataLogin struct {
	//in: body
	Body DataLogin
}

// swagger:parameters authorization
type SwaggerAuthorization struct {
	// type: apiKey
	// in: header
	// Required: true
	Authorization string
}

// Information from login for refresh token
// swagger:parameters refreshRequest
type SwaggerRefreshRequest struct {
	// in: body
	Body struct {
		// Required: true
		AccessToken  string `json:"access_token,omitempty"`
		// Required: true
		RefreshToken string `json:"refresh_token,omitempty"`
	}
}
