package table

import (
	"fmt"
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
	secret     string
	createBy   string
	startedBy  string
	minPlayers int
	maxPlayers int
	allowBots  bool
	players    *safemap.SafeMap[string, *player.Player]
	groups     *safemap.SafeMap[int, *group.Group]
	winners    []int
	state      table.State
	gameMode   gm.Mode
	game       game.IGame
	createdAt  time.Time
	startedAt  time.Time
	done       chan bool
}

func New(owner *player.Player, secret string, minPlayers, maxPlayers int, allowBots bool, gameMode gm.Mode) (*Table, error) {
	game, err := newGame(gameMode)
	if err != nil {
		return nil, err
	}

	groups := safemap.New[int, *group.Group]()
	for i := 1; i <= game.GetMaxGroups(); i++ {
		group := group.New(i, game.GetMinPlayersGroup(), game.GetMaxPlayersGroup())
		groups.Insert(i, group)
	}

	table := &Table{
		id:         utils.NewUUID(consts.TABLE_PREFIX_ID),
		createBy:   owner.UUID,
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
		createdAt:  time.Now(),
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
	player.TableID = t.id

	return nil
}

func (t *Table) AddGroupPlayer(group int, player *player.Player) error {
	// basic validation
	if group == 0 || !t.groups.HasKey(group) {
		return errors.ErrNotFoundGroup
	}
	if player == nil {
		return errors.ErrNotFoundPlayer
	}
	if t.state != table.StateStart {
		return errors.ErrStartedGame
	}
	// remove from previous group
	if player.GroupID > 0 {
		g, err := t.groups.GetOneValue(player.GroupID, false)
		if err == nil || err == errors.ErrMinPlayers {
			// TODO (@dev-rodrigobaliza) should stop the game ???
			g.DelPlayer(player.UUID)
		}

	}
	// add player to group
	g, _ := t.groups.GetOneValue(group, false)
	return g.AddPlayer(player)
}

func (t *Table) CheckSecret(secret string) bool {
	return t.secret == "" || t.secret == secret
}

func (t *Table) DelGroupPlayer(group int, player string) error {
	// basic validation
	if group == 0 || !t.groups.HasKey(group) {
		return errors.ErrNotFoundGroup
	}
	if player == "" {
		return errors.ErrNotFoundPlayer
	}
	// remove player from group
	g, _ := t.groups.GetOneValue(group, false)

	return g.DelPlayer(player)
}

func (t *Table) DelPlayer(player string) error {
	// basic validation
	if player == "" {
		return errors.ErrNotFoundPlayer
	}
	// get player group
	groupID := 0
	if t.players.HasKey(player) {
		p, _ := t.players.GetOneValue(player, false)
		groupID = p.GroupID
	}
	// del player
	err := t.players.DeleteKey(player)
	if err != nil {
		return errors.ErrNotFoundPlayer
	}
	// adjust owner
	if t.createBy == player {
		t.createBy = ""
	}
	// remove from group
	if groupID > 0 {
		t.groups.DeleteKey(groupID)
		// TODO (@dev-rodrigobaliza) check if game is running and permits incomplete groups playing
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

func (t *Table) GetGroup(groupID int) (*group.Group, error) {
	// basic validation
	if groupID == 0 || !t.groups.HasKey(groupID) {
		return nil, errors.ErrNotFoundGroup
	}
	group, _ := t.groups.GetOneValue(groupID, false)

	return group, nil
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
	return t.createBy
}

func (t *Table) GetPlayers() []string {
	return t.players.GetAllKeys()
}

func (t *Table) GetPlayersCount() int {
	return t.players.Size()
}

func (t *Table) HasGroup(groupID int) bool {
	return t.groups.HasKey(groupID)
}

func (t *Table) HasPlayer(playerID string) bool {
	return t.players.HasKey(playerID)
}

func (t *Table) IsPrivate() bool {
	return t.secret != ""
}

func (t *Table) Start() error {
	playersCount := t.players.Size()
	if playersCount < t.minPlayers {
		return errors.ErrNotEnoughPlayers
	}
	if t.state != table.StateStart {
		return errors.ErrStartedGame
	}
	// TODO (@dev-rodrigobaliza) more game conditions to start ???
	// TODO (@dev-rodrigobaliza) set players order

	t.startedAt = time.Now()
	t.state = table.StatePlay

	// send the news to all players
	response := make(map[string]interface{})
	response["table_game"] = "game started"
	t.sendAllPlayersResponse("info", "game status", response)

	// start the game loop
	go t.loop(playersCount)
	
	return nil
}

func (t *Table) Stop(force bool) {
	// TODO (@dev-rodrigobaliza) mayde do a little housekeeping here
	if t.state == table.StatePlay {
		t.done <- true
	}
}

func (t *Table) ToResponse() *response.Table {
	spectators := make([]*response.Player, 0)
	players := t.players.GetAllValues()
	for _, player := range players {
		if player.GroupID == 0 {
			spectators = append(spectators, player.ToResponse(false))
		}
	}

	groups := make([]*response.Group, 0)
	grs := t.groups.GetAllValues()
	for _, group := range grs {
		if group.GetPlayersCount() > 0 {
			groups = append(groups, group.ToResponse(false))
		}
	}

	winners := make([]*response.Group, 0)
	for _, w := range t.winners {
		winner, err := t.groups.GetOneValue(w, false)
		if err == nil {
			winners = append(winners, winner.ToResponse(false))
		}
	}

	created := fmt.Sprintf("%d", t.createdAt.UnixMilli())
	var started string
	if t.startedAt.After(t.createdAt) {
		started = fmt.Sprintf("%d", t.startedAt.UnixMilli())
	}
	ta := response.NewTable(t.id, t.gameMode.String(), t.createBy, t.startedBy, created, started, t.IsPrivate(), t.players.Size(), spectators, groups, winners)

	return ta
}

func (t *Table) loop(playersCount int) {
	// pass only groups with players
	groups := make([]*group.Group, 0)
	grs := t.groups.GetAllValues()
	for _, g := range grs {
		if g.GetPlayersCount() > 0 {
			groups = append(groups, g)
		}
	}

	err := t.game.Start(groups)
	if err != nil {
		// TODO (@dev-rodrigobaliza) log this error
		return
	}

	// send cards to all players
	for _, g := range grs {
		for _, p := range g.GetPlayers() {
			t.sendPlayerResponse(p, "cards", g)
		}
	}

	ticker := time.NewTicker(time.Millisecond * consts.TABLE_INTERVAL_LOOP)

loop:
	for {
		select {
		case <-t.done:
			return

		case <-ticker.C:
			ticker.Stop()
			finished, err := t.game.Loop()
			if err != nil {
				p, _ := t.game.GetActivePlayer()
				g := t.game.GetActiveGroup()

				switch err {
				case errors.ErrSendPlayerCards:
					t.sendPlayerResponse(p.UUID, "cards", g)

				case errors.ErrSendPlayerLoose:
					t.sendPlayerResponse(p.UUID, "loose", g)

				case errors.ErrSendPlayerWin:
					t.sendPlayerResponse(p.UUID, "win", g)
				}
			}

			if finished {
				// TODO (@dev-rodrigobaliza) game finished
				break loop
			}

			ticker.Reset(time.Millisecond * consts.TABLE_INTERVAL_LOOP)
		}
	}

	ticker.Stop()
	t.state = table.StateFinish

	if err == nil {
		grs := t.groups.GetAllValues()
		for _, g := range grs {
			p, err := g.GetNextPlayer()
			if err != nil {
				continue
			}

			t.sendPlayerResponse(p.UUID, p.Action, g)
		}
	}

	err = t.game.Stop()
	if err != nil {
		// TODO (@dev-rodrigobaliza) log this error
	}
}

func (t *Table) sendAllPlayersResponse(status, message string, response map[string]interface{}) {

}

func (t *Table) sendPlayerResponse(playerID, status string, group *group.Group) {
	pl, err := t.players.GetOneValue(playerID, false)
	if err != nil {
		// TODO (@dev-rodrigobaliza) log this error
		return
	}

	var message string
	response := make(map[string]interface{})

	switch status {
	case "cards":
		d, err := group.GetPlayerDeck(playerID)
		if err != nil {
			// TODO (@dev-rodrigobaliza) how do we handle this error?
			return
		}

		response["table"] = d.ToResponse(t.id, true)
		message = "player cards"

	case "loose", "win":
		s := table.Status{
			ID:     t.id,
			Status: status,
		}

		response["table"] = s
		message = "game status"

	default:
		return
	}

	if t.state == table.StateFinish {
		response["table_game"] = "game over"
	}

	// send response
	pl.SendResponse(nil, "info", message, response)

	// if status == "cards" {
	// 	response = make(map[string]interface{})
	// 	response["table"] = "waiting for player action"
	// 	pl.SendResponse(nil, "info", "game status", response)
	// }
}
