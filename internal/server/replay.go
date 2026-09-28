package server

import (
	"errors"
	"strings"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
)

func (s *Server) ReplayCommand(command protocol.Value) error {
	if command.Type != protocol.Array {
		return errors.New("AOF-Replay command not array error")
	}
	if len(command.Array) == 0 {
		return errors.New("AOF-REPLAY no command error")
	}

	cmd := command.Array[0]
	if cmd.Type != protocol.BulkString {
		return errors.New("AOF-REPLAY bulkstring error")
	}
	args := make([]string, len(command.Array)-1)

	for i, value := range command.Array[1:] {
		if value.Type != protocol.BulkString {
			return errors.New("AOF-REPLAY arguments must be bulk strings")
		}
		args[i] = value.Str
	}

	commandName := strings.ToUpper(cmd.Str)
	switch commandName {
	//Core
	case "SET":
		return s.replaySet(args)
	case "DEL":
		return s.replayDeleteCommand(args)
	case "EXPIRE":
		return s.replayExpire(args)
	case "PERSIST":
		return s.replayPersist(args)
	//Integers
	case "INCR":
		return s.replayIncr(args)
	case "DECR":
		return s.replayDecr(args)
	case "INCRBY":
		return s.replayIncrby(args)
	case "DECRBY":
		return s.replayDecrby(args)
	//Strings
	case "APPEND":
		return s.replayAppendCommand(args)
	case "SETNX":
		return s.replaySetnx(args)
	case "MSET":
		return s.replayMSet(args)
	//Lists
	case "LPUSH":
		return s.replayLPush(args)
	case "RPUSH":
		return s.replayRPush(args)
	case "LPOP":
		return s.replayLPop(args)
	case "RPOP":
		return s.replayRPop(args)
	case "LSET":
		return s.replayLSet(args)
	case "LTRIM":
		return s.replayLTrim(args)
	//Hashes
	case "HSET":
		return s.replayHSet(args)
	case "HDEL":
		return s.replayHDel(args)
	//ZSET
	case "ZADD":
		return s.replayZAdd(args)
	case "ZREM":
		return s.replayZRem(args)
	case "ZINCRBY":
		return s.replayZIncrby(args)
	default:
		return errors.New("AOF-REPLAY invalid command")
	}
}
