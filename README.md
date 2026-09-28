# Davids Redis Clone

A Redis-inspired in-memory key-value server written in Go. Built to explore TCP networking, RESP, backend server design, and concurrent fundamentals. This project was built from scratch, while using AI for code review.

![Status](https://github.com/DavidMWeaver4/Davids_Redis_Clone/actions/workflows/ci.yml/badge.svg)

>Project started on 2026/August/02

## Tech Stack

- Go 1.25.6
- TCP networking
- RESP
- sync.RWMutex
- bufio
- Skip lists
- File I/O
- Unit testing
- race detector
- go vet
- Staticcheck
- GitHub Actions



## Features
- TCP server
- Concurrent client handling
- Graceful shutdown
- Transactions with MULTI, EXEC, and DISCARD
- In-memory key-value storage
- Thread safe storage using sync.RWMutex
- String, integer, list, hash, and sorted set data types
- Sorted sets backed by a skip list
- RESP encoding / decoding
- Append-only file persistence
- AOF replay during server startup
- Redis like responses
- Unit testing on protocol, server, persistence, and storage layers
- Key expiration (lazy and active)
- Race detection
- Automated CI

## Architecture 
```
Client
  │
  ▼
TCP Connection
  │
  ▼
RESP Reader
  │
  ▼
Command Execution
  │
  ├──────────────► AOF Persistence ──────► AOF File
  │
  ▼
Store
  ├── Strings
  ├── Lists
  ├── Hashes
  └── Sorted Sets
        │
        ▼
     Skip List
  │
  ▼
RESP Writer
  │
  ▼
Client
```

The server separates protocol handling, command execution, persistence, and data storage so that each layer can be tested and refactored independently 

## Implemented Commands
### Core Commands
- PING
-	SET
-	GET
-	DEL
### Expiration/Key Management
-	EXISTS
- TTL
-	EXPIRE
-	PERSIST
Keys use both lazy expiration on access and active expiration.
### Integer Commands
-	INCR
-	DECR
-	INCRBY
- DECRBY
Integer operations include validation and overflow/underflow handling
### String Commands
-	APPEND
-	STRLEN
-	SETNX
-	MGET
-	MSET
### List Commands
-	LPUSH
-	RPUSH
-	LPOP
-	RPOP
-	LLEN
-	LRANGE
-	LINDEX
-	LSET
-	LTRIM
### Hash Commands
- HSET
- HGET
- HDEL
- HEXISTS
- HLEN
- HGETALL
### Sorted Set Commands
- ZADD
- ZSCORE
- ZCARD
- ZREM
- ZRANGE
- ZRANK
- ZINCRBY
- ZRANGEBYSCORE

Sorted sets use a skip list to maintain members ordered by score and support efficient ordered traversal
### Transactions
- MULTI
- EXEC
- DISCARD
Commands issued after MULTI are queued for execution by EXEC. DISCARD clears the queued commands without execution.
## RESP
Supported RESP values
- Simple strings
- Errors
- Integers
- Bulk strings
- Arrays
- Null Bulk strings
- Null Arrays
## Persistence
The server currently uses an Append-Only File for persistence
Write commands are encoded using RESP and appended to the AOF.
On startup:
AOF File
   ↓
RESP Reader
   ↓
Replay
   ↓
Server Store
The AOF implementation supports three fsync policies:

always — fsync after every write
everysec — fsync approximately once per second
no — rely on the operating system to flush writes

The default policy is everysec.

The persistence layer also detects asynchronous fsync failures. Once detected, persistence is flagged and becomes unavailable, any later write command requests are rejected because they cannot be safely stored. 

## Concurrency
The store layer uses sync.RWMutex to protect shared in-memory state while allowing concurrent read access. I wanted to implement lazy expiration so if an expired key is accessed, it will be deleted immediately, which requires Lock() instead of RLock() even on read operations



## Testing

- Unit tests for protocol, store, and command layers
- Race detector
- go vet
- staticcheck
- Automated CI

Run the complete local check with:
```
make check
```
make check runs:
- go fmt
- go vet
- staticcheck
- go test -race

Or run the test suite directly:
```
go test ./...
```
Race testing:
```
go test -race ./...
```
Generate Coverage:
```
make test-race-cover
```

## Project Tree
```
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── persistence/
│   │   ├── aof.go
│   │   ├── aof_replay.go
│   │   └── ...
│   ├── protocol/
│   │   ├── resp_reader.go
│   │   ├── resp_writer.go
│   │   ├── value.go
│   │   └── ...
│   ├── server/
│   │   ├── commands.go
│   │   ├── transactions.go
│   │   └── ...
│   └── store/
│       ├── store.go
│       ├── store_strings.go
│       ├── store_lists.go
│       ├── store_hashes.go
│       ├── zset.go
│       ├── store_expirations.go
│       └── skiplist/
├── Makefile
├── README.md
├── go.mod
└── mise.toml
```
## Running

Start the server with:
```
make run
```
Build the server with:
```
make build
```


> README last updated on 2026/September/28
