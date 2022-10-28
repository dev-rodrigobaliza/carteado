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
		"table_id": "gid-5fd5c570a0100cb"
	}
}

table enter (private):
{
	"request_id": 1,
	"service": "table",
	"resource": "enter",
	"data": {
		"table_id": "gid-5fd5c570a0100cb",
		"secret": "secret"
	}
}

table leave:
{
	"request_id": 1,
	"service": "table",
	"resource": "leave",
	"data": {
		"table_id": "gid-5fd571f700100cb"
	}
}

table remove:
{
	"request_id": 1,
	"service": "table",
	"resource": "remove",
	"data": {
		"table_id": "gid-5fd571f700100cb"
	}
}

table status:
{
	"request_id": 1,
	"service": "table",
	"resource": "status",
	"data": {
		"table_id": "gid-5fde6d1a300000a"
	}
}