package utils

import (
	"time"

	str2duration "github.com/xhit/go-str2duration/v2"
)

func StringToDuration(duration string) time.Duration {
	d, err := str2duration.ParseDuration(duration)
	if err != nil {
		return time.Hour * 24
	}

	return d
}
