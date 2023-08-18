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

func UInt2Str(i int64) string {
	return strconv.FormatInt(i, 10)
}
