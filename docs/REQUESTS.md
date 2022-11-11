# websocket requests

this are examples of how to consume the carteado websocket api

## auth login (admin)

``` json
{
 "request_id": 1,
 "service": "auth",
 "resource": "login",
 "data": {
  "token": "v4.local.KDBFxxUYtml5Afhb_amr0Yr48D5W5kHNX_IemOwIRp-0MyzeQvPXTNy5WVTsIfG5E-H8Y18Y9Itq6Bcglj22mqfss-iNE4L4semjZZUox4tpgQCJpTsSflxkGzaprUPkzj4ZGBybMNL5evHmP6xE0A_mAGB7"
 }
}
```

## auth login (user)

``` json
{
 "request_id": 1,
 "service": "auth",
 "resource": "login",
 "data": {
  "token": "v4.local.cEPMuKYqQlP7YfAO6Gey8ilVVWYHDYyC3OTG3LFEbw8JcGKg-p28KcJ2UbRKZEljtArFNepFT_DLIf-c2BXSiziJbJBxnGhA5Mcq_bYEPBo3apVghqHw3Zfk0HIw-3kIcKLy4cUFPGrEP54DIIx9TAwIHQ0"
 }
}
```

## admin status (admin only)

``` json
{
 "request_id": 1,
 "service": "admin",
 "resource": "status"
}
```

## table create (public)

``` json
{
 "request_id": 1,
 "service": "table",
 "resource": "create",
 "data": {
  "game_mode": "blackjack",
  "min_players": 1,
  "max_players": 5,
  "allow_bots": true
 }
}
```

## table create (private)

``` json
{
 "request_id": 1,
 "service": "table",
 "resource": "create",
 "data": {
  "game_mode": "blackjack",
  "min_players": 1,
  "max_players": 5,
  "allow_bots": true,
  "secret": "secret"
 }
}
```

## table enter (public)

``` json
{
 "request_id": 1,
 "service": "table",
 "resource": "enter",
 "data": {
  "table_id": "tid-5fe9db7f600000a"
 }
}
```

## table enter (private)

``` json
{
 "request_id": 1,
 "service": "table",
 "resource": "enter",
 "data": {
  "table_id": "tid-5fe9db7f600000a",
  "secret": "secret"
 }
}
```

## table leave

``` json
{
 "request_id": 1,
 "service": "table",
 "resource": "leave",
 "data": {
  "table_id": "tid-5fe9db7f600000a"
 }
}
```

## table remove

``` json
{
 "request_id": 1,
 "service": "table",
 "resource": "remove",
 "data": {
  "table_id": "tid-5fe9db7f600000a"
 }
}
```

## table status

``` json
{
 "request_id": 1,
 "service": "table",
 "resource": "status",
 "data": {
  "table_id": "tid-5feb5ab5100000a"
 }
}
```

## table group (enter)

``` json
{
 "request_id": 1,
 "service": "table",
 "resource": "group",
 "data": {
  "table_id": "tid-5ff3b8ffa00000a",
  "group_id": 1,
  "action": "enter"
 }
}
```

## table group (leave)

``` json
{
 "request_id": 1,
 "service": "table",
 "resource": "group",
 "data": {
  "table_id": "tid-5fe9db7f600000a",
  "group_id": 1,
  "action": "leave"
 }
}
```

## table group (status)

``` json
{
 "request_id": 1,
 "service": "table",
 "resource": "group",
 "data": {
  "table_id": "tid-5feb41b7200000a",
  "group_id": 1,
  "action": "status"
 }
}
```

## table group (bot)

``` json
{
 "request_id": 1,
 "service": "table",
 "resource": "group",
 "data": {
  "table_id": "tid-5ff3b8ffa00000a",
  "action": "bot",
  "quantity": 1
 }
}
```

## table game (start)

``` json
{
 "request_id": 1,
 "service": "table",
 "resource": "game",
 "data": {
  "table_id": "tid-5ff3b8ffa00000a",
  "action": "start"
 }
}
```

## table game (continue)

``` json
{
 "request_id": 1,
 "service": "table",
 "resource": "game",
 "data": {
  "table_id": "tid-5ff3b8ffa00000a",
  "action": "continue"
 }
}
```

## table game (discontinue)

``` json
{
 "request_id": 1,
 "service": "table",
 "resource": "game",
 "data": {
  "table_id": "tid-5ff3b8ffa00000a",
  "action": "discontinue"
 }
}
```

## table game (status)

``` json
{
 "request_id": 1,
 "service": "table",
 "resource": "game",
 "data": {
  "table_id": "tid-5ff3b8ffa00000a",
  "action": "status"
 }
}
```
