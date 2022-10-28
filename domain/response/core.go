package response

type Group struct {
	ID           string    `json:"id"`
	PlayersCount int       `json:"players_count"`
	Players      []*Player `json:"players"`
}

func NewGroup(id string, players []*Player) *Group {
	return &Group{
		ID:           id,
		PlayersCount: len(players),
		Players:      players,
	}
}

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
	GroupsCount  int       `json:"groups_count"`
	Groups       []*Group  `json:"groups"`
	Winners      []*Group  `json:"winners"`
}

func NewTable(id, mode, owner string, private bool, players []*Player, groups, winners []*Group) *Table {
	return &Table{
		ID:           id,
		Mode:         mode,
		Owner:        owner,
		Private:      private,
		PlayersCount: len(players),
		Players:      players,
		GroupsCount:  len(groups),
		Groups:       groups,
		Winners:      winners,
	}
}
