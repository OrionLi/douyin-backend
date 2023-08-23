package util

import (
	"strconv"
)

func StringToInt64(str string) int64 {
	parseInt, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return -1
	}
	return parseInt
}
