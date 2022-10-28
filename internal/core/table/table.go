package table

import (
	"time"

	"github.com/dev-rodrigobaliza/carteado/consts"
	gm "github.com/dev-rodrigobaliza/carteado/domain/core/game"
	"github.com/dev-rodrigobaliza/carteado/domain/core/table"
	"github.com/dev-rodrigobaliza/carteado/domain/response"
	"github.com/dev-rodrigobaliza/carteado/errors"
	"github.com/dev-rodrigobaliza/carteado/internal/core/game"
	"github.com/dev-rodrigobaliza/carteado/internal/core/group"
	"github.com/dev-rodrigobaliza/carteado/internal/core/player"
	"github.com/dev-rodrigobaliza/carteado/pkg/safemap"
	"github.com/dev-rodrigobaliza/carteado/utils"
)

type Table struct {
	id         string
	owner      string
	secret     string
	minPlayers int
	maxPlayers int
	allowBots  bool
	players    *safemap.SafeMap[string, *player.Player]
	groups     *safemap.SafeMap[int, *group.Group]
	winners    []int
	state      table.State
	gameMode   gm.Mode
	game       game.IGame
	done       chan bool
}

func New(owner *player.Player, secret string, minPlayers, maxPlayers int, allowBots bool, gameMode gm.Mode) (*Table, error) {
	game, err := newGame(gameMode)
	if err != nil {
		return nil, err
	}

	groups := safemap.New[int, *group.Group]()
	for i := 0; i < game.GetMaxGroups(); i++ {
		group := group.New(game.GetMinPlayersGroup(), game.GetMaxPlayersGroup())
		groups.Insert(i, group)
	}

	table := &Table{
		id:         utils.NewUUID(consts.TABLE_PREFIX_ID),
		owner:      owner.UUID,
		secret:     secret,
		minPlayers: minPlayers,
		maxPlayers: maxPlayers,
		allowBots:  allowBots,
		players:    safemap.New[string, *player.Player](),
		groups:     groups,
		winners:    make([]int, 0),
		state:      table.StateStart,
		gameMode:   gameMode,
		game:       game,
		done:       make(chan bool),
	}
	table.AddPlayer(owner, secret)

	return table, nil
}

func newGame(gameMode gm.Mode) (game.IGame, error) {
	var g game.IGame

	switch gameMode {
	case gm.ModeBlackJack:
		g = game.NewBlackJack()

	default:
		return nil, errors.ErrInvalidEmail
	}

	return g, nil
}

func (t *Table) AddPlayer(player *player.Player, secret string) error {
	// basic validation
	if player == nil {
		return errors.ErrNotFoundPlayer
	}
	if !t.CheckSecret(secret) {
		return errors.ErrInvalidPassword
	}
	// table validation
	if t.HasPlayer(player.UUID) {
		return errors.ErrExistsPlayer
	}

	if t.players.Size() >= t.maxPlayers {
		return errors.ErrMaxPlayers
	}
	// add player
	t.players.Insert(player.UUID, player)

	return nil
}

func (t *Table) CheckSecret(secret string) bool {
	return t.secret == "" || t.secret == secret
}

func (t *Table) DelPlayer(player string) error {
	// basic validation
	if player == "" {
		return errors.ErrNotFoundPlayer
	}
	// get player group
	groupID := 0
	if t.players.HasKey(player) {
		p, _ := t.players.GetOneValue(player)
		groupID = p.GroupID
	}
	// del player
	err := t.players.Delete(player)
	if err != nil {
		return errors.ErrNotFoundPlayer
	}
	// adjust owner
	if t.owner == player {
		t.owner = ""
	}
	// remove from group
	if groupID > 0 {
		t.groups.Delete(groupID)
		// TODO (@dev-rodrigobaliza) check if game is running and permits incomplet groups playing
	}
	// validate table rules
	if t.players.IsEmpty() {
		return errors.ErrEmptyTable
	}
	if t.state != table.StateStart && t.players.Size() < t.minPlayers {
		return errors.ErrMinPlayers
	}

	return nil
}

func (t *Table) GetAllowBots() bool {
	return t.allowBots
}

func (t *Table) GetID() string {
	return t.id
}

func (t *Table) GetMaxPlayers() int {
	return t.maxPlayers
}

func (t *Table) GetMinPlayers() int {
	return t.minPlayers
}

func (t *Table) GetOwner() string {
	return t.owner
}

func (t *Table) GetPlayers() []string {
	return t.players.GetAllKeys()
}

func (t *Table) GetPlayersCount() int {
	return t.players.Size()
}

func (t *Table) HasPlayer(playerID string) bool {
	return t.players.HasKey(playerID)
}

func (t *Table) IsPrivate() bool {
	return t.secret != ""
}

func (t *Table) Play() error {
	playersCount := t.players.Size()
	if playersCount < t.minPlayers {
		return errors.ErrNotEnoughPlayers
	}
	if t.state != table.StateStart {
		return errors.ErrStartedGame
	}
	// TODO (@dev-rodrigobaliza) more game conditions to start ???
	// TODO (@dev-rodrigobaliza) set players order
	// start the game loop
	go t.loop()

	return nil
}

func (t *Table) Stop(force bool) {
	// TODO (@dev-rodrigobaliza) mayde do a little housekeeping here
	if t.state == table.StatePlay {
		t.done <- true
	}
}

func (t *Table) ToResponse() *response.Table {
	pls := make([]*response.Player, 0)
	players := t.players.GetAllValues()
	for _, player := range players {
		pls = append(pls, player.ToResponse())
	}

	grs := make([]*response.Group, 0)
	groups := t.groups.GetAllValues()
	for _, group := range groups {
		if group.GetPlayersCount() > 0 {
			grs = append(grs, group.ToResponse())
		}
	}

	wis := make([]*response.Group, 0)
	for _, w := range t.winners {
		winner, err := t.groups.GetOneValue(w)
		if err == nil {
			wis = append(wis, winner.ToResponse())
		}
	}

	ta := response.NewTable(t.id, t.gameMode.String(), t.owner, t.IsPrivate(), pls, grs, wis)

	return ta
}

func (t *Table) loop() {
	err := t.game.Start()
	if err != nil {
		// TODO (@dev-rodrigobaliza) log this error
		return
	}

	ticker := time.NewTicker(time.Millisecond * consts.TABLE_INTERVAL_LOOP)

loop:
	for {
		select {
		case <-t.done:
			return

		case <-ticker.C:
			ok, err := t.game.Loop()
			if err != nil {
				// TODO (@dev-rodrigobaliza) log this error
				break loop
			}
			if !ok {
				// TODO (@dev-rodrigobaliza) game finished
				break loop
			}

		}
	}

	err = t.game.Stop()
	if err != nil {
		// TODO (@dev-rodrigobaliza) log this error
	}
}
