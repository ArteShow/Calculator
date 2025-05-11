package application

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"

	config "github.com/ArteShow/Calculator/pkg/Config"
	database "github.com/ArteShow/Calculator/pkg/Database"
	MyJWT "github.com/ArteShow/Calculator/pkg/JWT"
	user "github.com/ArteShow/Calculator/proto"
	"google.golang.org/grpc"
)

type User struct {
	Username string `json:"login"`
	Password string `json:"password"`
	UserId   int    `json:"userId"`
}

type Calculations struct {
	Calculations []Calculation `json:"expression"`
}

type Calculation struct {
	Expression string `json:"expression"`
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
	log.Printf("User name: %s 😊", user.Username)

	// Get the next user ID
	UserId, err := database.GetMaxId(config.GetDatabasePath(), "users")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get max userId %v", err), http.StatusInternalServerError)
		return
	}
	UserId++ // Increment the max ID to get the new user ID

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

	log.Printf("User saved successfully: %v 🎉", UserMap)
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
	log.Printf("Login name: %s 🔑", login.Login)

	// Get the user information based on the login (username) from the database
	user, err := database.GetUserByUsername(config.GetDatabasePath(), login.Login)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Check if the password matches the stored password
	if user.Password != login.Password {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// Generate the JWT token using the user ID
	tokenString, err := MyJWT.CreateJWT(user.UserId, "user", MyJWT.GetJWTKey())
	if err != nil {
		http.Error(w, "Failed to create JWT", http.StatusInternalServerError)
		return
	}

	// Send the token back to the client in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	log.Printf("JWT token generated successfully for user ID: %d 🎉", user.UserId)
}

func GetUserIdFromToken(w http.ResponseWriter, r *http.Request, tokenstring string) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return 0, fmt.Errorf("missing Authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return 0, fmt.Errorf("invalid Authorization header format")
	}

	tokenString := parts[1]
	token, err := MyJWT.ParseJWT(tokenString, []byte(MyJWT.GetJWTKey()))
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusUnauthorized)
		return 0, fmt.Errorf("invalid claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "user_id not found", http.StatusUnauthorized)
		return 0, fmt.Errorf("user_id not found")
	}
	return int(userIDFloat), nil
}

func Calculate(w http.ResponseWriter, r *http.Request) {
	var calculation Calculation
	err := json.NewDecoder(r.Body).Decode(&calculation)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Get userId from token
	userID, err := GetUserIdFromToken(w, r, w.Header().Get("Authorization"))
	if err != nil {
		http.Error(w, "Failed to get userId from token", http.StatusUnauthorized)
		return
	}
	log.Printf("User ID from token: %d 💡", userID)

	// gRPC call to send the user data and calculation
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		http.Error(w, "Failed to connect to gRPC server", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := user.NewUserServiceClient(conn)

	// Prepare the request with actual userId and calculation data
	req := &user.UserDataRequest{
		UserId:   int32(userID),
		Calculation: &user.Calculation{
			Expression: calculation.Expression,
		},
	}

	// Set a timeout for the request to the gRPC server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Send the request
	res, err := client.SendUserData(ctx, req)
	if err != nil {
		http.Error(w, "Failed to send user data to gRPC server", http.StatusInternalServerError)
		return
	}

	// Log the response from the gRPC server
	log.Printf("Server says: %s 🗣️", res.Message)

	// Send the response back to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": res.Message})
}

func GetExpressions(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIdFromToken(w, r, w.Header().Get("Authorization"))
	if err != nil {
		http.Error(w, "Failed to get userId from token", http.StatusUnauthorized)
		return
	}
	log.Printf("User ID from token: %d 📊", userID)

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		http.Error(w, "Failed to connect to gRPC server", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := user.NewUserServiceClient(conn)

	// Send gRPC request
	res, err := client.GetUserCalculations(context.Background(), &user.UserIdRequest{
		UserId: int32(userID),
	})
	if err != nil {
		http.Error(w, "Failed to get calculations", http.StatusInternalServerError)
		return
	}

	// Convert gRPC response to your `Calculations` struct
	calcs := Calculations{}
	for _, c := range res.Calculations {
		calcs.Calculations = append(calcs.Calculations, Calculation{
			Expression: c.Expression,
		})
	}

	// Respond to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calcs)
}

func GetExpressionById(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIdFromToken(w, r, w.Header().Get("Authorization"))
	if err != nil {
		http.Error(w, "Failed to get userId from token", http.StatusUnauthorized)
		return
	}
	log.Printf("User ID from token: %d 📖", userID)

	expressionID := strings.TrimPrefix(r.URL.Path, "/api/v1/expression/")
	expressionIDInt, err := strconv.Atoi(expressionID)
	if err != nil {
		http.Error(w, "Invalid expression ID", http.StatusBadRequest)
		return
	}
	log.Printf("Expression ID: %d ✨", expressionIDInt)

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) // Replace with your server address
	if err != nil {
		log.Fatalf("❌ Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Create a new client instance for the UserService
	client := user.NewUserServiceClient(conn)

	// Create the request object with the userId and customId
	request := &user.UserDataRequest{
		UserId:   int32(userID),
		CustomId: int32(expressionIDInt),
	}

	// Send the request using the SendUserData method
	response, err := client.SendUserData(context.Background(), request)
	if err != nil {
		log.Fatalf("❌ Failed to send data: %v", err)
	}

	// Print the response message
	log.Printf("Response from server: %s 💬", response.GetMessage())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": response.GetMessage()})
}

func StartApplicationServer() {
	http.HandleFunc("/api/v1/register", SaveRegUser)
	http.HandleFunc("/api/v1/login", LoginUser)
	http.HandleFunc("/api/v1/calculate", Calculate)
	http.HandleFunc("/api/v1/expressions", GetExpressions)
	http.HandleFunc("/api/v1/expression/", GetExpressionById)
	log.Println("Server started at http://localhost:8082 🚀")
	http.ListenAndServe(":8082", nil)
}
