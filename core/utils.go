package core

import "strconv"

func string2int(s string) int {
	i, _ := strconv.ParseInt(s, 0, 64)
	return int(i)
}
