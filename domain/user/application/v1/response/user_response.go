package response

type UserResponse struct {
	ID           string     `json:"id,omitempty"`
	Names        string     `json:"names,omitempty"`
	LastNames    string     `json:"last_names,omitempty"`
	Email        string     `json:"email,omitempty"`
}

// UserResponse It is the response of the all users information
// swagger:response SwaggerAllUserResponse
type SwaggerAllUserResponse struct {
	//in: body
	Body []UserResponse
}

// UserResponse It is the response of the user's information
// swagger:response SwaggerUserResponse
type SwaggerUserResponse struct {
	//in: body
	Body UserResponse
}

