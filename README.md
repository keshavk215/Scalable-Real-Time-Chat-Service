# Scalable Real-Time Chat System (Go, Redis, PostgreSQL)

## 📌 Project Overview

A high-performance, real-time chat backend designed to handle high concurrency and horizontal scaling. The system uses **Go** for the backend logic, **Redis Pub/Sub** for distributing messages across server instances, and **PostgreSQL** for durable data storage.

This project demonstrates the transition from a simple in-memory chat application to a distributed system architecture.

## 🚀 Key Features

* **High Concurrency:** Implemented using Go **Goroutines** and **Channels** to manage thousands of active WebSocket connections efficiently.
* **Low Latency:** Uses **WebSockets** for bi-directional, real-time communication (unlike HTTP polling).
* **Horizontal Scaling:** Integrated **Redis Pub/Sub** as a message broker. This allows multiple Go server instances to talk to each other, preventing "siloed" chat rooms.
* **Data Persistence:** Designed a **PostgreSQL** schema with indexing (`CREATE INDEX`) to optimize the retrieval of chat history (O(log n) performance).
* **Containerization:** Uses **Docker Compose** to orchestrate the database and message broker services.

## 🛠️ Tech Stack

* **Language:** Go (Golang)
* **Real-time Protocol:** WebSockets (`gorilla/websocket`)
* **Message Broker:** Redis
* **Database:** PostgreSQL (`lib/pq`)
* **Infrastructure:** Docker

## ⚙️ How to Run

### 1. Prerequisites

* Go 1.18+
* Docker & Docker Compose

### 2. Start Infrastructure

Spin up the Redis and PostgreSQL containers:

```
docker-compose up -d
```

### 3. Run the Server

Run the Go application (this handles DB migrations automatically):

```
go run main.go hub.go client.go database.go
```

### 4. Usage

* Open `http://localhost:8080` in your browser.
* Open a second tab (or a different browser) to the same URL.
* Messages sent in one tab will appear in the other instantly (Real-time).
* Refreshing the page will load previous messages (Persistence).

## 📐 Architecture Diagram

```
graph TD
    User1[Client A] -->|WebSocket| LB[Go Server Instance 1]
    User2[Client B] -->|WebSocket| LB2[Go Server Instance 2]
  
    LB -->|Publish/Subscribe| Redis[(Redis Pub/Sub)]
    LB2 -->|Publish/Subscribe| Redis
  
    LB -->|Read/Write| DB[(PostgreSQL)]
    LB2 -->|Read/Write| DB
```

## 🔍 Code Highlights

* **`hub.go`** : Manages the "Room" state and handles the Redis Pub/Sub integration.
* **`client.go`** : Handles individual WebSocket connections using dedicated Goroutines for Reading and Writing.
* **`database.go`** : manages SQL connections and ensures indexes are created on startup.
