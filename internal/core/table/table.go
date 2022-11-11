package table

import (
	"fmt"
	"log"
	"time"

	"github.com/dev-rodrigobaliza/carteado/consts"
	gm "github.com/dev-rodrigobaliza/carteado/domain/core/game"
	"github.com/dev-rodrigobaliza/carteado/domain/core/table"
	"github.com/dev-rodrigobaliza/carteado/domain/response"
	"github.com/dev-rodrigobaliza/carteado/errors"
	"github.com/dev-rodrigobaliza/carteado/internal/core/deck"
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
	game, err := game.New(gameMode)
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

func (t *Table) AddGroupPlayer(group int, player *player.Player) error {
	// basic validation
	if group == 0 || !t.groups.HasKey(group) {
		return errors.ErrNotFoundGroup
	}
	if player == nil {
		return errors.ErrNotFoundPlayer
	}
	if t.state != table.StateStart {
		return errors.ErrGameStart
	}
	if player.GroupID == group {
		return errors.ErrExistsPlayerGroup
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

func (t *Table) AddGroupPlayerBot(player *player.Player) error {
	// basic validation
	if player == nil {
		return errors.ErrNotFoundPlayer
	}
	if t.state != table.StateStart {
		return errors.ErrGameStart
	}
	// find next group available
	grs := t.groups.GetAllValues()
	for _, g := range grs {
		if g.GetPlayersCount() == 0 {
			// add player to group
			return g.AddPlayer(player)
		}
	}

	return errors.ErrNotFoundAvailableGroup
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
	// add player
	t.players.Insert(player.UUID, player)
	player.TableID = t.id

	return nil
}

func (t *Table) AddPlayerBot(player *player.Player) error {
	// basic validation
	if player == nil {
		return errors.ErrNotFoundPlayer
	}
	// table validation
	if !t.allowBots {
		return errors.ErrNotAllowedBotPlayers
	}
	if t.HasPlayer(player.UUID) {
		return errors.ErrExistsPlayer
	}
	if t.GetGroupPlayersCount() >= t.maxPlayers {
		return errors.ErrMaxPlayers
	}
	// add player
	t.players.Insert(player.UUID, player)
	player.TableID = t.id

	return t.AddGroupPlayerBot(player)
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
	if t.state != table.StateStart && t.GetGroupPlayersCount() < t.minPlayers {
		return errors.ErrMinPlayers
	}

	return nil
}

func (t *Table) GetAllowBots() bool {
	return t.allowBots
}

func (t *Table) GetGameStatus() *response.Game {
	return t.game.Response()
}

func (t *Table) GetGroup(groupID int) (*group.Group, error) {
	// basic validation
	if groupID == 0 || !t.groups.HasKey(groupID) {
		return nil, errors.ErrNotFoundGroup
	}
	group, _ := t.groups.GetOneValue(groupID, false)

	return group, nil
}

func (t *Table) GetGroupPlayersCount() int {
	var total int

	grs := t.groups.GetAllValues()
	for _, g := range grs {
		total += g.GetPlayersCount()
	}

	return total
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
	playersCount := t.GetGroupPlayersCount()
	if playersCount < t.minPlayers {
		return errors.ErrNotEnoughPlayers
	}
	if t.state != table.StateStart {
		return errors.ErrGameStart
	}
	// TODO (@dev-rodrigobaliza) more game conditions to start ???
	// TODO (@dev-rodrigobaliza) set players order

	t.startedAt = time.Now()
	t.state = table.StatePlay

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
		return err
	}

	go func() {
		time.Sleep(time.Millisecond * consts.TABLE_INTERVAL_START_GAME_RESPONSE)
		// send the news to all players
		t.sendResponseAllPlayers("start", "game status")
		// send cards to all players
		t.sendResponseInGamePlayers("card", "")
		// start the game loop
		t.loop(playersCount)
	}()

	return nil
}

func (t *Table) Stop(force bool) error {
	// TODO (@dev-rodrigobaliza) mayde do a little housekeeping here
	if t.state == table.StatePlay {
		t.done <- true

		return nil
	}

	return errors.ErrGameStop
}

func (t *Table) Response(admin bool) *response.Table {
	spectators := make([]*response.Player, 0)
	players := t.players.GetAllValues()
	for _, player := range players {
		if player.GroupID == 0 {
			spectators = append(spectators, player.Response(admin, false))
		}
	}

	groups := make([]*response.Group, 0)
	grs := t.groups.GetAllValues()
	for _, group := range grs {
		if group.GetPlayersCount() > 0 {
			groups = append(groups, group.Response(admin))
		}
	}

	winners := make([]*response.Group, 0)
	for _, w := range t.winners {
		winner, err := t.groups.GetOneValue(w, false)
		if err == nil {
			winners = append(winners, winner.Response(admin))
		}
	}

	game := t.game.Response()

	created := fmt.Sprintf("%d", t.createdAt.UnixMilli())
	var started string
	if t.startedAt.After(t.createdAt) {
		started = fmt.Sprintf("%d", t.startedAt.UnixMilli())
	}
	ta := response.NewTable(t.id, t.gameMode.String(), t.createBy, t.startedBy, created, started, t.IsPrivate(), t.GetGroupPlayersCount(), spectators, groups, winners, game)

	return ta
}

func (t *Table) getInGamePlayers() []*player.Player {
	playersInGame := make([]*player.Player, 0)
	grs := t.groups.GetAllValues()
	for _, g := range grs {
		p, err := g.GetNextPlayer()
		if err != nil {
			continue
		}

		playersInGame = append(playersInGame, p)
	}

	return playersInGame
}

func (t *Table) getPlayerDeck(player *player.Player) (*deck.Deck, error) {
	// get group from player
	g, err := t.groups.GetOneValue(player.GroupID, false)
	if err != nil {
		// TODO (@dev-rodrigobaliza) how do we handle this error?
		return nil, err
	}
	// get deck from group
	d, err := g.GetPlayerDeck(player.UUID)
	if err != nil {
		// TODO (@dev-rodrigobaliza) how do we handle this error?
		return nil, err
	}

	return d, nil
}

func (t *Table) loop(playersCount int) {
	// init the main loop ticker
	ticker := time.NewTicker(time.Millisecond * consts.TABLE_INTERVAL_LOOP)

loop:
	for {
		select {
		case <-t.done:
			return

		case <-ticker.C:
			ticker.Stop()
			playerID, finished, err := t.game.Loop()
			if err != nil {
				switch err {
				case errors.ErrSendBotPlayerContinue:
					t.sendResponseAllPlayers("bot-continue", "")

				case errors.ErrSendBotPlayerDiscontinue:
					t.sendResponseAllPlayers("bot-discontinue", "")

				case errors.ErrSendBotPlayerGotCard:
					t.sendResponseAllPlayers("bot-got-card", "")

				case errors.ErrSendPlayerAction:
					t.sendPlayerResponse(playerID, "action", "")

				case errors.ErrSendPlayerCards:
					t.sendPlayerResponse(playerID, "card", "")

				case errors.ErrSendPlayerGotCard:
					t.sendResponseAllPlayers("got-card", "")

				case errors.ErrSendPlayerLoose:
					t.sendResponseAllPlayers("discontinue", "")
					t.sendPlayerResponse(playerID, "loose", "")

				case errors.ErrSendPlayerWin:
					t.sendPlayerResponse(playerID, "win", "")
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
	t.sendResponseAllPlayers("stop", "")

	err := t.game.Stop()
	if err != nil {
		log.Printf("error stopping game: %s", err.Error())
	}
}

func (t *Table) sendResponseAllPlayers(status, message string) {
	// message for all players, inclusive spectators
	pls := t.players.GetAllValues()
	for _, p := range pls {
		if p.IsBot {
			continue
		}

		t.sendPlayerResponse(p.UUID, status, message)
	}
}

func (t *Table) sendResponseInGamePlayers(status, message string) {
	// message only for players in game
	inGamePlayers := t.getInGamePlayers()
	for _, p := range inGamePlayers {
		p.Action = status
		if p.IsBot {
			continue
		}

		t.sendPlayerResponse(p.UUID, status, message)
	}
}

func (t *Table) sendPlayerResponse(playerID, status, message string) {
	response := make(map[string]interface{})
	// get player
	p, err := t.players.GetOneValue(playerID, false)
	if err != nil {
		// TODO (@dev-rodrigobaliza) log this error
		return
	}
	if p.IsBot {
		return
	}
	// make response
	switch status {
	case "action":
		message = "waiting for player action"

	case "card":
		message = "player cards"

	case "start":
		response["status"] = "game started"

	case "stop":
		response["status"] = "game over"

	case "loose", "win":
		message = "game status"
		response["status"] = status

	case "continue", "discontinue", "got-card", "bot-continue", "bot-discontinue", "bot-got-card":
		message = "game status"
		response["status"] = status
		ap, _ := t.game.GetActivePlayer()
		response["player"] = ap.Response(p.User.IsAdmin, false)
		if p.User.IsAdmin {
			d, _ := t.getPlayerDeck(ap)
			response["deck"] = d.Response(t.id, false)
		}

	default:
		return
	}
	// get player deck (if pertinent)
	_, ok := response["deck"]
	if !ok {
		d, err := t.getPlayerDeck(p)
		if err == nil {
			response["deck"] = d.Response(t.id, !p.User.IsAdmin)
		}
	}
	// send response
	p.SendResponse(nil, "info", message, response)
}
