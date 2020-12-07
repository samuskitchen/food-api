package model

import (
	"strings"
	"time"
)

// Data of Food
// swagger:model
type Food struct {
	ID          string     `json:"id,omitempty"`
	// Required: true
	UserID      string     `json:"user_id,omitempty"`
	// Required: true
	Title       string     `json:"title,omitempty"`
	// Required: true
	Description string     `json:"description,omitempty"`
	FoodImage   string     `json:"food_image,omitempty"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `json:"-"`
}

// Information from food
// swagger:parameters getAllFood
type SwaggerAllFoodRequest struct {
	// type: apiKey
	// in: header
	// Required: true
	Authorization string
}


// Information from food
// swagger:parameters foodRequest
type SwaggerFoodRequest struct {
	// type: apiKey
	// in: header
	// Required: true
	Authorization string

	//in: body
	Body Food
}

// Information from food for update
// swagger:parameters foodUpdateRequest
type SwaggerFoodUpdateRequest struct {
	// type: apiKey
	// in: header
	// Required: true
	Authorization string

	// in: path
	// Required: true
	ID string

	// in: body
	Body struct{
		// Required: true
		Title       string     `json:"title,omitempty"`
		// Required: true
		Description string     `json:"description,omitempty"`
		FoodImage   string     `json:"food_image,omitempty"`
	}
}

// swagger:parameters idFoodPath
type SwaggerFoodPathId struct {
	// type: apiKey
	// in: header
	// Required: true
	Authorization string

	// in: path
	// Required: true
	ID string
}

// swagger:parameters idFoodDeletePath
type SwaggerFoodDeletePathId struct {
	// type: apiKey
	// in: header
	// Required: true
	Authorization string

	// in: path
	// Required: true
	ID string
}

// swagger:parameters idFoodByUserPath
type SwaggerFoodPathUser struct {
	// type: apiKey
	// in: header
	// Required: true
	Authorization string

	// in: path
	// Required: true
	UserID string
}

func (f *Food) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)

	switch strings.ToLower(action) {
	case "update":
		if f.Title == "" || f.Title == "null" {
			errorMessages["title_required"] = "title is required"
		}

		if f.Description == "" || f.Description == "null" {
			errorMessages["desc_required"] = "description is required"
		}

	default:
		if f.Title == "" || f.Title == "null" {
			errorMessages["title_required"] = "title is required"
		}

		if f.Description == "" || f.Description == "null" {
			errorMessages["desc_required"] = "description is required"
		}
	}
	return errorMessages
}
