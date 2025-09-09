package core

import (
	"time"

	"redigo/constant"
)

func ActiveDeleteExpiredKeys() {
	for {
		var expiredCount = 0
		var sampleCountRemain = constant.ActiveExpireSampleSize
		for key, expiredTime := range dictStore.GetKeyExpiredStore() {
			sampleCountRemain--
			if sampleCountRemain < 0 {
				break
			}

			if expiredTime < uint64(time.Now().UnixMilli()) {
				dictStore.Del(key)
				print("delete key: ", key)
				expiredCount++
			}
		}

		// continue if exceeding threshold
		if float64(expiredCount/constant.ActiveExpireSampleSize) <= constant.ActiveExpireThreshold {
			break
		}
	}
}
