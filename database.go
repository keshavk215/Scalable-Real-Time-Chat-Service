package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq" // Postgres driver
)

// Message struct defines our data shape
type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

var db *sql.DB

func initDB(dataSourceName string) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	// Instead of panicking immediately, we try to connect 10 times.
	// This handles the case where the Docker container is still booting up.
	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			log.Println("Successfully connected to PostgreSQL!")
			break
		}
		log.Printf("Database not ready yet (Attempt %d/10)... Retrying in 2 seconds...", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Panicf("Could not connect to database after retries: %v", err)
	}
	
	createTables()
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	-- Index for fast retrieval of latest messages
	CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at DESC);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Panic(err)
	}
}

func saveMessage(content string) {
	_, err := db.Exec("INSERT INTO messages (content) VALUES ($1)", content)
	if err != nil {
		log.Printf("Error saving message: %v", err)
	}
}

func getRecentMessages() []string {
	// Resume Point: "Indexing for fast query performance"
	// We fetch the last 50 messages. The Index helps this sort happen instantly.
	rows, err := db.Query("SELECT content FROM messages ORDER BY created_at DESC LIMIT 50")
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	var messages []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			log.Println(err)
			continue
		}
		// Prepend to slice to keep order (Oldest -> Newest)
		messages = append([]string{content}, messages...)
	}
	return messages
}