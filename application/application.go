package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	config "github.com/ArteShow/Calculator/pkg/Config"
	database "github.com/ArteShow/Calculator/pkg/Database"
	jwt "github.com/ArteShow/Calculator/pkg/JWT"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserId   int    `json:"userId"`
}

type Login struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func SaveRegUser(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the user data
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Get the next user ID
	UserId, err := database.GetMaxId(config.GetDatabasePath(), "users")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get max userId %v", err), http.StatusInternalServerError)
		return
	}
	if UserId == 1 {
		log.Println("UserId is 1")
	} else {
		UserId += 1 // Correct increment
	}

	// Prepare the user data to be saved
	UserMap := map[string]interface{}{
		"username": user.Username,
		"password": user.Password,
		"id":       UserId,
	}

	// Open the database and save user data
	db, err := database.OpenDatabase(config.GetDatabasePath())
	if err != nil {
		http.Error(w, "Failed to open database", http.StatusInternalServerError)
		return
	}
	err = database.InsertData(db, "users", UserMap)
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	log.Printf("User saved successfully: %v", UserMap)
	w.WriteHeader(http.StatusCreated) // Only write the status once here
	json.NewEncoder(w).Encode(UserMap)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming JSON request body
	var login Login
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Get the user ID based on the login (username) from the database
	userId, err := database.GetUserID(config.GetDatabasePath(), login.Login)
	if err != nil {
		http.Error(w, "Failed to get userId", http.StatusInternalServerError)
		return
	}

	// Convert the user ID to an integer
	intId, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, "Failed to convert userId to int", http.StatusInternalServerError)
		return
	}

	// Generate the JWT token using the user ID
	tokenString, err := jwt.CreateJWT(intId, "user", jwt.GetJWTKey())
	if err != nil {
		http.Error(w, "Failed to create JWT", http.StatusInternalServerError)
		return
	}

	// Send the token back to the client in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})

	//hear the grpc part
}

func StartApplicationServer() {
	http.HandleFunc("/api/v1/register", SaveRegUser)
	http.HandleFunc("/api/v1/login", LoginUser)

	http.ListenAndServe(":8082", nil)
}
