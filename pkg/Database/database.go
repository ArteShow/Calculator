package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserId   int    `json:"userId"`
}

func CreateDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Printf("❌ Failed to open DB at path %s: %v", path, err)
		return nil, err
	}
	log.Println("✅ Database connection created 😊")
	return db, nil
}

func CreateTable(db *sql.DB, tableName string, columns map[string]string) error {
	query := "CREATE TABLE IF NOT EXISTS " + tableName + " ("
	for col, typ := range columns {
		query += col + " " + typ + ", "
	}
	query = query[:len(query)-2] + ");"

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("❌ Failed to create table %s: %v", tableName, err)
		return err
	}
	log.Printf("✅ Table %s created (or already exists) 😎", tableName)
	return nil
}

func InsertData(db *sql.DB, tableName string, data map[string]interface{}) error {
	query := "INSERT OR IGNORE INTO " + tableName + " ("
	values := "("
	args := []interface{}{}

	for col, val := range data {
		query += col + ", "
		values += "?, "
		args = append(args, val)
	}

	query = query[:len(query)-2] + ") VALUES " + values[:len(values)-2] + ");"

	log.Printf("📝 Query: %s", query)
	log.Printf("🧾 Args: %v", args)

	_, err := db.Exec(query, args...)
	if err != nil {
		log.Printf("❌ Failed to insert data: %v", err)
		return fmt.Errorf("failed to insert data: %w", err)
	}

	log.Println("✅ Insert (or ignore) successful 🚀")
	return nil
}
func DeleteData(db *sql.DB, tableName string, condition string) error {
	query := "DELETE FROM " + tableName + " WHERE " + condition + ";"
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("❌ Failed to delete data from %s: %v", tableName, err)
		return err
	}
	log.Println("🗑️ Data deleted successfully!")
	return nil
}

func OpenDatabase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Printf("❌ Failed to open database: %v", err)
		return nil, err
	}
	log.Println("📂 Database opened!")
	return db, nil
}

func UpdateData(db *sql.DB, tableName string, data map[string]interface{}, condition string) error {
	query := "UPDATE " + tableName + " SET "
	args := []interface{}{}
	for col, val := range data {
		query += col + " = ?, "
		args = append(args, val)
	}
	query = query[:len(query)-2] + " WHERE " + condition + ";"

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("❌ Failed to prepare update: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		log.Printf("❌ Failed to update data: %v", err)
		return err
	}
	log.Println("✏️ Data updated successfully!")
	return nil
}

func GetUserID(path, username string) (string, error) {
	db, err := OpenDatabase(path)
	if err != nil {
		log.Printf("❌ Failed to open DB: %v", err)
		return "", err
	}
	defer db.Close()

	query := "SELECT id FROM users WHERE username = ? LIMIT 1"
	var userID string
	err = db.QueryRow(query, username).Scan(&userID)
	if err != nil {
		log.Printf("❌ Failed to get user ID for username '%s': %v", username, err)
		return "", err
	}
	log.Printf("📤 Got user ID for username '%s': %s", username, userID)
	return userID, nil
}

func GetUserByUsername(path, username string) (*User, error) {
	db, err := OpenDatabase(path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT id, username, password FROM users WHERE username = ?"
	var user User
	err = db.QueryRow(query, username).Scan(&user.UserId, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}


func GetMaxId(path string, table string) (int, error) {
	db, err := OpenDatabase(path)
	if err != nil {
		log.Printf("❌ Failed to open DB: %v", err)
		return 0, err
	}
	defer db.Close()

	query := "SELECT COALESCE(MAX(id), 0) FROM " + table
	var result int
	err = db.QueryRow(query).Scan(&result)
	if err != nil {
		log.Printf("❌ Failed to get max id: %v", err)
		return 0, err
	}
	if result == 0 {
		log.Println("ℹ️ No entries yet. Returning default ID = 1")
		return 1, nil
	}
	log.Printf("📈 Max ID from table %s: %d", table, result)
	return result, nil
}
