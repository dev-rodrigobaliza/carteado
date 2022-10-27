package errors

import "errors"

var (
	ErrExistsPlayer = errors.New("player already exists")

	ErrMaxPlayers = errors.New("players count over maximum")
	ErrMinPlayers = errors.New("players count above minimum")

	ErrNotFoundPlayer = errors.New("player not found")
	ErrNotFoundUser   = errors.New("user not found")

	ErrNotImplemented = errors.New("not implemented")

	ErrInvalidAccessToken  = errors.New("access token invalid")
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
)
