package core

type IGame interface {
	GetID() string
	GetOwner() string
	GetStatus() *GameStatus
	GetTable() *Table
	Enter(player, secret string) error
	Leave(player string) error
	Play() error
	Stop(force bool) error
}
