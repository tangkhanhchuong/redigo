package core

import (
	"errors"
)

const CRLF string = "\r\n"

// +OK\r\n => OK
func readSimpleString(data []byte) (string, int, error) {
	pos := 1
	for data[pos] != '\r' {
		pos++
	}
	return string(data[1:pos]), pos + 2, nil
}

// :123\r\n => 123
func readInt64(data []byte) (int64, int, error) {
	pos := 1
	var sign int64 = 1
	if data[pos] == '-' {
		sign = -1
		pos++
	}
	if data[pos] == '+' {
		pos++
	}

	var res int64 = 0
	for data[pos] != '\r' {
		res = res*10 + int64(data[pos]-'0')
		pos++
	}
	return res * sign, pos + 2, nil
}

func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

func readLen(data []byte) (int, int) {
	res, pos, _ := readInt64(data)
	return int(res), pos
}

// $5\r\nhello\r\n => "hello"
func readBulkString(data []byte) (string, int, error) {
	len, pos := readLen(data)
	return string(data[pos:(pos + len)]), pos + len + 2, nil
}

// *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n => {"hello", "world"}
func readArray(data []byte) ([]interface{}, int, error) {
	len, pos := readLen(data)
	var res []interface{} = make([]interface{}, len)

	for i := range res {
		ele, delta, err := decodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		res[i] = ele
		pos += delta
	}
	return res, pos, nil
}

func decodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}

	switch data[0] {
	case '+':
		return readSimpleString(data)
	case ':':
		return readInt64(data)
	case '-':
		return readError(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	default:
		return nil, 0, errors.New("unknown resp type")
	}
}

func Decode(data []byte) (interface{}, error) {
	res, _, err := decodeOne(data)
	return res, err
}

func ReadRESPCommand(data []byte) (*RedigoCmd, error) {
	val, err := Decode(data)
	if err != nil {
		return nil, err
	}

	arr, ok := val.([]interface{})
	if !ok {
		return nil, errors.New("command is not an array")
	}

	tokens := make([]string, len(arr))
	for i := range arr {
		str, ok := arr[i].(string)
		if !ok {
			return nil, errors.New("command element is not a string")
		}
		tokens[i] = str
	}

	return &RedigoCmd{Cmd: tokens[0], Args: tokens[1:]}, nil
}
