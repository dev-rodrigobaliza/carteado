package utils

import "time"

func StringToDuration(duration string) time.Duration {
	d, _ := time.ParseDuration(duration)

	return d
}
