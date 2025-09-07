# Scalable Real-Time Chat Service in Go

A high-performance, real-time chat service built in Go, designed for massive concurrency and horizontal scalability. This project demonstrates advanced backend engineering principles, including real-time communication with WebSockets, state management, and distributed messaging with Redis.

## Key Features

- **Real-Time Messaging**: Low-latency, bi-directional communication between thousands of clients using WebSockets.
- **High Concurrency**: Leverages Go's native concurrency model (Goroutines and Channels) to efficiently handle a large number of simultaneous connections.
- **Horizontal Scalability**: Uses a Redis Pub/Sub backend as a message broker, allowing the service to be deployed across multiple server instances without losing state.
- **Multiple Chat Rooms**: Supports the creation of and subscription to different chat rooms.
- **Persistent Data**: Stores user profiles, chat history, and room information in a PostgreSQL database.

## Tech Stack

- **Language**: Go
- **Real-Time Communication**: WebSockets (`gorilla/websocket` library)
- **State Management & Messaging**: Redis (Pub/Sub)
- **Database**: PostgreSQL (`pq` driver)
- **Containerization**: Docker

## Architecture Overview

The system is designed to be horizontally scalable. When a user sends a message via their WebSocket connection to one server instance, that instance does not broadcast the message directly to other connected clients. Instead, it publishes the message to a specific Redis channel.

All server instances are subscribed to the Redis channels. When a message is published, Redis broadcasts it to *all* server instances. Each instance then checks which of its locally connected clients are subscribed to that chat room and forwards the message to them via their WebSocket connections. This architecture decouples the servers and allows the system to scale to handle an immense number of users.

## Setup and Installation (Placeholder)

1. Clone the repository.
2. Ensure you have Go, Docker, Redis, and PostgreSQL installed.
3. Install Go dependencies: `go mod tidy`
4. Build the application: `go build ./cmd/main.go`

## Usage (Placeholder)

Run the compiled application: `./main`

The server will start on `localhost:8080`. Clients can connect to the WebSocket endpoint at `ws://localhost:8080/ws`.
