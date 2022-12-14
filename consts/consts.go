package consts

const (
	APP_CONFIG_FILENAME = "config.json"

	BOT_PREFIX_ID = "bid-"

	CARD_PREFIX_ID = "cid-"

	DECK_DEFAULT_CARD_SIZE = 52

	GAME_BLACKJACK_MAX_GROUPS        = 16
	GAME_BLACKJACK_MAX_PLAYERS_GROUP = 1
	GAME_BLACKJACK_MIN_PLAYERS_GROUP = 1
	GAME_BLACKJACK_WINNING_SCORE     = 21

	SALOON_MAX_PLAYERS        = 1024
	SALOON_MESSAGE_STACK_SIZE = SALOON_MAX_PLAYERS * PLAYER_MESSAGE_STACK_SIZE

	PLAYER_MESSAGE_STACK_SIZE = 64
	PLAYER_PREFIX_ID          = "pid-"

	TABLE_INTERVAL_LOOP                = 500 // milliseconds
	TABLE_INTERVAL_START_GAME_RESPONSE = 500 // milliseconds
	TABLE_MAX_PLAYERS                  = 16
	TABLE_PREFIX_ID                    = "tid-"
)
