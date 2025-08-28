package constant

var RespNil = []byte("$-1\r\n")
var RespOk = []byte("$2\r\nOK\r\n")
var TtlKeyNotExistCode = -2
var TtlKeyExistNoExpiryCode = -1
