package table

import (
	"fmt"
	"time"

	"github.com/dev-rodrigobaliza/carteado/consts"
	gm "github.com/dev-rodrigobaliza/carteado/domain/core/game"
	"github.com/dev-rodrigobaliza/carteado/domain/core/table"
	"github.com/dev-rodrigobaliza/carteado/errors"
	"github.com/dev-rodrigobaliza/carteado/internal/core/game"
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
	players    *safemap.SafeMap[string, *Player]
	winners    []string
	state      table.State
	gameMode   gm.Mode
	game       game.IGame
	done       chan bool
}

func NewTable(owner *Player, secret string, minPlayers, maxPlayers int, allowBots bool, gameMode gm.Mode) (*Table, error) {
	game, err := newGame(gameMode)
	if err != nil {
		return nil, err
	}

	id := "gid-" + utils.NewUUID()

	table := &Table{
		id:         id,
		owner:      owner.uuid,
		secret:     secret,
		minPlayers: minPlayers,
		maxPlayers: maxPlayers,
		allowBots:  allowBots,
		players:    safemap.New[string, *Player](),
		winners:    make([]string, 0),
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

func (t *Table) AddPlayer(player *Player, secret string) error {
	// basic validation
	if player == nil {
		return errors.ErrNotFoundPlayer
	}
	if !t.CheckSecret(secret) {
		return errors.ErrInvalidPassword
	}
	// table validation
	if t.HasPlayer(player.uuid) {
		return errors.ErrExistsPlayer
	}

	if t.players.Size() >= t.maxPlayers {
		return errors.ErrMaxPlayers
	}
	// add player
	t.players.Insert(player.uuid, player)

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
	// del player
	err := t.players.Delete(player)
	if err != nil {
		return errors.ErrNotFoundPlayer
	}
	// adjust owner
	if t.owner == player {
		t.owner = ""
	}
	// validate table rules
	if t.players.Size() < t.minPlayers {
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

func (t *Table) GetPlayers() []string {
	return t.players.GetAllKeys()
}

func (t *Table) GetPlayersCount() int {
	return t.players.Size()
}

func (t *Table) GetStatus() *table.Status {
	return &table.Status{
		ID:          t.id,
		Owner:       t.owner,
		Winners:     t.winners,
		State:       t.state,
		PlayerCount: t.players.Size(),
		MinPlayers:  t.minPlayers,
		MaxPlayers:  t.maxPlayers,
		AllowBots:   t.allowBots,
		Private:     t.IsPrivate(),
		GameMode:    t.gameMode,
		GameState:   t.game.GetState(),
		GameRound:   t.game.GetRound(),
	}
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
		return fmt.Errorf("%d players not enough, expect min %d players", playersCount, t.minPlayers)
	}
	if t.state != table.StateStart {
		return fmt.Errorf("game already started")
	}
	// TODO (@dev-rodrigobaliza) more game conditions to start ???
	// TODO (@dev-rodrigobaliza) set players order
	// start the game loop
	go t.loop()

	return nil
}

func (t *Table) Stop(force bool) {
	t.done <- true
}

func (t *Table) loop() {
	err := t.game.Start()
	if err != nil {
		// TODO (@dev-rodrigobaliza) log this error
		return
	}

	ticker := time.NewTicker(time.Millisecond * consts.LOOP_INTERVAL)

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
