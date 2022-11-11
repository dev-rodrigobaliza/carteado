# gameState

it's the game state machine (players and table management)

## start

* game room
* game table open (set partners, set position, set first player)
* chat
* can finish gameState (force - not admin only empty)
* can play gameState

## play

* table room
* game table lock
* chat
* can finish gameState (force - not admin only empty)
* runs the table state machine

## finish

* game room
* game table lock
* chat
* when last player leaves remove the game
