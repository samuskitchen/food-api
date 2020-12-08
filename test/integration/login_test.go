package integration

import (
	"bytes"
	"encoding/json"
	"food-api/domain/user/domain/model"
	authModel "food-api/infrastructure/auth/model"
	"food-api/test/integration/seed"
	"net/http"
	"net/http/httptest"
	"testing"
)

var dataResponseLogin authModel.DataLogin

// testIntegrationLogin necessary method to perform integration tests and also performs the login test
func testIntegrationLogin(t *testing.T, expectedUser []model.User) authModel.DataLogin {

	tests := []struct {
		Name         string
		RequestBody  model.User
		ExpectedCode int
	}{
		{
			Name:         "Login User Successful",
			RequestBody:  expectedUser[0],
			ExpectedCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		func() {
			var b bytes.Buffer
			if err := json.NewEncoder(&b).Encode(test.RequestBody); err != nil {
				t.Errorf("error encoding request body: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, "/api/login", &b)
			if err != nil {
				t.Errorf("error login request: %v", err)
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

			var responseLogin authModel.DataLogin
			if test.ExpectedCode != http.StatusConflict {

				if err := json.NewDecoder(w.Body).Decode(&responseLogin); err != nil {
					t.Errorf("error decoding responseLogin body: %v", err)
				}

				if e, a := test.RequestBody.ID, responseLogin.ID; e != a {
					t.Errorf("expected user ID: %v, got user ID: %v", e, a)
				}
			}

			dataResponseLogin = authModel.DataLogin{
				ID:           responseLogin.ID,
				Names:        responseLogin.Names,
				LastNames:    responseLogin.LastNames,
				AccessToken:  responseLogin.AccessToken,
				RefreshToken: responseLogin.RefreshToken,
			}

		}()
	}

	return dataResponseLogin
}

// testIntegrationLogout method needed to perform integration tests and also perform logout test
func testIntegrationLogout(t *testing.T, dataLogin authModel.DataLogin) {

	tests := []struct {
		Name         string
		ExpectedCode int
	}{
		{
			Name:         "Login User Successful",
			ExpectedCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		func() {

			req, err := http.NewRequest(http.MethodPost, "/api/logout", nil)
			if err != nil {
				t.Errorf("error logout request: %v", err)
			}
			req.Header.Set("Authorization", "Bearer "+dataLogin.AccessToken)

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			var responseLogout string
			if test.ExpectedCode != http.StatusOK {

				if err := json.NewDecoder(w.Body).Decode(&responseLogout); err != nil {
					t.Errorf("error decoding responseLogout body: %v", err)
				}

				if responseLogout == "Successfully logged out" {
					t.Errorf("expected no data to be returned, got %v data", responseLogout)
				}
			}

		}()
	}
}

// TestIntegration_Refresh method needed to perform integration tests and also perform refresh test
func TestIntegration_Refresh(t *testing.T) {
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
	responseRefresh := authModel.DataLogin{
		ID:           dataLogin.ID,
		Names:        dataLogin.Names,
		LastNames:    dataLogin.LastNames,
		AccessToken:  dataLogin.AccessToken,
		RefreshToken: dataLogin.RefreshToken,
	}

	tests := []struct {
		Name         string
		RequestBody  authModel.DataLogin
		ExpectedBody authModel.DataLogin
		ExpectedCode int
	}{
		{
			Name:         "Refresh Token Successful",
			RequestBody:  dataLogin,
			ExpectedBody: responseRefresh,
			ExpectedCode: http.StatusOK,
		},
		{
			Name:         "No Data Refresh",
			RequestBody:  authModel.DataLogin{},
			ExpectedBody: authModel.DataLogin{},
			ExpectedCode: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		func() {
			var b bytes.Buffer
			if err := json.NewEncoder(&b).Encode(test.RequestBody); err != nil {
				t.Errorf("error encoding request body: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, "/api/refresh", &b)
			if err != nil {
				t.Errorf("error refresh request: %v", err)
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

			var responseLogin authModel.DataLogin
			if test.ExpectedCode != http.StatusOK {

				if err := json.NewDecoder(w.Body).Decode(&responseLogin); err != nil {
					t.Errorf("error decoding responseLogin body: %v", err)
				}

				if e, a := test.RequestBody.ID, responseLogin.ID; e != a {
					t.Errorf("expected user ID: %v, got user ID: %v", e, a)
				}
			}

		}()
	}
}
