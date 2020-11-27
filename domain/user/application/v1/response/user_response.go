package response

type UserResponse struct {
	ID           string     `json:"id,omitempty"`
	Names        string     `json:"names,omitempty"`
	LastNames    string     `json:"last_names,omitempty"`
	Email        string     `json:"email,omitempty"`
}