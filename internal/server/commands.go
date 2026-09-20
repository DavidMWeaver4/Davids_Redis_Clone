package server

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
)

type commandHandler func(*Server, *Client, []string) protocol.Value

var commandHandlers map[string]commandHandler

func init() {
	commandHandlers = map[string]commandHandler{
		//Core
		"PING":   ping,
		"SET":    set,
		"GET":    get,
		"DEL":    deleteCommand,
		"EXISTS": exists,
		//Transactions
		"MULTI":   multi,
		"EXEC":    exec,
		"DISCARD": discard,
		//Expiration
		"TTL":     ttl,
		"EXPIRE":  expire,
		"PERSIST": persist,
		//Integers
		"INCR":   incr,
		"DECR":   decr,
		"INCRBY": incrby,
		"DECRBY": decrby,
		//Strings
		"APPEND": appendCommand,
		"STRLEN": strlen,
		"SETNX":  setnx,
		"MGET":   mget,
		"MSET":   mset,
		//Lists
		"LPUSH":  lpush,
		"RPUSH":  rpush,
		"LPOP":   lpop,
		"RPOP":   rpop,
		"LLEN":   llen,
		"LRANGE": lrange,
		"LINDEX": lindex,
		"LSET":   lset,
		"LTRIM":  ltrim,
		//Hashes
		"HSET":    hset,
		"HGET":    hget,
		"HDEL":    hdel,
		"HEXISTS": hexists,
		"HLEN":    hlen,
		"HGETALL": hgetall,
		//ZSET
		"ZADD":          zadd,
		"ZSCORE":        zscore,
		"ZCARD":         zcard,
		"ZREM":          zrem,
		"ZRANGE":        zrange,
		"ZRANK":         zrank,
		"ZINCRBY":       zincrby,
		"ZRANGEBYSCORE": zrangebyscore,
	}
}

var (
	ErrInvalidInteger = errors.New("invalid integer entered")
	ErrInvalidFloat   = errors.New("invalid float entered")
)

func (s *Server) execute(client *Client, command protocol.Value) protocol.Value {
	if command.Type != protocol.Array {
		return protocol.NewError("command must be an array")
	}
	if len(command.Array) == 0 {
		return protocol.NewError("no command entered")
	}

	cmd := command.Array[0]
	if cmd.Type != protocol.BulkString {
		return protocol.NewError("command must be a bulk string")
	}
	args := make([]string, len(command.Array)-1)

	for i, value := range command.Array[1:] {
		if value.Type != protocol.BulkString {
			return protocol.NewError("arguments must be bulk strings")
		}
		args[i] = value.Str
	}

	commandName := strings.ToUpper(cmd.Str)
	handler, ok := commandHandlers[commandName]
	if !ok {
		return protocol.NewError("invalid command")
	}
	if isWriteCommand(commandName) && s.persistenceUnavailable() {
		return protocol.NewError("persistence unavailable")
	}
	if client.inTransaction &&
		commandName != "MULTI" &&
		commandName != "EXEC" &&
		commandName != "DISCARD" {

		client.queue = append(client.queue, command)
		return protocol.NewSimpleString("QUEUED")
	}
	response := handler(s, client, args)
	if isWriteCommand(commandName) && response.Type != protocol.Error {
		err := s.appendAOF(command)
		if err != nil {
			s.markPersistenceFailed(err)
			return protocol.NewError("persistence failure")
		}
	}
	return response
}

/*
||\\
||	\\
|| 	 //
||	//
||//
*/

func parseIntHelper(value string) (int, error) {
	value64, err := strconv.ParseInt(value, 10, strconv.IntSize)
	if err != nil {
		return 0, ErrInvalidInteger
	}
	return int(value64), nil
}
func parseFloat64Helper(value string) (float64, error) {
	value64, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, ErrInvalidFloat
	}
	return value64, nil
}
func (s *Server) appendAOF(command protocol.Value) error {
	if s.aof == nil {
		return nil
	}
	return s.aof.Append(command)
}
func isWriteCommand(commandName string) bool {
	switch commandName {
	case "SET", "DEL", "EXPIRE", "PERSIST",
		"INCR", "DECR", "INCRBY", "DECRBY",
		"APPEND", "SETNX", "MSET",
		"LPUSH", "RPUSH", "LPOP", "RPOP",
		"LSET", "LTRIM",
		"HSET", "HDEL",
		"ZADD", "ZREM", "ZINCRBY":
		return true
	default:
		return false
	}
}
func (s *Server) persistenceUnavailable() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.persistenceFailed
}

func (s *Server) markPersistenceFailed(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.persistenceFailed {
		log.Printf("AOF append failed: %v", err)
		s.persistenceFailed = true
	}
}
