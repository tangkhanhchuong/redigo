package core

import (
	"log"
)

func ExecuteCommand(cmd *RedigoCommand) []byte {
	log.Printf("parsed command: %+v\n", cmd)

	var res []byte
	switch cmd.Cmd {
	case "PING":
		res = cmdPING(cmd.Args)
	case "SET":
		res = cmdSET(cmd.Args)
	case "GET":
		res = cmdGET(cmd.Args)
	case "TTL":
		res = cmdTTL(cmd.Args)
	default:
		res = []byte("-unknown command '" + cmd.Cmd + "'\r\n")
	}

	return res
}
