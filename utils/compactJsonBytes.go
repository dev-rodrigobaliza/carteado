package utils

import (
	"bytes"
	"encoding/json"
)

func CompactJson(buffer []byte) string {
	buf := new(bytes.Buffer)
	err := json.Compact(buf, buffer)
	if err != nil {
		return ""
	}

	return buf.String()
}
