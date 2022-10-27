package response

import "github.com/goccy/go-json"

type Player struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewPlayer(id string, name string) *Player {
	if name == "" {
		name = "* unregistered *"
	}

	return &Player{
		ID:   id,
		Name: name,
	}

}

type Table struct {
	ID           string    `json:"id"`
	Mode         string    `json:"mode"`
	Owner        string    `json:"owner"`
	Private      bool      `json:"private"`
	PlayersCount int       `json:"players_count"`
	Players      []*Player `json:"players"`
}

func NewTable(id, mode, owner string, private bool, players []*Player) *Table {
	return &Table{
		ID:           id,
		Mode:         mode,
		Owner:        owner,
		Private:      private,
		PlayersCount: len(players),
		Players:      players,
	}
}

type WSResponse struct {
	RequestID uint64                 `json:"request_id,omitempty"`
	Status    string                 `json:"status"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// FromBytes converts []byte to Input
func (w *WSResponse) FromBytes(buffer []byte) error {
	return json.Unmarshal(buffer, w)
}

// ToBytes converts Input to []byte
func (w *WSResponse) ToBytes() []byte {
	buffer, _ := json.Marshal(w)
	return buffer
}
