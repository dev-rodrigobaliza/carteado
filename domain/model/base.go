package model

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID         uint64         `json:"id,omitempty" gorm:"primarykey"`
	CreatedAt  time.Time      `json:"-" gorm:"autoCreateTime;not null;default:current_timestamp"`
	ModifiedAt time.Time      `json:"-" gorm:"autoCreateTime;autoUpdateTime;not null;default:current_timestamp ON update current_timestamp"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
