package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	// Messages from local clients (to be published to Redis)
	broadcast chan []byte

	// Messages from Redis (to be sent to local clients)
	redisInbound chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		redisInbound: make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// A dedicated Goroutine (Infinite Loop)
func (h *Hub) run(rdb *redis.Client) {

	// 1. Start a background Goroutine to listen to Redis
	// Pumps messages from Redis -> h.redisInbound
	go func() {
		ctx := context.Background()
		pubsub := rdb.Subscribe(ctx, "chat_room")
		defer pubsub.Close()

		ch := pubsub.Channel()

		for msg := range ch {
			// When Redis sends a message, push it to main loop
			h.redisInbound <- []byte(msg.Payload)
		}
	}()

	// 2. Main Loop
	for {
		select {
		// 1. Someone joined
		case client := <-h.register:
			h.clients[client] = true
		
		// 2. Someone left
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		
		// CASE A: Local User sent a message
		// ACTION: Publish it to Redis (don't send to clients yet!)
		case message := <-h.broadcast:
			// 1. SAVE TO DB (Persistence Layer)
            saveMessage(string(message))

            // 2. PUBLISH TO REDIS (Distribution Layer)
            ctx := context.Background()
            err := rdb.Publish(ctx, "chat_room", message).Err()
            if err != nil {
                log.Printf("Redis publish error: %v", err)
            }

		// CASE B: Redis sent a message
		// ACTION: Send to all local clients
		case message := <-h.redisInbound:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}