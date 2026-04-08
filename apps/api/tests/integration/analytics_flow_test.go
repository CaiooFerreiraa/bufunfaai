package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAnalyticsOverviewAndGoalsFlow(t *testing.T) {
	engine := newTestEngine()
	accessToken := registerAnalyticsUser(t, engine)

	overviewRecorder := httptest.NewRecorder()
	overviewRequest := httptest.NewRequest(http.MethodGet, "/v1/analytics/overview", nil)
	overviewRequest.Header.Set("Authorization", "Bearer "+accessToken)
	engine.ServeHTTP(overviewRecorder, overviewRequest)

	if overviewRecorder.Code != http.StatusOK {
		t.Fatalf("expected overview status %d, got %d", http.StatusOK, overviewRecorder.Code)
	}

	var overviewBody map[string]any
	if err := json.Unmarshal(overviewRecorder.Body.Bytes(), &overviewBody); err != nil {
		t.Fatalf("failed to decode overview response: %v", err)
	}

	overview := overviewBody["data"].(map[string]any)["overview"].(map[string]any)
	if overview["score"] == nil {
		t.Fatalf("expected overview score")
	}

	createGoalPayload := []byte(`{
		"title":"Reserva de emergencia",
		"goal_type":"save_amount",
		"target_amount_cents":300000,
		"current_amount_cents":50000,
		"due_date":"2026-12-31"
	}`)

	createGoalRecorder := httptest.NewRecorder()
	createGoalRequest := httptest.NewRequest(http.MethodPost, "/v1/analytics/goals", bytes.NewReader(createGoalPayload))
	createGoalRequest.Header.Set("Content-Type", "application/json")
	createGoalRequest.Header.Set("Authorization", "Bearer "+accessToken)
	engine.ServeHTTP(createGoalRecorder, createGoalRequest)

	if createGoalRecorder.Code != http.StatusCreated {
		t.Fatalf("expected create goal status %d, got %d", http.StatusCreated, createGoalRecorder.Code)
	}

	var createGoalBody map[string]any
	if err := json.Unmarshal(createGoalRecorder.Body.Bytes(), &createGoalBody); err != nil {
		t.Fatalf("failed to decode create goal response: %v", err)
	}

	goal := createGoalBody["data"].(map[string]any)["goal"].(map[string]any)
	goalID := goal["id"].(string)

	updateGoalPayload := []byte(`{
		"current_amount_cents":90000,
		"status":"in_progress"
	}`)

	updateGoalRecorder := httptest.NewRecorder()
	updateGoalRequest := httptest.NewRequest(http.MethodPatch, "/v1/analytics/goals/"+goalID, bytes.NewReader(updateGoalPayload))
	updateGoalRequest.Header.Set("Content-Type", "application/json")
	updateGoalRequest.Header.Set("Authorization", "Bearer "+accessToken)
	engine.ServeHTTP(updateGoalRecorder, updateGoalRequest)

	if updateGoalRecorder.Code != http.StatusOK {
		t.Fatalf("expected update goal status %d, got %d", http.StatusOK, updateGoalRecorder.Code)
	}

	goalsRecorder := httptest.NewRecorder()
	goalsRequest := httptest.NewRequest(http.MethodGet, "/v1/analytics/goals", nil)
	goalsRequest.Header.Set("Authorization", "Bearer "+accessToken)
	engine.ServeHTTP(goalsRecorder, goalsRequest)

	if goalsRecorder.Code != http.StatusOK {
		t.Fatalf("expected goals status %d, got %d", http.StatusOK, goalsRecorder.Code)
	}
}

func registerAnalyticsUser(t *testing.T, engine http.Handler) string {
	t.Helper()

	registerPayload := []byte(`{
		"full_name":"Analytics User",
		"email":"analytics@example.com",
		"password":"SenhaSegura123",
		"phone":"71999999999"
	}`)

	registerRecorder := httptest.NewRecorder()
	registerRequest := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewReader(registerPayload))
	registerRequest.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(registerRecorder, registerRequest)

	if registerRecorder.Code != http.StatusCreated {
		t.Fatalf("expected register status %d, got %d", http.StatusCreated, registerRecorder.Code)
	}

	var registerBody map[string]any
	if err := json.Unmarshal(registerRecorder.Body.Bytes(), &registerBody); err != nil {
		t.Fatalf("failed to decode register response: %v", err)
	}

	registerData := registerBody["data"].(map[string]any)
	session := registerData["session"].(map[string]any)
	return session["access_token"].(string)
}
