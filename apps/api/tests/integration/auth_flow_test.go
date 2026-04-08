package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthRegisterRefreshAndMeFlow(t *testing.T) {
	engine := newTestEngine()

	registerPayload := []byte(`{
		"full_name":"Maria Silva",
		"email":"maria@example.com",
		"password":"SenhaSegura123",
		"phone":"71999999999",
		"device":{
			"device_name":"iPhone",
			"platform":"ios",
			"app_version":"1.0.0",
			"fingerprint":"device-fingerprint"
		}
	}`)

	registerResponse := httptest.NewRecorder()
	registerRequest := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewReader(registerPayload))
	registerRequest.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(registerResponse, registerRequest)

	if registerResponse.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, registerResponse.Code)
	}

	var registerBody map[string]any
	if err := json.Unmarshal(registerResponse.Body.Bytes(), &registerBody); err != nil {
		t.Fatalf("failed to decode register response: %v", err)
	}

	data, ok := registerBody["data"].(map[string]any)
	if !ok {
		t.Fatalf("missing data payload")
	}

	session, ok := data["session"].(map[string]any)
	if !ok {
		t.Fatalf("missing session payload")
	}

	accessToken, _ := session["access_token"].(string)
	refreshToken, _ := session["refresh_token"].(string)
	if accessToken == "" || refreshToken == "" {
		t.Fatalf("expected access and refresh tokens")
	}

	meResponse := httptest.NewRecorder()
	meRequest := httptest.NewRequest(http.MethodGet, "/v1/users/me", nil)
	meRequest.Header.Set("Authorization", "Bearer "+accessToken)
	engine.ServeHTTP(meResponse, meRequest)

	if meResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, meResponse.Code)
	}

	refreshPayload := []byte(`{"refresh_token":"` + refreshToken + `"}`)
	refreshResponse := httptest.NewRecorder()
	refreshRequest := httptest.NewRequest(http.MethodPost, "/v1/auth/refresh", bytes.NewReader(refreshPayload))
	refreshRequest.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(refreshResponse, refreshRequest)

	if refreshResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, refreshResponse.Code)
	}
}
