package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Verify the connection is active
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Successfully connected to local SQLite database (database.db)!")

	// Create tables if they do not exist (SQLite3 compatible syntax)
	createTablesQuery := `
	CREATE TABLE IF NOT EXISTS users (
		email TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		password TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL,
		user_email TEXT,
		FOREIGN KEY (user_email) REFERENCES users(email) ON DELETE CASCADE
	);`

	_, err = DB.Exec(createTablesQuery)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
	fmt.Println("Database schema is ready!")
}
