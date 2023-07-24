package api

import "strconv"

// atoi converts a string to an integer, returning 0 in case of an error.
func atoi(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return val
}
