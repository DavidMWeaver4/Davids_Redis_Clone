# Davids_Redis_Clone

A Redis-inspired in-memory key-value server written in Go. Built to explore TCP networking, RESP, backend server design, and concurrent fundamentals. 

![Status](https://github.com/DavidMWeaver4/Davids_Redis_Clone/actions/workflows/ci.yml/badge.svg)

>Project started on 2026/August/02

## Tech Stack

Go, TCP networking, RESP, sync.RWMutex, bufio, race detector


## Features
- TCP server
- Concurrent clients
- RESP encoding
- RESP decoding
- Redis like responses
- Unit testing
- Support for strings and lists

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

## Testing

- Unit tests for protocol, store, and command layers
- Race detector
- go vet
- staticcheck
- Automated CI

## Current Tree
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
```
> Readme last updated on 2026/August/25
