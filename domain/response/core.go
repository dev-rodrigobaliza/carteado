package response

type Group struct {
	ID           int       `json:"id"`
	MinPlayers   int       `json:"min_players"`
	MaxPlayers   int       `json:"max_players"`
	PlayersCount int       `json:"players_count"`
	Players      []*Player `json:"players,omitempty"`
}

func NewGroup(id, minPlayers, maxPlayers int, players []*Player) *Group {
	return &Group{
		ID:           id,
		MinPlayers:   minPlayers,
		MaxPlayers:   maxPlayers,
		PlayersCount: len(players),
		Players:      players,
	}
}

type Player struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	TableID string `json:"table_id,omitempty"`
	GroupID string `json:"group_id,omitempty"`
}

func NewPlayer(id, name, tableID, groupID string) *Player {
	if name == "" {
		name = "* unregistered *"
	}

	player := &Player{
		ID:   id,
		Name: name,
	}
	if tableID != "" {
		player.TableID = tableID
		if groupID == "0" {
			player.GroupID = ""
		} else {
			player.GroupID = groupID
		}
	}

	return player
}

type Table struct {
	ID              string    `json:"id"`
	Mode            string    `json:"mode"`
	Owner           string    `json:"owner"`
	Private         bool      `json:"private"`
	SpectatorsCount int       `json:"spectators_count"`
	Spectators      []*Player `json:"spectators,omitempty"`
	PlayersCount    int       `json:"players_count"`
	GroupsCount     int       `json:"groups_count"`
	Groups          []*Group  `json:"groups,omitempty"`
	Winners         []*Group  `json:"winners,omitempty"`
}

func NewTable(id, mode, owner string, private bool, playersCount int, spectators []*Player, groups, winners []*Group) *Table {
	return &Table{
		ID:              id,
		Mode:            mode,
		Owner:           owner,
		Private:         private,
		SpectatorsCount: len(spectators),
		Spectators:      spectators,
		PlayersCount:    playersCount,
		GroupsCount:     len(groups),
		Groups:          groups,
		Winners:         winners,
	}
}
