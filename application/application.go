package application

import (
	"encoding/json"
	"net/http"
	"strconv"

	config "github.com/ArteShow/Calculator/pkg/Config"
	database "github.com/ArteShow/Calculator/pkg/Database"
	jwt "github.com/ArteShow/Calculator/pkg/JWT"
)

type User struct{
	Username string `json:"username"`
	Password string `json:"password"`
	UserId int `json:"userId"`
}

type Login struct{
	Login string `json:"login"`
	Password string `json:"password"`
}

func SaveRegUser(w http.ResponseWriter, r *http.Request){
	// Parse the request body to get the user data
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	UserMap := map[string]interface{}{
		"username": user.Username,
		"password": user.Password,
		"userId": user.UserId,
	}
	db, err := database.OpenDatabase("/db/database.db")
	if err != nil {
		http.Error(w, "Failed to open database", http.StatusInternalServerError)
		return
	}
	// Save the user data to the database
	err = database.InsertData(db, "users", UserMap)
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func LogMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := database.GetValueByField(config.GetDatabasePath(), "userId")
		if err != nil {
			http.Error(w, "Failed to get userId", http.StatusInternalServerError)
			return
		}
		intId, errerr := strconv.Atoi(userId)
		if errerr != nil {
			http.Error(w, "Failed to convert userId to int", http.StatusInternalServerError)
			return
		}
		tokenString, err := jwt.CreateJWT(intId, "user", jwt.GetJWTKey())
		if err != nil {
			http.Error(w, "Failed to create JWT", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
		next.ServeHTTP(w, r)
	})
}

func StartApplicationServer(){
	http.HandleFunc("/api/v1/register", SaveRegUser)
	http.HandleFunc("/api/v1/login", SaveRegUser)
}