package websocket

type WSMessage struct {
	Player *Player
	Data   []byte
}

func NewWSMessage(player *Player, data []byte) *WSMessage {
	return &WSMessage{
		Player: player,
		Data:   data,
	}
}
