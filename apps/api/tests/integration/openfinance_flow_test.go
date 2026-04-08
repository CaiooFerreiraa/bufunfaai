package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOpenFinanceConsentAndSyncFlow(t *testing.T) {
	engine := newTestEngine()

	registerPayload := []byte(`{
		"full_name":"Open Finance User",
		"email":"openfinance@example.com",
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
	accessToken := session["access_token"].(string)

	institutionsRecorder := httptest.NewRecorder()
	institutionsRequest := httptest.NewRequest(http.MethodGet, "/v1/open-finance/institutions", nil)
	institutionsRequest.Header.Set("Authorization", "Bearer "+accessToken)
	engine.ServeHTTP(institutionsRecorder, institutionsRequest)

	if institutionsRecorder.Code != http.StatusOK {
		t.Fatalf("expected institutions status %d, got %d", http.StatusOK, institutionsRecorder.Code)
	}

	var institutionsBody map[string]any
	if err := json.Unmarshal(institutionsRecorder.Body.Bytes(), &institutionsBody); err != nil {
		t.Fatalf("failed to decode institutions response: %v", err)
	}

	institutions := institutionsBody["data"].(map[string]any)["institutions"].([]any)
	institution := institutions[0].(map[string]any)
	institutionID := institution["id"].(string)

	createConsentPayload := []byte(`{
		"institution_id":"` + institutionID + `",
		"purpose":"Consolidacao financeira pessoal",
		"permissions":["ACCOUNTS_READ","BALANCES_READ","TRANSACTIONS_READ"],
		"redirect_uri":"https://app.bufunfa.ai/open-finance/callback"
	}`)

	createConsentRecorder := httptest.NewRecorder()
	createConsentRequest := httptest.NewRequest(http.MethodPost, "/v1/open-finance/consents", bytes.NewReader(createConsentPayload))
	createConsentRequest.Header.Set("Content-Type", "application/json")
	createConsentRequest.Header.Set("Authorization", "Bearer "+accessToken)
	engine.ServeHTTP(createConsentRecorder, createConsentRequest)

	if createConsentRecorder.Code != http.StatusCreated {
		t.Fatalf("expected consent status %d, got %d", http.StatusCreated, createConsentRecorder.Code)
	}

	var createConsentBody map[string]any
	if err := json.Unmarshal(createConsentRecorder.Body.Bytes(), &createConsentBody); err != nil {
		t.Fatalf("failed to decode consent response: %v", err)
	}

	consent := createConsentBody["data"].(map[string]any)["consent"].(map[string]any)
	consentID := consent["id"].(string)

	authorizeRecorder := httptest.NewRecorder()
	authorizeRequest := httptest.NewRequest(http.MethodPost, "/v1/open-finance/consents/"+consentID+"/authorize", nil)
	authorizeRequest.Header.Set("Authorization", "Bearer "+accessToken)
	engine.ServeHTTP(authorizeRecorder, authorizeRequest)

	if authorizeRecorder.Code != http.StatusOK {
		t.Fatalf("expected authorize status %d, got %d", http.StatusOK, authorizeRecorder.Code)
	}

	var authorizeBody map[string]any
	if err := json.Unmarshal(authorizeRecorder.Body.Bytes(), &authorizeBody); err != nil {
		t.Fatalf("failed to decode authorize response: %v", err)
	}

	authorizationURL := authorizeBody["data"].(map[string]any)["authorization_url"].(string)
	state := extractQueryValue(authorizationURL, "state")
	code := extractQueryValue(authorizationURL, "code")
	if state == "" || code == "" {
		t.Fatalf("expected state and code in authorization url")
	}

	callbackRecorder := httptest.NewRecorder()
	callbackRequest := httptest.NewRequest(http.MethodGet, "/v1/open-finance/callback?state="+state+"&code="+code, nil)
	engine.ServeHTTP(callbackRecorder, callbackRequest)

	if callbackRecorder.Code != http.StatusOK {
		t.Fatalf("expected callback status %d, got %d", http.StatusOK, callbackRecorder.Code)
	}

	var callbackBody map[string]any
	if err := json.Unmarshal(callbackRecorder.Body.Bytes(), &callbackBody); err != nil {
		t.Fatalf("failed to decode callback response: %v", err)
	}

	connection := callbackBody["data"].(map[string]any)["connection"].(map[string]any)
	connectionID := connection["id"].(string)

	syncRecorder := httptest.NewRecorder()
	syncRequest := httptest.NewRequest(http.MethodPost, "/v1/open-finance/connections/"+connectionID+"/sync", nil)
	syncRequest.Header.Set("Authorization", "Bearer "+accessToken)
	engine.ServeHTTP(syncRecorder, syncRequest)

	if syncRecorder.Code != http.StatusOK {
		t.Fatalf("expected sync status %d, got %d", http.StatusOK, syncRecorder.Code)
	}

	syncStatusRecorder := httptest.NewRecorder()
	syncStatusRequest := httptest.NewRequest(http.MethodGet, "/v1/open-finance/connections/"+connectionID+"/sync-status", nil)
	syncStatusRequest.Header.Set("Authorization", "Bearer "+accessToken)
	engine.ServeHTTP(syncStatusRecorder, syncStatusRequest)

	if syncStatusRecorder.Code != http.StatusOK {
		t.Fatalf("expected sync status response %d, got %d", http.StatusOK, syncStatusRecorder.Code)
	}

	reconcileRecorder := httptest.NewRecorder()
	reconcileRequest := httptest.NewRequest(http.MethodPost, "/internal/open-finance/reconcile?limit=5", nil)
	reconcileRequest.Header.Set("X-Internal-Secret", "test-cron-secret")
	engine.ServeHTTP(reconcileRecorder, reconcileRequest)

	if reconcileRecorder.Code != http.StatusOK {
		t.Fatalf("expected reconcile status %d, got %d", http.StatusOK, reconcileRecorder.Code)
	}
}

func extractQueryValue(url string, key string) string {
	parts := strings.Split(url, "?")
	if len(parts) != 2 {
		return ""
	}

	queryParts := strings.Split(parts[1], "&")
	for _, queryPart := range queryParts {
		pair := strings.SplitN(queryPart, "=", 2)
		if len(pair) == 2 && pair[0] == key {
			return pair[1]
		}
	}

	return ""
}
