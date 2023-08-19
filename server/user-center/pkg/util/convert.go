package util

import (
	"fmt"
	"strconv"
)

func StrToUint(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	return i
}

func UintToStr(i uint) string {
	return strconv.FormatInt(int64(i), 10)
}
