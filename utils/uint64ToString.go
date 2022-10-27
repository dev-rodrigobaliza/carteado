package utils

import "strconv"

func Uint64ToString(u uint64) string {
	return strconv.FormatUint(u, 10)
}
