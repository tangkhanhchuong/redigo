package core

import (
	"errors"
	"fmt"
	"strconv"

	"redigo/constant"
	"redigo/internal/data_structure"
)

func getSortedSet(key string) (*data_structure.SortedSet, error) {
	val := dictStore.Get(key)
	if val == nil {
		return nil, nil
	}

	parsed, ok := val.(*data_structure.SortedSet)
	if !ok {
		return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	return parsed, nil
}

func cmdZADD(args []string) []byte {
	if len(args) < 3 {
		return Encode(errors.New("ERR wrong number of arguments for 'zadd' command"), false)
	}

	key := args[0]
	scoreIndex := 1
	numScoreEleArgs := len(args) - scoreIndex
	if numScoreEleArgs%2 == 1 || numScoreEleArgs == 0 {
		return Encode(fmt.Errorf("ERR wrong number of arguments for 'zadd' command"), false)
	}

	sortedSet, err := getSortedSet(key)
	if err != nil {
		return Encode(err, false)
	}
	if sortedSet == nil {
		sortedSet = data_structure.NewSortedSet(key)
		dictStore.Set(key, sortedSet)
	}

	count := 0
	for i := scoreIndex; i < len(args); i += 2 {
		score, err := strconv.ParseFloat(args[i], 64)
		if err != nil {
			return Encode(errors.New("ERR value is not a valid float"), false)
		}
		member := args[i+1]

		res := sortedSet.Add(score, member)
		if res != 1 {
			return Encode(errors.New("ERR adding element"), false)
		}
		count++
	}

	return Encode(count, false)
}

// TODO: Need to support float score
func cmdZSCORE(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'zscore' command"), false)
	}
	key, member := args[0], args[1]
	sortedSet, err := getSortedSet(key)
	if err != nil {
		return Encode(err, false)
	}
	if sortedSet == nil {
		return constant.RespNil
	}

	ret, score := sortedSet.GetScore(member)
	if ret != 0 {
		return constant.RespNil
	}
	return Encode(int(score), false)
}

func cmdZRANK(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'zrank' command"), false)
	}

	key, member := args[0], args[1]
	sortedSet, err := getSortedSet(key)
	if err != nil {
		return Encode(err, false)
	}
	if sortedSet == nil {
		return constant.RespNil
	}

	rank, _ := sortedSet.GetRank(member, false)
	return Encode(rank, false)
}
