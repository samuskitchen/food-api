package response

type FoodResponse struct {
	ID          string `json:"id,omitempty"`
	UserID      string `json:"user_id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	FoodImage   string `json:"food_image,omitempty"`
}


// FoodResponse It is the response of the all food information
// swagger:response SwaggerAllFoodResponse
type SwaggerAllFoodResponse struct {
	// in: body
	Body []FoodResponse
}

// FoodResponse It is the response of the food information
// swagger:response SwaggerFoodResponse
type SwaggerFoodResponse struct {
	// in: body
	Body FoodResponse
}