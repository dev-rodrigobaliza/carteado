package response

type Card struct {
	Card  string `json:"card"`
	Value int    `json:"value"`
}

func NewCard(card string, value int) *Card {
	return &Card{
		Card:  card,
		Value: value,
	}
}

type Deck struct {
	TableID string  `json:"table_id"`
	Cards   []*Card `json:"cards,omitempty"`
}

func NewDeck(tableID string, cards []*Card) *Deck {
	return &Deck{
		TableID: tableID,
		Cards:   cards,
	}
}

type Game struct {
	State     string `json:"state,omitempty"`
	CreatedAt string `json:"created_at"`
	StartedAt string `json:"started_at,omitempty"`
}

func NewGame(createdAt, startedAt string) *Game {
	game := &Game{
		CreatedAt: createdAt,
		StartedAt: startedAt,
	}

	return game
}

type Group struct {
	ID           int       `json:"id"`
	MinPlayers   int       `json:"min_players"`
	MaxPlayers   int       `json:"max_players"`
	PlayersCount int       `json:"players_count"`
	Players      []*Player `json:"players,omitempty"`
	CreatedAt    string    `json:"created_at,omitempty"`
}

func NewGroup(id, minPlayers, maxPlayers int, players []*Player, createdAt string) *Group {
	return &Group{
		ID:           id,
		MinPlayers:   minPlayers,
		MaxPlayers:   maxPlayers,
		PlayersCount: len(players),
		Players:      players,
		CreatedAt:    createdAt,
	}
}

type Player struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	TableID   string `json:"table_id,omitempty"`
	GroupID   string `json:"group_id,omitempty"`
	CreatedAt string `json:"created_at"`
	LoggedAt  string `json:"logged_at,omitempty"`
}

func NewPlayer(id, name, tableID, groupID, createdAt, loggedAt string) *Player {
	player := &Player{
		ID:        id,
		Name:      name,
		TableID:   tableID,
		GroupID:   groupID,
		CreatedAt: createdAt,
		LoggedAt:  loggedAt,
	}

	return player
}

type Table struct {
	ID              string    `json:"id"`
	Mode            string    `json:"mode"`
	CreatedBy       string    `json:"created_by"`
	StartedBy       string    `json:"started_by,omitempty"`
	Private         bool      `json:"private"`
	SpectatorsCount int       `json:"spectators_count"`
	Spectators      []*Player `json:"spectators,omitempty"`
	PlayersCount    int       `json:"players_count"`
	GroupsCount     int       `json:"groups_count"`
	Groups          []*Group  `json:"groups,omitempty"`
	Winners         []*Group  `json:"winners,omitempty"`
	Game            *Game     `json:"game"`
	CreatedAt       string    `json:"created_at"`
	StartedAt       string    `json:"started_at,omitempty"`
}

func NewTable(id, mode, createdBy, startedBy, createdAt, startedAt string, private bool, playersCount int, spectators []*Player, groups, winners []*Group, game *Game) *Table {
	return &Table{
		ID:              id,
		Mode:            mode,
		CreatedBy:       createdBy,
		StartedBy:       startedBy,
		Private:         private,
		SpectatorsCount: len(spectators),
		Spectators:      spectators,
		PlayersCount:    playersCount,
		GroupsCount:     len(groups),
		Groups:          groups,
		Winners:         winners,
		Game:            game,
		CreatedAt:       createdAt,
		StartedAt:       startedAt,
	}
}
