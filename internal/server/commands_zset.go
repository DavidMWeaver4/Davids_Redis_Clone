package server

import (
	"errors"
	"math"
	"strconv"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
)

func zadd(s *Server, args []string) protocol.Value {
	if len(args) != 3 {
		return protocol.NewError("need 3 arguments in 'ZADD'")
	}
	score, err := parseScore(args[1])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	statusCode, err := s.store.ZAdd(args[0], score, args[2])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(int64(statusCode))
}

func zscore(s *Server, args []string) protocol.Value {
	if len(args) != 2 {
		return protocol.NewError("need 2 arguments in 'ZSCORE'")
	}
	score, found, err := s.store.ZScore(args[0], args[1])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	if !found {
		return protocol.NewNullBulkString()
	}
	return scoreValue(score)
}
func zcard(s *Server, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("need 1 argument in 'ZCARD'")
	}
	length, err := s.store.ZCard(args[0])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(int64(length))
}
func zrem(s *Server, args []string) protocol.Value {
	if len(args) < 2 {
		return protocol.NewError("need at least 2 arguments in 'ZREM'")
	}
	removed, err := s.store.ZRem(args[0], args[1:]...)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(int64(removed))
}
func zrange(s *Server, args []string) protocol.Value {
	if len(args) != 3 {
		return protocol.NewError("need 3 arguments in 'ZRANGE'")
	}
	start, err := parseIntHelper(args[1])
	if err != nil {
		return protocol.NewError(ErrInvalidInteger.Error())
	}
	end, err := parseIntHelper(args[2])
	if err != nil {
		return protocol.NewError(ErrInvalidInteger.Error())
	}
	payload, err := s.store.ZRange(args[0], start, end)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	results := make([]protocol.Value, 0, len(payload))
	for _, value := range payload {
		results = append(results, protocol.NewBulkString(value))
	}
	return protocol.NewArray(results)
}

/*
 *
 *
 *
 *
 */

func parseScore(raw string) (float64, error) {
	score, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, errors.New("value is not a valid float")
	}
	if math.IsNaN(score) {
		return 0, errors.New("value is not a valid float")
	}
	return score, nil
}

func scoreValue(score float64) protocol.Value {
	return protocol.NewBulkString(
		strconv.FormatFloat(score, 'g', -1, 64),
	)
}
