package database

import (
	"database/sql"
)

func CreateDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTable(db *sql.DB, tableName string, columns map[string]string) error {
	query := "CREATE TABLE IF NOT EXISTS " + tableName + " ("
	for col, typ := range columns {
		query += col + " " + typ + ", "
	}
	query = query[:len(query)-2] + ");" // Remove the last comma and space, then close the parenthesis

	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func InsertData(db *sql.DB, tableName string, data map[string]interface{}) error {
	query := "INSERT INTO " + tableName + " ("
	values := "("
	for col, _ := range data {
		query += col + ", "
		values += "?, "
	}
	query = query[:len(query)-2] + ") VALUES " + values[:len(values)-2] + ");" // Remove the last comma and space, then close the parenthesis

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(data)
	if err != nil {
		return err
	}
	return nil
}

func DeleteData(db *sql.DB, tableName string, condition string) error {
	query := "DELETE FROM " + tableName + " WHERE " + condition + ";"
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func UpdateData(db *sql.DB, tableName string, data map[string]interface{}, condition string) error {
	query := "UPDATE " + tableName + " SET "
	for col, _ := range data {
		query += col + " = ?, "
	}
	query = query[:len(query)-2] + " WHERE " + condition + ";" // Remove the last comma and space, then add the condition

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(data)
	if err != nil {
		return err
	}
	return nil
}