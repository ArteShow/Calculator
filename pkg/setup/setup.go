package setup

import (
	"database/sql"
	"fmt"
	"log"

	config "github.com/ArteShow/Calculator/pkg/Config"
	database "github.com/ArteShow/Calculator/pkg/Database"
	jwt "github.com/ArteShow/Calculator/pkg/JWT"

	_ "modernc.org/sqlite"
)

func Setup() {
	// Open the database
	db, err := database.OpenDatabase(config.GetDatabasePath())
	if err != nil {
		log.Fatalf("❌ Failed to open database: %v", err)
	}
	defer db.Close()

	log.Printf("📂 Database opened at path: %s", config.GetDatabasePath())

	// Define the tables and their columns
	tables := map[string]map[string]string{
		"users": {
			"id":       "INTEGER PRIMARY KEY AUTOINCREMENT",
			"username": "TEXT NOT NULL UNIQUE",
			"password": "TEXT NOT NULL",
		},
		"jwt": {
			"id":  "INTEGER PRIMARY KEY AUTOINCREMENT",
			"key": "TEXT NOT NULL",
		},
		"calculations": {
			"id":          "INTEGER PRIMARY KEY AUTOINCREMENT",
			"userId":      "INTEGER NOT NULL",
			"calculation": "TEXT NOT NULL",
			"result":      "TEXT NOT NULL",
		},
	}

	// Create tables in the database
	for tableName, columns := range tables {
		err = CreateTable(db, tableName, columns)
		if err != nil {
			log.Fatalf("❌ Failed to create table %s: %v", tableName, err)
		}
	}

	// Generate the JWT key using the function
	jwtKey, err := jwt.GenerateJWTKey(32) // Generate a 32-byte key
	if err != nil {
		log.Fatalf("❌ Failed to generate JWT key: %v", err)
	}

	// Insert the generated JWT key into the 'jwt' table
	err = jwt.InsertJWTKey(db, jwtKey)
	if err != nil {
		log.Fatalf("❌ Failed to insert JWT key: %v", err)
	}

	log.Println("✅ JWT key inserted successfully 🚀")
}

// InsertJWTKey inserts the generated JWT key into the 'jwt' table


func InsertIfNotExists(db *sql.DB, table string, data map[string]interface{}, uniqueField string) error {
	var exists int
	checkQuery := fmt.Sprintf("SELECT 1 FROM %s WHERE %s = ? LIMIT 1", table, uniqueField)
	err := db.QueryRow(checkQuery, data[uniqueField]).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("❌ error checking existing data: %w", err)
	}
	if exists == 1 {
		log.Printf("ℹ️ User with %s '%v' already exists. Skipping insert 😊", uniqueField, data[uniqueField])
		return nil
	}

	columns := ""
	placeholders := ""
	values := []interface{}{}
	for col, val := range data {
		columns += col + ", "
		placeholders += "?, "
		values = append(values, val)
	}
	columns = columns[:len(columns)-2]
	placeholders = placeholders[:len(placeholders)-2]

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, columns, placeholders)
	_, err = db.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("❌ failed to insert test data into %s: %w", table, err)
	}
	return nil
}


func CreateTable(db *sql.DB, tableName string, columns map[string]string) error {
	query := "CREATE TABLE IF NOT EXISTS " + tableName + " ("
	for col, typ := range columns {
		query += col + " " + typ + ", "
	}
	query = query[:len(query)-2] + ");"
	_, err := db.Exec(query)
	return err
}
