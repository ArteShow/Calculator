package application

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// GetExpressionById handles requests to get an expression by its ID
func GetExpressionById(w http.ResponseWriter, r *http.Request){
	log.Println("📩 Request at:", r.URL.Path)

	// Extract the ID from the URL
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/expression/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		w.WriteHeader(http.StatusNotFound)
		log.Println("❌ Error: Invalid URL")
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	// Convert the ID to an integer
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("❌ Error: ID is not a number:", parts[0])
		http.Error(w, "The Id is not a string", http.StatusBadRequest)
		return
	}

	// Fetch the expression from the internal server
	url := "http://localhost:8083/internal/expression/" + strconv.Itoa(id)
	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error retrieving the expression:", err)
		http.Error(w, "Error while connecting to server 8083", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error reading the response:", err)
		http.Error(w, "Error while reading the response body", http.StatusInternalServerError)
		return
	}

	// Return the expression to the client
	log.Println("✅ Answer from internal:", string(body))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// GetExpressionsList handles requests to get a list of all expressions
func GetExpressionsList(w http.ResponseWriter, r *http.Request){
	// Fetch the list of expressions from the internal server
	resp, err := http.Get("http://localhost:8083/internal/expression/list")
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error reading the response:", err)
		http.Error(w, "Error while retrieving expression list", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error retrieving the expression:", err)
		http.Error(w, "Error while retrieving expression list", http.StatusInternalServerError)
		return
	}

	// Return the list of expressions to the client
	log.Println("✅ Answer from internal:", string(body))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// SendExpression handles requests to send an expression to the internal server
func SendExpression(w http.ResponseWriter, r *http.Request) {
	// Read the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error reading the body:", err)
		http.Error(w, "Error while reading the request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Send the expression to the internal server
	log.Println("📡 Sending expression to internal:", string(body))
	_, err2 := http.Post("http://localhost:8083/internal/task", "application/json", bytes.NewBuffer(body))

	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error sending:", err)
		http.Error(w, "Error while sending the expression", http.StatusInternalServerError)
		return
	}

	log.Println("✅ Expression saved, now requesting an ID...")

	// Request a new ID from the internal server
	idResp, err := http.Post("http://localhost:8083/internal/expression", "application/json", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error requesting the ID:", err)
		http.Error(w, "Error while requesting the ID", http.StatusInternalServerError)
		return
	}
	defer idResp.Body.Close()

	// Read the ID response
	idBody, err := io.ReadAll(idResp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error reading the ID response:", err)
		http.Error(w, "Error while reading the ID response", http.StatusInternalServerError)
		return
	}

	// Debugging - Show the content of the ID response
	log.Printf("📜 Received ID response: %s\n", string(idBody))

	// Check if the ID response is empty
	if len(idBody) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("❌ Error: No ID received")
		http.Error(w, "Error: No ID received", http.StatusInternalServerError)
		return
	}

	// Return the ID to the client
	log.Println("✅ ID received:", string(idBody))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(idBody)
}

// RunServer starts the API server
func RunServer() {
	log.Println("🌍 API server started on port 8082")
	// Define the routes for the API
	http.HandleFunc("/api/v1/calculate", SendExpression)
	http.HandleFunc("/api/v1/expression/", GetExpressionById)
	http.HandleFunc("/api/v1/expressions", GetExpressionsList)

	// Start the server
	log.Fatal(http.ListenAndServe(":8082", nil))
}
