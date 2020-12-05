package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"food-api/domain/user/application/v1/response"
	responseUser "food-api/domain/user/application/v1/response"
	"food-api/domain/user/domain/model"
	"food-api/test/integration/seed"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// dataUser is data for test
func dataUser() []model.User {
	now := time.Now().Truncate(time.Second).Truncate(time.Millisecond).Truncate(time.Microsecond)

	return []model.User{
		{
			ID:        uuid.New().String(),
			Names:     "Daniel",
			LastNames: "De La Pava Suarez",
			Email:     "daniel.delapava@jikkosoft.com",
			Password:  "123456",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        uuid.New().String(),
			Names:     "Rebecca",
			LastNames: "Romero",
			Email:     "rebecca.romero@jikkosoft.com",
			Password:  "123456",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func TestIntegration_GetAllUser(t *testing.T) {
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

		req, err := http.NewRequest(http.MethodGet, "/api/v1/users/", nil)
		if err != nil {
			tt.Errorf("error creating request: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+dataLogin.AccessToken)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if e, a := http.StatusNotFound, w.Code; e != a {
			tt.Errorf("expected status code: %v, got status code: %v", e, a)
		}

		var responseData response.UserResponse
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

	t.Run("Ok (database has been seeded)", func(tt *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/users/", nil)
		if err != nil {
			tt.Errorf("error creating request: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+dataLogin.AccessToken)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if e, a := http.StatusOK, w.Code; e != a {
			tt.Errorf("expected status code: %v, got status code: %v", e, a)
		}

		var users []responseUser.UserResponse
		if err := json.NewDecoder(w.Body).Decode(&users); err != nil {
			tt.Errorf("error decoding response body: %v", err)
		}

		if d := cmp.Diff(expectedUsers[0].ID, users[0].ID); d != "" {
			tt.Errorf("unexpected difference in response body:\n%v", d)
		}
	})

	testIntegrationLogout(t, dataLogin)
}

func TestIntegration_GetOneHandler(t *testing.T) {
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
	dataResponse := responseUser.UserResponse{
		ID: expectedUsers[0].ID,
		Names: expectedUsers[0].Names,
		LastNames: expectedUsers[0].LastNames,
		Email: expectedUsers[0].Email,
	}

	tests := []struct {
		Name         string
		UserID       string
		ExpectedBody responseUser.UserResponse
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
			ExpectedBody: responseUser.UserResponse{},
			ExpectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s", test.UserID), nil)
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
				var userResponse responseUser.UserResponse

				if err := json.NewDecoder(w.Body).Decode(&userResponse); err != nil {
					t.Errorf("error decoding userResponse body: %v", err)
				}

				if e, a := test.ExpectedBody.ID, userResponse.ID; e != a {
					t.Errorf("expected user ID: %v, got user ID: %v", e, a)
				}
			}
		}

		t.Run(test.Name, fn)
	}

	testIntegrationLogout(t, dataLogin)
}

func TestIntegration_CreateHandler(t *testing.T) {
	defer func() {
		if err := seed.Truncate(dataConnection.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	}()

	tests := []struct {
		Name         string
		RequestBody  model.User
		ExpectedCode int
	}{
		{
			Name:         "Create User Successful",
			RequestBody:  dataUser()[0],
			ExpectedCode: http.StatusCreated,
		},
		{
			Name:         "Break Unique UserName Constraint",
			RequestBody:  dataUser()[0],
			ExpectedCode: http.StatusConflict,
		},
		{
			Name:         "No Data User",
			RequestBody:  model.User{},
			ExpectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			var b bytes.Buffer
			if err := json.NewEncoder(&b).Encode(test.RequestBody); err != nil {
				t.Errorf("error encoding request body: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, "/api/v1/users/", &b)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			defer func() {
				if err := req.Body.Close(); err != nil {
					t.Errorf("error encountered closing request body: %v", err)
				}
			}()

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedCode != http.StatusConflict {
				var userResponse responseUser.UserResponse

				if err := json.NewDecoder(w.Body).Decode(&userResponse); err != nil {
					t.Errorf("error decoding userResponse body: %v", err)
				}

				if e, a := test.RequestBody.Names, userResponse.Names; e != a {
					t.Errorf("expected user UserName: %v, got user UserName: %v", e, a)
				}
			}
		}

		t.Run(test.Name, fn)
	}
}

func TestIntegration_UpdateHandler(t *testing.T) {

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
	dataResponse := responseUser.UserResponse{
		ID: expectedUsers[0].ID,
		Names: expectedUsers[0].Names,
		LastNames: expectedUsers[0].LastNames,
		Email: expectedUsers[0].Email,
	}

	tests := []struct {
		Name         string
		UserID       string
		RequestBody  model.User
		ExpectedBody responseUser.UserResponse
		ExpectedCode int
	}{
		{
			Name:         "Update User Successful",
			UserID:       expectedUsers[0].ID,
			RequestBody:  expectedUsers[0],
			ExpectedBody: dataResponse,
			ExpectedCode: http.StatusOK,
		},
		{
			Name:         "Break Unique UserName Constraint",
			UserID:       expectedUsers[1].ID,
			RequestBody:  expectedUsers[0],
			ExpectedBody: dataResponse,
			ExpectedCode: http.StatusConflict,
		},
		{
			Name:         "No Data User",
			UserID:       expectedUsers[0].ID,
			RequestBody:  model.User{},
			ExpectedBody: responseUser.UserResponse{},
			ExpectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			var b bytes.Buffer
			if err := json.NewEncoder(&b).Encode(test.RequestBody); err != nil {
				t.Errorf("error encoding request body: %v", err)
			}

			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s", test.UserID), &b)
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
				var userResponse responseUser.UserResponse

				if err := json.NewDecoder(w.Body).Decode(&userResponse); err != nil {
					t.Errorf("error decoding userResponse body: %v", err)
				}

				if e, a := test.RequestBody.Names, userResponse.Names; e != a {
					t.Errorf("expected user UserName: %v, got user UserName: %v", e, a)
				}
			}
		}

		t.Run(test.Name, fn)
	}

	testIntegrationLogout(t, dataLogin)
}
