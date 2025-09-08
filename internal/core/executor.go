package core

import (
	"log"
)

// TODO: Implement ZRANGE, ZCARD, ZCOUNT, ZREM
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
	case "DEL":
		res = cmdDEL(cmd.Args)
	case "TTL":
		res = cmdTTL(cmd.Args)
	case "SADD":
		res = cmdSADD(cmd.Args)
	case "SREM":
		res = cmdSREM(cmd.Args)
	case "SMEMBERS":
		res = cmdSMEMBERS(cmd.Args)
	case "SISMEMBER":
		res = cmdSISMEMBER(cmd.Args)
	case "SCARD":
		res = cmdSCARD(cmd.Args)
	case "ZADD":
		res = cmdZADD(cmd.Args)
	case "ZSCORE":
		res = cmdZSCORE(cmd.Args)
	case "ZRANK":
		res = cmdZRANK(cmd.Args)
	default:
		res = []byte("-unknown command '" + cmd.Cmd + "'\r\n")
	}

	return res
}
