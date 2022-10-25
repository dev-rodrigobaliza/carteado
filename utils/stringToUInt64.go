package utils

import "strconv"

func StringToUint64(s string) (uint64, error) {
	if s == "" {
		return 0, nil
	}

	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return i, nil
}
