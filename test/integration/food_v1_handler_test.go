package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"food-api/domain/food/application/v1/response"
	responseFood "food-api/domain/food/application/v1/response"
	"food-api/domain/food/domain/model"
	"food-api/test/integration/seed"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// dataFood is data for test
func dataFood() []model.Food {
	now := time.Now()
	userId := uuid.New().String()

	return []model.Food{
		{
			ID:          uuid.New().String(),
			UserID:      userId,
			Title:       "Title",
			Description: "Description",
			FoodImage:   "/profile-photos/food_api/309-3092053_gopher-link-transparent-cartoons-gopher-link.png",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          uuid.New().String(),
			UserID:      userId,
			Title:       "Title",
			Description: "Description",
			FoodImage:   "/profile-photos/food_api/309-3092053_gopher-link-transparent-cartoons-gopher-link.png",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
}

// dataFoodResponse is data for test
func dataFoodResponse() []responseFood.FoodResponse {
	userId := uuid.New().String()

	return []responseFood.FoodResponse{
		{
			ID:          uuid.New().String(),
			UserID:      userId,
			Title:       "Title",
			Description: "Description",
			FoodImage:   "/profile-photos/food_api/309-3092053_gopher-link-transparent-cartoons-gopher-link.png",
		},
		{
			ID:          uuid.New().String(),
			UserID:      userId,
			Title:       "Title",
			Description: "Description",
			FoodImage:   "/profile-photos/food_api/309-3092053_gopher-link-transparent-cartoons-gopher-link.png",
		},
	}
}

func TestFoodRouter_GetAllFood(t *testing.T) {

	defer func() {
		if err := seed.Truncate(dataConnection.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	}()

	expectedUsers, err := seed.UsersSeed(dataConnection.DB)
	if err != nil {
		t.Fatalf("error seeding users: %v", err)
	}

	dataLogin := testIntegrationLogin(t, expectedUsers)

	t.Run("No Content (no seed data)", func(tt *testing.T) {
		if err := seed.Truncate(dataConnection.DB); err != nil {
			tt.Errorf("error truncating test database tables: %v", err)
		}

		req, err := http.NewRequest(http.MethodGet, "/api/v1/foods/", nil)
		if err != nil {
			tt.Errorf("error creating request: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+dataLogin.AccessToken)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if e, a := http.StatusNotFound, w.Code; e != a {
			tt.Errorf("expected status code: %v, got status code: %v", e, a)
		}

		var responseData response.FoodResponse
		if err := json.Unmarshal([]byte(w.Body.String()), &responseData); err != nil {
			tt.Errorf("error decoding response body: %v", err)
		}

		if responseData.ID != "" {
			tt.Errorf("expected no data to be returned, got %v data", responseData)
		}
	})

	expectedUsers, err = seed.UsersSeed(dataConnection.DB)
	if err != nil {
		t.Fatalf("error seeding users: %v", err)
	}

	expectedFoods, err := seed.FoodsSeed(dataConnection.DB, expectedUsers)
	if err != nil {
		t.Fatalf("error seeding foods: %v", err)
	}

	t.Run("Ok (database has been seeded)", func(tt *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/foods/", nil)
		if err != nil {
			tt.Errorf("error creating request: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+dataLogin.AccessToken)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if e, a := http.StatusOK, w.Code; e != a {
			tt.Errorf("expected status code: %v, got status code: %v", e, a)
		}

		var foodResponses []responseFood.FoodResponse
		if err := json.NewDecoder(w.Body).Decode(&foodResponses); err != nil {
			tt.Errorf("error decoding response body: %v", err)
		}

		if d := cmp.Diff(expectedFoods[0].ID, foodResponses[0].ID); d != "" {
			tt.Errorf("unexpected difference in response body:\n%v", d)
		}
	})

	testIntegrationLogout(t, dataLogin)
}

func TestFoodRouter_GetOneByUserHandler(t *testing.T) {

	defer func() {
		if err := seed.Truncate(dataConnection.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	}()

	expectedUsers, err := seed.UsersSeed(dataConnection.DB)
	if err != nil {
		t.Fatalf("error seeding users: %v", err)
	}

	expectedFoods, err := seed.FoodsSeed(dataConnection.DB, expectedUsers)
	if err != nil {
		t.Fatalf("error seeding foods: %v", err)
	}

	dataLogin := testIntegrationLogin(t, expectedUsers)

	dataResponse := responseFood.FoodResponse{
		ID:          expectedFoods[0].ID,
		UserID:      expectedUsers[0].ID,
		Title:       expectedFoods[0].Title,
		Description: expectedFoods[0].Description,
		FoodImage:   expectedFoods[0].FoodImage,
	}

	tests := []struct {
		Name         string
		UserID       string
		ExpectedBody responseFood.FoodResponse
		ExpectedCode int
	}{
		{
			Name:         "Get One User Successful",
			UserID:       expectedUsers[0].ID,
			ExpectedBody: dataResponse,
			ExpectedCode: http.StatusOK,
		},
		{
			Name:         "User Not Found",
			UserID:       "0",
			ExpectedBody: responseFood.FoodResponse{},
			ExpectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/foods/user/%s", test.UserID), nil)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			req.Header.Set("Authorization", "Bearer "+dataLogin.AccessToken)
			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedCode != http.StatusNotFound {
				var foodResponse responseFood.FoodResponse

				if err := json.NewDecoder(w.Body).Decode(&foodResponse); err != nil {
					t.Errorf("error decoding foodResponse body: %v", err)
				}

				if e, a := test.ExpectedBody.ID, foodResponse.ID; e != a {
					t.Errorf("expected user ID: %v, got user ID: %v", e, a)
				}
			}
		}

		t.Run(test.Name, fn)
	}

	testIntegrationLogout(t, dataLogin)
}

func TestFoodRouter_GetOneHandler(t *testing.T) {

	defer func() {
		if err := seed.Truncate(dataConnection.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	}()

	expectedUsers, err := seed.UsersSeed(dataConnection.DB)
	if err != nil {
		t.Fatalf("error seeding users: %v", err)
	}

	expectedFoods, err := seed.FoodsSeed(dataConnection.DB, expectedUsers)
	if err != nil {
		t.Fatalf("error seeding foods: %v", err)
	}

	dataLogin := testIntegrationLogin(t, expectedUsers)

	dataResponse := responseFood.FoodResponse{
		ID:          expectedFoods[0].ID,
		UserID:      expectedUsers[0].ID,
		Title:       expectedFoods[0].Title,
		Description: expectedFoods[0].Description,
		FoodImage:   expectedFoods[0].FoodImage,
	}

	tests := []struct {
		Name         string
		ID           string
		ExpectedBody responseFood.FoodResponse
		ExpectedCode int
	}{
		{
			Name:         "Get One User Successful",
			ID:           expectedFoods[0].ID,
			ExpectedBody: dataResponse,
			ExpectedCode: http.StatusOK,
		},
		{
			Name:         "User Not Found",
			ID:           "0",
			ExpectedBody: responseFood.FoodResponse{},
			ExpectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/foods/%s", test.ID), nil)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			req.Header.Set("Authorization", "Bearer "+dataLogin.AccessToken)
			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedCode != http.StatusNotFound {
				var foodResponse responseFood.FoodResponse

				if err := json.NewDecoder(w.Body).Decode(&foodResponse); err != nil {
					t.Errorf("error decoding foodResponse body: %v", err)
				}

				if e, a := test.ExpectedBody.ID, foodResponse.ID; e != a {
					t.Errorf("expected user ID: %v, got user ID: %v", e, a)
				}
			}
		}

		t.Run(test.Name, fn)
	}

	testIntegrationLogout(t, dataLogin)

}

func TestFoodRouter_CreateHandler(t *testing.T) {

	defer func() {
		if err := seed.Truncate(dataConnection.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	}()

	expectedUsers, err := seed.UsersSeed(dataConnection.DB)
	if err != nil {
		t.Fatalf("error seeding users: %v", err)
	}

	dataLogin := testIntegrationLogin(t, expectedUsers)

	dataTest := model.Food{
		ID:          dataFood()[0].ID,
		UserID:      expectedUsers[0].ID,
		Title:       dataFood()[0].Title,
		Description: dataFood()[0].Description,
		FoodImage:   dataFood()[0].FoodImage,
	}

	tests := []struct {
		Name         string
		RequestBody  model.Food
		ExpectedCode int
	}{
		{
			Name:         "Create Food Successful",
			RequestBody:  dataTest,
			ExpectedCode: http.StatusCreated,
		},
		{
			Name:         "Break Unique UserName Constraint",
			RequestBody:  dataTest,
			ExpectedCode: http.StatusConflict,
		},
		{
			Name:         "No Data Food",
			RequestBody:  model.Food{},
			ExpectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			var b bytes.Buffer
			if err := json.NewEncoder(&b).Encode(test.RequestBody); err != nil {
				t.Errorf("error encoding request body: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, "/api/v1/foods/", &b)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			defer func() {
				if err := req.Body.Close(); err != nil {
					t.Errorf("error encountered closing request body: %v", err)
				}
			}()

			req.Header.Set("Authorization", "Bearer "+dataLogin.AccessToken)
			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedCode != http.StatusConflict {
				var foodResponse responseFood.FoodResponse

				if err := json.NewDecoder(w.Body).Decode(&foodResponse); err != nil {
					t.Errorf("error decoding foodResponse body: %v", err)
				}

				if e, a := test.RequestBody.Title, foodResponse.Title; e != a {
					t.Errorf("expected user Title: %v, got user Title: %v", e, a)
				}
			}
		}

		t.Run(test.Name, fn)
	}

	testIntegrationLogout(t, dataLogin)
}

func TestFoodRouter_UpdateHandler(t *testing.T) {

	defer func() {
		if err := seed.Truncate(dataConnection.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	}()

	expectedUsers, err := seed.UsersSeed(dataConnection.DB)
	if err != nil {
		t.Fatalf("error seeding users: %v", err)
	}

	expectedFoods, err := seed.FoodsSeed(dataConnection.DB, expectedUsers)
	if err != nil {
		t.Fatalf("error seeding foods: %v", err)
	}

	dataLogin := testIntegrationLogin(t, expectedUsers)

	dataTest := model.Food{
		ID:          expectedFoods[0].ID,
		UserID:      expectedUsers[0].ID,
		Title:       expectedFoods[0].Title,
		Description: expectedFoods[0].Description,
		FoodImage:   expectedFoods[0].FoodImage,
	}

	dataResponse := responseFood.FoodResponse{
		ID:          expectedFoods[0].ID,
		UserID:      expectedUsers[0].ID,
		Title:       expectedFoods[0].Title,
		Description: expectedFoods[0].Description,
		FoodImage:   expectedFoods[0].FoodImage,
	}

	tests := []struct {
		Name         string
		ID           string
		RequestBody  model.Food
		ExpectedBody responseFood.FoodResponse
		ExpectedCode int
	}{
		{
			Name:         "Update Food Successful",
			ID:           expectedFoods[0].ID,
			RequestBody:  dataTest,
			ExpectedBody: dataResponse,
			ExpectedCode: http.StatusOK,
		},
		{
			Name:         "No Data User",
			ID:           expectedUsers[0].ID,
			RequestBody:  model.Food{},
			ExpectedBody: responseFood.FoodResponse{},
			ExpectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			var b bytes.Buffer
			if err := json.NewEncoder(&b).Encode(test.RequestBody); err != nil {
				t.Errorf("error encoding request body: %v", err)
			}

			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/foods/%s", test.ID), &b)
			if err != nil {
				t.Errorf("error updating request: %v", err)
			}

			defer func() {
				if err := req.Body.Close(); err != nil {
					t.Errorf("error encountered closing request body: %v", err)
				}
			}()

			req.Header.Set("Authorization", "Bearer "+dataLogin.AccessToken)
			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedCode != http.StatusConflict {
				var foodResponse responseFood.FoodResponse

				if err := json.NewDecoder(w.Body).Decode(&foodResponse); err != nil {
					t.Errorf("error decoding foodResponse body: %v", err)
				}

				if e, a := test.RequestBody.Title, foodResponse.Title; e != a {
					t.Errorf("expected user Title: %v, got user Title: %v", e, a)
				}
			}
		}

		t.Run(test.Name, fn)
	}

	testIntegrationLogout(t, dataLogin)
}

func TestFoodRouter_DeleteHandler(t *testing.T) {

	defer func() {
		if err := seed.Truncate(dataConnection.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	}()

	expectedUsers, err := seed.UsersSeed(dataConnection.DB)
	if err != nil {
		t.Fatalf("error seeding users: %v", err)
	}

	expectedFoods, err := seed.FoodsSeed(dataConnection.DB, expectedUsers)
	if err != nil {
		t.Fatalf("error seeding foods: %v", err)
	}

	dataLogin := testIntegrationLogin(t, expectedUsers)

	tests := []struct {
		Name         string
		ID           string
		ExpectedCode int
	}{
		{
			Name:         "Delete User Successful",
			ID:           expectedFoods[0].ID,
			ExpectedCode: http.StatusNoContent,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/foods/%s", test.ID), nil)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			req.Header.Set("Authorization", "Bearer "+dataLogin.AccessToken)
			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}
		}

		t.Run(test.Name, fn)
	}

	testIntegrationLogout(t, dataLogin)

}
