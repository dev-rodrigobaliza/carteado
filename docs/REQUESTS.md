auth login (admin):
{
	"request_id": 1,
	"service": "auth",
	"resource": "login",
	"data": {
		"token": "v4.local.S0SQOjr_sCExqvT8IVhMTLaXtl-5GUFUc14bTFFKP3CperTtffwchOMJUcsijZZtgkRZJlGqmRdEE3Ux-U1i2Iz8_91dvzG3z-wQPjJeGBmHIihI2eWdR-MP8RRoEk_3bKWMDcpICYDtSkb5MRcs7Rp3qxpC"
	}
}

auth login (user):
{
	"request_id": 1,
	"service": "auth",
	"resource": "login",
	"data": {
		"token": "v4.local.cEPMuKYqQlP7YfAO6Gey8ilVVWYHDYyC3OTG3LFEbw8JcGKg-p28KcJ2UbRKZEljtArFNepFT_DLIf-c2BXSiziJbJBxnGhA5Mcq_bYEPBo3apVghqHw3Zfk0HIw-3kIcKLy4cUFPGrEP54DIIx9TAwIHQ0"
	}
}

admin status (admin only):
{
	"request_id": 1,
	"service": "admin",
	"resource": "status",
	"data": {
		"authenticated_only": false
	}
}

table create (public):
{
	"request_id": 1,
	"service": "table",
	"resource": "create",
	"data": {
		"game_mode": "blackjack",
		"min_players": 5,
		"max_players": 5,
		"allow_bots": true
	}
}

table create (private):
{
	"request_id": 1,
	"service": "table",
	"resource": "create",
	"data": {
		"game_mode": "blackjack",
		"min_players": 5,
		"max_players": 5,
		"allow_bots": true,
		"secret": "secret"
	}
}

table enter (public):
{
	"request_id": 1,
	"service": "table",
	"resource": "enter",
	"data": {
		"table_id": "tid-5fe9db7f600000a"
	}
}

table enter (private):
{
	"request_id": 1,
	"service": "table",
	"resource": "enter",
	"data": {
		"table_id": "tid-5fe9db7f600000a",
		"secret": "secret"
	}
}

table leave:
{
	"request_id": 1,
	"service": "table",
	"resource": "leave",
	"data": {
		"table_id": "tid-5fe9db7f600000a"
	}
}

table remove:
{
	"request_id": 1,
	"service": "table",
	"resource": "remove",
	"data": {
		"table_id": "tid-5fe9db7f600000a"
	}
}

table status:
{
	"request_id": 1,
	"service": "table",
	"resource": "status",
	"data": {
		"table_id": "tid-5feb5ab5100000a"
	}
}

table group (enter):
{
	"request_id": 1,
	"service": "table",
	"resource": "group",
	"data": {
		"table_id": "tid-5feb60e8d00000a",
		"group_id": 1,
		"action": "enter"
	}
}

table group (leave):
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

table group (status):
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