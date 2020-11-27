package response

type FoodResponse struct {
	ID          string `json:"id,omitempty"`
	UserID      string `json:"user_id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	FoodImage   string `json:"food_image,omitempty"`
}