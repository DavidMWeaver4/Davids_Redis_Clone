package server

import (
	"errors"
	"math"
	"strconv"
	"strings"

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
	if len(args) != 3 && len(args) != 4 {
		return protocol.NewError("need 3 or 4 arguments in 'ZRANGE'")
	}
	start, err := parseIntHelper(args[1])
	if err != nil {
		return protocol.NewError(ErrInvalidInteger.Error())
	}
	end, err := parseIntHelper(args[2])
	if err != nil {
		return protocol.NewError(ErrInvalidInteger.Error())
	}
	withScores := false
	if len(args) == 4 {
		if strings.ToUpper(args[3]) != "WITHSCORES" {
			return protocol.NewError("syntax error 'WITHSCORES'")
		}
		withScores = true
	}
	payload, err := s.store.ZRange(args[0], start, end)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	results := make([]protocol.Value, 0, len(payload))
	for _, value := range payload {
		results = append(results, protocol.NewBulkString(value.Member))
		if withScores {
			results = append(results, scoreValue(value.Score))
		}
	}
	return protocol.NewArray(results)
}
func zrank(s *Server, args []string) protocol.Value {
	if len(args) != 2 {
		return protocol.NewError("need 2 arguments in 'ZRANK'")
	}
	rank, found, err := s.store.ZRank(args[0], args[1])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	if !found {
		return protocol.NewNullBulkString()
	}
	return protocol.NewInteger(int64(rank))
}
func zincrby(s *Server, args []string) protocol.Value {
	if len(args) != 3 {
		return protocol.NewError("need 3 arguments in 'ZINCRBY'")
	}
	inc, err := parseFloat64Helper(args[1])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	newScore, err := s.store.ZIncrby(args[0], inc, args[2])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return scoreValue(newScore)
}
func zrangebyscore(s *Server, args []string) protocol.Value {
	if len(args) != 5 {
		return protocol.NewError("need 5 arguments in 'ZRANGEBYSCORE'")
	}
	minScore, err := parseFloat64Helper(args[1])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	maxScore, err := parseFloat64Helper(args[2])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	offset, err := parseIntHelper(args[3])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	count, err := parseIntHelper(args[4])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	payload, err := s.store.ZRangeByScore(args[0], minScore, maxScore, offset, count)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	results := make([]protocol.Value, 0, len(payload))
	for _, value := range payload {
		results = append(results, protocol.NewBulkString(value.Member))
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
