package util

import (
	"fmt"
	"strconv"
)

func StrToInt64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	return i
}

func UintToStr(i uint) string {
	return strconv.FormatInt(int64(i), 10)
}

func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}
