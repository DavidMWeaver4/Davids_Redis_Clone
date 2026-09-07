# Davids Redis Clone

A Redis-inspired in-memory key-value server written in Go. Built to explore TCP networking, RESP, backend server design, and concurrent fundamentals. 

![Status](https://github.com/DavidMWeaver4/Davids_Redis_Clone/actions/workflows/ci.yml/badge.svg)

>Project started on 2026/August/02

## Tech Stack

- Go
- TCP networking
- RESP
- sync.RWMutex
- bufio
- Skip lists
- Unit testing
- race detector
- go vet
- Staticcheck
- GitHub Actions



## Features
- TCP server
- Concurrent client handling
- Transactions
- In-memory key-value storage
- Sorted sets backed by a skip list
- RESP encoding / decoding
- Redis like responses
- Unit testing on protocol, server, and storage layers
- Key expiration (lazy and active)
- Race detection
- Automated CI

## Architecture 
```
Client
  ↓
TCP connection
  ↓
RESP Reader
  ↓
Command Handler
  ↓
Store
  ↓
RESP Writer
```

The server separates protocol handling, command execution, and data storage so that each layer can be tested and refactored independently 

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
### Integer Commands
-	INCR
-	DECR
-	INCRBY
- DECRBY
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
## RESP
Supported RESP values
- Simple strings
- Errors
- Integers
- Bulk strings
- Arrays
- Null Bulk strings
- Null Arrays

## Concurrency
The store layer uses sync.RWMutex to protect shared in-memory state while allowing concurrent read access. I wanted to implement lazy expiration so if an expired key is accessed, it will be deleted immiedately, which requires Lock() instead of RLock() even on read operations



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
Or run the test suite directly:
```
go test ./...
```
Race testing:
```
go test -race ./...
```

## Project Tree
```
├── Makefile
├── README.md
├── cmd/
│   └── server/
│       └── main.go
├── go.mod
└── internal/
    ├── protocol/
    ├── server/
    └── store/
        └── skiplist/
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


> README last updated on 2026/September/7
