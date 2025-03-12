package internal

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	calculate "github.com/ArteShow/Calculation_Service/pkg/Calculation"
)

var (
	id                   int64
	Calculations         = map[string]map[string]Expressions{"expression": {}}
	GenerateIdBool       bool
	expression           string
	ExpressionByID       = map[int]string{}
	TIME_ADDITION_MS     time.Duration
	TIME_SUBTRACTION_MS  time.Duration
	TIME_MULTIPLICATIONS_MS time.Duration
	TIME_DIVISIONS_MS    time.Duration
)

type Expression struct {
	Expression string `json:"expression"`
}

type Id struct {
	Id int `json:"id"`
}

type Expressions struct {
	ID     int     `json:"id"`
	Status int     `json:"status"`
	Result float64 `json:"result"`
	Error  error
}

// SendExpressionsList sends the entire list of expressions as a response
func SendExpressionsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(Calculations)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error creating the response:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// GenerateID generates a new unique ID for an expression and adds it to the list
func GenerateID(w http.ResponseWriter, r *http.Request) {
	if GenerateIdBool {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Println("⚠️ Request ignored: ID generation already in progress")
		http.Error(w, "ID generation is already in progress", http.StatusConflict)
		return
	}

	GenerateIdBool = true
	defer func() { GenerateIdBool = false }()

	log.Println("🔢 Generating new ID...")

	newID := atomic.AddInt64(&id, 1)
	Calculations["expression"][expression] = Expressions{
		ID:     int(newID),
		Status: 0,
		Result: 0,
		Error:  nil,
	}
	ExpressionByID[int(newID)] = expression
	log.Println("✅ New ID generated:", newID)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Id{Id: int(newID)})
	Calculate(expression)
}

// SendExpressionById retrieves and sends the expression by its ID
func SendExpressionById(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/internal/expression/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		w.WriteHeader(http.StatusNotFound)
		log.Println("❌ Error: Invalid URL")
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("❌ Error: ID is not a number:", parts[0])
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Search for the expression based on the ID
	exprStr, found := ExpressionByID[id]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		log.Println("❌ Error: Expression not found")
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	expression, found := Calculations["expression"][exprStr]

	// Send the JSON response with the expression
	response, err := json.Marshal(expression)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error creating the response:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	w.WriteHeader(http.StatusOK)
}

// StoreExpression stores the expression sent in the request
func StoreExpression(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error reading the expression:", err)
		http.Error(w, "Error reading the expression", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var expr Expression
	err = json.Unmarshal(body, &expr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error unmarshaling:", err)
		log.Println("📜 Received expression:", string(body))
		http.Error(w, "Error unmarshaling", http.StatusBadRequest)
		return
	}

	log.Println("✅ Expression received:", expr.Expression)

	w.WriteHeader(http.StatusCreated)
	expression = expr.Expression
}

// Calculate performs the calculation for a given expression
func Calculate(expression string) {
	//////////////////////
	TIME_ADDITION_MS = 5 * time.Millisecond
	TIME_SUBTRACTION_MS = 1 * time.Millisecond
	TIME_MULTIPLICATIONS_MS = 5 * time.Millisecond
	TIME_MULTIPLICATIONS_MS = 5 * time.Millisecond
	//////////////////////

	var wg sync.WaitGroup
	var mu sync.Mutex
	var statusCode int
	var finalResult float64
	var finalError error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	parts := strings.Fields(expression)
	wg.Add(len(parts))
	for _, expr := range parts {
		go func(expr string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				log.Println("🛑 Timeout: Calculation aborted")
				return
			default:
				result, err, code := calculate.Calc(expr)
				mu.Lock()
				if err != nil {
					log.Println("❌ Error in calculation:", err)
					finalError = err
				} else {
					log.Printf("✅ Calculation: %s = %f, StatusCode: %d\n", expr, result, code)
					finalResult += result
					statusCode = code
				}
				mu.Unlock()
			}
		}(expr)
	}

	wg.Wait()

	log.Println("✅ All calculations completed!")
	mu.Lock()
	Calculations["expression"][expression] = Expressions{
		ID:     Calculations["expression"][expression].ID,
		Status: statusCode,
		Result: finalResult,
		Error:  finalError,
	}
	mu.Unlock()

}

// RunServerAgent starts the internal server to handle requests
func RunServerAgent() {
	log.Println("🚀 Internal server started on port 8083")
	http.HandleFunc("/internal/task", StoreExpression)
	http.HandleFunc("/internal/expression/", SendExpressionById)
	http.HandleFunc("/internal/expression", GenerateID)
	http.HandleFunc("/internal/expression/list", SendExpressionsList)  // Correct: list route for expressions

	log.Fatal(http.ListenAndServe(":8083", nil))
}
