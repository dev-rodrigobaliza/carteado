package model

type Logout uint8

const (
	LogoutUnknown Logout = iota
	LogoutAutomatic
	LogoutManual
)

type Login struct {
	Base   `gorm:"embedded"`
	UserID uint64 `gorm:"index;default:0"`
	Token  string `gorm:"type:varchar(200);index;not null"`
	Logout Logout `gorm:"default:0;not null"`
	User   User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
