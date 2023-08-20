package util

import (
	"fmt"
	"strconv"
)

// StrToUint 字符串转uint
func StrToUint(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	return i
}

// UintToStr uint转字符串
func UintToStr(i uint) string {
	return strconv.FormatInt(int64(i), 10)
}
