package application

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSaveRegUser(t *testing.T) {
	payload := `{"login": "testuser", "password": "testpass"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	SaveRegUser(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", res.StatusCode)
	}
}

func TestLoginUser(t *testing.T) {
	payload := `{"login": "testuser", "password": "testpass"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	LoginUser(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	var result map[string]string
	json.Unmarshal(body, &result)

	if _, ok := result["token"]; !ok {
		t.Errorf("expected token in response")
	}
}

func TestCalculate_InvalidToken(t *testing.T) {
	payload := `{"expression": "2+2"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(payload))
	req.Header.Set("Authorization", "Bearer invalidtoken")
	w := httptest.NewRecorder()

	Calculate(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 for invalid token, got %d", res.StatusCode)
	}
}

func TestGetExpressions_InvalidToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/expressions", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	w := httptest.NewRecorder()

	GetExpressions(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 for invalid token, got %d", res.StatusCode)
	}
}

func TestGetExpressionById_InvalidToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/expression/1", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	w := httptest.NewRecorder()

	GetExpressionById(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 for invalid token, got %d", res.StatusCode)
	}
}
