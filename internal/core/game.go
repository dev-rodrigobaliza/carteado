package core

type IGame interface {
	GetStatus() *GameStatus
	GetTable() *Table
	Start() error
}
