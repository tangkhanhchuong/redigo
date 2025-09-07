package core

import (
	"errors"

	"redigo/internal/data_structure"
)

func getSet(key string) (*data_structure.SimpleSet, error) {
	val := dictStore.Get(key)
	if val == nil {
		return nil, nil
	}

	parsed, ok := val.(*data_structure.SimpleSet)
	if !ok {
		return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	return parsed, nil
}

func cmdSADD(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'sadd' command"), false)
	}

	key := args[0]
	set, err := getSet(key)
	if err != nil {
		return Encode(err, false)
	}
	if set == nil {
		set = data_structure.NewSimpleSet(key)
		dictStore.Set(key, set)
	}

	count := set.Add(args[1:]...)
	return Encode(count, false)
}

func cmdSREM(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'srem' command"), false)
	}

	key := args[0]
	set, err := getSet(key)
	if err != nil {
		return Encode(err, false)
	}
	if set == nil {
		return Encode(0, false)
	}

	count := set.Remove(args[1:]...)
	return Encode(count, false)
}

func cmdSMEMBERS(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'smembers' command"), false)
	}

	key := args[0]
	set, err := getSet(key)
	if err != nil {
		return Encode(err, false)
	}
	if set == nil {
		return Encode(0, false)
	}

	return Encode(set.Members(), false)
}

func cmdSISMEMBER(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'sismember' command"), false)
	}

	key := args[0]
	set, err := getSet(key)
	if err != nil {
		return Encode(err, false)
	}
	if set == nil {
		return Encode(0, false)
	}

	return Encode(set.IsMember(args[1]), false)
}

func cmdSCARD(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'scard' command"), false)
	}

	key := args[0]
	set, err := getSet(key)
	if err != nil {
		return Encode(err, false)
	}
	if set == nil {
		return Encode(0, false)
	}

	return Encode(set.Card(), false)
}
