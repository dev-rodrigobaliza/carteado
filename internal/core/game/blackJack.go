package game

import (
	"github.com/dev-rodrigobaliza/carteado/consts"
	"github.com/dev-rodrigobaliza/carteado/domain/core/game"
	gru "github.com/dev-rodrigobaliza/carteado/domain/core/group"
	"github.com/dev-rodrigobaliza/carteado/errors"
	"github.com/dev-rodrigobaliza/carteado/internal/core/deck"
	"github.com/dev-rodrigobaliza/carteado/internal/core/group"
	"github.com/dev-rodrigobaliza/carteado/internal/core/player"
)

type BlackJack struct {
	maxGroups       int
	maxPlayersGroup int
	minPlayersGroup int
	groups          []*group.Group
	state           game.State
	round           uint64
	group           int
	dealerDeck      *deck.Deck
}

// This line is for get feedback in case we are not implementing the interface correctly
var _ IGame = (*BlackJack)(nil)

func NewBlackJack() *BlackJack {
	return &BlackJack{
		maxGroups:       consts.GAME_BLACKJACK_MAX_GROUPS,
		maxPlayersGroup: consts.GAME_BLACKJACK_MAX_PLAYERS_GROUP,
		minPlayersGroup: consts.GAME_BLACKJACK_MIN_PLAYERS_GROUP,
		state:           game.StateWaiting,
		round:           0,
	}
}

func (g *BlackJack) GetActiveGroup() *group.Group {
	return g.groups[g.group]
}

func (g *BlackJack) GetActivePlayer() (*player.Player, error) {
	grp := g.groups[g.group]
	player, err := grp.GetNextPlayer()
	if err != nil {
		return nil, err
	}

	return player, nil
}

func (g *BlackJack) GetActivePlayerDeck() (*deck.Deck, error) {
	grp := g.groups[g.group]
	deck, err := grp.GetNextDeck()
	if err != nil {
		return nil, err
	}

	return deck, nil
}

func (g *BlackJack) GetMaxGroups() int {
	return g.maxGroups
}

func (g *BlackJack) GetMaxPlayersGroup() int {
	return g.maxPlayersGroup
}

func (g *BlackJack) GetMinPlayersGroup() int {
	return g.minPlayersGroup
}

func (g *BlackJack) GetRound() uint64 {
	return g.round
}

func (g *BlackJack) GetState() game.State {
	return g.state
}

func (g *BlackJack) Loop() (bool, error) {
	switch g.state {
	case game.StateDealing:
		return g.deal()

	case game.StateBetting:
		return g.bet()

	case game.StatePlaying:
		return g.play()

	case game.StateWaiting:
		return g.wait()
	}

	return false, errors.ErrInvalidGameState
}

func (g *BlackJack) SetState(gameState game.State) {
	g.state = gameState
}

// TODO (@dev-rodrigobaliza) table calls start, here change the state to dealing, if not occurs error table will call loop until the game is finished
func (g *BlackJack) Start(groups []*group.Group) error {
	if g.state != game.StateWaiting {
		return errors.ErrInvalidGameState
	}
	groupsCount := len(groups)
	if g.minPlayersGroup < groupsCount {
		return errors.ErrNotEnoughPlayers
	}

	g.groups = groups
	err := g.initGroups()
	if err != nil {
		return err
	}
	err = g.initDecks()
	if err != nil {
		return err
	}

	g.state = game.StatePlaying
	g.round = 0
	g.group = 0

	return nil
}

func (g *BlackJack) Stop() error {
	if g.state != game.StatePlaying {
		return errors.ErrInvalidGameState
	}
	// TODO (@dev-rodrigobaliza) finish here all the loose ends

	return nil
}

func (g *BlackJack) bet() (bool, error) {
	return false, errors.ErrNotImplemented
}

func (g *BlackJack) checkHighScore() (bool, error) {
	highPlayers := make([]string, 0)
	highScore := 0

	for _, group := range g.groups {
		if group.GetPlayersCount() > 0 {
			if group.GetGroupScore() > highScore {
				highScore = group.GetGroupScore()
				highPlayers = append(make([]string, 0), group.GetPlayers()[0])
			} else if group.GetGroupScore() == highScore {
				highPlayers = append(highPlayers, group.GetPlayers()[0])
			}
		}
	}

	for _, group := range g.groups {
		if group.GetPlayersCount() > 0 {
			group.NextPlayer = 0
			player, err := group.GetNextPlayer()
			if err != nil {
				continue
			}

			if group.GetGroupScore() <= consts.GAME_BLACKJACK_WINNING_SCORE {
				for i := 0; i < len(highPlayers); i++ {
					if player.UUID == highPlayers[i] {
						if len(highPlayers) >= 1 {
							player.Action = "win"
						}
					} else {
						player.Action = "loose"
					}
				}
			} else {
				player.Action = "loose"
			}
		}
	}

	return true, nil
}

func (g *BlackJack) checkWinner(group *group.Group, player *player.Player, state gru.State) (bool, error) {
	player.Action = ""

	points := group.GetGroupScore()
	if points < consts.GAME_BLACKJACK_WINNING_SCORE {
		group.State = state
		return false, nil
	} else if points == consts.GAME_BLACKJACK_WINNING_SCORE {
		group.State = gru.StateFinish
		return true, errors.ErrSendPlayerWin
	}

	group.State = gru.StateFinish
	return false, errors.ErrSendPlayerLoose
}

func (g *BlackJack) deal() (bool, error) {
	return false, errors.ErrNotImplemented
}

func (g *BlackJack) giveCard(grp *group.Group) error {
	c, err := g.dealerDeck.GetFirstCard(true)
	if err != nil {
		return errors.ErrEmptyDeck
	}
	c.Visible = true
	grp.AddCard("", c)

	return nil
}

func (g *BlackJack) initDecks() error {
	// create dealer deck
	cardSet := []string{
		"1h", "2h", "3h", "4h", "5h", "6h", "7h", "jh", "qh", "kh",
		"1d", "2d", "3d", "4d", "5d", "6d", "7d", "jd", "qd", "kd",
		"1c", "2c", "3c", "4c", "5c", "6c", "7c", "jc", "qc", "kc",
		"1s", "2s", "3s", "4s", "5s", "6s", "7s", "js", "qs", "ks",
	}
	valueSet := []int{
		1, 2, 3, 4, 5, 6, 7, 10, 10, 10,
		1, 2, 3, 4, 5, 6, 7, 10, 10, 10,
		1, 2, 3, 4, 5, 6, 7, 10, 10, 10,
		1, 2, 3, 4, 5, 6, 7, 10, 10, 10,
	}
	var cards []string
	var values []int
	// one card set for each four groups
	for i := 0; i <= (len(g.groups) / 4); i++ {
		cards = append(cards, cardSet...)
		values = append(values, valueSet...)
	}
	d, err := deck.NewCustom(cards, values, false)
	if err != nil {
		return err
	}
	d.Shuffle(1)
	g.dealerDeck = d
	// initially give two cards visible for each player
	for _, gr := range g.groups {
		for i := 0; i < 2; i++ {
			g.giveCard(gr)
		}
	}

	return nil
}

func (g *BlackJack) initGroups() error {
	// basic validation
	if len(g.groups) == 0 {
		return errors.ErrEmptyTable
	}
	// set initial group state for each group
	for _, gr := range g.groups {
		gr.State = gru.StateAction
		gr.NextPlayer = 0
	}

	return nil
}

func (g *BlackJack) play() (bool, error) {
	// basic validation
	groupsCount := len(g.groups)
	if groupsCount == 0 {
		return true, errors.ErrEmptyTable
	}
	// check if all groups has been played
	if groupsCount == g.group {
		return g.checkHighScore()
	}
	// check if actual group has player
	group := g.groups[g.group]
	if group.GetPlayersCount() == 0 {
		// empty group, remove from the game
		g.groups = append(g.groups[:g.group], g.groups[g.group+1:]...)
		g.group++

		return false, nil
	}
	// get the active player
	player, err := group.GetNextPlayer()
	if err != nil {
		return g.checkHighScore()
	}
	// finite state machine for the game rules
	switch group.State {
	case gru.StateReady:
		g.giveCard(group)
		group.State = gru.StateAction
		return false, errors.ErrSendPlayerCards

	case gru.StateAction:
		switch player.Action {
		case "continue":
			return g.checkWinner(group, player, gru.StateReady)

		case "stop":
			return g.checkWinner(group, player, gru.StateFinish)

		case "":
			return false, nil

		default:
			return false, errors.ErrInvalidAction // send invalid action to player
		}

	case gru.StateStop:
		return g.checkWinner(group, player, gru.StateFinish)

	case gru.StateFinish:
		g.group++
		return false, nil
	}

	return false, errors.ErrInvalidAction
}

func (g *BlackJack) wait() (bool, error) {
	return false, errors.ErrNotImplemented
}
