package core

type IGame interface {
	EnterGame(player, secret string) error
	GetStatus() *GameStatus
	GetTable() *Table
	LeaveGame(player string) error
	Start() error
	isPrivate() bool
}
