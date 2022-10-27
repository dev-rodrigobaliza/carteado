package utils

import (
	"fmt"

	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake

func init() {
	sf = sonyflake.NewSonyflake(sonyflake.Settings{})
	if sf == nil {
		panic("sonyflake not created")
	}
}

func NewUUID(prefix string) string {
	id, err := sf.NextID()
	if err != nil {
		return RandomString(15)
	}

	return fmt.Sprintf("%s%x", prefix, id)
}
