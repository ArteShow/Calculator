package application

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test for GetExpressionById
func TestGetExpressionById(t *testing.T) {
	// Create a new GET request for the URL
	req, err := http.NewRequest("GET", "/api/v1/expression/1", nil)
	if err != nil {
		t.Fatal(err) // If there is an error creating the request, stop the test
	}

	// Create a new ResponseRecorder to test the response
	rr := httptest.NewRecorder()

	// Set the handler and execute the request
	handler := http.HandlerFunc(GetExpressionById)
	handler.ServeHTTP(rr, req)

	// Check the status code of the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, status)
	}

	// Check if the Content-Type of the response is correct
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected Content-Type to be 'application/json', but got '%s'", contentType)
	}

	// Here you could further check the response body, depending on its content
}

// Test for GetExpressionsList
func TestGetExpressionsList(t *testing.T) {
	// Create a new GET request for the list of expressions
	req, err := http.NewRequest("GET", "/api/v1/expressions", nil)
	if err != nil {
		t.Fatal(err) // If there is an error creating the request, stop the test
	}

	// Create a new ResponseRecorder to test the response
	rr := httptest.NewRecorder()

	// Set the handler and execute the request
	handler := http.HandlerFunc(GetExpressionsList)
	handler.ServeHTTP(rr, req)

	// Check the status code of the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, status)
	}
}

// Test for SendExpression
func TestSendExpression(t *testing.T) {
	// Example body for the POST request
	body := `{"expression": "2+2"}`

	// Create a new POST request with the example body
	req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Fatal(err) // If there is an error creating the request, stop the test
	}
	req.Header.Set("Content-Type", "application/json") // Set the correct content type header

	// Create a new ResponseRecorder to test the response
	rr := httptest.NewRecorder()

	// Set the handler and execute the request
	handler := http.HandlerFunc(SendExpression)
	handler.ServeHTTP(rr, req)

	// Check the status code of the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, status)
	}
}
