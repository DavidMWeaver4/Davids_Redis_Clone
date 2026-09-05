package server

import (
	"errors"
	"strconv"
	"strings"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
)

type commandHandler func(*Server, []string) protocol.Value

var commandHandlers = map[string]commandHandler{
	//Core
	"PING":   ping,
	"SET":    set,
	"GET":    get,
	"DEL":    deleteCommand,
	"EXISTS": exists,
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
	"ZADD":    zadd,
	"ZSCORE":  zscore,
	"ZCARD":   zcard,
	"ZREM":    zrem,
	"ZRANGE":  zrange,
	"ZRANK":   zrank,
	"ZINCRBY": zincrby,
}
var (
	ErrInvalidInteger = errors.New("invalid integer entered")
	ErrInvalidFloat   = errors.New("invalid float entered")
)

func (s *Server) execute(command protocol.Value) protocol.Value {
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
	handler, ok := commandHandlers[strings.ToUpper(cmd.Str)]
	if !ok {
		return protocol.NewError("invalid command")
	}
	return handler(s, args)
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
