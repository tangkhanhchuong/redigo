package core

import (
	"errors"
	"strconv"
	"time"

	"redigo/constant"
)

func cmdPING(args []string) []byte {
	if len(args) > 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'ping' command"), false)
	}

	if len(args) == 0 {
		return Encode("PONG", false)
	}
	return Encode(args[0], false)
}

func cmdSET(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'set' command"), false)
	}
	if len(args) == 3 || len(args) > 4 {
		return Encode(errors.New("ERR syntax error"), false)
	}

	var ttlMs int64 = -1
	var key, val string = args[0], args[1]
	if len(args) > 2 {
		ttlSec, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return Encode(errors.New("ERR value is not an integer or out of range"), false)
		}
		ttlMs = ttlSec * 1000
	}

	dictStore.Set(key, val)
	dictStore.SetExpiry(key, ttlMs)

	return constant.RespOk

}

func cmdGET(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'get' command"), false)
	}

	key := args[0]
	val := dictStore.Get(key)
	if val == nil {
		return constant.RespNil
	}

	if dictStore.HasExpired(key) {
		return constant.RespNil
	}

	return Encode(val, false)
}

func cmdDEL(args []string) []byte {
	if len(args) == 0 {
		return Encode(errors.New("ERR wrong number of arguments for 'del' command"), false)
	}

	count := 0
	for _, key := range args {
		dictStore.Del(key)
		count += 1
	}

	return Encode(count, false)
}

func cmdTTL(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'ttl' command"), false)
	}

	key := args[0]
	val := dictStore.Get(key)
	if val == nil {
		return Encode(constant.TtlKeyNotExistCode, false)
	}

	exp, existed := dictStore.GetExpiry(key)
	if !existed {
		return Encode(constant.TtlKeyExistNoExpiryCode, false)
	}

	remain := int64(exp) - int64(time.Now().UnixMilli())
	if remain < 0 {
		return Encode(constant.TtlKeyNotExistCode, false)
	}

	return Encode(remain/1000, false)
}
