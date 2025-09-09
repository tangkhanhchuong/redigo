package constant

import "time"

var RespNil = []byte("$-1\r\n")
var RespOk = []byte("$2\r\nOK\r\n")
var TtlKeyNotExistCode = -2
var TtlKeyExistNoExpiryCode = -1

const ActiveExpireFrequency = 100 * time.Millisecond
const ActiveExpireSampleSize = 20
const ActiveExpireThreshold = 0.1
