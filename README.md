# Davids_Redis_Clone

This is a Redis clone written in Go. Built to explore TCP networking, RESP, backend server design, and concurrent fundamentals. 

![Status](https://github.com/DavidMWeaver4/Davids_Redis_Clone/actions/workflows/ci.yml/badge.svg)

>Project started on 2026/August/02

## Tech Stack

Go, TCP, RESP, sync.RWMutex, bufio 


## Features
- TCP server
- Concurrent clients
- RESP encoding
- RESP decoding
- Commands: PING, SET, GET, DEL, EXISTS, TTL, EXPIRE, PERSIST
- Unit testing

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
> Readme last updated on 2026/August/19
