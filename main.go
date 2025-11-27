package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	// 1. Connect to Redis (For Pub/Sub / Scaling)
	// Check for an Environment Variable first, else default to localhost
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis is running here
	})

	// 2. Connect to PostgreSQL (For History / Persistence)
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		// Default for local testing (matches the docker-compose settings below)
		dbConnStr = "postgres://postgres:password@localhost:5433/chat_db?sslmode=disable"
	}
	
	// Initialize the DB 
	initDB(dbConnStr)

	// 3. Start Hub (Pass Redis connection)
	hub := newHub()
	go hub.run(rdb) // Start the Hub in a separate Goroutine

	// 4. Register Routes
	// Serve the static HTML file
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// Handle WebSocket connections
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	log.Printf("Server started on %s", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}