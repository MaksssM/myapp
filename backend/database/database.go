package database

import (
	"database/sql"
	"log"
)

var Db *sql.DB // Экспортируемая переменная

func InitDatabase() {
	// Реализация функции инициализации базы данных
	var err error
	Db, err = sql.Open("sqlite", "file:subscriptions.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = Db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	createSubsTable := `CREATE TABLE IF NOT EXISTS subscriptions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		subscribed_to INTEGER NOT NULL
	)`
	if _, err := Db.Exec(createSubsTable); err != nil {
		log.Fatalf("Failed to create subscriptions table: %v", err)
	}

	createPostsTable := `CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		author TEXT NOT NULL,
		content TEXT NOT NULL,
		media TEXT,
		date TEXT NOT NULL
	)`
	if _, err := Db.Exec(createPostsTable); err != nil {
		log.Fatalf("Failed to create posts table: %v", err)
	}
}