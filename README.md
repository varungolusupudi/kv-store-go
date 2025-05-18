# kv-store-go

A simple in-memory Redis-like key-value store built from scratch in Go.  
Supports concurrent clients, TTL-based expiration, and a minimal CLI interface.

## Features
- TCP server handling multiple clients concurrently
- Commands: `SET`, `GET`, `DEL`, `EXPIRE`
- TTL-based expiration with background cleanup
- Concurrency-safe using `sync.RWMutex`
- Lightweight CLI client

## How to Run

### Start the server:
```bash
cd server
go run server.go
```

### Open another terminal and start a client:
```bash
cd client/
go run client.go
```

## Example Usage (from client)
```
SET name varun
OK

GET name
varun

EXPIRE name 10
OK

GET name
varun

# (after 10 seconds)
GET name
Key expired
```