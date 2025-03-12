package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Test for Calculate
func TestCalculate(t *testing.T) {
	// ID generation
	req, err := http.NewRequest("POST", "/internal/expression", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GenerateID)
	handler.ServeHTTP(rr, req)

	// Test request for storing the expression
	expressionData := `{"expression": "2 + 3 * 2"}`
	req2, err := http.NewRequest("POST", "/internal/task", bytes.NewBuffer([]byte(expressionData)))
	if err != nil {
		t.Fatal(err)
	}
	req2.Header.Set("Content-Type", "application/json")
	rr2 := httptest.NewRecorder()
	handler2 := http.HandlerFunc(StoreExpression)
	handler2.ServeHTTP(rr2, req2)

	// Testing the calculation
	// We give the calculation some time (via the `sync.WaitGroup` in the `Calculate` function)
	time.Sleep(1 * time.Second)

	// Fetching the calculation results
	req3, err := http.NewRequest("GET", "/internal/expression/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr3 := httptest.NewRecorder()
	handler3 := http.HandlerFunc(SendExpressionById)
	handler3.ServeHTTP(rr3, req3)

	if status := rr3.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, status)
	}

	// Checking the calculation results
	var expr Expressions
	err = json.NewDecoder(rr3.Body).Decode(&expr)
	if err != nil {
		t.Fatal("Error decoding response:", err)
	}

	if expr.Result != 8 {
		t.Errorf("Expected result 8, but got %f", expr.Result)
	}

	// Optional: Checking the status code and the ID
	if expr.Status != 200 {
		t.Errorf("Expected status 200, but got %d", expr.Status)
	}
}

// Test for SendExpressionsList
func TestSendExpressionsList(t *testing.T) {
	req, err := http.NewRequest("GET", "/internal/expression/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(SendExpressionsList)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, status)
	}
}

// Test for GenerateID
func TestGenerateID(t *testing.T) {
	// Initially setting that ID generation is not running
	GenerateIdBool = false

	// Test request for ID generation
	req, err := http.NewRequest("POST", "/internal/expression", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(GenerateID)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, but got %d", http.StatusCreated, status)
	}

	// Checking if the ID was returned
	var idResponse Id
	err = json.NewDecoder(rr.Body).Decode(&idResponse)
	if err != nil {
		t.Fatal("Error decoding response:", err)
	}

	if idResponse.Id <= 0 {
		t.Errorf("Expected a valid ID, but got %d", idResponse.Id)
	}
}

// Test for SendExpressionById
func TestSendExpressionById(t *testing.T) {
	// Pre-simulating ID generation
	expression = "2+2"
	GenerateIdBool = false
	req, err := http.NewRequest("POST", "/internal/expression", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GenerateID)
	handler.ServeHTTP(rr, req)

	// Now testing the SendExpressionById function
	req2, err := http.NewRequest("GET", "/internal/expression/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr2 := httptest.NewRecorder()

	handler2 := http.HandlerFunc(SendExpressionById)
	handler2.ServeHTTP(rr2, req2)

	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, status)
	}
}

// Test for StoreExpression
func TestStoreExpression(t *testing.T) {
	// Example expression to be stored
	expressionData := `{"expression": "2+3"}`

	req, err := http.NewRequest("POST", "/internal/task", bytes.NewBuffer([]byte(expressionData)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(StoreExpression)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, but got %d", http.StatusCreated, status)
	}

	// Ensuring that the expression was stored
	if expression != "2+3" {
		t.Errorf("Expected expression '2+3', but got %s", expression)
	}
}
