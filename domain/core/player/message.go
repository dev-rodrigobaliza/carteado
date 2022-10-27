package player

type Message[T any] struct {
	Player T
	Data   []byte
}
