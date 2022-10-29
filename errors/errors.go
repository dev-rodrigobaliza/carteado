package errors

import "errors"

var (
	ErrEmptyTable = errors.New("table is empty")

	ErrExistsPlayer = errors.New("player already exists")

	ErrFailedOpenFileConfig  = errors.New("failed to open config file")
	ErrFailedParseFileConfig = errors.New("failed to parse config file")

	ErrInvalidAccessToken  = errors.New("access token invalid")
	ErrInvalidCard         = errors.New("card invalid")
	ErrInvalidCardDeck     = errors.New("card deck invalid")
	ErrInvalidCardFace     = errors.New("card face invalid")
	ErrInvalidCardSuit     = errors.New("card suit invalid")
	ErrInvalidCardValue    = errors.New("card value invalid")
	ErrInvalidDatabaseType = errors.New("database type invalid")
	ErrInvalidEmail        = errors.New("email invalid")
	ErrInvalidGameMode     = errors.New("game mode invalid")
	ErrInvalidGameState    = errors.New("game state invalid")
	ErrInvalidIP           = errors.New("ip address invalid")
	ErrInvalidLogin        = errors.New("login invalid")
	ErrInvalidName         = errors.New("name invalid")
	ErrInvalidPassword     = errors.New("password invalid")
	ErrInvalidUser         = errors.New("user invalid")
	ErrInvalidUserID       = errors.New("user id invalid")
	ErrInvalidUserIDEmail  = errors.New("user id and/or email invalid")

	ErrMaxPlayers = errors.New("players count over maximum")
	ErrMinPlayers = errors.New("players count above minimum")

	ErrNotEnoughPlayers = errors.New("players not enough")

	ErrNotFoundCard   = errors.New("card not found")
	ErrNotFoundGroup  = errors.New("group not found")
	ErrNotFoundPlayer = errors.New("player not found")
	ErrNotFoundUser   = errors.New("user not found")

	ErrNotImplemented = errors.New("not implemented")

	ErrStartedGame = errors.New("game already started")

	ErrUnauthorized = errors.New("unauthorized")
)
