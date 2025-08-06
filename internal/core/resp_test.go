package core_test

import (
	"fmt"
	"testing"

	"redigo/internal/core"
)

func TestSimpleStringDecode(t *testing.T) {
	cases := map[string]string{
		"+OK\r\n": "OK",
	}

	for input, expected := range cases {
		t.Run(fmt.Sprintf("SimpleString:%q", input), func(t *testing.T) {
			val, err := core.Decode([]byte(input))
			if err != nil {
				t.Fatalf("unexpected error :%v", err)
			}

			if str, ok := val.(string); !ok || str != expected {
				t.Errorf("expected %q, got %T(%v)", expected, val, val)
			}
		})
	}
}

func TestErrorDecode(t *testing.T) {
	cases := map[string]string{
		"-Error message\r\n": "Error message",
	}

	for input, expected := range cases {
		t.Run(fmt.Sprintf("Error:%q", input), func(t *testing.T) {
			value, err := core.Decode([]byte(input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if str, ok := value.(string); !ok || str != expected {
				t.Errorf("expected %q, got %T(%v)", expected, value, value)
			}
		})
	}
}

func TestInt64Decode(t *testing.T) {
	cases := map[string]int64{
		":0\r\n":    0,
		":1000\r\n": 1000,
	}

	for input, expected := range cases {
		t.Run(fmt.Sprintf("Int:%q", input), func(t *testing.T) {
			value, err := core.Decode([]byte(input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if num, ok := value.(int64); !ok || num != expected {
				t.Errorf("expected %d, got %T(%v)", expected, value, value)
			}
		})
	}
}

func TestBulkStringDecode(t *testing.T) {
	cases := map[string]string{
		"$5\r\nhello\r\n": "hello",
		"$0\r\n\r\n":      "",
	}

	for input, expected := range cases {
		t.Run(fmt.Sprintf("BulkString:%q", input), func(t *testing.T) {
			value, err := core.Decode([]byte(input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if str, ok := value.(string); !ok || str != expected {
				t.Errorf("expected %q, got %T(%v)", expected, value, value)
			}
		})
	}
}

func TestArrayDecode(t *testing.T) {
	cases := map[string][]interface{}{
		"*0\r\n":                                                   {},
		"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n":                     {"hello", "world"},
		"*3\r\n:1\r\n:2\r\n:3\r\n":                                 {int64(1), int64(2), int64(3)},
		"*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$5\r\nhello\r\n":            {int64(1), int64(2), int64(3), int64(4), "hello"},
		"*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Hello\r\n-World\r\n": {[]interface{}{int64(1), int64(2), int64(3)}, []interface{}{"Hello", "World"}},
	}

	for input, expected := range cases {
		t.Run(fmt.Sprintf("Array:%q", input), func(t *testing.T) {
			value, err := core.Decode([]byte(input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			array, ok := value.([]interface{})
			if !ok {
				t.Fatalf("expected []interface{}, got %T(%v)", value, value)
			}
			if len(array) != len(expected) {
				t.Fatalf("expected length %d, got %d", len(expected), len(array))
			}
			for i := range array {
				if fmt.Sprintf("%v", array[i]) != fmt.Sprintf("%v", expected[i]) {
					t.Errorf("at index %d: expected %v, got %v", i, expected[i], array[i])
				}
			}
		})
	}
}

func TestParseCmd(t *testing.T) {
	cases := map[string]core.RedigoCmd{
		"*3\r\n$3\r\nput\r\n$5\r\nhello\r\n$5\r\nworld\r\n": {
			Cmd:  "put",
			Args: []string{"hello", "world"},
		},
		"*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n": {
			Cmd:  "ECHO",
			Args: []string{"hello"},
		},
	}

	for input, expected := range cases {
		t.Run(expected.Cmd, func(t *testing.T) {
			cmd, err := core.ReadRESPCommand([]byte(input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cmd.Cmd != expected.Cmd {
				t.Errorf("expected Cmd=%q, got %q", expected.Cmd, cmd.Cmd)
			}

			if len(cmd.Args) != len(expected.Args) {
				t.Errorf("expected %d args, got %d", len(cmd.Args), len(expected.Args))
			}

			for i := range cmd.Args {
				if cmd.Args[i] != expected.Args[i] {
					t.Errorf("arg[%d]: expected %q, got %q", i, expected.Args[i], cmd.Args[i])
				}
			}
		})
	}
}
